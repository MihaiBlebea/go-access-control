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

	"github.com/MihaiBlebea/go-access-control/user"
	uhandlers "github.com/MihaiBlebea/go-access-control/user/handlers"
)

const prefix = "/api/v1/"

func New(userService user.Service, logger *logrus.Logger) {
	r := mux.NewRouter()

	api := r.PathPrefix(prefix).Subrouter()

	// Handle api calls
	api.Handle("/health-check", healthHandler()).
		Methods(http.MethodGet)

	api.Handle("/login", uhandlers.LoginHandler(userService)).
		Methods(http.MethodPost)

	api.Handle("/register", uhandlers.RegisterHandler(userService)).
		Methods(http.MethodPost)

	api.Handle("/authorize", uhandlers.AuthorizeHandler(userService)).
		Methods(http.MethodPost)

	api.Handle("/refresh", uhandlers.RefreshHandler(userService)).
		Methods(http.MethodPost)

	api.Handle("/remove", uhandlers.RemoveHandler(userService)).
		Methods(http.MethodDelete)

	r.Use(loggerMiddleware(logger))

	srv := &http.Server{
		Handler:      cors.AllowAll().Handler(r),
		Addr:         fmt.Sprintf("0.0.0.0:%s", os.Getenv("HTTP_PORT")),
		WriteTimeout: 15 * time.Second,
		ReadTimeout:  15 * time.Second,
	}

	logger.Info(fmt.Sprintf("Started server on port %s", os.Getenv("HTTP_PORT")))

	log.Fatal(srv.ListenAndServe())
}
