package uhandlers

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/MihaiBlebea/go-access-control/user"
)

type RemoveRequest struct {
	AccessToken string `json:"access_token"`
}

type RemoveResponse struct {
	ID      int    `json:"id,omitempty"`
	Success bool   `json:"success"`
	Message string `json:"message,omitempty"`
}

func RemoveHandler(s user.Service) http.Handler {
	validate := func(r *http.Request) (*RemoveRequest, error) {
		request := RemoveRequest{}

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
		response := RefreshResponse{}

		request, err := validate(r)
		if err != nil {
			response.Message = err.Error()
			sendResponse(w, response, http.StatusBadRequest)
			return
		}

		id, err := s.RemoveUser(request.AccessToken)
		if err != nil {
			response.Message = err.Error()
			sendResponse(w, response, http.StatusBadRequest)
			return
		}

		response.Success = true
		response.ID = id

		sendResponse(w, response, http.StatusOK)
	})
}
