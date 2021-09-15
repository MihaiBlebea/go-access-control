package http

import (
	"errors"
	"net/http"
	"strings"
)

func extractToken(r *http.Request) (string, error) {
	reqToken := r.Header.Get("Authorization")
	if reqToken == "" {
		return "", errors.New("bearer token is invalid")
	}

	splitToken := strings.Split(reqToken, "Bearer ")
	if len(splitToken) < 2 {
		return "", errors.New("bearer token is invalid")
	}

	reqToken = splitToken[1]

	if reqToken == "" {
		return "", errors.New("bearer token is invalid")
	}

	return reqToken, nil
}

func unauthorizedResponse(w http.ResponseWriter) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.Header().Set("X-Content-Type-Options", "nosniff")
	w.WriteHeader(401)

	w.Write([]byte{})
}
