package testing

import (
	"fmt"
	"github.com/openconfig/ygot/ygot"
	oc "ovs-gnxi/generated/ocstruct"
	"testing"
)

func TestGenerateConfig(t *testing.T) {
	// Create a new device which is named according to the fake root specified above. To generate
	// the fakeroot then generate_fakeroot should be specified. This entity corresponds to the
	// root of the YANG schema tree. The fakeroot name is the CamelCase version of the name
	// supplied by the fakeroot_name argument.
	d := &oc.Device{}

	// Create a new interface called "eth0"
	i, err := d.Interfaces.NewInterface("eth0")
	fmt.Printf("%v", i)

	// Set the fields that are within the struct.
	i.State.AdminStatus = oc.OpenconfigInterfaces_Interfaces_Interface_State_AdminStatus_UP
	i.Config.Mtu = ygot.Uint16(1500)
	i.Config.Description = ygot.String("An Interface")

	if err := d.Interfaces.Interface["eth0"].Validate(); err != nil {
		panic(fmt.Sprintf("Interface validation failed: %v", err))
	}

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
		t.Fatalf("error in generating json: %v", err)
	}

	fmt.Printf(j)
}
