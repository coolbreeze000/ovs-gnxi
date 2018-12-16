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
	Desc      string
	XPaths    []string
	Extractor func(n []*gnmi.Notification) string
	ExpResp   interface{}
}{{
	Desc:      "get system config hostname",
	XPaths:    []string{"/system/config/hostname"},
	Extractor: ExtractSingleStringValueFromResponse,
	ExpResp:   "ovs.gnxi.lan",
}}

func ExtractSingleStringValueFromResponse(n []*gnmi.Notification) string {
	return n[0].Update[0].Val.GetStringVal()
}
