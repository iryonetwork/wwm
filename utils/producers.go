package utils

import (
	"net/http"

	"github.com/go-openapi/runtime"
	"github.com/go-openapi/runtime/middleware"
)

type producer int8

// Producers
const (
	_ producer = iota
	JSONProducer
	TextProducer
	BinProducer
	FileProducer
)

// UseProducer is used to override which producer will be used in response
func UseProducer(responder middleware.Responder, p producer) middleware.Responder {
	return middleware.ResponderFunc(func(rw http.ResponseWriter, pr runtime.Producer) {
		switch p {
		case JSONProducer:
			rw.Header().Set(runtime.HeaderContentType, "application/json; charset=utf-8")
			responder.WriteResponse(rw, runtime.JSONProducer())

		case TextProducer:
			rw.Header().Set(runtime.HeaderContentType, "text/plain")
			responder.WriteResponse(rw, runtime.TextProducer())

		case BinProducer:
			rw.Header().Set(runtime.HeaderContentType, "application/octet-stream")
			responder.WriteResponse(rw, runtime.ByteStreamProducer())

		case FileProducer:
			// content type is expected to be set already
			responder.WriteResponse(rw, runtime.ByteStreamProducer())

		default:
			responder.WriteResponse(rw, pr)
		}
	})
}
