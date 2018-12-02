package main

import (
	"flag"
	"fmt"
	"github.com/google/gnxi/gnmi/modeldata"
	"github.com/google/gnxi/utils/credentials"
	pb "github.com/openconfig/gnmi/proto/gnmi"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
	"net"
	"net/http"
	"os"
	"ovs-gnxi/generated/ocstruct"
	"ovs-gnxi/target/gnxi"
	"ovs-gnxi/target/gnxi/gnmi"
	"ovs-gnxi/target/logging"
	"ovs-gnxi/target/ovs"
	"reflect"
	"sync"
)

var log = logging.New("ovs-gnxi")

func RunPrometheus(instance *PrometheusMonitoringInstance) {
	instance.StartPrometheus()
	log.Error("Prometheus exit")
}

func RunOVSClient(wg *sync.WaitGroup, client *ovs.Client) {
	defer client.Connection.Disconnect()
	defer wg.Done()
	client.StartMonitorAll()
	log.Error("OVS Client exit")
}

func CreateOVSClient(broker *SystemBroker, config *ovs.Config) *ovs.Client {
	log.Info("Initializing OVS Client...")

	client, err := ovs.NewClient("ovs.gnxi.lan", "tcp", "6640", "certs/target.key", "certs/target.crt", "certs/ca.crt", config)
	if err != nil {
		log.Errorf("Unable to initialize OVS Client: %v", err)
		os.Exit(1)
	}

	return client
}

func setGNMIServerFlags() {
	err := flag.Set("ca", "certs/ca.crt")
	if err != nil {
		log.Fatalf("Unable to set ca flag: %v", err)
	}
	err = flag.Set("cert", "certs/target.crt")
	if err != nil {
		log.Fatalf("Unable to set cert flag: %v", err)
	}
	err = flag.Set("key", "certs/target.key")
	if err != nil {
		log.Fatalf("Unable to set key flag: %v", err)
	}

	flag.Parse()
}

func CreateGNMIServer(client *ovs.Client, broker *SystemBroker) *gnmi.Server {
	<-client.Config.Initialized

	setGNMIServerFlags()

	model := gnmi.NewModel(modeldata.ModelData,
		reflect.TypeOf((*ocstruct.Device)(nil)),
		ocstruct.SchemaTree["Device"],
		ocstruct.Unmarshal,
		ocstruct.ΛEnum)

	config, err := gnxi.GenerateConfig(client.Config)
	if err != nil {
		log.Fatalf("Unable to generate gNMI Config: %v", err)
	}

	log.Info(fmt.Sprintf("%s", config))

	s, err := gnmi.NewServer(model, []byte(config), nil)
	if err != nil {
		log.Fatalf("Error on creating gNMI target: %v", err)
	}

	return s
}

func RunGNMIServer(wg *sync.WaitGroup, server *gnmi.Server) {
	defer wg.Done()

	opts := credentials.ServerCredentials()
	g := grpc.NewServer(opts...)

	pb.RegisterGNMIServer(g, server)
	reflection.Register(g)

	log.Infof("Starting to listen")
	listen, err := net.Listen("tcp", ":10161")
	if err != nil {
		log.Fatalf("Failed to listen: %v", err)
	}

	log.Info("Starting to serve")
	if err := g.Serve(listen); err != nil {
		log.Fatalf("Failed to serve: %v", err)
	}

	log.Error("GNMI Server exit")
}

func main() {
	defer os.Exit(0)
	defer log.Info("Exiting Open vSwitch gNXI interface\n")

	log.Info("Starting Open vSwitch gNXI interface\n")

	setGNMIServerFlags()

	var wg sync.WaitGroup
	wg.Add(2)

	prometheusInstance, err := NewPrometheusMonitoringInstance("0.0.0.0", "8080")
	if err != nil {
		log.Errorf("Unable to configure Prometheus Monitoring: %v", err)
		os.Exit(1)
	}

	go RunPrometheus(prometheusInstance)

	broker := NewSystemBroker()

	ovsConfig := ovs.NewConfig(nil)

	client := CreateOVSClient(broker, ovsConfig)

	go RunOVSClient(&wg, client)

	gnmiServer := CreateGNMIServer(client, broker)

	go RunGNMIServer(&wg, gnmiServer)

	broker.GNMIServer = gnmiServer
	broker.OVSClient = client

	gnmiServer.OverwriteCallback(broker.GNMIConfigChangeCallback)
	ovsConfig.OverwriteCallback(broker.OVSConfigChangeCallback)

	wg.Wait()
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
