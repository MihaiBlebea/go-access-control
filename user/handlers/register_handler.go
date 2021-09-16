package uhandlers

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/MihaiBlebea/go-access-control/event"
	"github.com/MihaiBlebea/go-access-control/user"
	"github.com/gorilla/context"
)

type RegisterRequest struct {
	FirstName         string `json:"first_name"`
	LastName          string `json:"last_name"`
	Email             string `json:"email"`
	Password          string `json:"password"`
	ConfirmSuccessURL string `json:"confirm_success_url"`
	ConfirmFailURL    string `json:"confirm_fail_url"`
	ConfirmWebhook    string `json:"confirm_webhook"`
}

type RegisterResponse struct {
	ID           int    `json:"id,omitempty"`
	AccessToken  string `json:"access_token,omitempty"`
	RefreshToken string `json:"refresh_token,omitempty"`
	Success      bool   `json:"success"`
	Message      string `json:"message,omitempty"`
}

func RegisterHandler(s user.Service, es event.Service) http.Handler {
	validate := func(r *http.Request) (*RegisterRequest, error) {
		request := RegisterRequest{}

		err := json.NewDecoder(r.Body).Decode(&request)
		if err != nil {
			return &request, err
		}

		if request.FirstName == "" {
			return &request, errors.New("invalid request param first_name")
		}

		if request.LastName == "" {
			return &request, errors.New("invalid request param last_name")
		}

		if request.Email == "" {
			return &request, errors.New("invalid request param email")
		}

		if request.Password == "" {
			return &request, errors.New("invalid request param password")
		}

		return &request, nil
	}

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		response := RegisterResponse{}

		request, err := validate(r)
		if err != nil {
			response.Message = err.Error()
			sendResponse(w, response, http.StatusBadRequest)
			return
		}

		projectID := context.Get(r, "project_id").(int)

		user, err := s.Register(
			projectID,
			request.FirstName,
			request.LastName,
			request.Email,
			request.Password,
			request.ConfirmSuccessURL,
			request.ConfirmFailURL,
			request.ConfirmWebhook,
		)
		if err != nil {
			response.Message = err.Error()
			sendResponse(w, response, http.StatusBadRequest)
			return
		}

		es.StoreEvent(user.ID, "user:registered")

		response.Success = true
		response.ID = user.ID
		response.AccessToken = user.AccessToken
		response.RefreshToken = user.RefreshToken

		sendResponse(w, response, http.StatusOK)
	})
}
