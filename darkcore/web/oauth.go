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

	"github.com/darklab8/fl-darkstat/darkcore/settings/logus"
	"github.com/darklab8/fl-darkstat/darkcore/web/registry"
	statsettings "github.com/darklab8/fl-darkstat/darkstat/settings"
	"github.com/darklab8/go-utils/typelog"
	"github.com/darklab8/go-utils/utils/utils_settings"
)

// I have now implemented roughly this
// You redirect unauthenticated users to https://discoverygc.com/forums/oauth/?client_id=darkstat_dev&redirect_url=https://darkstat-dev.dd84ai.com/oauth/redirect
// I check that they're in the right group, and redirect them to https://darkstat-dev.dd84ai.com/oauth/redirect?code=${some_assigned_token}
// You POST {"code":"${code}"} to https://discoverygc.com/forums/oauth/access_token
// I return either {"error": "nope"} or {"access_token": "yes they're a dev"}
// if you get access_token "yes they're a dev" you can grant them a session cookie permitting them to view info from the private development repo
// let me know if any of this is wrong

var DiscoOauthSiteUrl = "https://discoverygc.com"

func NewOauthStart(w *Web) *registry.Endpoint {
	return &registry.Endpoint{
		Url: "GET /oauth",
		Handler: func(w http.ResponseWriter, r *http.Request) {
			redirect_url := fmt.Sprintf("%s/forums/oauth/?client_id=darkstat_dev&redirect_url=%s/oauth/redirect", DiscoOauthSiteUrl, statsettings.Env.SiteUrl)
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
				_, err := fmt.Fprintln(w, "failed oauth procedure.", err)
				Log.CheckError(err, "failed oauth procedure responce")
			}

			tempus_value := NewTempusToken()
			fmt.Println("setting tempus cookie for succesful oauth login", "host=", r.Host)
			http.SetCookie(w, &http.Cookie{Name: "tempus", Value: tempus_value, Expires: time.Now().Add(1 * time.Hour), Path: "/", HttpOnly: true})

			// http.Redirect(w, r, statsettings.Env.SiteUrl, http.StatusSeeOther)
			// redirect with delay instead
			buf := bytes.NewBuffer([]byte{})
			err = RedirectPage(
				"Succesfully oauth authentificated, u will be redirected in 3 seconds to main darkstat page",
				"/").Render(context.Background(), buf)
			logus.Log.CheckError(err, "failed to redirect oauth response")
			_, err = fmt.Fprint(w, buf.String())
			logus.Log.CheckError(err, "failed to print into response")
		},
	}
}

// I return either {"error": "nope"} or {"access_token": "yes they're a dev"}
type OauthAnswer struct {
	Error       *string `json:"error"`
	AccessToken *string `json:"access_token"`
}

const AccesTokenIsDev = "yes they're a dev"

func validateCode(code string) (bool, error) {
	var err error
	var client = &http.Client{}
	var answer OauthAnswer
	var param = url.Values{}
	param.Set("code", code)
	var payload = bytes.NewBufferString(param.Encode())
	request, err := http.NewRequest("POST", fmt.Sprintf("%s/forums/oauth/access_token", DiscoOauthSiteUrl), payload)
	if err != nil {
		return false, err
	}
	request.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	response, err := client.Do(request)
	if utils_settings.Envs.UserAgent != "" {
		request.Header.Set("User-Agent", utils_settings.Envs.UserAgent)
	}
	if err != nil {
		return false, err
	}
	defer func() {
		err = response.Body.Close()
		Log.CheckError(err, "failed to close body")
	}()
	err = json.NewDecoder(response.Body).Decode(&answer)
	if err != nil {
		return false, err
	}
	logus.Log.Info("got answer", typelog.Struct(answer))

	if answer.AccessToken == nil {
		return false, errors.New("access token is nil")
	}

	is_dev := *answer.AccessToken == AccesTokenIsDev
	if !is_dev {
		return false, errors.New("access token is not equal 'yes they're a dev'")
	}

	return true, nil
}
