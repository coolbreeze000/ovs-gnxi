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

package ovs

import (
	"github.com/openconfig/ygot/ygot"
	"os"
	"ovs-gnxi/shared"
	oc "ovs-gnxi/shared/gnmi/modeldata/generated/ocstruct"
	gnxi "ovs-gnxi/target/gnxi/service"
	"strings"
)

const (
	ovsAddress  = "target.gnxi.lan"
	ovsProtocol = "tcp"
	ovsPort     = "6640"
)

type SystemBroker struct {
	certs                *shared.ServerCertificates
	GNXIService          *gnxi.Service
	OVSClient            *Client
	startOVSClientChan   chan bool
	startGNXIServiceChan chan bool
	stopOVSClientChan    chan bool
	stopGNXIServiceChan  chan bool
}

func NewSystemBroker(gnxiService *gnxi.Service, certs *shared.ServerCertificates) *SystemBroker {
	var err error
	s := &SystemBroker{GNXIService: gnxiService, certs: certs}

	log.Info("Initializing OVS Client...")

	s.OVSClient, err = NewClient(ovsAddress, ovsProtocol, ovsPort)
	if err != nil {
		log.Errorf("Unable to initialize OVS Client: %v", err)
		os.Exit(1)
	}

	return s
}

func (s *SystemBroker) RegisterWatchdogChannels(startOVSClientChan, startGNXIServiceChan, stopOVSClientChan, stopGNXIServiceChan chan bool) {
	s.startOVSClientChan = startOVSClientChan
	s.startGNXIServiceChan = startGNXIServiceChan
	s.stopOVSClientChan = stopOVSClientChan
	s.stopGNXIServiceChan = stopGNXIServiceChan
}

func (s *SystemBroker) GenerateConfig(config *Config) ([]byte, error) {
	d := &oc.Device{
		System: &oc.System{
			Hostname: ygot.String(config.ObjCache.System.Hostname),
			Openflow: &oc.System_Openflow{},
		},
	}

	v, err := d.NewComponent("os")
	if err != nil {
		return []byte(""), err
	}

	v.Type = &oc.Component_Type_Union_E_OpenconfigPlatformTypes_OPENCONFIG_SOFTWARE_COMPONENT{
		E_OpenconfigPlatformTypes_OPENCONFIG_SOFTWARE_COMPONENT: oc.OpenconfigPlatformTypes_OPENCONFIG_SOFTWARE_COMPONENT_OPERATING_SYSTEM,
	}
	v.Description = ygot.String(config.ObjCache.System.Version)

	for _, i := range config.ObjCache.Interfaces {
		o, err := d.NewInterface(i.Name)
		if err != nil {
			return []byte(""), err
		}

		switch adminStatus := i.AdminStatus; adminStatus {
		case "up":
			o.AdminStatus = oc.OpenconfigInterfaces_Interface_AdminStatus_UP
		case "down":
			o.AdminStatus = oc.OpenconfigInterfaces_Interface_AdminStatus_DOWN
		default:
			o.AdminStatus = oc.OpenconfigInterfaces_Interface_AdminStatus_UNSET
		}

		switch linkStatus := i.LinkStatus; linkStatus {
		case "up":
			o.OperStatus = oc.OpenconfigInterfaces_Interface_OperStatus_UP
		case "down":
			o.OperStatus = oc.OpenconfigInterfaces_Interface_OperStatus_DOWN
		default:
			o.OperStatus = oc.OpenconfigInterfaces_Interface_OperStatus_UNSET
		}

		o.Mtu = ygot.Uint16(i.MTU)

		o.Counters = &oc.Interface_Counters{
			InPkts:      ygot.Uint64(i.Statistics.ReceivedPackets),
			InErrors:    ygot.Uint64(i.Statistics.ReceivedErrors),
			InDiscards:  ygot.Uint64(i.Statistics.ReceivedDropped),
			OutPkts:     ygot.Uint64(i.Statistics.TransmittedPackets),
			OutErrors:   ygot.Uint64(i.Statistics.TransmittedErrors),
			OutDiscards: ygot.Uint64(i.Statistics.TransmittedDropped),
		}

		if err := d.Interface[i.Name].Validate(); err != nil {
			return []byte(""), err
		}
	}

	for _, i := range config.ObjCache.Controllers {
		c, err := d.System.Openflow.NewController(i.Name)
		if err != nil {
			return []byte(""), err
		}

		n, err := c.NewConnection(0)
		n.Address = ygot.String(i.Target.Address)
		n.Port = ygot.Uint16(i.Target.Port)
		n.Connected = ygot.Bool(i.Connected)

		switch protocol := strings.ToLower(i.Target.Protocol); protocol {
		case "tcp":
			n.Transport = oc.OpenconfigOpenflow_Transport_TCP
		case "tls":
			n.Transport = oc.OpenconfigOpenflow_Transport_TLS
		default:
			n.Transport = oc.OpenconfigOpenflow_Transport_UNSET
		}
	}

	j, err := ygot.EmitJSON(d, &ygot.EmitJSONConfig{
		Format: ygot.RFC7951,
		Indent: "  ",
		RFC7951Config: &ygot.RFC7951JSONConfig{
			AppendModuleName: true,
		},
	})
	if err != nil {
		return []byte(""), err
	}

	return []byte(j), nil
}

