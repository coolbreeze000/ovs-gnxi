package gnxi

import (
	"github.com/openconfig/ygot/ygot"
	oc "ovs-gnxi/generated/ocstruct"
	"ovs-gnxi/target/logging"
	"ovs-gnxi/target/ovs"
)

var log = logging.New("ovs-gnxi")

func GenerateConfig(config *ovs.Config) (string, error) {
	d := &oc.Device{
		System: &oc.System{
			Hostname: ygot.String(config.Cache.System.Hostname),
			Openflow: &oc.System_Openflow{},
		},
	}

	s, err := d.NewComponent("os")
	if err != nil {
		return "", err
	}

	s.Type = &oc.Component_Type_Union_E_OpenconfigPlatformTypes_OPENCONFIG_SOFTWARE_COMPONENT{
		E_OpenconfigPlatformTypes_OPENCONFIG_SOFTWARE_COMPONENT: oc.OpenconfigPlatformTypes_OPENCONFIG_SOFTWARE_COMPONENT_OPERATING_SYSTEM,
	}
	s.Description = ygot.String(config.Cache.System.Version)

	for _, i := range config.Cache.Interfaces {
		o, err := d.NewInterface(i.Name)
		if err != nil {
			return "", err
		}

		switch status := i.Status; status {
		case "up":
			o.AdminStatus = oc.OpenconfigInterfaces_Interface_AdminStatus_UP
		case "down":
			o.AdminStatus = oc.OpenconfigInterfaces_Interface_AdminStatus_DOWN
		default:
			o.AdminStatus = oc.OpenconfigInterfaces_Interface_AdminStatus_UNSET
		}

		o.Mtu = ygot.Uint16(i.MTU)

		if err := d.Interface[i.Name].Validate(); err != nil {
			return "", err
		}
	}

	for _, i := range config.Cache.Controllers {
		c, err := d.System.Openflow.NewController(i.Name)
		if err != nil {
			return "", err
		}
		n, err := c.NewConnection(0)
		n.Address = ygot.String(i.Address)
		n.Port = ygot.Uint16(i.Port)

		switch protocol := i.Protocol; protocol {
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
		return "", err
	}

	return j, nil
}
