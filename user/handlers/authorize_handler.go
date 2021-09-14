package uhandlers

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/MihaiBlebea/go-access-control/user"
)

type AuthorizeRequest struct {
	Token string `json:"token"`
}

type AuthorizeResponse struct {
	ID      int          `json:"id,omitempty"`
	Token   string       `json:"token,omitempty"`
	User    *UserDetails `json:"user,omitempty"`
	Success bool         `json:"success"`
	Message string       `json:"message,omitempty"`
}

func AuthorizeHandler(s user.Service) http.Handler {
	validate := func(r *http.Request) (*AuthorizeRequest, error) {
		request := AuthorizeRequest{}

		err := json.NewDecoder(r.Body).Decode(&request)
		if err != nil {
			return &request, err
		}

		if request.Token == "" {
			return &request, errors.New("invalid request param token")
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

		user, err := s.Authorize(request.Token)
		if err != nil {
			response.Message = err.Error()
			sendResponse(w, response, http.StatusBadRequest)
			return
		}

		response.Success = true
		response.ID = user.ID
		response.Token = user.Token
		response.User = toUserDetails(user)

		sendResponse(w, response, http.StatusOK)
	})
}
