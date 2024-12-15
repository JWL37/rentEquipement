package middleware

import (
	"log"
	"net/http"
	"rentEquipement/internal/session"
)

type AuthMiddleware struct {
	SM *session.SessionsManager
}

func (am *AuthMiddleware) Auth(next http.Handler) http.Handler {

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log.Println("auth middleware", r.URL.Path)

		sess, err := am.SM.Check(r)
		if err != nil {
			log.Println("no auth", r.URL.Path, r.Method)

			http.Redirect(w, r, "/", http.StatusFound)
			return
		}
		ctx := session.ContextWithSession(r.Context(), sess)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
