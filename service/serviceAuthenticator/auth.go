// Package serviceAuthenticator provides implementation of runtime.ClientAuthInfoWriter
// interface for internal services using cert and key to authorize its requests to other
// internal services.
package serviceAuthenticator

import (
	"crypto/rsa"
	"crypto/tls"
	"fmt"
	"time"

	jwt "github.com/dgrijalva/jwt-go"
	"github.com/go-openapi/runtime"
	"github.com/go-openapi/strfmt"
	"github.com/rs/zerolog"
	"golang.org/x/crypto/acme"

	"github.com/iryonetwork/wwm/log/errorChecker"
	"github.com/iryonetwork/wwm/service/authenticator"
)

type storageRequestAuthenticator struct {
	pk     *rsa.PrivateKey
	logger zerolog.Logger
}

const servicePrincipal = "__service__"

var tokenExpiersIn = time.Duration(15) * time.Minute

// AuthenticateRequest authenticates API request
func (a *storageRequestAuthenticator) AuthenticateRequest(r runtime.ClientRequest, f strfmt.Registry) error {
	token, err := a.createToken()
	if err != nil {
		a.logger.Error().Err(err).Msg("failed to create token")
		return err
	}

	errorChecker.LogError(r.SetHeaderParam("Authorization", token))
	return nil
}

// New returns new instance of storageRequestAuthenticator that implements runtime.ClientAuthInfoWriter
func New(certFile, keyFile string, logger zerolog.Logger) (runtime.ClientAuthInfoWriter, error) {
	logger = logger.With().Str("component", "service/serviceAuthenticator").Logger()

	logger.Debug().Msg("Initialize service's request authenticator")
	cert, err := tls.LoadX509KeyPair(certFile, keyFile)
	if err != nil {
		return nil, err
	}

	pk, ok := cert.PrivateKey.(*rsa.PrivateKey)
	if !ok {
		return nil, fmt.Errorf("Certificate doesn't contain rsa key")
	}

	return &storageRequestAuthenticator{
		pk:     pk,
		logger: logger,
	}, nil
}

func (a *storageRequestAuthenticator) createToken() (string, error) {
	thumb, err := acme.JWKThumbprint(a.pk.Public())
	if err != nil {
		return "", err
	}

	claims := &authenticator.Claims{
		KeyID: thumb,
		StandardClaims: jwt.StandardClaims{
			Subject:   servicePrincipal + thumb,
			IssuedAt:  time.Now().Unix(),
			ExpiresAt: time.Now().Add(tokenExpiersIn).Unix(),
		},
	}

	// create the token
	return jwt.NewWithClaims(jwt.SigningMethodRS256, claims).SignedString(a.pk)
}
