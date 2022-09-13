package middleware

import (
	"net/http"

	"github.com/Ollub/user_service/internal/session"
)

const AutenticationHeader = "x-authentication-token"

var (
	noAuthUrls = map[string]struct{}{
		"/register": {},
		"/login":    {},
	}
)

func Authentication(sm *session.SessionsJWTVer) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if _, ok := noAuthUrls[r.URL.Path]; ok {
				next.ServeHTTP(w, r)
				return
			}
			token := r.Header.Get(AutenticationHeader)
			if token == "" {
				http.Error(w, "Missing authentication header", http.StatusUnauthorized)
				return
			}
			sess, err := sm.Check(r.Context(), token)
			if err != nil {
				http.Error(w, "No auth", http.StatusUnauthorized)
				return
			}
			ctx := session.ToContext(r.Context(), sess)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}
