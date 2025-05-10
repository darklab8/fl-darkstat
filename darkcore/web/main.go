package web

/*
Entrypoint for front and for dev web server?
*/

import (
	"errors"
	"fmt"
	"log"
	"net"
	"net/http"
	"os"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/darklab8/fl-darkstat/darkcore/metrics"
	"github.com/darklab8/fl-darkstat/darkstat/front/urls"
	"github.com/darklab8/go-typelog/typelog"
	"github.com/darklab8/go-utils/utils/regexy"

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

	site_root string
}

func (w *Web) GetMux() *http.ServeMux { return w.mux }

type WebOpt func(w *Web)

func WithAppData(data *appdata.AppData) WebOpt {
	return func(w *Web) {
		w.data = data
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

// getIP returns the ip address from the http request
func getIP(r *http.Request) (string, error) {
	ips := r.Header.Get("X-Forwarded-For")
	splitIps := strings.Split(ips, ",")

	if len(splitIps) > 0 {
		// get last IP in list since ELB prepends other user defined IPs, meaning the last one is the actual client IP.
		netIP := net.ParseIP(splitIps[len(splitIps)-1])
		if netIP != nil {
			return netIP.String(), nil
		}
	}

	ip, _, err := net.SplitHostPort(r.RemoteAddr)
	if err != nil {
		return "", err
	}

	netIP := net.ParseIP(ip)
	if netIP != nil {
		ip := netIP.String()
		if ip == "::1" {
			return "127.0.0.1", nil
		}
		return ip, nil
	}

	return "", errors.New("IP not found")
}

var UrlGeneralizer *regexp.Regexp = regexy.InitRegex(`-[\w0-9]*`)

type statusRecorder struct {
	http.ResponseWriter
	status int
}

func (rec *statusRecorder) WriteHeader(code int) {
	rec.status = code
	rec.ResponseWriter.WriteHeader(code)
}

func prometheusMidleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Gathering all info pre request
		pattern := r.Pattern
		url := r.URL.Path[1:]
		if pattern == "" && (strings.Contains(url, urls.Index.ToString()) ||
			strings.Contains(url, urls.DarkIndex.ToString()) ||
			strings.Contains(url, urls.VanillaIndex.ToString()) ||
			strings.Contains(url, "cdn")) {
			pattern = UrlGeneralizer.ReplaceAllString(url, "-{item_id}")
			logus.Log.Debug("generalized url to pattern", typelog.String("pattern", pattern))
		}
		if pattern == "" {
			pattern = "unknown"
		}

		ip, err := getIP(r)
		logus.Log.CheckError(err, "not found ip in prometheus middleware incoming request")

		// Metrics pre request
		metrics.HttpRequestByPatternStartedTotal.WithLabelValues(pattern).Inc()

		// Request
		rec := statusRecorder{w, 200}
		time_start := time.Now()
		next.ServeHTTP(&rec, r)

		if rec.status >= 400 && rec.status < 500 && !strings.Contains(r.URL.Path, "favicon.ico") {
			Log.Warn("finished request", typelog.String("pattern", pattern), typelog.Int("status_code", rec.status), typelog.String("url", r.URL.Path))
		} else if rec.status >= 500 {
			Log.Error("finished request", typelog.String("pattern", pattern), typelog.Int("status_code", rec.status), typelog.String("url", r.URL.Path))
		} else {
			Log.Info("finished request", typelog.String("pattern", pattern), typelog.Int("status_code", rec.status), typelog.String("url", r.URL.Path))
		}

		time_finish := time.Since(time_start).Seconds()

		// Metrics after request
		metrics.HttpRequestByPatternFinishedTotal.WithLabelValues(pattern, strconv.Itoa(rec.status)).Inc()
		metrics.HttpRequestByPatternDurationSum.WithLabelValues(pattern, strconv.Itoa(rec.status)).Add(time_finish)
		metrics.HttpRequestByPatternDurationHist.WithLabelValues(pattern, strconv.Itoa(rec.status)).Observe(time_finish)

		metrics.HttpRequestByIpFinishedTotal.WithLabelValues(ip, strconv.Itoa(rec.status)).Inc()
		metrics.HttpRequestByIpDurationSum.WithLabelValues(ip, strconv.Itoa(rec.status)).Add(time_finish)

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
		sock_server = http.Server{Handler: hander}
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
