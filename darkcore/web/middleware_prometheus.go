package web

import (
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/darklab8/fl-darkstat/darkcore/metrics"
	"github.com/darklab8/fl-darkstat/darkcore/settings/logus"
	"github.com/darklab8/fl-darkstat/darkstat/front/urls"
	"github.com/darklab8/go-typelog/typelog"
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
		next.ServeHTTP(&rec, r)
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

		Logger := Log.WithFields(
			typelog.String("pattern", pattern),
			typelog.Int("status_code", rec.status),
			typelog.String("url", r.URL.Path),
			typelog.Float64("duration", time_finish),
			typelog.Int("body_size", rec.body_size),
		)
		if rec.status >= 400 && rec.status < 500 && !strings.Contains(r.URL.Path, "favicon.ico") {
			Logger.Warn("finished request")
		} else if rec.status >= 500 {
			Logger.Error("finished request")
		} else {
			Logger.Info("finished request")
		}

		// Metrics after request
		metrics.HttpResponseByPatternDurationHist.WithLabelValues(pattern, strconv.Itoa(rec.status)).Observe(time_finish)

		metrics.HttpResponseByPatternBodySizeHist.WithLabelValues(pattern, strconv.Itoa(rec.status)).Observe(float64(rec.body_size))

		metrics.HttpResponseByIpFinishedTotal.WithLabelValues(ip, strconv.Itoa(rec.status)).Inc()
		metrics.HttpResponseByIpDurationSum.WithLabelValues(ip, strconv.Itoa(rec.status)).Add(time_finish)

	})
}
