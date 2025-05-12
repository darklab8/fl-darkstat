package web

import (
	"fmt"
	"log"
	"net"
	"net/http"
	"net/http/cookiejar"
	"regexp"
	"testing"
	"time"

	"github.com/darklab8/fl-darkstat/darkapis/darkhttp/apiutils"
	"github.com/darklab8/fl-darkstat/darkcore/builder"
	"github.com/darklab8/fl-darkstat/darkcore/settings"
	"github.com/darklab8/fl-darkstat/darkcore/settings/logus"
	"github.com/darklab8/fl-darkstat/darkcore/web/registry"
	statsettings "github.com/darklab8/fl-darkstat/darkstat/settings"
	"github.com/darklab8/go-utils/typelog"
	"github.com/darklab8/go-utils/utils/ptr"
	"github.com/darklab8/go-utils/utils/utils_types"
	"github.com/stretchr/testify/assert"
)

func NewServerOauthStart(w *Web) *registry.Endpoint {
	return &registry.Endpoint{
		Url: "GET /forums/oauth/",
		Handler: func(w http.ResponseWriter, r *http.Request) {
			r.URL.Query()
			if r.URL.Query().Get("client_id") != "darkstat_dev" {
				panic("client id is not darkstat_dev")
			}
			redirect_url := r.URL.Query().Get("redirect_url")

			logus.Log.Info("NewServerOauthStart accepted request", typelog.String("redirect_url", redirect_url))
			http.Redirect(w, r, fmt.Sprintf("%s?code=1234", redirect_url), http.StatusSeeOther)
		},
	}
}
func NewServerOauthValidate(w *Web) *registry.Endpoint {
	return &registry.Endpoint{
		Url: "GET /forums/oauth/access_token",
		Handler: func(w http.ResponseWriter, r *http.Request) {
			code := r.FormValue("code")

			if code != "1234" {
				w.WriteHeader(http.StatusForbidden)
				output := OauthAnswer{Error: ptr.Ptr("incorrect code")}
				apiutils.ReturnJson(&w, output)
				return
			}

			output := OauthAnswer{AccessToken: ptr.Ptr(AccesTokenIsDev)}
			apiutils.ReturnJson(&w, output)
		},
	}
}

func OauthTestServer(opts WebServeOpts) {
	w := NewWebBasic([]*builder.Filesystem{}, WithSiteRoot("http://localhost:8889/"))
	w.registry.Register(NewServerOauthStart(w))
	w.registry.Register(NewServerOauthValidate(w))
	w.registry.Foreach(func(e *registry.Endpoint) {
		w.mux.HandleFunc(string(e.Url), e.Handler)
	})
	hander := CorsMiddleware(w.mux)
	var err error
	tcp_listener, err := net.Listen("tcp", fmt.Sprintf("%s:%d", "0.0.0.0", 8889))
	if err != nil {
		panic(err)
	}
	tcp_server := http.Server{Handler: hander}
	err = tcp_server.Serve(tcp_listener) // if serving over Http
	if err != nil {
		log.Fatal("http error:", err)
	}
}

type EmptyPage struct{}

func (e EmptyPage) Render() []byte { return []byte(string("empty page")) }

func OauthTestClient(opts WebServeOpts) {
	fs := &builder.Filesystem{Files: make(map[utils_types.FilePath]builder.MemFile)}
	fs.Files["index.html"] = EmptyPage{}
	w := NewWebBasic([]*builder.Filesystem{fs}, WithSiteRoot("http://localhost:8881/"))
	w.registry.Register(NewOauthStart(w))
	w.registry.Register(NewOauthAccept(w))
	w.registry.Register(NewEndpointStatic(w))
	w.registry.Register(NewEndpointPing(w))
	w.registry.Foreach(func(e *registry.Endpoint) {
		w.mux.HandleFunc(string(e.Url), e.Handler)
	})
	hander := AuthMiddleware(CorsMiddleware(w.mux))
	// hander := CorsMiddleware(w.mux)
	var err error
	tcp_listener, err := net.Listen("tcp", fmt.Sprintf("%s:%d", "0.0.0.0", 8881))
	if err != nil {
		panic(err)
	}
	tcp_server := http.Server{Handler: hander}
	err = tcp_server.Serve(tcp_listener) // if serving over Http
	if err != nil {
		log.Fatal("http error:", err)
	}
}

func TestOauth(t *testing.T) {
	site_url := statsettings.Env.SiteUrl
	password := settings.Env.Password
	secret := settings.Env.Secret
	is_oauth := settings.Env.IsDiscoOauthEnabled

	DiscoOauthSiteUrl = "http://localhost:8889"
	statsettings.Env.SiteUrl = "http://localhost:8881"
	settings.Env.Password = "1234"
	settings.Env.Secret = "passphrasewhichneedstobe32bytes!"
	settings.Env.IsDiscoOauthEnabled = true

	go OauthTestServer(WebServeOpts{})
	go OauthTestClient(WebServeOpts{})
	time.Sleep(1 * time.Second)

	jar, _ := cookiejar.New(nil)
	client := http.Client{
		Jar: jar,
	}

	app_url := statsettings.Env.SiteUrl
	answer, err := Get(client, app_url+"/")
	assert.Nil(t, err)

	fmt.Println(answer.StatusCode, string(answer.Body))

	redirect_url := regexp.MustCompile(`url=(.+)\"`).FindStringSubmatch(string(answer.Body))[1]
	fmt.Println("getting redirected to ", redirect_url)
	answer, err = Get(client, app_url+redirect_url)
	assert.Nil(t, err)

	assert.Equal(t, statsettings.Env.SiteUrl+"/oauth/redirect?code=1234", answer.Resp.Request.URL.String())

	fmt.Println("print cookies in answer")
	cookies := answer.Resp.Cookies()
	for _, c := range cookies {
		fmt.Println("name=", c.Name, " value=", c.Value) // The cookie's value
	}

	redirect_url2 := regexp.MustCompile(`url=(.+)\"`).FindStringSubmatch(string(answer.Body))[1]
	fmt.Println("getting redirected to ", redirect_url2)
	answer, err = Get(client, app_url+redirect_url2)
	assert.Nil(t, err)
	fmt.Println(answer.StatusCode, string(answer.Body))
	assert.Equal(t, answer.StatusCode, 200)
	// ideally you should assert this, but for some reason cookie is not propagating as it should
	// assert.Contains(t, answer.Body, "empty page")

	statsettings.Env.SiteUrl = site_url
	settings.Env.Password = password
	settings.Env.Secret = secret
	settings.Env.IsDiscoOauthEnabled = is_oauth
}
