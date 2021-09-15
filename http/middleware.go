package http

import (
	"fmt"
	"net/http"

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

func getProjectID(projectService proj.Service, logger *logrus.Logger) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

			token, err := extractToken(r)
			if err != nil {
				logger.Info(fmt.Sprintf("Error %s request %s: %s", r.Method, r.URL.Path, err.Error()))
				unauthorizedResponse(w)
				return
			}

			project, err := projectService.GetProject(token)
			if err != nil {
				logger.Info(fmt.Sprintf("Error %s request %s: %s", r.Method, r.URL.Path, err.Error()))
				unauthorizedResponse(w)
				return
			}

			context.Set(r, "project_id", project.ID)

			next.ServeHTTP(w, r)
		})
	}
}
