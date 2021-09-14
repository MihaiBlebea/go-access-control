package uhandlers

import (
	"encoding/json"
	"net/http"

	"github.com/MihaiBlebea/go-access-control/user"
)

type UserDetails struct {
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
	Email     string `json:"email"`
}

func sendResponse(w http.ResponseWriter, resp interface{}, code int) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.Header().Set("X-Content-Type-Options", "nosniff")
	w.WriteHeader(code)

	b, _ := json.Marshal(resp)

	w.Write(b)
}

func toUserDetails(u *user.User) *UserDetails {
	return &UserDetails{u.FirstName, u.LastName, u.Email}
}
