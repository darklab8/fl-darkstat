package darkgrpc

import (
	"net/http"
	"strconv"
	"time"

	"github.com/darklab8/fl-darkstat/darkcore/metrics"
	"github.com/darklab8/fl-darkstat/darkcore/settings/logus"
	"github.com/darklab8/go-utils/typelog"
)

type logResponseWriter struct {
	http.ResponseWriter
	statusCode int
}

func (rsp *logResponseWriter) WriteHeader(code int) {
	rsp.statusCode = code
	rsp.ResponseWriter.WriteHeader(code)
}

// Unwrap returns the original http.ResponseWriter. This is necessary
// to expose Flush() and Push() on the underlying response writer.
func (rsp *logResponseWriter) Unwrap() http.ResponseWriter {
	return rsp.ResponseWriter
}

func newLogResponseWriter(w http.ResponseWriter) *logResponseWriter {
	return &logResponseWriter{w, http.StatusOK}
}

// MiddlewarePrometheusForAPIGateway is inspired by
// how to write middlewares for api gateways
// from over there https://github.com/grpc-ecosystem/grpc-gateway/blob/115075f9df32cc3a4fe3325ee2dc274d5569876f/docs/docs/operations/logging.md?plain=1#L10
func MiddlewarePrometheusForAPIGateway(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		lw := newLogResponseWriter(w)

		// BODY READING EXAMPLE
		// Note that buffering the entire request body could consume a lot of memory.
		// body, err := io.ReadAll(r.Body)
		// if err != nil {
		// 	http.Error(w, fmt.Sprintf("failed to read body: %v", err), http.StatusBadRequest)
		// 	return
		// }
		// clonedR := r.Clone(r.Context())
		// clonedR.Body = io.NopCloser(bytes.NewReader(body))
		// h.ServeHTTP(lw, clonedR)

		time_start := time.Now()
		h.ServeHTTP(lw, r)
		time_finish := time.Since(time_start).Seconds()
		pattern := r.URL.Path
		metrics.HttpResponseByPatternDurationHist.WithLabelValues(pattern, strconv.Itoa(lw.statusCode)).Observe(time_finish)

		Logger := logus.Log.WithFields(
			typelog.Any("pattern", pattern),
			typelog.Any("status_code", lw.statusCode),
		)
		if lw.statusCode >= 400 && lw.statusCode < 500 {
			Logger.Warn("finished api gateway request")
		} else if lw.statusCode >= 500 {
			Logger.Error("finished api gateway request")
		} else {
			Logger.Info("finished api gateway request")
		}
	})
}
