package main

import (
	"github.com/google/gnxi/gnmi"
	"github.com/google/gnxi/gnmi/modeldata"
	"github.com/google/gnxi/gnmi/modeldata/gostruct"
	"github.com/google/gnxi/utils/credentials"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
	"net"
	"os"
	"reflect"

	"fmt"
	"github.com/op/go-logging"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"net/http"

	pb "github.com/openconfig/gnmi/proto/gnmi"
)

var (
	logModule = "ovs-target"
	log       = logging.MustGetLogger(logModule)
)

type server struct {
	*gnmi.Server
}

func newServer(model *gnmi.Model, config []byte) (*server, error) {
	s, err := gnmi.NewServer(model, config, nil)
	if err != nil {
		return nil, err
	}
	return &server{Server: s}, nil
}

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

	client, err := NewOVSClient("127.0.0.1", "6640")
	if err != nil {
		log.Fatal("Unable to initialize OVS Client\n")
	}

	log.Info(client.String())

	//
	// Test https://github.com/google/gnxi/blob/master/gnmi_target/gnmi_target.go
	//

	model := gnmi.NewModel(modeldata.ModelData,
		reflect.TypeOf((*gostruct.Device)(nil)),
		gostruct.SchemaTree["Device"],
		gostruct.Unmarshal,
		gostruct.ΛEnum)

	opts := credentials.ServerCredentials()
	g := grpc.NewServer(opts...)

	configData := []byte(
		`"{
		"openconfig-openflow:openflow": {
			"controllers": {
				"controller": [
					{
						"config": {
							"name": "main"
						},
						"connections": {
							"connection": [
								{
									"config": {
										"address": "127.0.0.1"
									},
									"state": {
										"address": "127.0.0.1"
									}
								}
							]
						},
						"name": "main"
					}
				]
			}
		}
	}"`)
	s, err := newServer(model, configData)
	if err != nil {
		log.Fatalf("Error on creating gNMI target: %v", err)
	}
	pb.RegisterGNMIServer(g, s)
	reflection.Register(g)

	log.Infof("Starting to listen")
	listen, err := net.Listen("tcp", "0.0.0.0")
	if err != nil {
		log.Fatalf("Failed to listen: %v", err)
	}

	log.Info("Starting to serve")
	if err := g.Serve(listen); err != nil {
		log.Fatalf("Failed to serve: %v", err)
	}

	//
	//
	//

	run()
}

func run() {
	select {}
}

type PrometheusMonitoringInstance struct {
	IPAddress         string
	Port              string
	ErrorsGaugeMetric prometheus.Gauge
}

func NewPrometheusMonitoringInstance(ipAddress, port string) (*PrometheusMonitoringInstance, error) {
	p := PrometheusMonitoringInstance{IPAddress: ipAddress, Port: port}
	p.InitializeMetrics()
	p.RegisterMetrics()
	return &p, nil
}

func (p *PrometheusMonitoringInstance) InitializeMetrics() {
	p.ErrorsGaugeMetric = prometheus.NewGauge(prometheus.GaugeOpts{
		Name: "ovsgnxi_errors",
		Help: "The number of errors ovs-target experienced during runtime.",
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
