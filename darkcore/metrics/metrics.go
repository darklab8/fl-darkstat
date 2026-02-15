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

// Example with OTLP metrics
//	func Init[T any](result T, err error) T {
//		logus.Log.CheckPanic(err, "failed to init metric")
//		return result
//	}
// var (
// 	meter                      = otel.Meter("darkstat")
// 	upTime metric.Int64Counter = Init(meter.Int64Counter(
// 		"darkstat_uptime_seconds",
// 		metric.WithDescription("Up time in seconds by otlp")),
// 	)
// )

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
)

type Metronom struct {
	Reg *prometheus.Registry
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
		HttpResponseByPatternAndUserAgentTotal,
		HttpResponseByPatternDurationHist,
		HttpResponseByIpFinishedTotal,
		HttpResponseByIpDurationSum,
		HttpResponseByPatternBodySizeHist,
		collectors.NewGoCollector(),
		collectors.NewProcessCollector(collectors.ProcessCollectorOpts{}),
		version.NewCollector("darkstat"),
	)

	mux.Handle(
		"/metrics", promhttp.HandlerFor(
			newreg,
			promhttp.HandlerOpts{
				EnableOpenMetrics: true,
			}),
	)

	return &Metronom{
		Reg: newreg,
	}
}

func (m *Metronom) Run() {
	for {
		seconds := 60
		upTime.Add(float64(seconds))
		time.Sleep(time.Second * time.Duration(seconds))
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
