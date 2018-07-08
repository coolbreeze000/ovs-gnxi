package main

import (
	"os"

	"github.com/op/go-logging"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"net/http"
	"fmt"
)

var (
	logModule = "ovs-gnxi"
	log       = logging.MustGetLogger(logModule)
)

func main() {
	defer os.Exit(0)
	defer log.Info("Exiting Open vSwitch gNXI interface\n")

	log.Info("Starting Open vSwitch gNXI interface\n")

	prometheusInstance, err := NewPrometheusMonitoringInstance("0.0.0.0", "8080")
	if err != nil {
		log.Errorf("Unable to configure Prometheus Monitoring: %v", err)
		os.Exit(1)
	}
	go prometheusInstance.StartPrometheus()

	run()
}

func run() {
	select {}
}

type PrometheusMonitoringInstance struct {
	IPAddress			string
	Port				string
	ErrorsGaugeMetric	prometheus.Gauge
}

func NewPrometheusMonitoringInstance(ipAddress, port string) (*PrometheusMonitoringInstance, error) {
	p := PrometheusMonitoringInstance{IPAddress: ipAddress, Port: port}
	p.InitializeMetrics()
	p.RegisterMetrics()
	return &p, nil
}

func (p *PrometheusMonitoringInstance) InitializeMetrics() {
	p.ErrorsGaugeMetric = prometheus.NewGauge(prometheus.GaugeOpts{
		Name: "cod_errors",
		Help: "The number of errors CoD experienced during runtime.",
	})
	p.ErrorsGaugeMetric.Set(0)
}

func (p *PrometheusMonitoringInstance) RegisterMetrics() {
	prometheus.MustRegister(p.ErrorsGaugeMetric)
}

func (p *PrometheusMonitoringInstance) StartPrometheus() {
	log.Infof("Starting prometheus on %v:%v...", p.IPAddress, p.Port)
	http.Handle("/metrics", promhttp.Handler())
	log.Fatal(http.ListenAndServe(fmt.Sprintf("%v:%v", p.IPAddress, p.Port), nil))
}