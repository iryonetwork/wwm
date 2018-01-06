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

// type authService interface {
// 	Login(context.Context, string, string) (string, error)
// 	Validate(context.Context, string, string)
// }

// type handler struct {
// 	svc
// }

// func startServer() error {
// 	h := handler{}

// 	http.HandleFunc("/login", h.login)
// 	http.HandleFunc("/validate", h.validate)

// 	return http.ListenAndServe(":8080", nil)
// }

// func (h *handler) login(w http.ResponseWriter, r *http.Request) {
// 	//
// }

// func (h *handler) validate(w http.ResponseWriter, r *http.Request) {
// 	//
// }
