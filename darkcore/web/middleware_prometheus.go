package web

import (
	"bytes"
	"errors"
	"io"
	"net"
	"net/http"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/darklab8/fl-darkstat/darkcore/metrics"
	"github.com/darklab8/fl-darkstat/darkcore/settings/logus"
	"github.com/darklab8/fl-darkstat/darkcore/settings/traces"
	"github.com/darklab8/go-utils/typelog"
	"github.com/darklab8/go-utils/utils/regexy"
	"github.com/prometheus/client_golang/prometheus"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/codes"
	"go.opentelemetry.io/otel/trace"
)

type statusRecorder struct {
	http.ResponseWriter
	status      int
	body_size   int
	wroteHeader bool
}

func (rec *statusRecorder) WriteHeader(code int) {
	if rec.wroteHeader {
		return // net/http ignores repeated WriteHeader calls; don't overwrite status
	}
	rec.wroteHeader = true
	rec.status = code
	rec.ResponseWriter.WriteHeader(code)
}

func (rec *statusRecorder) Write(b []byte) (int, error) {
	rec.wroteHeader = true // an implicit 200 happens on first Write
	n, err := rec.ResponseWriter.Write(b)
	rec.body_size += n // accumulate
	return n, err
}

// Keep streaming/flush support, which middleware like gzhttp relies on.
func (rec *statusRecorder) Flush() {
	if f, ok := rec.ResponseWriter.(http.Flusher); ok {
		f.Flush()
	}
}

// Lets http.ResponseController reach the underlying writer.
func (rec *statusRecorder) Unwrap() http.ResponseWriter {
	return rec.ResponseWriter
}
func prometheusMidleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		metrics.HttpActiveRequests.Inc()
		defer metrics.HttpActiveRequests.Dec()

		// Request
		rec := statusRecorder{ResponseWriter: w, status: http.StatusOK}
		time_start := time.Now()

		opts := []trace.SpanStartOption{
			trace.WithSpanKind(trace.SpanKindServer),
		}
		ctx, span := traces.Tracer.Start(r.Context(), "web-request", opts...)
		defer span.End()

		var bodyBytes []byte
		if r.Body != nil && r.Body != http.NoBody {
			var err error

			body_to_read := http.MaxBytesReader(w, r.Body, 10<<20) // 10Mb Read limit
			bodyBytes, err = io.ReadAll(body_to_read)
			if err != nil {
				http.Error(w, "cannot read body", http.StatusBadRequest)
				return
			}
			_ = r.Body.Close()
			r.Body = io.NopCloser(bytes.NewReader(bodyBytes))
		}

		new_r := r.WithContext(ctx)
		next.ServeHTTP(&rec, new_r)
		span.AddEvent("ServeHTTP is made")
		time_finish := time.Since(time_start).Seconds()

		// Gathering request info
		pattern := new_r.Pattern
		url := r.URL.Path[1:]

		// pattern recognition should work only for route_all.go
		if pattern == "/" && r.URL.Path != "/" {
			if rec.status != http.StatusNotFound {
				if strings.Contains(url, "static/") {
					pattern = "static/{file}"
				} else {
					pattern = UrlGeneralizer.ReplaceAllString(url, "-{item_id}")
					logus.Log.Debug("generalized url to pattern", typelog.String("url_pattern", pattern))
				}
			} else {
				pattern = "unknown"
			}
		}

		ip, err := getIP(r)
		logus.Log.CheckError(err, "not found ip in prometheus middleware incoming request")

		user_agent := r.Header.Get("User-Agent")

		Logger := Log.WithFields(
			typelog.String("url_pattern", pattern),
			typelog.Int("status_code", rec.status),
			typelog.String("url", r.URL.Path),
			typelog.Float64("duration", time_finish),
			typelog.Int("body_size", rec.body_size),
			typelog.String("user_agent", user_agent),
			typelog.String("method", r.Method),
			typelog.Any("body", string(bodyBytes)),
		)
		span.SetAttributes(
			attribute.String("pattern", pattern),
			attribute.Int("status_code", rec.status),
			attribute.String("url", r.URL.Path),
			attribute.Float64("duration", time_finish),
			attribute.Int("body_size", rec.body_size),
			attribute.String("user_agent", user_agent),
			attribute.String("method", r.Method),
		)
		span.SetName("web-request " + pattern)

		if rec.status >= 400 && rec.status < 500 && !strings.Contains(r.URL.Path, "favicon.ico") {
			Logger.WarnCtx(ctx, "finished request")
			span.SetStatus(codes.Error, "4xx status code")
		} else if rec.status >= 500 {
			Logger.ErrorCtx(ctx, "finished request")
			span.SetStatus(codes.Error, "5xx status code")
		} else {
			Logger.InfoCtx(ctx, "finished request")
		}

		var trace_id string
		if span.SpanContext().HasSpanID() {
			trace_id = span.SpanContext().TraceID().String()
		}

		// confirm it is present
		// curl -H 'Accept: application/openmetrics-text' localhost:8000/metrics | grep "darkstat_http_by_pattern_duration_seconds_hist"
		metrics.HttpResponseByPatternDurationHist.WithLabelValues(
			pattern,
			strconv.Itoa(rec.status),
		).(prometheus.ExemplarObserver).ObserveWithExemplar(
			time_finish, prometheus.Labels{"traceID": trace_id},
		)

		metrics.HttpResponseByPatternBodySizeHist.WithLabelValues(pattern, strconv.Itoa(rec.status)).Observe(float64(rec.body_size))

		metrics.HttpResponseByIpFinishedTotal.WithLabelValues(ip, strconv.Itoa(rec.status)).Inc()

		metrics.HttpResponseByPatternAndUserAgentTotal.WithLabelValues(pattern, user_agent).Inc()
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
