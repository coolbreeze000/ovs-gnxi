package watchdog

import (
	"ovs-gnxi/shared/logging"
	"ovs-gnxi/target/gnxi"
)

var log = logging.New("ovs-gnxi")

type Watchdog struct {
	gnxiServer           *gnxi.Server
	startOVSClientChan   chan bool
	startGNXIServiceChan chan bool
	stopOVSClientChan    chan bool
	stopGNXIServiceChan  chan bool
}

func NewWatchdog(gnxiServer *gnxi.Server) *Watchdog {
	log.Info("Initializing Watchdog...")
	w := &Watchdog{gnxiServer: gnxiServer, startOVSClientChan: make(chan bool), startGNXIServiceChan: make(chan bool), stopOVSClientChan: make(chan bool), stopGNXIServiceChan: make(chan bool)}
	w.gnxiServer.SystemBroker.RegisterWatchdogChannels(w.startOVSClientChan, w.startGNXIServiceChan, w.stopOVSClientChan, w.stopGNXIServiceChan)
	return w
}

func (w *Watchdog) RunServices() {
	w.runOVSClientService()
	w.gnxiServer.InitializeService()
	w.runGNXIService()

	for {
		select {
		case _ = <-w.startOVSClientChan:
			w.runOVSClientService()
		case _ = <-w.stopOVSClientChan:
			w.gnxiServer.SystemBroker.OVSClient.StopClient()
		case _ = <-w.startGNXIServiceChan:
			w.gnxiServer.InitializeService()
			w.runGNXIService()
		case _ = <-w.stopGNXIServiceChan:
			w.gnxiServer.Service.StopService()
		}
	}
}

func (w *Watchdog) runOVSClientService() {
	go w.gnxiServer.SystemBroker.OVSClient.StartClient(w.gnxiServer.CertManager.GetActivePackage().KeySystemPath, w.gnxiServer.CertManager.GetActivePackage().CertificateSystemPath, w.gnxiServer.CertManager.GetActivePackage().CASystemPath)
}

func (w *Watchdog) runGNXIService() {
	go w.gnxiServer.Service.StartService()
}
