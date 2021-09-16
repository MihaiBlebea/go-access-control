package uhandlers

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/MihaiBlebea/go-access-control/event"
	"github.com/MihaiBlebea/go-access-control/user"
)

type LoginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type LoginResponse struct {
	ID           int    `json:"id,omitempty"`
	AccessToken  string `json:"access_token,omitempty"`
	RefreshToken string `json:"refresh_token,omitempty"`
	Success      bool   `json:"success"`
	Message      string `json:"message,omitempty"`
}

func LoginHandler(s user.Service, es event.Service) http.Handler {
	validate := func(r *http.Request) (*LoginRequest, error) {
		request := LoginRequest{}

		err := json.NewDecoder(r.Body).Decode(&request)
		if err != nil {
			return &request, err
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
		response := LoginResponse{}

		request, err := validate(r)
		if err != nil {
			response.Message = err.Error()
			sendResponse(w, response, http.StatusBadRequest)
			return
		}

		user, err := s.Login(request.Email, request.Password)
		if err != nil {
			response.Message = err.Error()
			sendResponse(w, response, http.StatusBadRequest)
			return
		}

		es.StoreEvent(user.ID, "user:loggedin")

		response.Success = true
		response.ID = user.ID
		response.AccessToken = user.AccessToken
		response.RefreshToken = user.RefreshToken

		sendResponse(w, response, http.StatusOK)
	})
}
