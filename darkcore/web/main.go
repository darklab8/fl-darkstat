package web

/*
Entrypoint for front and for dev web server?
*/

import (
	"fmt"
	"log"
	"net"
	"net/http"
	"os"

	"github.com/darklab8/fl-darkstat/configs/cfg"
	"github.com/darklab8/fl-darkstat/darkcore/builder"
	"github.com/darklab8/fl-darkstat/darkcore/settings/logus"
	"github.com/darklab8/fl-darkstat/darkcore/web/registry"
)

type Mutex interface {
	Lock()
	Unlock()
}

type Web struct {
	filesystems  []*builder.Filesystem
	registry     *registry.Registion
	mux          *http.ServeMux
	AppDataMutex Mutex

	site_root string
}

func (w *Web) GetMux() *http.ServeMux { return w.mux }

type WebOpt func(w *Web)

func WithMutexableData(app_data_mutex Mutex) WebOpt {
	return func(w *Web) {
		w.AppDataMutex = app_data_mutex
	}
}

func WithSiteRoot(site_root string) WebOpt {
	return func(w *Web) {
		w.site_root = site_root
	}
}

func NewWeb(filesystems []*builder.Filesystem, opts ...WebOpt) *Web {
	w := &Web{
		filesystems: filesystems,
		registry:    registry.NewRegister(),
		mux:         http.NewServeMux(),
	}

	for _, opt := range opts {
		opt(w)
	}

	fmt.Println("site_root", w.site_root)

	w.registry.Register(NewEndpointStatic(w))

	w.registry.Register(NewEndpointPing(w))
	w.registry.Register(NewOauthStart(w))
	w.registry.Register(NewOauthAccept(w))

	return w
}

type WebServeOpts struct {
	Port        *int
	SockAddress string
}

func enableCors(w *http.ResponseWriter) {
	(*w).Header().Set("Access-Control-Allow-Origin", "*")
	(*w).Header().Set("Access-Control-Allow-Headers", "*")
}

func CorsMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		enableCors(&w)
		next.ServeHTTP(w, r)
	})
}

func (w *Web) Serve(opts WebServeOpts) ServerClose {
	w.registry.Foreach(func(e *registry.Endpoint) {
		w.mux.HandleFunc(string(e.Url), e.Handler)
	})

	ip := "0.0.0.0"
	var port int = 8000
	if opts.Port != nil {
		port = *opts.Port
	}

	fmt.Printf("launching web server, visit http://localhost:%d to check it!\n", port)
	hander := AuthMiddleware(CorsMiddleware(w.mux))

	var sock_listener net.Listener
	var sock_server http.Server
	if cfg.IsLinux && opts.SockAddress != "" {
		sock_server = http.Server{Handler: hander}
		var err error
		os.Mkdir("/tmp/darkstat", 0777)
		err = os.Remove(opts.SockAddress)
		logus.Log.CheckWarn(err, "attempting to remove socket")

		fmt.Println("starting to listen to sock ", opts.SockAddress)
		sock_listener, err = net.Listen("unix", opts.SockAddress)
		if err != nil {
			panic(err)
		}
		go sock_server.Serve(sock_listener)
	}

	tcp_server := http.Server{Handler: hander}

	var err error
	tcp_listener, err := net.Listen("tcp", fmt.Sprintf("%s:%d", ip, port))
	if err != nil {
		panic(err)
	}
	go func() {
		err := tcp_server.Serve(tcp_listener) // if serving over Http
		if err != nil {
			log.Fatal("http error:", err)
		}
	}()
	return ServerClose{
		sock_adrr: opts.SockAddress,
	}
}

type ServerClose struct {
	sock_adrr string
}

func (s ServerClose) Close() {
	os.Remove(s.sock_adrr)
}

const DarkstatAPISock = "/tmp/darkstat/api.sock"
