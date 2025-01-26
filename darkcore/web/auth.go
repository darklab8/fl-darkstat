package web

import (
	"net/http"
	"time"

	"github.com/darklab8/fl-darkstat/darkcore/settings"
	"github.com/darklab8/fl-darkstat/darkcore/settings/logus"
	"github.com/darklab8/go-typelog/typelog"
)

func AuthMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var query_password string
		var cookie_password string

		// Auth is not active
		if settings.Env.Password == "" {
			logus.Log.Debug("Server is without auth. Passing through")
			next.ServeHTTP(w, r)
			return
		}

		// Auth is Active.
		if password_cookie, err := r.Cookie(settings.Env.Password); err == nil {
			if password_cookie.Value == "true" {
				cookie_password = settings.Env.Password
				logus.Log.Debug("Found password in cookie. acquired.", typelog.Any("password", cookie_password))
			}
		}
		r.URL.Query()
		if password_query := r.URL.Query().Get("password"); password_query != "" {
			query_password = password_query
			logus.Log.Debug("Found password in query param. acquired.", typelog.Any("password", password_query))
		}

		var tempus_token string
		if tempus_cookie, err := r.Cookie("tempus"); err == nil {
			tempus_token = tempus_cookie.Value
		}

		logus.Log.Debug("check auth",
			typelog.Any("query_password", query_password),
			typelog.Any("cookie_password", cookie_password),
			typelog.Any("tempus", tempus_token),
		)
		if query_password == settings.Env.Password || cookie_password == settings.Env.Password || IsTempusValid(tempus_token) {
			if query_password == settings.Env.Password || cookie_password == settings.Env.Password {
				logus.Log.Debug("Valid password. Access Granted")

			}
			if IsTempusValid(tempus_token) {
				logus.Log.Debug("Valid tempus. Access Granted")
			}

			expiration := time.Now().Add(24 * time.Hour)
			cookie := &http.Cookie{Name: settings.Env.Password, Value: "true", Expires: expiration}
			logus.Log.Debug("setting password cookie")
			http.SetCookie(w, cookie)

			tempus_cookie := &http.Cookie{Name: "tempus", Value: NewTempusToken(), Expires: time.Now().Add(1 * time.Hour)}
			logus.Log.Debug("setting tempus cookie")
			http.SetCookie(w, tempus_cookie)

			next.ServeHTTP(w, r)
			return
		}

		http.Error(w, "Password is incorrect", http.StatusForbidden)
		logus.Log.Debug("Password is incorrect. Forbidden")
	})
}
