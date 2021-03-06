package main

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestGenerateJWT(t *testing.T) {
	token, err := GenerateJWT()
	if err != nil {
		t.Fatalf("Test Failed: Cannot generate JWT token: %v", err)
	}

	if len(token) != 152 {
		t.Fatalf("Test Failed: Incorrect token generated %v", len(token))
	}
}

func TestRouting(t *testing.T) {
	srv := httptest.NewServer(handler())
	defer srv.Close()

	tt := []struct {
		name     string
		endpoint string
		status   interface{}
		isToken  bool
	}{
		{name: "Test getToken", endpoint: "getToken", status: http.StatusOK, isToken: true},
		{name: "Test isAuthorized", endpoint: "isAuthorized", status: http.StatusForbidden},
	}

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			res, err := http.Get(fmt.Sprintf("%s/"+tc.endpoint, srv.URL))

			if err != nil {
				t.Fatalf("could not set GET request: %v", err)
			}

			defer res.Body.Close()

			b, err := ioutil.ReadAll(res.Body)

			if err != nil {
				t.Fatalf("could not read response: %v", err)
			}

			if tc.isToken == true {
				if len(b) != 152 {
					t.Fatalf("Incorrect token generated got %v", string(b))
				}
			}

			if res.StatusCode != tc.status {
				t.Fatalf("Received incorrect status, expected %v got %v", tc.status, res.Status)
			}
		})
	}
}

func Test_isAuthorized(t *testing.T) {
	token, err := GenerateJWT()
	if err != nil {
		t.Fatalf("Could not generate token %v", err)
	}

	tt := []struct {
		name   string
		token  string
		result string
	}{
		{name: "valid token", token: token, result: ""},
		{name: "invalid token", token: "invalidtoken", result: "token contains an invalid number of segments"},
		{name: "empty token", token: "", result: "token contains an invalid number of segments"},
	}

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			req := httptest.NewRequest("GET", "https://localhost:8080", nil)
			req.Header.Set("Token", tc.token)
			w := httptest.NewRecorder()

			err = isAuthorized(w, req)

			if err != nil && err.Error() != tc.result {
				t.Fatalf("Could not authorize because %v", err)
			}
		})
	}
}
