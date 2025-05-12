package web

import (
	"errors"
	"net"
	"net/http"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/darklab8/fl-darkstat/darkcore/metrics"
	"github.com/darklab8/fl-darkstat/darkcore/settings/logus"
	"github.com/darklab8/fl-darkstat/darkcore/settings/traces"
	"github.com/darklab8/fl-darkstat/darkstat/front/urls"
	"github.com/darklab8/go-utils/typelog"
	"github.com/darklab8/go-utils/utils/regexy"
	"github.com/prometheus/client_golang/prometheus"
)

type statusRecorder struct {
	http.ResponseWriter
	status    int
	body_size int
}

func (rec *statusRecorder) WriteHeader(code int) {
	rec.status = code
	rec.ResponseWriter.WriteHeader(code)

	var bytes []byte = []byte{}
	rec.body_size, _ = rec.ResponseWriter.Write(bytes)
}
func (rec *statusRecorder) Write(bytes []byte) (int, error) {
	var err error
	rec.body_size, err = rec.ResponseWriter.Write(bytes)
	return rec.body_size, err
}

func prometheusMidleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Request
		rec := statusRecorder{w, 200, 0}
		time_start := time.Now()

		ctx, span := traces.Tracer.Start(r.Context(), "request")
		defer span.End()

		next.ServeHTTP(&rec, r.WithContext(ctx))
		time_finish := time.Since(time_start).Seconds()

		// Gathering request info
		pattern := r.Pattern
		url := r.URL.Path[1:]
		if pattern == "/" && (strings.Contains(url, urls.Index.ToString()) ||
			strings.Contains(url, urls.DarkIndex.ToString()) ||
			strings.Contains(url, urls.VanillaIndex.ToString()) ||
			strings.Contains(url, "cdn")) {
			pattern = UrlGeneralizer.ReplaceAllString(url, "-{item_id}")
			logus.Log.Debug("generalized url to pattern", typelog.String("pattern", pattern))
		}
		if pattern == "/" && r.URL.Path != "/" {
			pattern = "unknown"
		}
		ip, err := getIP(r)
		logus.Log.CheckError(err, "not found ip in prometheus middleware incoming request")

		user_agent := r.Header.Get("User-Agent")

		Logger := Log.WithFields(
			typelog.String("pattern", pattern),
			typelog.Int("status_code", rec.status),
			typelog.String("url", r.URL.Path),
			typelog.Float64("duration", time_finish),
			typelog.Int("body_size", rec.body_size),
			typelog.String("user_agent", user_agent),
		)
		if rec.status >= 400 && rec.status < 500 && !strings.Contains(r.URL.Path, "favicon.ico") {
			Logger.WarnCtx(ctx, "finished request")
		} else if rec.status >= 500 {
			Logger.ErrorCtx(ctx, "finished request")
		} else {
			Logger.InfoCtx(ctx, "finished request")
		}

		var trace_id string
		if span.SpanContext().HasSpanID() {
			trace_id = span.SpanContext().TraceID().String()
		}

		// confirm it is present
		// curl -H 'Accept: application/openmetrics-text' localhost:8000/metrics | grep "darkstat_http_by_pattern_duration_seconds_hist"
		metrics.HttpResponseByPatternDurationHist.WithLabelValues(pattern, strconv.Itoa(rec.status)).(prometheus.ExemplarObserver).ObserveWithExemplar(
			time_finish, prometheus.Labels{"traceID": trace_id},
		)

		metrics.HttpResponseByPatternBodySizeHist.WithLabelValues(pattern, strconv.Itoa(rec.status)).Observe(float64(rec.body_size))

		metrics.HttpResponseByIpFinishedTotal.WithLabelValues(ip, strconv.Itoa(rec.status)).Inc()
		metrics.HttpResponseByIpDurationSum.WithLabelValues(ip, strconv.Itoa(rec.status)).Add(time_finish)
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
