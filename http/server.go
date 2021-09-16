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

	"github.com/MihaiBlebea/go-access-control/event"
	proj "github.com/MihaiBlebea/go-access-control/project"
	"github.com/MihaiBlebea/go-access-control/user"
	uhandlers "github.com/MihaiBlebea/go-access-control/user/handlers"
	"github.com/MihaiBlebea/go-access-control/webapp"
)

const prefix = "/api/v1/"

func New(
	userService user.Service,
	projService proj.Service,
	eventService event.Service,
	logger *logrus.Logger) {

	r := mux.NewRouter()

	api := r.PathPrefix(prefix).Subrouter()

	userApi := api.PathPrefix("/user/").Subrouter()

	// Handle api calls
	api.Handle("/health-check", healthHandler()).
		Methods(http.MethodGet)

	// user endpoints
	userApi.Handle("/login", uhandlers.LoginHandler(userService, eventService)).
		Methods(http.MethodPost)

	userApi.Handle("/register", uhandlers.RegisterHandler(userService, eventService)).
		Methods(http.MethodPost)

	userApi.Handle("/authorize", uhandlers.AuthorizeHandler(userService, eventService)).
		Methods(http.MethodPost)

	userApi.Handle("/refresh", uhandlers.RefreshHandler(userService, eventService)).
		Methods(http.MethodPost)

	userApi.Handle("/remove", uhandlers.RemoveHandler(userService, eventService)).
		Methods(http.MethodDelete)

	r.Handle("/confirm", uhandlers.ConfirmHandler(userService, eventService)).
		Methods(http.MethodGet)

	// project endpoints
	api.Handle("/project", proj.ProjectHandler(projService)).
		Methods(http.MethodPost)

	// webapp endpoints
	r.Handle("/project/{slug}", webapp.ProjectGetHandler(projService, userService, eventService)).
		Methods(http.MethodGet)

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
