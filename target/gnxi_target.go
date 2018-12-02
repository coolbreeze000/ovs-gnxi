package main

import (
	"flag"
	"fmt"
	"github.com/google/gnxi/gnmi"
	"github.com/google/gnxi/gnmi/modeldata"
	"github.com/google/gnxi/utils/credentials"
	pb "github.com/openconfig/gnmi/proto/gnmi"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/reflection"
	"google.golang.org/grpc/status"
	"net"
	"net/http"
	"os"
	"ovs-gnxi/generated/ocstruct"
	"ovs-gnxi/target/gnxi"
	"ovs-gnxi/target/logging"
	"ovs-gnxi/target/ovs"
	"reflect"
	"sync"
)

var log = logging.New("ovs-gnxi")

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

// Get overrides the Get func of gnmi.Target to provide user auth.
func (s *server) Get(ctx context.Context, req *pb.GetRequest) (*pb.GetResponse, error) {
	msg, ok := credentials.AuthorizeUser(ctx)
	if !ok {
		log.Infof("denied a Get request: %v", msg)
		return nil, status.Error(codes.PermissionDenied, msg)
	}
	log.Infof("allowed a Get request: %v", msg)
	return s.Server.Get(ctx, req)
}

// Set overrides the Set func of gnmi.Target to provide user auth.
func (s *server) Set(ctx context.Context, req *pb.SetRequest) (*pb.SetResponse, error) {
	msg, ok := credentials.AuthorizeUser(ctx)
	if !ok {
		log.Infof("denied a Set request: %v", msg)
		return nil, status.Error(codes.PermissionDenied, msg)
	}
	log.Infof("allowed a Set request: %v", msg)
	return s.Server.Set(ctx, req)
}

func RunPrometheus(wg *sync.WaitGroup, instance *PrometheusMonitoringInstance) {
	defer wg.Done()
	instance.StartPrometheus()
	log.Error("Prometheus exit")
}

func RunOVSClient(wg *sync.WaitGroup, client *ovs.Client) {
	defer wg.Done()
	client.StartMonitorAll()
	log.Error("OVS Client exit")
}

func RunGNMIServer(wg *sync.WaitGroup, client *ovs.Client) {
	defer wg.Done()
	<-client.Config.Initialized

	//
	//
	//
	// https://github.com/google/link022/blob/9dbee2acc0d1e02987f94a2c9ced45aa7fbfe91b/agent/gnmi/handler.go
	// https://github.com/google/link022/blob/master/agent/gnmi/server.go

	flag.Set("ca", "certs/ca.crt")
	flag.Set("cert", "certs/target.crt")
	flag.Set("key", "certs/target.key")

	//
	// Test https://github.com/google/gnxi/blob/master/gnmi_target/gnmi_target.go
	//

	model := gnmi.NewModel(modeldata.ModelData,
		reflect.TypeOf((*ocstruct.Device)(nil)),
		ocstruct.SchemaTree["Device"],
		ocstruct.Unmarshal,
		ocstruct.ΛEnum)

	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, "Supported models:\n")
		for _, m := range model.SupportedModels() {
			fmt.Fprintf(os.Stderr, "  %s\n", m)
		}
		fmt.Fprintf(os.Stderr, "\n")
		fmt.Fprintf(os.Stderr, "Usage of %s:\n", os.Args[0])
		flag.PrintDefaults()
	}

	flag.Parse()

	opts := credentials.ServerCredentials()
	g := grpc.NewServer(opts...)

	config, err := gnxi.GenerateConfig(client.Config)
	if err != nil {
		log.Fatalf("Unable to generate gNMI Config: %v", err)
	}

	log.Info(config)

	s, err := newServer(model, []byte(config))
	if err != nil {
		log.Fatalf("Error on creating gNMI target: %v", err)
	}
	pb.RegisterGNMIServer(g, s)
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

	//
	//
	//

	log.Error("GNMI Server exit")
}

func main() {
	defer os.Exit(0)
	defer log.Info("Exiting Open vSwitch gNXI interface\n")

	log.Info("Starting Open vSwitch gNXI interface\n")

	var wg sync.WaitGroup
	wg.Add(3)

	prometheusInstance, err := NewPrometheusMonitoringInstance("0.0.0.0", "8080")
	if err != nil {
		log.Errorf("Unable to configure Prometheus Monitoring: %v", err)
		os.Exit(1)
	}

	go RunPrometheus(&wg, prometheusInstance)

	log.Info("Initializing OVS Client...")

	client, err := ovs.NewClient("ovs.gnxi.lan", "tcp", "6640", "certs/target.key", "certs/target.crt", "certs/ca.crt")
	if err != nil {
		log.Errorf("Unable to initialize OVS Client: %v", err)
		os.Exit(1)
	}

	defer client.Connection.Disconnect()

	go RunOVSClient(&wg, client)

	go RunGNMIServer(&wg, client)

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
