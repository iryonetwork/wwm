package authorizer

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	jwt "github.com/dgrijalva/jwt-go"
	"github.com/go-openapi/swag"
	"github.com/iryonetwork/wwm/gen/auth/models"
	"github.com/rs/zerolog"
)

func TestGetPrincipalFromToken(t *testing.T) {
	service := New("http://doesnt.matter", zerolog.New(os.Stdout))

	token1, _ := jwt.NewWithClaims(jwt.SigningMethodHS256, &jwt.StandardClaims{Subject: "abc"}).SignedString([]byte("key"))
	token2, _ := jwt.NewWithClaims(jwt.SigningMethodHS256, &jwt.StandardClaims{Subject: ""}).SignedString([]byte("key"))
	token3, _ := jwt.NewWithClaims(jwt.SigningMethodHS256, &jwt.StandardClaims{Subject: "test"}).SignedString([]byte("key"))

	testData := []struct {
		token     string
		principal string
		err       error
	}{
		{
			token:     token1,
			principal: "abc",
			err:       nil,
		},
		{
			token:     "wrong token",
			principal: "",
			err:       fmt.Errorf(ErrInvalidToken),
		},
		{
			token:     token2,
			principal: "",
			err:       fmt.Errorf(ErrInvalidToken),
		},
		{
			token:     token3,
			principal: "test",
			err:       nil,
		},
	}

	for i, test := range testData {
		principal, err := service.GetPrincipalFromToken(test.token)

		if err != nil && test.err == nil {
			t.Errorf("#%d GetPrincipalFromToken(%s) err = %s; expected no error", i, test.token, err)
		} else if err == nil && test.err != nil {
			t.Errorf("#%d GetPrincipalFromToken(%s) err = nil; expected error to be %s", i, test.token, test.err)
		} else if err != nil && test.err != nil && err.Error() != test.err.Error() {
			t.Errorf("#%d GetPrincipalFromToken(%s) err = %s; expected error to be %s", i, test.token, err, test.err)
		}

		if *principal != test.principal {
			t.Errorf("#%d GetPrincipalFromToken(%s) = %s; want %s", i, test.token, *principal, test.principal)
		}
	}
}

func TestAuthorizer(t *testing.T) {
	testData := []struct {
		requestPath   string
		requestMethod string
		responseCode  int
		responseData  string
		err           error
	}{
		{
			requestPath:   "/storage",
			requestMethod: http.MethodPost,
			responseCode:  http.StatusOK,
			responseData:  `[{"result": true}]`,
			err:           nil,
		},
		{
			requestPath:   "/storage",
			requestMethod: http.MethodPost,
			responseCode:  http.StatusOK,
			responseData:  `[{"result": false}]`,
			err:           fmt.Errorf(ErrUnauthorized),
		},
		{
			requestPath:   "/storage",
			requestMethod: http.MethodGet,
			responseCode:  http.StatusInternalServerError,
			responseData:  `{"message": "Server Error", "code": "server_error"}`,
			err:           fmt.Errorf("Server Error"),
		},
		{
			requestPath:   "/storage/other",
			requestMethod: http.MethodPut,
			responseCode:  http.StatusInternalServerError,
			responseData:  `{"message": "Server Error", "code": "server_err`,
			err:           fmt.Errorf(`invalid character '\n' in string literal`),
		},
		{
			requestPath:   "/something/else",
			requestMethod: http.MethodGet,
			responseCode:  http.StatusOK,
			responseData:  `[{"test": false}]`,
			err:           fmt.Errorf(ErrUnauthorized),
		},
		{
			requestPath:   "/something/else",
			requestMethod: http.MethodGet,
			responseCode:  http.StatusOK,
			responseData:  `completely wrong`,
			err:           fmt.Errorf(`invalid character 'c' looking for beginning of value`),
		},
	}

	for i, test := range testData {
		authToken := fmt.Sprintf("token%d", i)

		ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			defer r.Body.Close()
			body, _ := ioutil.ReadAll(r.Body)

			pairs := models.PostValidateParamsBody{}
			swag.ReadJSON(body, &pairs)

			if *pairs[0].Actions != methodToAction(test.requestMethod) {
				t.Errorf("#%d Authorize(method: %s, path: %s) action = %d; expected %d", i, test.requestMethod, test.requestPath, *pairs[0].Actions, methodToAction(test.requestMethod))
			}

			if *pairs[0].Resource != "/api"+test.requestPath {
				t.Errorf("#%d Authorize(method: %s, path: %s) path = %s; expected /api%s", i, test.requestMethod, test.requestPath, *pairs[0].Resource, test.requestPath)
			}

			w.Header().Set("Content-Type", "application/json")
			if r.Header.Get("Authorization") != authToken {
				err := models.Error{
					Message: "Token is wrong",
					Code:    "unauthorized",
				}

				w.WriteHeader(http.StatusUnauthorized)
				body, _ := err.MarshalBinary()
				w.Write(body)
				return
			}

			w.WriteHeader(test.responseCode)
			fmt.Fprintln(w, test.responseData)
		}))
		defer ts.Close()

		authorizer := New(ts.URL, zerolog.New(os.Stdout)).Authorizer()

		req, _ := http.NewRequest(test.requestMethod, test.requestPath, nil)
		req.Header.Add("Authorization", authToken)
		err := authorizer.Authorize(req, nil)

		if err != nil && test.err == nil {
			t.Errorf("#%d Authorize(method: %s, path: %s) err = %s; expected no error", i, test.requestMethod, test.requestPath, err)
		} else if err == nil && test.err != nil {
			t.Errorf("#%d Authorize(method: %s, path: %s) err = nil; expected error to be %s", i, test.requestMethod, test.requestPath, test.err)
		} else if err != nil && test.err != nil && err.Error() != test.err.Error() {
			t.Errorf("#%d Authorize(method: %s, path: %s) err = %s; expected error to be %s", i, test.requestMethod, test.requestPath, err, test.err)
		}
	}

}
