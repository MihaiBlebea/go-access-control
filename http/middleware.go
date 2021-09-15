package http

import (
	"errors"
	"fmt"
	"net/http"
	"strings"

	proj "github.com/MihaiBlebea/go-access-control/project"
	"github.com/gorilla/context"
	"github.com/sirupsen/logrus"
)

func loggerMiddleware(logger *logrus.Logger) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			logger.Info(fmt.Sprintf("Incoming %s request %s", r.Method, r.URL.Path))
			next.ServeHTTP(w, r)
		})
	}
}

func projMiddleware(projectService proj.Service, logger *logrus.Logger) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

			token, err := extractToken(r)
			if err != nil {
				logger.Info(fmt.Sprintf("Error %s request %s: %s", r.Method, r.URL.Path, err.Error()))
				unauthorizedResponse(w)
				return
			}

			id, err := projectService.GetProjectID(token)
			if err != nil {
				logger.Info(fmt.Sprintf("Error %s request %s: %s", r.Method, r.URL.Path, err.Error()))
				unauthorizedResponse(w)
				return
			}

			context.Set(r, "project_id", id)

			next.ServeHTTP(w, r)
		})
	}
}

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
