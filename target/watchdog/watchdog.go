/* Copyright 2019 Google Inc.

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
	go w.gnxiServer.SystemBroker.OVSClient.StartClient(w.gnxiServer.CertManager.GetActivePackageKeyPath(), w.gnxiServer.CertManager.GetActivePackageCertPath(), w.gnxiServer.CertManager.GetActivePackageCAPath())
}

func (w *Watchdog) runGNXIService() {
	go w.gnxiServer.Service.StartService()
}
