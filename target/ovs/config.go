package ovs

import (
	"fmt"
	"github.com/socketplane/libovsdb"
	"reflect"
)

type System struct {
	uuid     string
	Version  string
	Hostname string
}

func (s *System) String() string {
	return fmt.Sprintf("System(uuid: \"%v\", Version: \"%v\", Hostname: \"%v\")", s.uuid, s.Version, s.Hostname)
}

type OpenFlowController struct {
	uuid     string
	Name     string
	Address  string
	Protocol string
	Port     uint16
}

type Interface struct {
	uuid       string
	Name       string
	MTU        uint16
	Status     string
	Statistics string
}

func (c *OpenFlowController) String() string {
	return fmt.Sprintf("OpenFlowController(uuid: \"%v\", Name: \"%v\", Address: \"%v\", Protocol: \"%v\", Port: \"%v\")", c.uuid, c.Name, c.Address, c.Protocol, c.Port)
}

type Config struct {
	rawCache map[string]map[string]libovsdb.Row
	Cache    struct {
		System      *System
		Controllers []*OpenFlowController
		Interfaces  []*Interface
	}
}

func NewConfig() *Config {
	c := Config{rawCache: make(map[string]map[string]libovsdb.Row)}
	return &c
}

func (c *Config) SyncCache(updates libovsdb.TableUpdates) {
	for tableName, tableUpdate := range updates.Updates {
		c.initializeCacheTableIfNotExists(tableName)

		for uuid, row := range tableUpdate.Rows {
			empty := libovsdb.Row{}
			if !reflect.DeepEqual(row.New, empty) {
				c.rawCache[tableName][uuid] = row.New
			} else {
				delete(c.rawCache[tableName], uuid)
			}
		}
	}
}

func (c *Config) initializeCacheTableIfNotExists(tableName string) {
	if _, ok := c.rawCache[tableName]; !ok {
		c.rawCache[tableName] = make(map[string]libovsdb.Row)
	}
}

func (c *Config) DumpRawCache() {
	log.Debug(c.rawCache)
}
