package http

import (
	"fmt"
	"log"

	"net/http"
	"os"
	"time"

	"github.com/gorilla/mux"
	"github.com/rs/cors"
	"github.com/sirupsen/logrus"

	proj "github.com/MihaiBlebea/go-access-control/project"
	"github.com/MihaiBlebea/go-access-control/user"
	uhandlers "github.com/MihaiBlebea/go-access-control/user/handlers"
)

const prefix = "/api/v1/"

func New(userService user.Service, projService proj.Service, logger *logrus.Logger) {
	r := mux.NewRouter()

	api := r.PathPrefix(prefix).Subrouter()

	userApi := api.PathPrefix("/user/").Subrouter()

	// Handle api calls
	api.Handle("/health-check", healthHandler()).
		Methods(http.MethodGet)

	// user endpoints
	userApi.Handle("/login", uhandlers.LoginHandler(userService)).
		Methods(http.MethodPost)

	userApi.Handle("/register", uhandlers.RegisterHandler(userService)).
		Methods(http.MethodPost)

	userApi.Handle("/authorize", uhandlers.AuthorizeHandler(userService)).
		Methods(http.MethodPost)

	userApi.Handle("/refresh", uhandlers.RefreshHandler(userService)).
		Methods(http.MethodPost)

	userApi.Handle("/remove", uhandlers.RemoveHandler(userService)).
		Methods(http.MethodDelete)

	r.Handle("/confirm", uhandlers.ConfirmHandler(userService)).
		Methods(http.MethodGet)

	// project endpoints
	api.Handle("/project", proj.ProjectHandler(projService)).
		Methods(http.MethodPost)

	r.Use(loggerMiddleware(logger))

	userApi.Use(getProjectID(projService, logger))

	srv := &http.Server{
		Handler:      cors.AllowAll().Handler(r),
		Addr:         fmt.Sprintf("0.0.0.0:%s", os.Getenv("HTTP_PORT")),
		WriteTimeout: 15 * time.Second,
		ReadTimeout:  15 * time.Second,
	}

	logger.Info(fmt.Sprintf("Started server on port %s", os.Getenv("HTTP_PORT")))

	log.Fatal(srv.ListenAndServe())
}
