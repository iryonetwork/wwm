package filesDataExporter

import (
	"bytes"
	"context"
	"crypto/aes"
	"crypto/cipher"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"os"
	"reflect"
	"testing"

	"github.com/pkg/errors"
	"github.com/rs/zerolog"
)

var (
	testEncryptionKey = []byte{0x47, 0x41, 0xf2, 0xd6, 0xa5, 0xa0, 0x46, 0x33, 0x87, 0xa8, 0x2, 0xa2, 0x52, 0xd0, 0x20, 0xed, 0x8c, 0x65, 0x76, 0x50, 0x77, 0x79, 0x4c, 0x7c, 0x76, 0x49, 0xbf, 0x85, 0x18, 0x27, 0x4f, 0xa0}
	noErrors          = false
	withErrors        = true
)

func TestSanitize(t *testing.T) {
	testCases := []struct {
		description             string
		fieldsToSanitize        []FieldToSanitize
		data                    string
		fieldToDecryptToVerify  []string
		expectedAfterDecryption map[string]string
		errorExpected           bool
	}{
		{
			"Succesfully removed fields, nothing to encrypt",
			[]FieldToSanitize{
				{
					Type:           "value",
					EhrPath:        "field1",
					Transformation: transformationRemove,
				},
				{
					Type:           "value",
					EhrPath:        "field3",
					Transformation: transformationRemove,
				},
				{
					Type:           "value",
					EhrPath:        "field_that_is_not_in_data",
					Transformation: transformationRemove,
				},
			},
			"{\"field1\": \"value1\", \"field2\": \"value2\", \"field3\": \"value3\"}",
			nil,
			map[string]string{
				"field2": "value2",
			},
			noErrors,
		},
		{
			"Succesfully encrypted fields, nothing to remove",
			[]FieldToSanitize{
				{
					Type:           "value",
					EhrPath:        "field1",
					Transformation: transformationEncrypt,
				},
				{
					Type:           "value",
					EhrPath:        "field3",
					Transformation: transformationEncrypt,
				},
				{
					Type:           "value",
					EhrPath:        "field_that_is_not_in_data",
					Transformation: transformationEncrypt,
				},
			},
			"{\"field1\": \"value1\", \"field2\": \"value2\", \"field3\": \"value3\"}",
			[]string{"field1", "field3"},
			map[string]string{
				"field1": "value1",
				"field2": "value2",
				"field3": "value3",
			},
			noErrors,
		},
		{
			"Succesfully removed and encrypted fields",
			[]FieldToSanitize{
				{
					Type:           "value",
					EhrPath:        "field1",
					Transformation: transformationEncrypt,
				},
				{
					Type:           "value",
					EhrPath:        "field1",
					Transformation: transformationRemove,
				},
				{
					Type:           "value",
					EhrPath:        "field3",
					Transformation: transformationEncrypt,
				},
				{
					Type:           "value",
					EhrPath:        "field_that_is_not_in_data",
					Transformation: transformationRemove,
				},
				{
					Type:           "value",
					EhrPath:        "field_that_is_not_in_data_2",
					Transformation: transformationEncrypt,
				},
			},
			"{\"field1\": \"value1\", \"field2\": \"value2\", \"field3\": \"value3\"}",
			[]string{"field3"},
			map[string]string{
				"field2": "value2",
				"field3": "value3",
			},
			noErrors,
		},
		{
			"Succesfully removed and encrypted array fields",
			[]FieldToSanitize{
				{
					Type:    "array",
					EhrPath: "array",
					Properties: []FieldToSanitize{
						{
							Type:           "value",
							EhrPath:        "field1",
							Transformation: transformationEncrypt,
						},
						{
							Type:           "value",
							EhrPath:        "field1",
							Transformation: transformationRemove,
						},
						{
							Type:           "value",
							EhrPath:        "field2",
							Transformation: transformationEncrypt,
						},
						{
							Type:           "value",
							EhrPath:        "field_that_is_not_in_data",
							Transformation: transformationRemove,
						},
						{
							Type:    "array",
							EhrPath: "internalArray",
							Properties: []FieldToSanitize{
								{
									Type:           "value",
									EhrPath:        "field3",
									Transformation: transformationEncrypt,
								},
							},
						},
					},
				},
			},
			"{\"array:0field1\": \"value1\", \"array:0field2\": \"value2\",  \"array:0fieldNotInSpec\": \"value3\", \"array:0internalArray:0field3\": \"value4\", \"array:0internalArray:0fieldNotInSpec2\": \"value5\", \"array:1field1\": \"value1\", \"array:1field2\": \"value2\",  \"array:1fieldNotInSpec\": \"value3\", \"array:1internalArray:0field3\": \"value4\", \"array:1internalArray:0fieldNotInSpec2\": \"value5\"}",
			[]string{"array:0field2", "array:1field2", "array:0internalArray:0field3", "array:1internalArray:0field3"},
			map[string]string{
				"array:0field2":                         "value2",
				"array:0fieldNotInSpec":                 "value3",
				"array:0internalArray:0field3":          "value4",
				"array:0internalArray:0fieldNotInSpec2": "value5",
				"array:1field2":                         "value2",
				"array:1fieldNotInSpec":                 "value3",
				"array:1internalArray:0field3":          "value4",
				"array:1internalArray:0fieldNotInSpec2": "value5",
			},
			noErrors,
		},
		{
			"Succesfully applied substring transformation",
			[]FieldToSanitize{
				{
					Type:           "value",
					EhrPath:        "field1",
					Transformation: transformationSubstring,
					TransformationParameters: map[string]interface{}{
						"start": -1.0,
						"end":   -1.0,
					},
				},
				{
					Type:           "value",
					EhrPath:        "field3",
					Transformation: transformationSubstring,
					TransformationParameters: map[string]interface{}{
						"start": 1.0,
						"end":   -1.0,
					},
				},
				{
					Type:           "value",
					EhrPath:        "field4",
					Transformation: transformationSubstring,
					TransformationParameters: map[string]interface{}{
						"start": 1.0,
						"end":   2.0,
					},
				},
				{
					Type:           "value",
					EhrPath:        "field5",
					Transformation: transformationSubstring,
					TransformationParameters: map[string]interface{}{
						"start": -1.0,
						"end":   2.0,
					},
				},
				{
					Type:           "value",
					EhrPath:        "field_that_is_not_in_data",
					Transformation: transformationSubstring,
				},
			},
			"{\"field1\": \"1234\", \"field2\": \"1234\", \"field3\": \"1234\", \"field4\": \"1234\", \"field5\": \"1234\"}",
			nil,
			map[string]string{
				"field1": "1234",
				"field2": "1234",
				"field3": "234",
				"field4": "2",
				"field5": "12",
			},
			noErrors,
		},
		{
			"Error, invalid input",
			[]FieldToSanitize{
				{
					Type:           "value",
					EhrPath:        "field1",
					Transformation: transformationEncrypt,
				},
				{
					Type:           "value",
					EhrPath:        "field1",
					Transformation: transformationRemove,
				},
				{
					Type:           "value",
					EhrPath:        "field3",
					Transformation: transformationEncrypt,
				},
				{
					Type:           "value",
					EhrPath:        "field_that_is_not_in_data",
					Transformation: transformationRemove,
				},
				{
					Type:           "value",
					EhrPath:        "field_that_is_not_in_data_2",
					Transformation: transformationEncrypt,
				},
			},
			"{\"this\": \"is\", \"not\": \"valid\" \"json\": \"string\"}",
			nil,
			nil,
			withErrors,
		},
	}

	for _, test := range testCases {
		t.Run(test.description, func(t *testing.T) {
			svc := getTestService(t, test.fieldsToSanitize)

			out, err := svc.Sanitize(context.Background(), []byte(test.data))

			// assert error
			if test.errorExpected && err == nil {
				t.Error("Expected error, got nil")
			} else if !test.errorExpected && err != nil {
				t.Errorf("Expected error to be nil, got %v", err)
			}

			var outMap map[string]string

			// check expected results
			if test.expectedAfterDecryption != nil {
				if out == nil {
					t.Fatalf("Expected %s, got nil", toJSON(test.expectedAfterDecryption))
				}

				err = json.Unmarshal(out, &outMap)
				if err != nil {
					t.Fatalf("Error %s on trying to unmarshal returned value %s.", err, out)
				}

				// decrypt returned values so comparison is possible
				for _, fieldToDecrypt := range test.fieldToDecryptToVerify {
					if val, ok := outMap[fieldToDecrypt]; ok {
						decrypted, err := decrypt(val, testEncryptionKey)
						if err != nil {
							t.Fatalf("Error %s on decrypting returned encrypted value %s of field %s.", err, val, fieldToDecrypt)
						}
						outMap[fieldToDecrypt] = string(decrypted)
					}
				}

			}

			if !reflect.DeepEqual(outMap, test.expectedAfterDecryption) {
				t.Errorf("Expected %s, got %s", toJSON(test.expectedAfterDecryption), toJSON(outMap))
			}
		})
	}
}

