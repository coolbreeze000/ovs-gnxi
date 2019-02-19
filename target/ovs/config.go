/* Copyright 2019 Google Inc.

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

package ovs

import (
	"errors"
	"fmt"
	"github.com/socketplane/libovsdb"
	"reflect"
	"strconv"
	"strings"
	"sync"
)

const (
	primaryControllerName = "main"
)

type ConfigCallback func(config *Config) error

type System struct {
	uuid     string
	Version  string
	Hostname string
}

func (s *System) Equal(comp *System) bool {
	switch {
	case s.Hostname != comp.Hostname:
		return false
	default:
		return true
	}
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
	uuid      string
	Name      string
	Connected bool
	Target    *OpenFlowControllerTarget
}

func (c *OpenFlowController) Equal(comp *OpenFlowController) bool {
	log.Debugf("COMPARE: %v against %v", c, comp)

	switch {
	case c.Name != comp.Name:
		return false
	case c.Target.Protocol != comp.Target.Protocol:
		return false
	case c.Target.Port != comp.Target.Port:
		return false
	case c.Target.Address != comp.Target.Address:
		return false
	default:
		return true
	}
}

func (c *OpenFlowController) String() string {
	return fmt.Sprintf("OpenFlowController(uuid: \"%v\", Name: \"%v\", Connected: \"%v\", Address: \"%v\", Protocol: \"%v\", Port: \"%v\")", c.uuid, c.Name, c.Connected, c.Target.Address, c.Target.Protocol, c.Target.Port)
}

func ParseOpenFlowControllerTarget(target string) (*OpenFlowControllerTarget, error) {
	s := strings.Split(target, ":")

	port, err := strconv.ParseUint(s[2], 10, 16)
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

func (i *Interface) Equal(comp *Interface) bool {
	switch {
	case i.Name != comp.Name:
		return false
	case i.MTU != comp.MTU:
		return false
	default:
		return true
	}
}

func (i *Interface) String() string {
	return fmt.Sprintf("Interface(uuid: \"%v\", Name: \"%v\", MTU: \"%v\", AdminStatus: \"%v\", LinkStatus: \"%v\", Statistics: \"%v\")", i.uuid, i.Name, i.MTU, i.AdminStatus, i.LinkStatus, i.Statistics)
}

type ObjectCache struct {
	System      *System
	Controllers map[string]*OpenFlowController
	Interfaces  map[string]*Interface
}
type Config struct {
	rawCache    map[string]map[string]libovsdb.Row
	ObjCache    *ObjectCache
	mu          sync.RWMutex
	callback    ConfigCallback
	Initialized chan struct{}
}

func (c *Config) getInterfaceByUUID(uuid string) *Interface {
	for _, i := range c.ObjCache.Interfaces {
		if i.uuid == uuid {
			return i
		}
	}

	return nil
}

func (c *Config) getControllerByUUID(uuid string) *OpenFlowController {
	for _, i := range c.ObjCache.Controllers {
		if i.uuid == uuid {
			return i
		}
	}

	return nil
}

func NewConfig() *Config {
	c := &Config{rawCache: make(map[string]map[string]libovsdb.Row), Initialized: make(chan struct{}),
		ObjCache: &ObjectCache{System: &System{}, Controllers: make(map[string]*OpenFlowController), Interfaces: make(map[string]*Interface)}}
	return c
}

func CopyConfigObjectCache(c *ObjectCache) *ObjectCache {
	cache := &ObjectCache{System: &System{}, Controllers: make(map[string]*OpenFlowController), Interfaces: make(map[string]*Interface)}

	cache.System = &System{
		uuid:     c.System.uuid,
		Version:  c.System.Version,
		Hostname: c.System.Hostname,
	}

	cache.Controllers[primaryControllerName] = &OpenFlowController{
		uuid:      c.Controllers[primaryControllerName].uuid,
		Name:      c.Controllers[primaryControllerName].Name,
		Connected: c.Controllers[primaryControllerName].Connected,
		Target: &OpenFlowControllerTarget{
			Address:  c.Controllers[primaryControllerName].Target.Address,
			Port:     c.Controllers[primaryControllerName].Target.Port,
			Protocol: c.Controllers[primaryControllerName].Target.Protocol,
		},
	}

	for _, i := range c.Interfaces {
		cache.Interfaces[i.Name] = &Interface{
			uuid:        i.uuid,
			Name:        i.Name,
			MTU:         i.MTU,
			AdminStatus: i.AdminStatus,
			LinkStatus:  i.LinkStatus,
			Statistics: &InterfaceStatistics{
				ReceivedPackets:    i.Statistics.ReceivedPackets,
				ReceivedErrors:     i.Statistics.ReceivedErrors,
				ReceivedDropped:    i.Statistics.ReceivedDropped,
				TransmittedPackets: i.Statistics.TransmittedPackets,
				TransmittedErrors:  i.Statistics.TransmittedErrors,
				TransmittedDropped: i.Statistics.TransmittedDropped,
			},
		}
	}

	return cache
}

func OverwriteObjectCacheWithJSON(cache *ObjectCache, jsonConfig map[string]interface{}) {
	for _, i := range jsonConfig["openconfig-platform:components"].(map[string]interface{})["component"].([]interface{}) {
		if i.(map[string]interface{})["config"].(map[string]interface{})["name"] == "os" {
			cache.System.Version = i.(map[string]interface{})["state"].(map[string]interface{})["description"].(string)
		}
	}

	cache.System.Hostname = jsonConfig["openconfig-system:system"].(map[string]interface{})["config"].(map[string]interface{})["hostname"].(string)

	for _, i := range jsonConfig["openconfig-system:system"].(map[string]interface{})["openconfig-openflow:openflow"].(map[string]interface{})["controllers"].(map[string]interface{})["controller"].([]interface{}) {
		name := i.(map[string]interface{})["config"].(map[string]interface{})["name"].(string)

		for _, j := range i.(map[string]interface{})["connections"].(map[string]interface{})["connection"].([]interface{}) {
			cache.Controllers[name].Target.Address = j.(map[string]interface{})["config"].(map[string]interface{})["address"].(string)
			cache.Controllers[name].Target.Port = j.(map[string]interface{})["config"].(map[string]interface{})["port"].(uint16)
			cache.Controllers[name].Target.Protocol = strings.ToLower(j.(map[string]interface{})["config"].(map[string]interface{})["transport"].(string))
			cache.Controllers[name].Connected = j.(map[string]interface{})["state"].(map[string]interface{})["connected"].(bool)
		}
	}

	for _, i := range jsonConfig["openconfig-interfaces:interfaces"].(map[string]interface{})["interface"].([]interface{}) {
		name := i.(map[string]interface{})["config"].(map[string]interface{})["name"].(string)

		cache.Interfaces[name].Name = name
		cache.Interfaces[name].MTU = i.(map[string]interface{})["config"].(map[string]interface{})["mtu"].(uint16)

		if _, ok := i.(map[string]interface{})["state"].(map[string]interface{})["admin-status"]; ok {
			cache.Interfaces[name].AdminStatus = i.(map[string]interface{})["state"].(map[string]interface{})["admin-status"].(string)
		}

		if _, ok := i.(map[string]interface{})["state"].(map[string]interface{})["oper-status"]; ok {
			cache.Interfaces[name].LinkStatus = i.(map[string]interface{})["state"].(map[string]interface{})["oper-status"].(string)
		}

		if inPkts, err := strconv.ParseUint(i.(map[string]interface{})["state"].(map[string]interface{})["counters"].(map[string]interface{})["in-pkts"].(string), 10, 64); err == nil {
			cache.Interfaces[name].Statistics.ReceivedPackets = inPkts
		}

		if inErrs, err := strconv.ParseUint(i.(map[string]interface{})["state"].(map[string]interface{})["counters"].(map[string]interface{})["in-errors"].(string), 10, 64); err == nil {
			cache.Interfaces[name].Statistics.ReceivedErrors = inErrs
		}

		if inDisc, err := strconv.ParseUint(i.(map[string]interface{})["state"].(map[string]interface{})["counters"].(map[string]interface{})["in-discards"].(string), 10, 64); err == nil {
			cache.Interfaces[name].Statistics.ReceivedDropped = inDisc
		}

		if outPkts, err := strconv.ParseUint(i.(map[string]interface{})["state"].(map[string]interface{})["counters"].(map[string]interface{})["out-pkts"].(string), 10, 64); err == nil {
			cache.Interfaces[name].Statistics.TransmittedPackets = outPkts
		}

		if outErrs, err := strconv.ParseUint(i.(map[string]interface{})["state"].(map[string]interface{})["counters"].(map[string]interface{})["out-errors"].(string), 10, 64); err == nil {
			cache.Interfaces[name].Statistics.TransmittedErrors = outErrs
		}

		if outDisc, err := strconv.ParseUint(i.(map[string]interface{})["state"].(map[string]interface{})["counters"].(map[string]interface{})["out-discards"].(string), 10, 64); err == nil {
			cache.Interfaces[name].Statistics.TransmittedDropped = outDisc
		}
	}
}

func (c *Config) InitializeCache(updates *libovsdb.TableUpdates) {
	c.SyncCache(updates)
	close(c.Initialized)
}

func (c *Config) SyncCache(updates *libovsdb.TableUpdates) {
	c.mu.Lock()
	defer c.mu.Unlock()

	log.Debug("Syncing config cache...")

	for tableName, tableUpdate := range updates.Updates {
		c.initializeCacheTableIfNotExists(tableName)

		for uuid, row := range tableUpdate.Rows {
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

	if c.callback != nil {
		if err := c.callback(c); err != nil {
			log.Errorf("Config callback error: %v", err)
		}
	}

	c.DumpRawCache()
	c.DumpObjectCache()

	log.Debug("Syncing config cache complete")
}

func (c *Config) UpdateObjectCacheEntry(tableName, uuid string, row libovsdb.Row) error {
	switch tableName {
	case SystemTable:
		c.ObjCache.System = &System{
			uuid:     uuid,
			Version:  row.Fields["ovs_version"].(string),
			Hostname: row.Fields["external_ids"].(libovsdb.OvsMap).GoMap["hostname"].(string),
		}
	case ControllerTable:
		target, err := ParseOpenFlowControllerTarget(row.Fields["target"].(string))
		if err != nil {
			return err
		}

		c.ObjCache.Controllers[primaryControllerName] = &OpenFlowController{
			uuid:      uuid,
			Name:      primaryControllerName,
			Connected: row.Fields["is_connected"].(bool),
			Target:    target,
		}
	case InterfaceTable:
		var mtu uint16

		switch row.Fields["mtu"].(type) {
		case float64:
			mtu = uint16(row.Fields["mtu"].(float64))
		case libovsdb.OvsSet:
			log.Errorf("Unable to perform correct type conversion OvsSet for interface mtu: %v", row)
		default:
			log.Errorf("Unable to perform correct type conversion for interface mtu: %v", row)
		}

		c.ObjCache.Interfaces[row.Fields["name"].(string)] = &Interface{
			uuid:        uuid,
			Name:        row.Fields["name"].(string),
			MTU:         mtu,
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
		if c.ObjCache.System.uuid == uuid {
			c.ObjCache.System = nil
		}
	case ControllerTable:
		if c.getControllerByUUID(uuid) != nil {
			delete(c.ObjCache.Controllers, uuid)
		}
	case InterfaceTable:
		if c.getInterfaceByUUID(uuid) != nil {
			delete(c.ObjCache.Interfaces, uuid)
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
	log.Debug(c.ObjCache)
}

func (c *Config) OverwriteObjectCache(cache *ObjectCache) {
	c.mu.Lock()
	defer c.mu.Unlock()

	c.ObjCache = cache
}

func (c *Config) OverwriteCallback(callback ConfigCallback) {
	c.mu.Lock()
	defer c.mu.Unlock()

	c.callback = callback
}
