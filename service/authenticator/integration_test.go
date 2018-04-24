package authenticator

import (
	"crypto/ecdsa"
	"crypto/tls"
	"net/http"
	"os"
	"testing"
	"time"

	jwt "github.com/dgrijalva/jwt-go"
	"golang.org/x/crypto/acme"
)

// if you want to run this test you must have API up and running (make up)
// run this test with INTEGRATION=1 go test github.com/iryonetwork/wwm/service/authenticator -run TestSyncAuthentication
// you will have to modify test cases (firs two should work, other two probably not)
func TestSyncAuthentication(t *testing.T) {
	_, shouldRun := os.LookupEnv("INTEGRATION")
	if shouldRun {

		tests := []struct {
			cert   string
			key    string
			url    string
			status int
		}{
			{
				cert:   "../../bin/tls/localAuthSync.pem",
				key:    "../../bin/tls/localAuthSync-key.pem",
				url:    "https://iryo.cloud/auth/renew",
				status: 403,
			},
			{
				cert:   "../../bin/tls/localAuthSync.pem",
				key:    "../../bin/tls/localAuthSync-key.pem",
				url:    "https://iryo.cloud/auth/database",
				status: 200,
			},
			{
				cert:   "../../bin/tls/localAuthSync.pem",
				key:    "../../bin/tls/localAuthSync-key.pem",
				url:    "https://iryo.local/waitlist",
				status: 200,
			},
			{
				cert:   "../../bin/tls/localAuthSync.pem",
				key:    "../../bin/tls/localAuthSync-key.pem",
				url:    "https://iryo.local/waitlist/e0fd94de-5203-4556-b775-01a49a8bedec",
				status: 403,
			},
		}

		for _, test := range tests {

			cert, err := tls.LoadX509KeyPair(test.cert, test.key)
			if err != nil {
				t.Fatal(err)
			}

			pk, _ := cert.PrivateKey.(*ecdsa.PrivateKey)

			thumb, err := acme.JWKThumbprint(pk.Public())
			if err != nil {
				t.Fatal(err)
			}

			claims := &Claims{
				KeyID: thumb,
				StandardClaims: jwt.StandardClaims{
					Subject:   "sync",
					IssuedAt:  time.Now().Unix(),
					ExpiresAt: time.Now().Add(tokenExpiersIn).Unix(),
				},
			}

			// create the token
			token, _ := jwt.NewWithClaims(jwt.SigningMethodRS256, claims).SignedString(pk)

			request, err := http.NewRequest(http.MethodGet, test.url, nil)
			request.Header.Add("Authorization", token)

			tr := &http.Transport{
				TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
			}

			netClient := &http.Client{
				Transport: tr,
				Timeout:   time.Second * 10,
			}
			response, err := netClient.Do(request)
			if err != nil {
				t.Fatal(err)
			}

			if response.StatusCode != test.status {
				t.Fatalf("Expected status %d; got %d", test.status, response.StatusCode)
			}
		}
	}
}