func (s *SystemBroker) OVSConfigChangeCallback(ovsConfig *Config) error {
	log.Debug("Received new change by OVS device")
	gnmiConfig, err := s.GenerateConfig(ovsConfig)
	if err != nil {
		log.Errorf("Unable to generate gNMI config from OVS config source: %v", err)
		return err
	}

	if s.GNXIService != nil {
		s.GNXIService.OverwriteConfig(gnmiConfig)

		select {
		case s.GNXIService.ConfigUpdate <- true: // Send Config Update Notification, unless one already pending.
		default:
		}

		log.Debugf("Using following config data: %s", gnmiConfig)
	}

	return nil
}

func (s *SystemBroker) GNMIConfigSetupCallback(new ygot.ValidatedGoStruct) error {
	log.Debug("Received initial config by gNMI target")

	jsonConfig, err := ygot.ConstructIETFJSON(new, &ygot.RFC7951JSONConfig{
		AppendModuleName: true,
	})
	if err != nil {
		log.Errorf("unable to generate JSON config from gNMI config source: %v", err)
		return err
	}

	cache := CopyConfigObjectCache(s.OVSClient.Config.ObjCache)
	OverwriteObjectCacheWithJSON(cache, jsonConfig)
	s.OVSClient.Config.OverwriteObjectCache(cache)

	return nil
}

func (s *SystemBroker) GNMIConfigChangeCallback(new ygot.ValidatedGoStruct) error {
	log.Debug("Received new change by gNMI target")

	jsonConfigNew, err := ygot.ConstructIETFJSON(new, &ygot.RFC7951JSONConfig{
		AppendModuleName: true,
	})
	if err != nil {
		log.Errorf("unable to generate JSON config from gNMI config source: %v", err)
		return err
	}

	prevCache := CopyConfigObjectCache(s.OVSClient.Config.ObjCache)
	newCache := CopyConfigObjectCache(s.OVSClient.Config.ObjCache)
	OverwriteObjectCacheWithJSON(newCache, jsonConfigNew)

	s.OVSClient.Config.OverwriteObjectCache(newCache)

	err = s.OVSClient.SyncChangesToRemote(prevCache, s.OVSClient.Config.ObjCache)
	if err != nil {
		log.Errorf("unable to sync changes to OVS system: %v", err)
		return err
	}

	return nil
}

func (s *SystemBroker) GNOIRebootCallback() error {
	log.Debug("Received OVS reboot request by GNOI target")
	s.OVSClient.StopMonitoring()
	err := s.OVSClient.RestartSystem()
	if err != nil {
		log.Errorf("unable to reboot OVS system: %v", err)
		return err
	}

	s.stopOVSClientChan <- true
	s.startOVSClientChan <- true
	s.stopGNXIServiceChan <- true
	s.startGNXIServiceChan <- true

	return nil
}

func (s *SystemBroker) GNOIRotateCertificatesCallback(certs *shared.ServerCertificates) error {
	return nil
}
