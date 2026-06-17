package metrics

import (
	"net/http"
	"time"

	"github.com/darklab8/fl-darkstat/darkcore/settings"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/collectors"
	"github.com/prometheus/client_golang/prometheus/collectors/version"
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

	metrics := append(
		web_metrics,
		collectors.NewGoCollector(),
		collectors.NewProcessCollector(collectors.ProcessCollectorOpts{}),
		version.NewCollector("darkstat"),
	)
	metrics = append(metrics, pob_metrics...)

	reg.MustRegister(
		metrics...,
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
