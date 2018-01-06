package main

import (
	"context"
	"net/http"
	"os"

	"github.com/go-kit/kit/log"
	httptransport "github.com/go-kit/kit/transport/http"
	"github.com/gogo/protobuf/jsonpb"
	"github.com/gogo/protobuf/proto"
	"github.com/iryonetwork/wwm/service/authenticator"
	"github.com/iryonetwork/wwm/specs"
	"github.com/iryonetwork/wwm/storage/auth"
)

func main() {
	logger := log.NewLogfmtLogger(os.Stderr)

	// initialize storage
	storage, err := auth.New()
	if err != nil {
		logger.Log("msg", "Failed to initialize auth storage", "error", err)
		return
	}

	// initialize the service
	authSvc := authenticator.New(storage)

	// setup login handler
	loginEndpoint := makeLoginEndpoint(authSvc)
	loginHandler := httptransport.NewServer(
		loginEndpoint,
		decodeLoginRequest,
		encodeResponse,
		httptransport.ServerErrorEncoder(errorEncoder))
	http.Handle("/login", loginHandler)

	validateEndpoint := makeValidateEndpoint(authSvc)
	validateHandler := httptransport.NewServer(
		validateEndpoint,
		decodeValidateRequest,
		encodeResponse,
		httptransport.ServerErrorEncoder(errorEncoder))
	http.Handle("/validate", validateHandler)

	// start the server
	logger.Log("msg", "HTTP start", "addr", ":80")
	logger.Log("err", http.ListenAndServe(":80", nil))
}

func decodeLoginRequest(_ context.Context, r *http.Request) (interface{}, error) {
	in := &specs.LoginRequest{}
	err := requestToProto(in, r)
	return in, err
}

func decodeValidateRequest(_ context.Context, r *http.Request) (interface{}, error) {
	in := &specs.ValidationRequest{}
	err := requestToProto(in, r)
	return in, err
}

func encodeResponse(_ context.Context, w http.ResponseWriter, response interface{}) error {
	return protoToResponse(w, response.(proto.Message))
}

func requestToProto(in proto.Message, r *http.Request) error {
	return jsonpb.Unmarshal(r.Body, in)
}

func protoToResponse(w http.ResponseWriter, data proto.Message) error {
	m := &jsonpb.Marshaler{}
	return m.Marshal(w, data)
}

func errorEncoder(_ context.Context, err error, w http.ResponseWriter) {
	protoToResponse(w, &specs.ErrorResponse{Error: err.Error()})
}
