package utils

import (
	"net/http"

	"github.com/go-openapi/runtime"
	"github.com/go-openapi/runtime/middleware"
)

type producers int8

// Producers
const (
	_ producers = iota
	JSONProducer
	TextProducer
	BinProducer
)

// UseProducer is used to override which producer will be used in response
func UseProducer(responder middleware.Responder, p producers) middleware.Responder {
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

		default:
			responder.WriteResponse(rw, pr)
		}
	})
}
