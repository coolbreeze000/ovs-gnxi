package main

import (
	"github.com/openconfig/ygot/ygot"
	"ovs-gnxi/target/gnxi"
	"ovs-gnxi/target/gnxi/gnmi"
	"ovs-gnxi/target/ovs"
	"sync"
)

type SystemBroker struct {
	OVSClient      *ovs.Client
	OVSConfigLock  sync.RWMutex
	GNMIServer     *gnmi.Server
	GNMIConfigLock sync.RWMutex
}

func NewSystemBroker() *SystemBroker {
	b := SystemBroker{}
	return &b
}

func (s *SystemBroker) OVSConfigChangeCallback(ovsConfig *ovs.Config) error {
	log.Debug("Received new change by OVS device")
	gnmiConfig, err := gnxi.GenerateConfig(ovsConfig)
	if err != nil {
		log.Errorf("Unable to generate gNMI config from OVS config source: %v", err)
		return err
	}
	s.GNMIServer.OverwriteConfig(gnmiConfig)

	return nil
}

func (s *SystemBroker) GNMIConfigChangeCallback(config ygot.ValidatedGoStruct) error {
	return nil
}
