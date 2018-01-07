package main

import (
	"context"

	"github.com/go-kit/kit/endpoint"
	"github.com/iryonetwork/wwm/service/authenticator"
	"github.com/iryonetwork/wwm/specs"
)

func makeLoginEndpoint(svc authenticator.Service) endpoint.Endpoint {
	return func(ctx context.Context, req interface{}) (interface{}, error) {
		r := req.(*specs.LoginRequest)
		token, err := svc.Login(ctx, r.Username, r.Password)
		return &specs.LoginResponse{Token: token}, err
	}
}

func makeValidateEndpoint(svc authenticator.Service) endpoint.Endpoint {
	return func(ctx context.Context, req interface{}) (interface{}, error) {
		r := req.(*specs.ValidationRequest)
		results, err := svc.Validate(ctx, r.Queries)
		return &specs.ValidationResponse{Results: results}, err
	}
}

func makeKeysEndpoint(svc authenticator.Service) endpoint.Endpoint {
	return func(ctx context.Context, req interface{}) (interface{}, error) {
		r := req.(*specs.KeyRequest)
		result, err := svc.GetPublicKey(ctx, r.KeyID)
		return &specs.KeyResponse{Key: result}, err
	}
}
