package discovery

import (
	"github.com/go-openapi/runtime"
	"github.com/go-openapi/strfmt"
	"github.com/iryonetwork/wwm/log/errorChecker"
)

type authWriter struct {
	authToken string
}

func newAuthWriter(authToken string) *authWriter {
	return &authWriter{authToken: authToken}
}

func (a *authWriter) AuthenticateRequest(r runtime.ClientRequest, _ strfmt.Registry) error {
	errorChecker.LogError(r.SetHeaderParam("Authorization", a.authToken))
	return nil
}
