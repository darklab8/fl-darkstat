package web

/*
Entrypoint for front and for dev web server?
*/

import (
	"context"
	"fmt"
	"log"
	"net"
	"net/http"
	"os"
	"strings"

	"github.com/darklab8/fl-darkstat/configs/cfg"
	"github.com/darklab8/fl-darkstat/darkcore/builder"
	"github.com/darklab8/fl-darkstat/darkcore/settings"
	"github.com/darklab8/fl-darkstat/darkcore/settings/logus"
	"github.com/darklab8/fl-darkstat/darkcore/web/registry"
	"github.com/darklab8/fl-darkstat/darkstat/appdata"
)

var (
	Log = logus.Log.WithScope("web")
)

type Mutex interface {
	Lock()
	Unlock()
	RLock()
	RUnlock()
}

type Web struct {
	data         *appdata.AppData
	filesystems  []*builder.Filesystem
	registry     *registry.Registion
	mux          *http.ServeMux
	AppDataMutex Mutex
	ctx          context.Context

	site_root string
}

func (w *Web) GetMux() *http.ServeMux { return w.mux }

type WebOpt func(w *Web)

func WithAppData(data *appdata.AppData) WebOpt {
	return func(w *Web) {
		w.data = data
	}
}
func WithCtx(ctx context.Context) WebOpt {
	return func(w *Web) {
		w.ctx = ctx
	}
}

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

func NewWebBasic(filesystems []*builder.Filesystem, opts ...WebOpt) *Web {
	w := &Web{
		filesystems: filesystems,
		registry:    registry.NewRegister(),
		mux:         http.NewServeMux(),
		ctx:         context.Background(),
	}

	for _, opt := range opts {
		opt(w)
	}
	return w
}

func NewWeb(filesystems []*builder.Filesystem, opts ...WebOpt) *Web {
	w := NewWebBasic(filesystems, opts...)
	fmt.Println("site_root", w.site_root)

	// w.registry.Register(NewBaseTravelRoutes(w)) // example how to write route for appdata
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

func setHeaders(r *http.Request, w *http.ResponseWriter) {
	(*w).Header().Set("Access-Control-Allow-Origin", "*")
	(*w).Header().Set("Access-Control-Allow-Headers", "*")
	(*w).Header().Set("Access-Control-Allow-Methods", "*")

	if settings.Env.CacheControl != "" {
		switch r.Method {
		case http.MethodGet:
			if r.URL == nil {
				return
			}
			if len(r.URL.Path) < 1 {
				return
			}
			requested := r.URL.Path[1:]
			if strings.Contains(requested, "info_") {
				(*w).Header().Set("Cache-Control", "max-age=1200")
				return
			} else if strings.Contains(requested, "route/route_") {
				(*w).Header().Set("Cache-Control", "max-age=1200")
				return
			} else if strings.Contains(requested, "id_rephacks_") {
				(*w).Header().Set("Cache-Control", "max-age=1200")
				return
			} else if strings.Contains(requested, "ships_details") {
				(*w).Header().Set("Cache-Control", "max-age=1200")
				return
			}
		}

		(*w).Header().Set("Cache-Control", "max-age=60")
	}
}

func CorsMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		setHeaders(r, &w)
		if r.Method == "OPTIONS" {
			w.WriteHeader(http.StatusOK)
		} else {
			next.ServeHTTP(w, r)
		}

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
	hander := prometheusMidleware(CorsMiddleware(AuthMiddleware(w.mux)))

	var sock_listener net.Listener
	var sock_server http.Server
	if cfg.IsLinux && opts.SockAddress != "" {
		sock_server = http.Server{
			BaseContext: func(_ net.Listener) context.Context { return w.ctx },
			Handler:     hander,
		}
		var err error
		_ = os.Mkdir("/tmp/darkstat", 0777)
		err = os.Remove(opts.SockAddress)
		logus.Log.CheckWarn(err, "attempting to remove socket")

		fmt.Println("starting to listen to sock ", opts.SockAddress)
		sock_listener, err = net.Listen("unix", opts.SockAddress)
		if err != nil {
			panic(err)
		}
		go func() {
			err = sock_server.Serve(sock_listener)
			Log.CheckError(err, "serve sock server")
		}()
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
	err := os.Remove(s.sock_adrr)
	Log.CheckError(err, "failed to close server sock addr")
}

const DarkstatHttpSock = "/tmp/darkstat/http.sock"
