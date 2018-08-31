package utils

import (
	"bytes"
	"testing"

	"github.com/go-openapi/strfmt"
	"github.com/satori/go.uuid"
)

func TestUUIDToBytes(t *testing.T) {
	uuid, _ := uuid.NewV4()

	in := strfmt.UUID(uuid.String())

	b, err := UUIDToBytes(in)

	if err != nil {
		t.Fatalf("error should be nil; got %v", err)
	}

	if !bytes.Equal(uuid.Bytes(), b) {
		t.Fatalf("bytes should equal")
	}
}
