package metrics

import (
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
)

var (
	PoBInfo *prometheus.GaugeVec = promauto.NewGaugeVec(prometheus.GaugeOpts{
		Name: "darkstat_pob_info",
		Help: "PoB Info",
	}, []string{"pob_nick", "pob_name"})
	PoBHealth *prometheus.GaugeVec = promauto.NewGaugeVec(prometheus.GaugeOpts{
		Name: "darkstat_pob_health",
		Help: "PoB health",
	}, []string{"pob_nick"})
	PoBCargoLeft *prometheus.GaugeVec = promauto.NewGaugeVec(prometheus.GaugeOpts{
		Name: "darkstat_pob_cargo_left",
		Help: "PoB cargo left",
	}, []string{"pob_nick"})
	PoBLevel *prometheus.GaugeVec = promauto.NewGaugeVec(prometheus.GaugeOpts{
		Name: "darkstat_pob_level",
		Help: "PoB level",
	}, []string{"pob_nick"})
	PoBMoney *prometheus.GaugeVec = promauto.NewGaugeVec(prometheus.GaugeOpts{
		Name: "darkstat_pob_money",
		Help: "PoB money",
	}, []string{"pob_nick"})
	PoBItemsAmount *prometheus.GaugeVec = promauto.NewGaugeVec(prometheus.GaugeOpts{
		Name: "darkstat_pob_items_amount",
		Help: "PoB items amount",
	}, []string{"pob_nick"})

	PoBGoodQuantity *prometheus.GaugeVec = promauto.NewGaugeVec(prometheus.GaugeOpts{
		Name: "darkstat_pob_good_quantity",
		Help: "PoB items amount",
	}, []string{"pob_nick", "good_category", "good_nick", "good_name"})

	pob_metrics = []prometheus.Collector{
		PoBInfo,
		PoBHealth,
		PoBCargoLeft,
		PoBLevel,
		PoBMoney,
		PoBItemsAmount,
		PoBGoodQuantity,
	}
)