func getTestService(t *testing.T, fieldsToSanitize []FieldToSanitize) Sanitizer {
	svc, err := NewSanitizer(fieldsToSanitize, testEncryptionKey, zerolog.New(os.Stdout))
	if err != nil {
		t.Fatalf("Failed to initialize sanitizer")
	}

	return svc
}

func decrypt(data string, encryptionKey []byte) ([]byte, error) {
	rawData, err := base64.StdEncoding.DecodeString(data)
	if err != nil {
		return nil, errors.Wrap(err, "Failed to base64 decode encrypted value")
	}

	if len(rawData) < nonceLength {
		return nil, fmt.Errorf("Invalid length of encrypted data")
	}

	block, err := aes.NewCipher(encryptionKey)
	if err != nil {
		return nil, err
	}

	aesgcm, err := cipher.NewGCM(block)
	if err != nil {
		return nil, err
	}

	decrypted, err := aesgcm.Open(nil, rawData[:nonceLength], rawData[nonceLength:], nil)
	if err != nil {
		return nil, err
	}
	if decrypted == nil {
		return []byte{}, nil
	}

	return decrypted, nil
}

func toJSON(in interface{}) string {
	buf := bytes.NewBuffer(nil)
	json.NewEncoder(buf).Encode(in)
	return buf.String()
}
