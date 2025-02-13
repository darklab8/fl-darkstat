package web

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"time"

	"github.com/darklab8/fl-darkstat/darkcore/settings"
	"github.com/darklab8/fl-darkstat/darkcore/settings/logus"
	"github.com/darklab8/fl-darkstat/darkcore/web/registry"
	statsettings "github.com/darklab8/fl-darkstat/darkstat/settings"
	"github.com/darklab8/go-typelog/typelog"
)

// I have now implemented roughly this
// You redirect unauthenticated users to https://discoverygc.com/forums/oauth/?client_id=darkstat_dev&redirect_url=https://darkstat-dev.dd84ai.com/oauth/redirect
// I check that they're in the right group, and redirect them to https://darkstat-dev.dd84ai.com/oauth/redirect?code=${some_assigned_token}
// You POST {"code":"${code}"} to https://discoverygc.com/forums/oauth/access_token
// I return either {"error": "nope"} or {"access_token": "yes they're a dev"}
// if you get access_token "yes they're a dev" you can grant them a session cookie permitting them to view info from the private development repo
// let me know if any of this is wrong

func NewOauthStart(w *Web) *registry.Endpoint {
	return &registry.Endpoint{
		Url: "GET /oauth",
		Handler: func(w http.ResponseWriter, r *http.Request) {
			redirect_url := fmt.Sprintf("https://discoverygc.com/forums/oauth/?client_id=darkstat_dev&redirect_url=%s/oauth/redirect", statsettings.Env.SiteUrl)
			logus.Log.Info("oauth started", typelog.String("redirect_url", redirect_url))
			http.Redirect(w, r, redirect_url, http.StatusSeeOther)
		},
	}
}

// for local testing, manually replace incoming incoming redirect to http://localhost:8000
func NewOauthAccept(w *Web) *registry.Endpoint {
	return &registry.Endpoint{
		Url: "GET /oauth/redirect",
		Handler: func(w http.ResponseWriter, r *http.Request) {

			r.URL.Query()
			var oauth_code string
			if code := r.URL.Query().Get("code"); code != "" {
				oauth_code = code
				logus.Log.Debug("Found code in query param. acquired.", typelog.Any("oath_code", code))
			}

			is_dev, err := validateCode(oauth_code)

			if !is_dev {
				fmt.Fprintf(w, fmt.Sprintln("failed oauth procedure.", err))
			}

			tempus_value := NewTempusToken()
			tempus_cookie := &http.Cookie{Name: "tempus", Value: tempus_value, Expires: time.Now().Add(1 * time.Hour)}
			fmt.Println("setting tempus cookie for succesful oauth login", "host=", r.Host)
			http.SetCookie(w, tempus_cookie)

			tempus_as_query_param := ""
			if settings.Env.IsDevEnv {
				tempus_as_query_param = fmt.Sprintf("/?tempus=%s", tempus_value)
			}

			// http.Redirect(w, r, statsettings.Env.SiteUrl, http.StatusSeeOther)
			// redirect with delay instead
			buf := bytes.NewBuffer([]byte{})
			RedirectPage(
				"Succesfully oauth authentificated, u will be redirected in 3 seconds to main darkstat page",
				statsettings.Env.SiteUrl+tempus_as_query_param).Render(context.Background(), buf)
			fmt.Fprint(w, buf.String())
		},
	}
}

// I return either {"error": "nope"} or {"access_token": "yes they're a dev"}
type OauthAnswer struct {
	Error       *string `json:"error"`
	AccessToken *string `json:"access_token"`
}

func validateCode(code string) (bool, error) {
	var err error
	var client = &http.Client{}
	var answer OauthAnswer
	var param = url.Values{}
	param.Set("code", code)
	var payload = bytes.NewBufferString(param.Encode())
	request, err := http.NewRequest("POST", "https://discoverygc.com/forums/oauth/access_token", payload)
	if err != nil {
		return false, err
	}
	request.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	response, err := client.Do(request)
	if err != nil {
		return false, err
	}
	defer response.Body.Close()
	err = json.NewDecoder(response.Body).Decode(&answer)
	if err != nil {
		return false, err
	}
	logus.Log.Info("got answer", typelog.Struct(answer))

	if answer.AccessToken == nil {
		return false, errors.New("access token is nil")
	}

	is_dev := *answer.AccessToken == "yes they're a dev"
	if !is_dev {
		return false, errors.New("access token is not equal 'yes they're a dev'")
	}

	return true, nil
}
