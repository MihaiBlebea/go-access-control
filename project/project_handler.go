package proj

import (
	"encoding/json"
	"errors"
	"net/http"
)

type ProjectRequest struct {
	Name string `json:"name"`
	Host string `json:"host"`
}

type ProjectResponse struct {
	ID      int    `json:"id,omitempty"`
	Name    string `json:"name,omitempty"`
	ApiKey  string `json:"api_key,omitempty"`
	Success bool   `json:"success"`
	Message string `json:"message,omitempty"`
}

func ProjectHandler(s Service) http.Handler {
	validate := func(r *http.Request) (*ProjectRequest, error) {
		request := ProjectRequest{}

		err := json.NewDecoder(r.Body).Decode(&request)
		if err != nil {
			return &request, err
		}

		if request.Name == "" {
			return &request, errors.New("invalid request param name")
		}

		// if request.Host == "" {
		// 	return &request, errors.New("invalid request param password")
		// }

		return &request, nil
	}

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		response := ProjectResponse{}

		request, err := validate(r)
		if err != nil {
			response.Message = err.Error()
			sendResponse(w, response, http.StatusBadRequest)
			return
		}

		p, err := s.CreateProject(request.Name, request.Host)
		if err != nil {
			response.Message = err.Error()
			sendResponse(w, response, http.StatusBadRequest)
			return
		}

		response.Success = true
		response.ID = p.ID
		response.Name = p.Name
		response.ApiKey = p.ApiKey

		sendResponse(w, response, http.StatusOK)
	})
}
