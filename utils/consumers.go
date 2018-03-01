package utils

import (
	"mime"

	"github.com/go-openapi/runtime"
)

var (
	JSONMediaType, _, _           = mime.ParseMediaType("application/json; charset=utf-8")
	JSONCollectionMediaType, _, _ = mime.ParseMediaType("application/x-collection+json")
	ByteStreamMediaType, _, _     = mime.ParseMediaType("application/octet-stream")
	XMLOpenEHRMediaType, _, _     = mime.ParseMediaType("text/openEhrXml")
	JSONOpenEHRMediaType, _, _    = mime.ParseMediaType("text/openEhrJson")
	TextMediaType, _, _           = mime.ParseMediaType("text/plain")
)

// Returns set of consumers for runtime HTTP client for storage API to correctly process responses in sync/storage services
func ConsumersForSync() map[string]runtime.Consumer {
	return map[string]runtime.Consumer{
		JSONMediaType:           runtime.JSONConsumer(),
		JSONCollectionMediaType: runtime.ByteStreamConsumer(),
		ByteStreamMediaType:     runtime.ByteStreamConsumer(),
		XMLOpenEHRMediaType:     runtime.ByteStreamConsumer(),
		JSONOpenEHRMediaType:    runtime.ByteStreamConsumer(),
		TextMediaType:           runtime.TextConsumer(),
	}
}
