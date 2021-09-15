package email

import (
	"net/http"

	"github.com/MihaiBlebea/go-access-control/user"
)

func EmailConfirmHandler(s user.Service) http.Handler {

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		// Update the state of the user

		// Send back html template page

		sendResponse(w, http.StatusOK)
	})
}

func sendResponse(w http.ResponseWriter, code int) {
	w.Header().Set("Content-Type", "html/text; charset=utf-8")
	w.Header().Set("X-Content-Type-Options", "nosniff")
	w.WriteHeader(code)

	w.Write([]byte("All good, email confirmed"))
}
