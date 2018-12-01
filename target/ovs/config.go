package ovs

import (
	"errors"
	"fmt"
	"github.com/socketplane/libovsdb"
	"reflect"
	"strconv"
	"strings"
)

type System struct {
	uuid     string
	Version  string
	Hostname string
}

func (s *System) String() string {
	return fmt.Sprintf("System(uuid: \"%v\", Version: \"%v\", Hostname: \"%v\")", s.uuid, s.Version, s.Hostname)
}

type OpenFlowControllerTarget struct {
	Address  string
	Protocol string
	Port     uint16
}

func (t *OpenFlowControllerTarget) String() string {
	return fmt.Sprintf("OpenFlowControllerTarget(Address: \"%v\", Protocol: \"%v\", Port: \"%v\")", t.Address, t.Protocol, t.Port)
}

type OpenFlowController struct {
	uuid   string
	Name   string
	Target *OpenFlowControllerTarget
}

func (c *OpenFlowController) String() string {
	return fmt.Sprintf("OpenFlowController(uuid: \"%v\", Name: \"%v\", Address: \"%v\", Protocol: \"%v\", Port: \"%v\")", c.uuid, c.Name, c.Target.Address, c.Target.Protocol, c.Target.Port)
}

func ParseOpenFlowControllerTarget(target string) (*OpenFlowControllerTarget, error) {
	s := strings.Split(target, ":")

	port, err := strconv.ParseUint(s[2], 16, 16)
	if err != nil {
		return nil, err
	}

	return &OpenFlowControllerTarget{Address: s[1], Protocol: s[0], Port: uint16(port)}, nil
}

type InterfaceStatistics struct {
	ReceivedPackets    uint64
	ReceivedErrors     uint64
	ReceivedDropped    uint64
	TransmittedPackets uint64
	TransmittedErrors  uint64
	TransmittedDropped uint64
}

func (s *InterfaceStatistics) String() string {
	return fmt.Sprintf("InterfaceStatistics(ReceivedPackets: \"%v\", ReceivedErrors: \"%v\", ReceivedDropped: \"%v\", TransmittedPackets: \"%v\", TransmittedErrors: \"%v\", TransmittedDropped: \"%v\")",
		s.ReceivedPackets, s.ReceivedErrors, s.ReceivedDropped, s.TransmittedPackets, s.TransmittedErrors, s.TransmittedDropped)
}

type Interface struct {
	uuid        string
	Name        string
	MTU         uint16
	AdminStatus string
	LinkStatus  string
	Statistics  *InterfaceStatistics
}

func (i *Interface) String() string {
	return fmt.Sprintf("Interface(uuid: \"%v\", Name: \"%v\", MTU: \"%v\", AdminStatus: \"%v\", LinkStatus: \"%v\", Statistics: \"%v\")", i.uuid, i.Name, i.MTU, i.AdminStatus, i.LinkStatus, i.Statistics)
}

type Config struct {
	rawCache    map[string]map[string]libovsdb.Row
	ObjectCache struct {
		System      *System
		Controllers map[string]*OpenFlowController
		Interfaces  map[string]*Interface
	}
	Initialized chan struct{}
}

func NewConfig() *Config {
	c := Config{rawCache: make(map[string]map[string]libovsdb.Row), Initialized: make(chan struct{}), ObjectCache: struct {
		System      *System
		Controllers map[string]*OpenFlowController
		Interfaces  map[string]*Interface
	}{System: &System{}, Controllers: make(map[string]*OpenFlowController), Interfaces: make(map[string]*Interface)}}
	return &c
}

func (c *Config) InitializeCache(updates *libovsdb.TableUpdates) {
	c.SyncCache(updates)
	close(c.Initialized)
}

func (c *Config) SyncCache(updates *libovsdb.TableUpdates) {
	log.Debug("Syncing config cache...")

	for tableName, tableUpdate := range updates.Updates {
		c.initializeCacheTableIfNotExists(tableName)

		for uuid, row := range tableUpdate.Rows {
			log.Debug(row)

			empty := libovsdb.Row{}
			if !reflect.DeepEqual(row.New, empty) {
				c.rawCache[tableName][uuid] = row.New
				err := c.UpdateObjectCacheEntry(tableName, uuid, c.rawCache[tableName][uuid])
				if err != nil {
					log.Error(err)
				}
			} else {
				delete(c.rawCache[tableName], uuid)
				err := c.DeleteObjectCacheEntry(tableName, uuid)
				log.Error(err)
			}
		}
	}

	c.DumpRawCache()
	c.DumpObjectCache()

	log.Debug("Syncing config cache complete")
}

func (c *Config) UpdateObjectCacheEntry(tableName, uuid string, row libovsdb.Row) error {
	switch tableName {
	case SystemTable:
		c.ObjectCache.System = &System{
			uuid:     uuid,
			Version:  row.Fields["ovs_version"].(string),
			Hostname: row.Fields["external_ids"].(libovsdb.OvsMap).GoMap["hostname"].(string),
		}
	case ControllerTable:
		target, err := ParseOpenFlowControllerTarget(row.Fields["target"].(string))
		if err != nil {
			return err
		}

		c.ObjectCache.Controllers[uuid] = &OpenFlowController{
			uuid:   uuid,
			Name:   uuid,
			Target: target,
		}
	case InterfaceTable:
		c.ObjectCache.Interfaces[uuid] = &Interface{
			uuid:        uuid,
			Name:        uuid,
			MTU:         uint16(row.Fields["mtu"].(float64)),
			AdminStatus: row.Fields["admin_state"].(string),
			LinkStatus:  row.Fields["link_state"].(string),
			Statistics: &InterfaceStatistics{
				ReceivedPackets:    uint64(row.Fields["statistics"].(libovsdb.OvsMap).GoMap["rx_packets"].(float64)),
				ReceivedErrors:     uint64(row.Fields["statistics"].(libovsdb.OvsMap).GoMap["rx_errors"].(float64)),
				ReceivedDropped:    uint64(row.Fields["statistics"].(libovsdb.OvsMap).GoMap["rx_dropped"].(float64)),
				TransmittedPackets: uint64(row.Fields["statistics"].(libovsdb.OvsMap).GoMap["tx_packets"].(float64)),
				TransmittedErrors:  uint64(row.Fields["statistics"].(libovsdb.OvsMap).GoMap["tx_errors"].(float64)),
				TransmittedDropped: uint64(row.Fields["statistics"].(libovsdb.OvsMap).GoMap["tx_dropped"].(float64)),
			},
		}
	default:
		return errors.New("unable to update unsupported table entry")
	}

	return nil
}

func (c *Config) DeleteObjectCacheEntry(tableName, uuid string) error {
	switch tableName {
	case SystemTable:
		if c.ObjectCache.System.uuid == uuid {
			c.ObjectCache.System = nil
		}
	case ControllerTable:
		if _, ok := c.ObjectCache.Controllers[uuid]; !ok {
			delete(c.ObjectCache.Controllers, uuid)
		}
	case InterfaceTable:
		if _, ok := c.ObjectCache.Interfaces[uuid]; !ok {
			delete(c.ObjectCache.Interfaces, uuid)
		}
	default:
		return errors.New("unable to delete unsupported table entry")
	}

	return nil
}

func (c *Config) initializeCacheTableIfNotExists(tableName string) {
	if _, ok := c.rawCache[tableName]; !ok {
		c.rawCache[tableName] = make(map[string]libovsdb.Row)
	}
}

func (c *Config) DumpRawCache() {
	log.Debug(c.rawCache)
}

func (c *Config) DumpObjectCache() {
	log.Debug(c.ObjectCache)
}
