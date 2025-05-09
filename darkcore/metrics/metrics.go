package metrics

import (
	"net/http"
	"time"

	"github.com/darklab8/fl-darkstat/darkcore/settings"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/collectors"
	"github.com/prometheus/client_golang/prometheus/collectors/version"
	"github.com/prometheus/client_golang/prometheus/promauto"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

var (
	upTime prometheus.Counter = promauto.NewCounter(prometheus.CounterOpts{
		Name: "darkstat_uptime_seconds",
		Help: "Up time in seconds",
	})
	HttpRequestByPatternStartedTotal *prometheus.CounterVec = promauto.NewCounterVec(prometheus.CounterOpts{
		Name: "darkstat_http_by_pattern_started_total",
		Help: "Started http requests",
	}, []string{"pattern"})
	HttpRequestByPatternFinishedTotal *prometheus.CounterVec = promauto.NewCounterVec(prometheus.CounterOpts{
		Name: "darkstat_http_by_pattern_finished_total",
		Help: "Finished http requests",
	}, []string{"pattern", "status_code"})
	HttpRequestByPatternDurationSum *prometheus.GaugeVec = promauto.NewGaugeVec(prometheus.GaugeOpts{
		Name: "darkstat_http_by_pattern_duration_seconds_sum",
		Help: "Duration sum of http request in seconds",
	}, []string{"pattern", "status_code"})
	HttpRequestByPatternDurationHist *prometheus.HistogramVec = promauto.NewHistogramVec(prometheus.HistogramOpts{
		Name: "darkstat_http_by_pattern_duration_seconds_hist",
		Help: "Duration histogram of http request in seconds",
	}, []string{"pattern", "status_code"})

	HttpRequestByIpFinishedTotal *prometheus.CounterVec = promauto.NewCounterVec(prometheus.CounterOpts{
		Name: "darkstat_http_by_ip_finished_total",
		Help: "Finished http requests by ip total",
	}, []string{"ip", "status_code"})
	HttpRequestByIpDurationSum *prometheus.GaugeVec = promauto.NewGaugeVec(prometheus.GaugeOpts{
		Name: "darkstat_http_by_ip_duration_seconds_sum",
		Help: "Duration sum of http request by ip in seconds",
	}, []string{"ip", "status_code"})
)

type Metronom struct {
}

func NewMetronom(mux *http.ServeMux) *Metronom {

	newreg := prometheus.NewRegistry()
	reg := prometheus.WrapRegistererWith(
		prometheus.Labels{
			"environment": settings.Env.Environment,
			"project":     "darkstat",
		}, newreg)
	reg.MustRegister(
		upTime,
		HttpRequestByPatternStartedTotal,
		HttpRequestByPatternFinishedTotal,
		HttpRequestByPatternDurationSum,
		HttpRequestByPatternDurationHist,
		HttpRequestByIpFinishedTotal,
		HttpRequestByIpDurationSum,
		collectors.NewGoCollector(),
		collectors.NewProcessCollector(collectors.ProcessCollectorOpts{}),
		version.NewCollector("darkstat"),
	)

	mux.Handle(
		"/metrics", promhttp.HandlerFor(
			newreg,
			promhttp.HandlerOpts{}),
	)

	return &Metronom{}
}

func (m *Metronom) Run() {
	for {

		upTime.Add(60)
		time.Sleep(time.Minute)
	}
}

// Example to extract curent value of metrics
// httpRequestDurationSum.WithLabelValues("123", "123", "123").Set(float64(counter))
// counter++
// metric := &dto.Metric{}
// httpRequestDurationSum.WithLabelValues("123", "123", "123").Write(metric)
// fmt.Println("described metrics value = ", metric.GetGauge().GetValue())
// metric = &dto.Metric{}
// httpRequestDurationSum.WithLabelValues("123", "123", "444").Write(metric)
// fmt.Println("described metrics value = ", metric.GetGauge().GetValue())
