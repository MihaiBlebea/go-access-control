package uhandlers

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/MihaiBlebea/go-access-control/event"
	"github.com/MihaiBlebea/go-access-control/user"
	"github.com/gorilla/context"
)

type AuthorizeRequest struct {
	AccessToken string `json:"access_token"`
}

type AuthorizeResponse struct {
	ID      int          `json:"id,omitempty"`
	User    *UserDetails `json:"user,omitempty"`
	Success bool         `json:"success"`
	Message string       `json:"message,omitempty"`
}

func AuthorizeHandler(s user.Service, es event.Service) http.Handler {
	validate := func(r *http.Request) (*AuthorizeRequest, error) {
		request := AuthorizeRequest{}

		err := json.NewDecoder(r.Body).Decode(&request)
		if err != nil {
			return &request, err
		}

		if request.AccessToken == "" {
			return &request, errors.New("invalid request param access_token")
		}

		return &request, nil
	}

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		response := AuthorizeResponse{}

		request, err := validate(r)
		if err != nil {
			response.Message = err.Error()
			sendResponse(w, response, http.StatusBadRequest)
			return
		}

		projectID := context.Get(r, "project_id").(int)

		user, err := s.Authorize(projectID, request.AccessToken)
		if err != nil {
			response.Message = err.Error()
			sendResponse(w, response, http.StatusBadRequest)
			return
		}

		es.StoreEvent(user.ID, "user:authorized")

		response.Success = true
		response.ID = user.ID
		response.User = toUserDetails(user)

		sendResponse(w, response, http.StatusOK)
	})
}
