package uhandlers

import (
	"net/http"

	"github.com/MihaiBlebea/go-access-control/user"
)

func ConfirmHandler(s user.Service) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		token := r.URL.Query().Get("token")
		if token == "" {
			return
		}

		u, claims, err := s.ConfirmUser(token)
		if err != nil {
			user.SendWebhook(claims.ConfirmWebhook, &user.Payload{Success: false, ID: u.ID})
			redirect(w, r, claims.ConfirmFailURL)
			return
		}

		user.SendWebhook(claims.ConfirmWebhook, &user.Payload{Success: true, ID: u.ID})
		redirect(w, r, claims.ConfirmSuccessURL)
	})
}

func redirect(w http.ResponseWriter, r *http.Request, redirectTo string) {
	if redirectTo == "" {
		return
	}

	http.Redirect(w, r, redirectTo, http.StatusSeeOther)
}
