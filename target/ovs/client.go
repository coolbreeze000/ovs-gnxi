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

// Contains code from source: https://github.com/socketplane/libovsdb/tree/4de3618546deba09d8875d719752db32bd4652c0

package ovs

import (
	"fmt"
	"github.com/google/go-cmp/cmp"
	"github.com/socketplane/libovsdb"
	"os/exec"
	"ovs-gnxi/shared/logging"
)

const (
	DefaultDatabase = "Open_vSwitch"
	SystemTable     = "Open_vSwitch"
	ControllerTable = "Controller"
	InterfaceTable  = "Interface"
	StartOVS        = "start_ovs.sh"
	StopOVS         = "stop_ovs.sh"
	RestartOVS      = "restart_ovs.sh"
)

var log = logging.New("ovs-gnxi")

var update chan *libovsdb.TableUpdates

type Client struct {
	Address    string
	Protocol   string
	Port       string
	Connection *libovsdb.OvsdbClient
	Database   string
	Notifier   *Notifier
	Config     *Config
	ErrorChan  chan error
	QuitChan   chan bool
}

func (o *Client) String() string {
	return fmt.Sprintf("OVSClient(Address: \"%v\", Protocol: \"%v\", Port: \"%v\")", o.Address, o.Protocol, o.Port)
}

func NewClient(address, protocol, port string) (*Client, error) {
	o := Client{Address: address, Protocol: protocol, Port: port, Database: DefaultDatabase, ErrorChan: make(chan error), Config: NewConfig()}
	return &o, nil
}

func (o *Client) StartClient(privateKeyPath, publicKeyPath, caPath string) {
	log.Info("Start OVS Client")

	var err error
	o.Connection, err = libovsdb.ConnectUsingProtocolWithTLS(o.Protocol, fmt.Sprintf("%v:%v", o.Address, o.Port), privateKeyPath, publicKeyPath, caPath)
	if err != nil {
		o.ErrorChan <- err
	}

	o.MonitorAll()
	o.Config.DumpRawCache()
}

func (o *Client) StopMonitoring() {
	o.QuitChan <- true
}

func (o *Client) StopClient() {
	log.Info("Stop OVS Client")
	o.Connection.Disconnect()
	o.Config = NewConfig()
}

func (o *Client) SetSystem(system *System) error {
	condition := libovsdb.NewCondition("_uuid", "==", libovsdb.UUID{GoUUID: system.uuid})

	row := make(map[string]interface{})
	row["hostname"] = system.Hostname

	updateOp := libovsdb.Operation{
		Op:    "update",
		Table: SystemTable,
		Where: []interface{}{condition},
		Row:   row,
	}

	log.Debug(updateOp)

	operations := []libovsdb.Operation{updateOp}
	reply, _ := o.Connection.Transact(o.Database, operations...)

	if len(reply) < len(operations) {
		log.Error("number of Replies should be at least equal to number of Operations")
	}
	ok := true
	for i, o := range reply {
		if o.Error != "" && i < len(operations) {
			log.Errorf("transaction failed due to an error :", o.Error, " details:", o.Details, " in ", operations[i])
			ok = false
		} else if o.Error != "" {
			log.Errorf("transaction failed due to an error :", o.Error)
			ok = false
		}
	}
	if ok {
		return nil
	}

	return fmt.Errorf("unable to set system information")
}

func (o *Client) SetOpenFlowController(controller *OpenFlowController) error {
	condition := libovsdb.NewCondition("_uuid", "==", libovsdb.UUID{GoUUID: controller.uuid})

	row := make(map[string]interface{})
	row["target"] = fmt.Sprintf("%v:%v:%v", controller.Target.Protocol, controller.Target.Address, controller.Target.Port)

	updateOp := libovsdb.Operation{
		Op:    "update",
		Table: ControllerTable,
		Where: []interface{}{condition},
		Row:   row,
	}

	log.Debug(updateOp)

	operations := []libovsdb.Operation{updateOp}
	reply, _ := o.Connection.Transact(o.Database, operations...)

	if len(reply) < len(operations) {
		log.Error("number of Replies should be at least equal to number of Operations")
	}
	ok := true
	for i, o := range reply {
		if o.Error != "" && i < len(operations) {
			log.Errorf("transaction failed due to an error :", o.Error, " details:", o.Details, " in ", operations[i])
			ok = false
		} else if o.Error != "" {
			log.Errorf("transaction failed due to an error :", o.Error)
			ok = false
		}
	}
	if ok {
		return nil
	}

	return fmt.Errorf("unable to set system information")
}

