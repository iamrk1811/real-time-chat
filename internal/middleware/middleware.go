package middleware

import (
	"context"
	"net/http"

	"github.com/iamrk1811/real-time-chat/config"
	"github.com/iamrk1811/real-time-chat/internal/repo"
	"github.com/iamrk1811/real-time-chat/utils"
)

func CorsMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// TODO allow orgin will dynamic
		origin := r.Header.Get("origin")
		w.Header().Set("Access-Control-Allow-Origin", origin)
		w.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE")
		w.Header().Set("Access-Control-Allow-Headers", "X-Requested-With, X-HTTP-Method-Override, Content-Type, Accept		, Authorization")
		w.Header().Set("Access-Control-Allow-Credentials", "true")

		// If this is a preflight request, then stop here
		if r.Method == "OPTIONS" {
			w.WriteHeader(http.StatusOK)
			return
		}
		next.ServeHTTP(w, r)
	})
}

func UserProtectionMiddleware(next http.Handler, repo *repo.CRUDRepo) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		authCookie, err := r.Cookie(string(config.SessionKey))
		if err != nil {
			utils.WriteResponse(w, http.StatusUnauthorized, "You don't have access to this resource", nil)
			return
		}

		sessionID := authCookie.Value
		session := repo.FetchUserBySessionID(r.Context(), sessionID)
		if session == nil {
			utils.WriteResponse(w, http.StatusUnauthorized, "You don't have access to this resource", nil)
			return
		}

		ctx := context.WithValue(r.Context(), config.SessionKey, session)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
