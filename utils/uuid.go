package utils

import (
	"encoding/hex"

	"github.com/go-openapi/strfmt"
	uuid "github.com/satori/go.uuid"
)

// UUIDToBytes converts string uuid to bytes
func UUIDToBytes(uuid strfmt.UUID) ([]byte, error) {
	str := string(uuid)
	out := make([]byte, 16)

	part, err := hex.DecodeString(str[0:8])
	if err != nil {
		return nil, err
	}
	copy(out[0:4], part)

	part, err = hex.DecodeString(str[9:13])
	if err != nil {
		return nil, err
	}
	copy(out[4:6], part)

	part, err = hex.DecodeString(str[14:18])
	if err != nil {
		return nil, err
	}
	copy(out[6:8], part)

	part, err = hex.DecodeString(str[19:23])
	if err != nil {
		return nil, err
	}
	copy(out[8:10], part)

	part, err = hex.DecodeString(str[24:])
	if err != nil {
		return nil, err
	}
	copy(out[10:], part)

	return out, nil
}

func NormalizeUUIDString(id string) (string, error) {
	uuid, err := uuid.FromString(id)
	if err != nil {
		return "", err
	}

	return uuid.String(), nil
}