func (o *Client) SetInterface(interf *Interface) error {
	condition := libovsdb.NewCondition("_uuid", "==", libovsdb.UUID{GoUUID: interf.uuid})

	row := make(map[string]interface{})
	row["name"] = interf.Name
	row["mtu"] = interf.MTU

	updateOp := libovsdb.Operation{
		Op:    "update",
		Table: InterfaceTable,
		Where: []interface{}{condition},
		Row:   row,
	}

	log.Debug(updateOp)

	operations := []libovsdb.Operation{updateOp}
	reply, _ := o.Connection.Transact(o.Database, operations...)

	if len(reply) < len(operations) {
		log.Error("number of Replies should be at least equal to number of Operations")
	}
	ok := true
	for i, o := range reply {
		if o.Error != "" && i < len(operations) {
			log.Errorf("transaction failed due to an error :", o.Error, " details:", o.Details, " in ", operations[i])
			ok = false
		} else if o.Error != "" {
			log.Errorf("transaction failed due to an error :", o.Error)
			ok = false
		}
	}
	if ok {
		return nil
	}

	return fmt.Errorf("unable to set system information")
}

func (o *Client) SyncChangesToRemote(prev, new *ObjectCache) error {
	if !cmp.Equal(prev.System, new.System) {
		log.Info("target is in inconsistent state with OVS device, syncing System")

		err := o.SetSystem(new.System)
		if err != nil {
			return err
		}
	}

	for _, controller := range new.Controllers {
		if prev.Controllers[controller.Name].uuid == controller.uuid {
			if !cmp.Equal(prev.Controllers[controller.Name], controller) {
				log.Info("target is in inconsistent state with OVS device, syncing Controller")

				err := o.SetOpenFlowController(controller)
				if err != nil {
					return err
				}
			}
		}
	}

	for _, interf := range new.Interfaces {
		if prev.Interfaces[interf.Name].uuid == interf.uuid {
			if !cmp.Equal(prev.Interfaces[interf.Name], interf) {
				log.Info("target is in inconsistent state with OVS device, syncing Interface")

				err := o.SetInterface(interf)
				if err != nil {
					return err
				}
			}
		}
	}

	return nil
}

type Notifier struct {
	config *Config
}

func (n Notifier) Update(context interface{}, tableUpdates libovsdb.TableUpdates) {
	n.config.SyncCache(&tableUpdates)
	update <- &tableUpdates
}
func (n Notifier) Locked([]interface{}) {
}
func (n Notifier) Stolen([]interface{}) {
}
func (n Notifier) Echo([]interface{}) {
}
func (n Notifier) Disconnected(client *libovsdb.OvsdbClient) {
}

func (o *Client) receivedMonitorUpdate() {
	for {
		select {
		case <-o.QuitChan:
			return
		case currUpdate := <-update:
			for tableName, tableUpdate := range currUpdate.Updates {
				log.Debugf("Received Table update for \"%v\" with content: %v", tableName, tableUpdate)
			}
		}
	}
}

func (o *Client) MonitorAll() {
	o.QuitChan = make(chan bool)
	update = make(chan *libovsdb.TableUpdates)

	request := libovsdb.MonitorRequest{
		Select: libovsdb.MonitorSelect{
			Initial: true,
			Insert:  true,
			Delete:  true,
			Modify:  true,
		},
	}

	requests := make(map[string]libovsdb.MonitorRequest)
	requests[SystemTable] = request
	requests[ControllerTable] = request
	requests[InterfaceTable] = request

	o.Notifier = &Notifier{config: o.Config}
	o.Connection.Register(o.Notifier)

	initial, err := o.Connection.Monitor(DefaultDatabase, "", requests)
	if err != nil {
		o.ErrorChan <- err
	}

	o.Config.InitializeCache(initial)

	go o.receivedMonitorUpdate()
}

func (o *Client) StartSystem() error {
	log.Debug("Starting OVS system...")

	cmd := exec.Command("/bin/sh", StartOVS)

	_, err := cmd.Output()
	if err != nil {
		return err
	}

	return nil
}

func (o *Client) StopSystem() error {
	log.Debug("Stopping OVS system...")

	cmd := exec.Command("/bin/sh", StopOVS)

	_, err := cmd.Output()
	if err != nil {
		return err
	}

	return nil
}

func (o *Client) RestartSystem() error {
	log.Debug("Restarting OVS system...")

	cmd := exec.Command("/bin/sh", RestartOVS)

	_, err := cmd.Output()
	if err != nil {
		return err
	}

	return nil
}
