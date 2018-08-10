package filesDataExporter

import (
	"context"
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"strconv"

	"github.com/rs/zerolog"
)

type sanitizer struct {
	fieldsToSanitize []FieldToSanitize
	encryptionKey    []byte
	aesgcm           cipher.AEAD
	logger           zerolog.Logger
}

const transformationRemove = "remove"
const transformationEncrypt = "encrypt"
const transformationSubstring = "substring"

const nonceLength = 12

// Sanitize sanitizes JSON string by encrypting values and/or removing certain keys
func (s *sanitizer) Sanitize(ctx context.Context, data []byte) ([]byte, error) {
	if len(s.fieldsToSanitize) == 0 {
		// there is nothing to sanitize, skip unmarshalling for performance
		return data, nil
	}

	var dataMap map[string]interface{}

	if err := json.Unmarshal(data, &dataMap); err != nil {
		return nil, err
	}

	_, err := s.sanitize(ctx, s.fieldsToSanitize, &dataMap, "")
	if err != nil {
		return nil, err
	}

	out, err := json.Marshal(dataMap)
	if err != nil {
		return nil, err
	}
	fmt.Println(string(out))

	return out, nil
}

func (s *sanitizer) sanitize(ctx context.Context, fieldsToSanitize []FieldToSanitize, data *map[string]interface{}, prefix string) (fieldFound bool, err error) {
	fieldFound = false

	for _, field := range fieldsToSanitize {
		switch field.Type {
		case "array":
			for i := 0; true; i++ {
				ff, err := s.sanitize(ctx, field.Properties, data, fmt.Sprintf("%s%s:%d", prefix, field.EhrPath, i))
				if err != nil {
					return fieldFound, err
				}
				if ff == false {
					break
				} else {
					fieldFound = true
				}
			}
		default:
			value, ok := (*data)[fmt.Sprintf("%s%s", prefix, field.EhrPath)]
			if ok {
				fieldFound = true
				switch field.Transformation {
				case transformationRemove:
					delete(*data, fmt.Sprintf("%s%s", prefix, field.EhrPath))
				case transformationEncrypt:
					transformed, err := s.encrypt(value)
					if err != nil {
						return fieldFound, err
					}

					(*data)[fmt.Sprintf("%s%s", prefix, field.EhrPath)] = transformed
				case transformationSubstring:
					transformed, err := s.substring(value, int(field.TransformationParameters["start"].(float64)), int(field.TransformationParameters["end"].(float64)))
					if err != nil {
						return fieldFound, err
					}

					(*data)[fmt.Sprintf("%s%s", prefix, field.EhrPath)] = transformed
				}
			}
		}
	}
	return fieldFound, nil
}

func (s *sanitizer) encrypt(value interface{}) (string, error) {
	// data is expected to be flat JSON with string, boolean & number values
	var stringToEncrypt string
	switch t := value.(type) {
	case string:
		stringToEncrypt = value.(string)
	case int:
		stringToEncrypt = strconv.Itoa(value.(int))
	case float32:
		stringToEncrypt = strconv.FormatFloat(float64(value.(float32)), 'E', -1, 64)
	case float64:
		stringToEncrypt = strconv.FormatFloat(value.(float64), 'E', -1, 64)
	case bool:
		stringToEncrypt = strconv.FormatBool(value.(bool))
	default:
		return "", fmt.Errorf("Invalid type %s of data field meant to be encrypted", t)
	}

	nonce := make([]byte, nonceLength)
	if _, err := io.ReadFull(rand.Reader, nonce); err != nil {
		return "", err
	}

	encrypted :=
		base64.StdEncoding.EncodeToString(append(nonce, s.aesgcm.Seal(nil, nonce, []byte(stringToEncrypt), nil)...))

	return encrypted, nil
}

func (s *sanitizer) substring(value interface{}, start, end int) (string, error) {
	stringValue, ok := value.(string)
	if !ok {
		return "", fmt.Errorf("Invalid type of data field meant for extractSlice transformation")
	}

	if start != -1 && end != -1 {
		return string([]rune(stringValue)[start:end]), nil
	} else if start == -1 && end != -1 {
		return string([]rune(stringValue)[:end]), nil
	} else if start != -1 && end == -1 {
		return string([]rune(stringValue)[start:]), nil
	}

	return stringValue, nil
}

func NewSanitizer(fieldsToSanitize []FieldToSanitize, encryptionKey []byte, logger zerolog.Logger) (Sanitizer, error) {
	if len(encryptionKey) != 32 {
		return nil, fmt.Errorf("Encryption key must be 32 bytes long")
	}

	logger = logger.With().Str("component", "reports/filesDataExporter/sanitizer").Logger()

	block, err := aes.NewCipher(encryptionKey)
	if err != nil {
		return nil, err
	}

	aesgcm, err := cipher.NewGCM(block)
	if err != nil {
		return nil, err
	}

	return &sanitizer{
		fieldsToSanitize: fieldsToSanitize,
		encryptionKey:    encryptionKey,
		aesgcm:           aesgcm,
		logger:           logger,
	}, nil
}
