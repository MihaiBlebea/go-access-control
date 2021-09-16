package uhandlers

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/MihaiBlebea/go-access-control/event"
	"github.com/MihaiBlebea/go-access-control/user"
)

type RefreshRequest struct {
	RefreshToken string `json:"refresh_token"`
}

type RefreshResponse struct {
	ID          int          `json:"id,omitempty"`
	AccessToken string       `json:"access_token,omitempty"`
	User        *UserDetails `json:"user,omitempty"`
	Success     bool         `json:"success"`
	Message     string       `json:"message,omitempty"`
}

func RefreshHandler(s user.Service, es event.Service) http.Handler {
	validate := func(r *http.Request) (*RefreshRequest, error) {
		request := RefreshRequest{}

		err := json.NewDecoder(r.Body).Decode(&request)
		if err != nil {
			return &request, err
		}

		if request.RefreshToken == "" {
			return &request, errors.New("invalid request param refresh_token")
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

		user, err := s.RefreshToken(request.RefreshToken)
		if err != nil {
			response.Message = err.Error()
			sendResponse(w, response, http.StatusBadRequest)
			return
		}

		es.StoreEvent(user.ID, "user:token_refreshed")

		response.Success = true
		response.ID = user.ID
		response.AccessToken = user.AccessToken
		response.User = toUserDetails(user)

		sendResponse(w, response, http.StatusOK)
	})
}
