package webapp

import (
	"errors"
	"net/http"

	"github.com/MihaiBlebea/go-access-control/event"
	proj "github.com/MihaiBlebea/go-access-control/project"
	"github.com/MihaiBlebea/go-access-control/user"
	"github.com/gorilla/mux"
)

type User struct {
	FirstName string
	LastName  string
	Email     string
	Confirmed bool
	CreatedAt string
	HtmlID    string
	Events    []Event
}

type Event struct {
	Action    string
	CreatedAt string
}

func ProjectGetHandler(ps proj.Service, us user.Service, es event.Service) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		params := mux.Vars(r)
		slug := params["slug"]
		if slug == "" {
			renderError(w, 500, errors.New("no slug provided"))
			return
		}

		project, err := ps.GetProjectBySlug(slug)
		if err != nil {
			renderError(w, 500, err)
			return
		}

		users, err := us.ProjectUsers(project.ID)
		if err != nil {
			renderError(w, 500, err)
		}

		usrs := []User{}
		for _, u := range users {

			events, err := es.UserEvents(u.ID)
			if err != nil {
				renderError(w, 500, err)
				return
			}

			evs := []Event{}
			for _, e := range events {
				evs = append(evs, Event{
					Action:    e.Action,
					CreatedAt: e.CreatedAt.String(),
				})
			}

			usrs = append(usrs, User{
				FirstName: u.FirstName,
				LastName:  u.LastName,
				Email:     u.Email,
				Confirmed: u.Confirmed,
				CreatedAt: u.CreatedAt.String(),
				HtmlID:    getHtmlID(),
				Events:    evs,
			})
		}

		data := struct {
			ProjectName string
			ApiKey      string
			Users       []User
		}{
			ProjectName: project.Name,
			ApiKey:      project.ApiKey,
			Users:       usrs,
		}

		if err := render(w, "project", &data); err != nil {
			renderError(w, 500, err)
		}
	})
}
