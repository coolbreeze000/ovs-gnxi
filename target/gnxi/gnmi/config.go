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

package gnmi

import (
	"github.com/openconfig/ygot/ygot"
	oc "ovs-gnxi/target/gnxi/gnmi/modeldata/generated/ocstruct"
	"ovs-gnxi/target/ovs"
)

func GenerateConfig(config *ovs.Config) ([]byte, error) {
	log.Info("Start generating initial gNMI config from OVS system source...")
	log.Debugf("Using following initial config data: %v", config.ObjectCache)

	d := &oc.Device{
		System: &oc.System{
			Hostname: ygot.String(config.ObjectCache.System.Hostname),
			Openflow: &oc.System_Openflow{},
		},
	}

	s, err := d.NewComponent("os")
	if err != nil {
		return []byte(""), err
	}

	s.Type = &oc.Component_Type_Union_E_OpenconfigPlatformTypes_OPENCONFIG_SOFTWARE_COMPONENT{
		E_OpenconfigPlatformTypes_OPENCONFIG_SOFTWARE_COMPONENT: oc.OpenconfigPlatformTypes_OPENCONFIG_SOFTWARE_COMPONENT_OPERATING_SYSTEM,
	}
	s.Description = ygot.String(config.ObjectCache.System.Version)

	for _, i := range config.ObjectCache.Interfaces {
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

	for _, i := range config.ObjectCache.Controllers {
		c, err := d.System.Openflow.NewController(i.Name)
		if err != nil {
			return []byte(""), err
		}
		n, err := c.NewConnection(0)
		n.Address = ygot.String(i.Target.Address)
		n.Port = ygot.Uint16(i.Target.Port)
		n.Connected = ygot.Bool(i.Connected)

		switch protocol := i.Target.Protocol; protocol {
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
