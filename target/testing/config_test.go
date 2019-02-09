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

package testing

import (
	"fmt"
	"github.com/openconfig/ygot/ygot"
	"os"
	oc "ovs-gnxi/shared/gnmi/modeldata/generated/ocstruct"
	"testing"
)

func TestGenerateConfig(t *testing.T) {
	// Create a new device which is named according to the fake root specified above. To generate
	// the fakeroot then generate_fakeroot should be specified. This entity corresponds to the
	// root of the YANG schema tree. The fakeroot name is the CamelCase version of the name
	// supplied by the fakeroot_name argument.
	d := &oc.Device{
		System: &oc.System{
			Hostname: ygot.String("target.gnxi.lan"),
			Openflow: &oc.System_Openflow{},
		},
	}

	s, err := d.NewComponent("os")
	s.Type = &oc.Component_Type_Union_E_OpenconfigPlatformTypes_OPENCONFIG_SOFTWARE_COMPONENT{
		E_OpenconfigPlatformTypes_OPENCONFIG_SOFTWARE_COMPONENT: oc.OpenconfigPlatformTypes_OPENCONFIG_SOFTWARE_COMPONENT_OPERATING_SYSTEM,
	}
	s.Description = ygot.String("2.9.0")

	// Create a new interface called "eth0"
	i, err := d.NewInterface("sw1")

	// Set the fields that are within the struct.
	i.AdminStatus = oc.OpenconfigInterfaces_Interface_AdminStatus_UP
	i.Mtu = ygot.Uint16(1500)

	if err := d.Interface["eth0"].Validate(); err != nil {
		panic(fmt.Sprintf("Interface validation failed: %v", err))
	}

	c, err := d.System.Openflow.NewController("053f3666-8e20-4efc-a6a6-dd54a86e360b")
	fmt.Printf("%v", c)
	n, err := c.NewConnection(0)
	n.Address = ygot.String("172.18.0.4")
	n.Port = ygot.Uint16(6653)
	n.Transport = oc.OpenconfigOpenflow_Transport_TCP

	// To render the device (which is currently empty) to JSON in RFC7951 format, then we
	// simply call the ygot.EmitJSON method with the relevant arguments.
	j, err := ygot.EmitJSON(d, &ygot.EmitJSONConfig{
		Format: ygot.RFC7951,
		Indent: "  ",
		RFC7951Config: &ygot.RFC7951JSONConfig{
			AppendModuleName: true,
		},
	})
	if err != nil {
		t.Fatalf("error on generating config json: %v", err)
	}

	if _, err := os.Stat("config.json"); err != nil {
		if os.IsNotExist(err) {
			f, err := os.Create("config.json")
			if err != nil {
				t.Fatalf("error on creating config file: %v", err)
			}

			_, err = f.WriteString(j)
			if err != nil {
				t.Fatalf("error on writing config json: %v", err)
			}

			defer f.Close()
		} else {
			t.Fatalf("error: %v", err)
		}
	}
}
