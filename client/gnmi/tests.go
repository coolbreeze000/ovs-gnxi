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
	"github.com/openconfig/gnmi/proto/gnmi"
	"github.com/openconfig/gnmi/proto/gnmi_ext"
)

var CapabilitiesTests = []struct {
	Desc          string
	ExpVersion    string
	ExpModels     []*gnmi.ModelData
	ExpEncodings  []gnmi.Encoding
	ExpExtensions []*gnmi_ext.Extension
}{{
	Desc:       "retrieve system capabilities",
	ExpVersion: "0.7.0",
	ExpModels: []*gnmi.ModelData{
		{Name: "openconfig-interfaces", Organization: "OpenConfig working group", Version: "2.0.0"},
		{Name: "openconfig-openflow", Organization: "OpenConfig working group", Version: "0.1.0"},
		{Name: "openconfig-platform", Organization: "OpenConfig working group", Version: "0.5.0"},
		{Name: "openconfig-system", Organization: "OpenConfig working group", Version: "0.2.0"},
	},
	ExpEncodings: []gnmi.Encoding{
		gnmi.Encoding(gnmi.Encoding_JSON),
		gnmi.Encoding(gnmi.Encoding_JSON_IETF),
	},
	ExpExtensions: []*gnmi_ext.Extension{},
}}

var GetTests = []struct {
	Desc            string
	XPaths          []string
	ExtractorString func(n []*gnmi.Notification) string
	ExpResp         interface{}
	ExtractorUInt   func(n []*gnmi.Notification) uint64
	MinResp         interface{}
}{
	{
		Desc:            "get system config hostname",
		XPaths:          []string{"/system/config/hostname"},
		ExtractorString: ExtractSingleStringValueFromResponse,
		ExpResp:         "ovs.gnxi.lan",
	},
	{
		Desc:            "get platform component os state version",
		XPaths:          []string{"/components/component[name=os]/state/description"},
		ExtractorString: ExtractSingleStringValueFromResponse,
		ExpResp:         "2.9.0",
	},
	{
		Desc:            "get system openflow controller connection config address",
		XPaths:          []string{"/system/openflow/controllers/controller[name=main]/connections/connection[aux-id=0]/config/address"},
		ExtractorString: ExtractSingleStringValueFromResponse,
		ExpResp:         "172.18.0.2",
	},
	{
		Desc:          "get system openflow controller connection config port",
		XPaths:        []string{"/system/openflow/controllers/controller[name=main]/connections/connection[aux-id=0]/config/port"},
		ExtractorUInt: ExtractSingleUintValueFromResponse,
		ExpResp:       uint64(6653),
	},
	{
		Desc:          "get interface state counters in-pkts",
		XPaths:        []string{"/interfaces/interface[name=sw1-eth1]/state/counters/in-pkts"},
		ExtractorUInt: ExtractSingleUintValueFromResponse,
		MinResp:       uint64(0),
	},
	{
		Desc:          "get interface state counters out-pkts",
		XPaths:        []string{"/interfaces/interface[name=sw1-eth1]/state/counters/out-pkts"},
		ExtractorUInt: ExtractSingleUintValueFromResponse,
		MinResp:       uint64(0),
	},
}

var SubscribeTests = []struct {
	Desc          string
	XPaths        []string
	ExtractorUInt func(n []*gnmi.Notification) uint64
	MinResp       interface{}
}{
	{
		Desc:          "subscribe to interface state counters in-pkts",
		XPaths:        []string{"/interfaces/interface[name=sw1-eth1]/state/counters/in-pkts"},
		ExtractorUInt: ExtractSingleUintValueFromResponse,
		MinResp:       uint64(0),
	},
	{
		Desc:          "subscribe to interface state counters out-pkts",
		XPaths:        []string{"/interfaces/interface[name=sw1-eth1]/state/counters/out-pkts"},
		ExtractorUInt: ExtractSingleUintValueFromResponse,
		MinResp:       uint64(0),
	},
}

func ExtractSingleStringValueFromResponse(n []*gnmi.Notification) string {
	return n[0].Update[0].Val.GetStringVal()
}

func ExtractSingleUintValueFromResponse(n []*gnmi.Notification) uint64 {
	return n[0].Update[0].Val.GetUintVal()
}