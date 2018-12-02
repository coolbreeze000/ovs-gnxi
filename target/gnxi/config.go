package gnxi

import (
	"github.com/openconfig/ygot/ygot"
	oc "ovs-gnxi/generated/ocstruct"
	"ovs-gnxi/target/logging"
	"ovs-gnxi/target/ovs"
)

var log = logging.New("ovs-gnxi")

func GenerateConfig(config *ovs.Config) (string, error) {
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
		return "", err
	}

	s.Type = &oc.Component_Type_Union_E_OpenconfigPlatformTypes_OPENCONFIG_SOFTWARE_COMPONENT{
		E_OpenconfigPlatformTypes_OPENCONFIG_SOFTWARE_COMPONENT: oc.OpenconfigPlatformTypes_OPENCONFIG_SOFTWARE_COMPONENT_OPERATING_SYSTEM,
	}
	s.Description = ygot.String(config.ObjectCache.System.Version)

	for _, i := range config.ObjectCache.Interfaces {
		o, err := d.NewInterface(i.Name)
		if err != nil {
			return "", err
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
			return "", err
		}
	}

	for _, i := range config.ObjectCache.Controllers {
		c, err := d.System.Openflow.NewController(i.Name)
		if err != nil {
			return "", err
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
		return "", err
	}

	return j, nil
}
