package uhandlers

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/MihaiBlebea/go-access-control/user"
)

type RegisterRequest struct {
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
	Email     string `json:"email"`
	Password  string `json:"password"`
}

type RegisterResponse struct {
	ID           int    `json:"id,omitempty"`
	AccessToken  string `json:"access_token,omitempty"`
	RefreshToken string `json:"refresh_token,omitempty"`
	Success      bool   `json:"success"`
	Message      string `json:"message,omitempty"`
}

func RegisterHandler(s user.Service) http.Handler {
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

		user, err := s.Register(
			request.FirstName,
			request.LastName,
			request.Email,
			request.Password,
		)
		if err != nil {
			response.Message = err.Error()
			sendResponse(w, response, http.StatusBadRequest)
			return
		}

		response.Success = true
		response.ID = user.ID
		response.AccessToken = user.AccessToken
		response.RefreshToken = user.RefreshToken

		sendResponse(w, response, http.StatusOK)
	})
}
