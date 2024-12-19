package middleware

import (
	"context"
	log "github.com/sirupsen/logrus"
	"net/http"
	"repo-app/pkg/jwt"
	"repo-app/pkg/types"
	"strings"
)

var unprotectedRoutes = map[string]struct{}{
	"/register": {},
	"/login":    {},
}

func Auth(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if _, ok := unprotectedRoutes[r.URL.Path]; ok {
			next.ServeHTTP(w, r)
			return
		}

		authHeader := r.Header.Get("Authorization")
		if authHeader == "" {
			http.Error(w, http.StatusText(http.StatusUnauthorized), http.StatusUnauthorized)
			log.WithField("route", r.RequestURI).Info(http.StatusText(http.StatusUnauthorized))
			return
		}

		token := strings.TrimPrefix(authHeader, "Bearer ")
		log.WithField("token", token).Debug(types.TokenRecieved)

		data, err := jwt.NewJWT(jwt.JWTSecret).Parse(token)
		if err != nil {
			http.Error(w, http.StatusText(http.StatusUnauthorized), http.StatusUnauthorized)
			log.Error("Error parsing token")
			return
		}
		log.WithField("data", data).Debug(types.TokenParsed)

		ctx := context.WithValue(r.Context(), types.UserDataKey, data)
		r = r.WithContext(ctx)

		next.ServeHTTP(w, r)
	})
}
