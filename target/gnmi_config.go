package main

import (
	"fmt"
	"github.com/openconfig/ygot/ygot"
	"os"
	oc "ovs-gnxi/generated/ocstruct"
)

func InitializeConfig() {
	// Create a new device which is named according to the fake root specified above. To generate
	// the fakeroot then generate_fakeroot should be specified. This entity corresponds to the
	// root of the YANG schema tree. The fakeroot name is the CamelCase version of the name
	// supplied by the fakeroot_name argument.
	d := &oc.Device{
		System: &oc.System{
			Hostname: ygot.String("ovs.gnxi.lan"),
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
		log.Fatalf("error on generating config json: %v", err)
	}

	if _, err := os.Stat("config.json"); err != nil {
		if os.IsNotExist(err) {
			f, err := os.Create("config.json")
			if err != nil {
				log.Fatalf("error on creating config file: %v", err)
			}

			_, err = f.WriteString(j)
			if err != nil {
				log.Fatalf("error on writing config json: %v", err)
			}

			defer f.Close()
		} else {
			log.Fatalf("error: %v", err)
		}
	}
}
