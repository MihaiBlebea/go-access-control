package uhandlers

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/MihaiBlebea/go-access-control/user"
)

type ConfirmRequest struct {
	ConfirmToken string `json:"confirm_token"`
}

type ConfirmResponse struct {
	ID          int          `json:"id,omitempty"`
	AccessToken string       `json:"access_token,omitempty"`
	User        *UserDetails `json:"user,omitempty"`
	Success     bool         `json:"success"`
	Message     string       `json:"message,omitempty"`
}

func ConfirmHandler(s user.Service) http.Handler {
	validate := func(r *http.Request) (*ConfirmRequest, error) {
		request := ConfirmRequest{}

		err := json.NewDecoder(r.Body).Decode(&request)
		if err != nil {
			return &request, err
		}

		if request.ConfirmToken == "" {
			return &request, errors.New("invalid request param confirm_token")
		}

		return &request, nil
	}

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		response := RefreshResponse{}

		request, err := validate(r)
		if err != nil {
			response.Message = err.Error()
			sendResponse(w, response, http.StatusBadRequest)
			return
		}

		user, err := s.ConfirmUser(request.ConfirmToken)
		if err != nil {
			response.Message = err.Error()
			sendResponse(w, response, http.StatusBadRequest)
			return
		}

		response.Success = true
		response.ID = user.ID
		response.User = toUserDetails(user)
		response.AccessToken = user.AccessToken

		sendResponse(w, response, http.StatusOK)
	})
}
