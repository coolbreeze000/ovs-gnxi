/* Copyright 2018 Google Inc.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    https://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package main

import (
	"flag"
	"fmt"
	"github.com/google/gnxi/utils/credentials"
	"github.com/openconfig/gnmi/cache"
	pb "github.com/openconfig/gnmi/proto/gnmi"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
	"net"
	"net/http"
	"os"
	"ovs-gnxi/shared/logging"
	"ovs-gnxi/target/gnxi"
	"ovs-gnxi/target/gnxi/gnmi"
	"ovs-gnxi/target/gnxi/gnmi/modeldata"
	"ovs-gnxi/target/gnxi/gnmi/modeldata/generated/ocstruct"
	"ovs-gnxi/target/gnxi/gnoi"
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
	err = flag.Set("username", "admin")
	if err != nil {
		log.Fatalf("Unable to set username flag: %v", err)
	}
	err = flag.Set("password", "testpassword")
	if err != nil {
		log.Fatalf("Unable to set password flag: %v", err)
	}

	flag.Parse()
}

func CreateGNMIService(client *ovs.Client) *gnmi.Service {
	<-client.Config.Initialized

	setGNMIServerFlags()

	model := gnmi.NewModel(modeldata.ModelData,
		reflect.TypeOf((*ocstruct.Device)(nil)),
		ocstruct.SchemaTree["Device"],
		ocstruct.Unmarshal,
		ocstruct.ΛEnum)

	config, err := gnmi.GenerateConfig(client.Config)
	if err != nil {
		log.Fatalf("Unable to generate gNMI Config: %v", err)
	}

	log.Info(fmt.Sprintf("%s", config))

	s, err := gnmi.NewService(model, []byte(config), nil, cache.New([]string{"*"}), 8)
	if err != nil {
		log.Fatalf("Error on creating gNMI service: %v", err)
	}

	return s
}

func CreateGNOIService() *gnoi.Service {
	s, err := gnoi.NewService(nil, nil)
	if err != nil {
		log.Fatalf("Error on creating gNOI service: %v", err)
	}

	return s
}

func RunGNMIService(wg *sync.WaitGroup, service *gnmi.Service) {
	defer wg.Done()

	opts := credentials.ServerCredentials()
	g := grpc.NewServer(opts...)

	pb.RegisterGNMIServer(g, service)
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

	gnmiService := CreateGNMIService(client)
	gnoiService := CreateGNOIService()

	go RunGNMIService(&wg, gnmiService)

	broker.GNXIServer = gnxi.NewServer(gnmiService, gnoiService)
	broker.OVSClient = client

	gnmiService.OverwriteCallback(broker.GNMIConfigChangeCallback)
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
