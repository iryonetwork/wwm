// Package authorizer priovides functions that will be used for authorizing users.
//
// To use it you must first initialize the service:
//  auth := authorizer.New("https://localAuth/auth/validate", logger)
//
// and then you can use its methods for your API:
//  api.TokenAuth = auth.GetPrincipalFromToken
//  api.APIAuthorizer = auth.Authorizer()
package authorizer

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"net/http"
	"time"

	jwt "github.com/dgrijalva/jwt-go"
	"github.com/go-openapi/runtime"
	"github.com/go-openapi/swag"
	"github.com/iryonetwork/wwm/gen/auth/models"
	"github.com/iryonetwork/wwm/service/authenticator"
	"github.com/rs/zerolog"
)

// Service describes the actions supported by the authorizer service
type Service interface {
	// Authorizer returns function that checks if logged in user has permission to do a request
	Authorizer() runtime.Authorizer

	// GetPrincipalFromToken returns user ID parsed from token
	GetPrincipalFromToken(tokenString string) (*string, error)
}

type authorizer struct {
	validateURL string
	logger      zerolog.Logger
}

// New returns new authorizer service
func New(validateURL string, logger zerolog.Logger) Service {
	return &authorizer{
		validateURL: validateURL,
		logger:      logger.With().Str("component", "service/authorizer").Logger(),
	}
}

// GetPrincipalFromToken returns principal parsed from token
// It DOES NOT check if token is properly signed
// Signing will be checked in authenticator service
func (a *authorizer) GetPrincipalFromToken(tokenString string) (*string, error) {
	principal := ""

	jwt.ParseWithClaims(tokenString, &authenticator.Claims{}, func(token *jwt.Token) (interface{}, error) {
		claims, ok := token.Claims.(*authenticator.Claims)
		if ok {
			principal = claims.Subject
		}
		return "", nil
	})

	if principal == "" {
		a.logger.Error().Str("cmd", "GetPrincipalFromToken").Msg("Token is invalid")
		return &principal, fmt.Errorf(ErrInvalidToken)
	}

	return &principal, nil
}

// Actions
const (
	Read   = 1
	Write  = 1 << 1
	Delete = 1 << 2
)

// Errors
const (
	ErrInvalidToken = "Token is invalid"
	ErrUnauthorized = "Unauthorized"
)

type responseAndError struct {
	r   *http.Response
	err error
}

// Authorizer checks if logged in user has permission to do a request
func (a *authorizer) Authorizer() runtime.Authorizer {
	logger := a.logger.With().Str("cmd", "Authorizer").Logger()
	return runtime.AuthorizerFunc(func(request *http.Request, principal interface{}) error {
		action := methodToAction(request.Method)
		resource := "/api" + request.URL.EscapedPath()
		pairs := models.PostValidateParamsBody{
			&models.ValidationPair{
				Actions:  &action,
				Resource: &resource,
			},
		}

		body, err := swag.WriteJSON(pairs)
		if err != nil {
			logger.Error().Err(err).Msg("WriteJSON failed")
			return err
		}

		r, err := http.NewRequest(http.MethodPost, a.validateURL, bytes.NewBuffer(body))
		if err != nil {
			logger.Error().Err(err).Msg("Initializing request failed")
			return err
		}
		r.Header.Add("Authorization", request.Header.Get("Authorization"))
		r.Header.Add("Content-Type", "application/json")

		transport := &http.Transport{}
		netClient := &http.Client{
			Transport: transport,
			Timeout:   time.Second * 10,
		}

		c := make(chan responseAndError)
		go func() {
			response, err := netClient.Do(r)
			c <- responseAndError{response, err}
		}()

		var response responseAndError

		select {
		case <-request.Context().Done():
			transport.CancelRequest(r)
			<-c // wait for canceld request
			logger.Error().Err(request.Context().Err()).Msg("Context was done")
			return fmt.Errorf("Context was done")
		case response = <-c:
		}

		if response.err != nil {
			logger.Error().Err(response.err).Msg("Making request failed")
			return response.err
		}
		defer response.r.Body.Close()

		responseBody, err := ioutil.ReadAll(response.r.Body)
		if err != nil {
			logger.Error().Err(err).Msg("Reading response failed")
			return err
		}

		if response.r.StatusCode == http.StatusOK {
			validationResponse := models.PostValidateOKBody{}
			err := swag.ReadJSON(responseBody, &validationResponse)
			if err != nil {
				logger.Error().Err(err).Msg("Parsing response failed")
				return err
			}

			if !validationResponse[0].Result {
				logger.Debug().Msg(ErrUnauthorized)
				return fmt.Errorf(ErrUnauthorized)
			}

			logger.Debug().Msg("Authorized successfully")
			return nil
		}

		jsonError := &models.Error{}
		err = jsonError.UnmarshalBinary(responseBody)
		if err != nil {
			logger.Error().Err(err).Msg("Parsing error response failed")
			return err
		}

		return fmt.Errorf(jsonError.Message)
	})
}

func methodToAction(method string) int64 {
	switch method {
	case http.MethodDelete:
		return Delete
	case http.MethodPost:
	case http.MethodPut:
		return Write
	}

	return Read
}
