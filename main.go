package main

import (
	"errors"
	"fmt"
	"log"
	"net/http"
	"time"

	jwt "github.com/dgrijalva/jwt-go"
)

var mySigningKey = []byte("supersecretkey")

func getToken(w http.ResponseWriter, r *http.Request) {
	validToken, err := GenerateJWT()
	if err != nil {
		fmt.Println("Failed to generate token")
	}

	fmt.Fprintf(w, validToken)
}

func authroizedPage(w http.ResponseWriter, r *http.Request) {
	err := isAuthorized(w, r)
	if err != nil {
		w.WriteHeader(http.StatusForbidden)
		fmt.Fprintf(w, err.Error())
	}

	fmt.Fprintf(w, "Hello you are authorized!")
	fmt.Println("Endpoint Hit: homePage")
}

// GenerateJWT token
func GenerateJWT() (string, error) {
	token := jwt.New(jwt.SigningMethodHS256)

	claims := token.Claims.(jwt.MapClaims)

	claims["authorized"] = true
	claims["client"] = "Kiril"
	claims["exp"] = time.Now().Add(time.Minute * 30).Unix()

	tokenString, err := token.SignedString(mySigningKey)

	if err != nil {
		err = fmt.Errorf("Something Went Wrong: %s", err.Error())
		return "", err
	}

	return tokenString, nil
}

func isAuthorized(w http.ResponseWriter, r *http.Request) error {
	if r.Header["Token"] != nil {
		token, err := jwt.Parse(r.Header["Token"][0], func(token *jwt.Token) (interface{}, error) {
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, fmt.Errorf("There was an error")
			}
			return mySigningKey, nil
		})

		if err != nil {
			return err
		}

		if token.Valid {
			return nil
		}
	}

	return errors.New("Not Authorized")
}

func handler() http.Handler {
	r := http.NewServeMux()
	r.HandleFunc("/isAuthorized", authroizedPage)
	r.HandleFunc("/getToken", getToken)
	return r
}

func main() {
	err := http.ListenAndServe(":8081", handler())
	if err != nil {
		log.Fatal(err)
	}
}
