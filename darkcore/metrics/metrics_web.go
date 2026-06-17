package metrics

import (
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
)

var (
	upTime prometheus.Counter = promauto.NewCounter(prometheus.CounterOpts{
		Name: "darkstat_uptime_seconds",
		Help: "Up time in seconds",
	})
	HttpResponseByPatternDurationHist *prometheus.HistogramVec = promauto.NewHistogramVec(prometheus.HistogramOpts{
		Name: "darkstat_http_by_pattern_duration_seconds_hist",
		Help: "Duration histogram of http request in seconds",
	}, []string{"pattern", "status_code"})

	HttpResponseByPatternAndUserAgentTotal *prometheus.CounterVec = promauto.NewCounterVec(prometheus.CounterOpts{
		Name: "darkstat_http_by_useragent_finished_total",
		Help: "Count of http requests by useragent and pattern",
	}, []string{"pattern", "useragent"})

	HttpResponseByIpFinishedTotal *prometheus.CounterVec = promauto.NewCounterVec(prometheus.CounterOpts{
		Name: "darkstat_http_by_ip_finished_total",
		Help: "Finished http requests by ip total",
	}, []string{"ip", "status_code"})
	HttpResponseByIpDurationSum *prometheus.GaugeVec = promauto.NewGaugeVec(prometheus.GaugeOpts{
		Name: "darkstat_http_by_ip_duration_seconds_sum",
		Help: "Duration sum of http request by ip in seconds",
	}, []string{"ip", "status_code"})

	// Technically u need only Sum and Count? May be buckets are overkill
	HttpResponseByPatternBodySizeHist *prometheus.HistogramVec = promauto.NewHistogramVec(prometheus.HistogramOpts{
		Name: "darkstat_http_by_pattern_body_size_bytes_hist",
		Help: "Body size histogram of http response in bytes",
		// Buckets: []float64{.005, .01, .025, .05, .1, .25, .5, 1, 2.5, 5, 10}, // default not usable for bytes
		Buckets: []float64{10_000, 25_000, 50_000, 100_000, 250_000, 500_000, 1_000_000, 2_500_000, 5_000_000, 10_000_000},
	}, []string{"pattern", "status_code"})

	web_metrics = []prometheus.Collector{
		upTime,
		HttpResponseByPatternAndUserAgentTotal,
		HttpResponseByPatternDurationHist,
		HttpResponseByIpFinishedTotal,
		HttpResponseByIpDurationSum,
		HttpResponseByPatternBodySizeHist,
	}
)
