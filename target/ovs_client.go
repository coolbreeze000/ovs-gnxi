package main

import (
	"errors"
	"fmt"
	"github.com/socketplane/libovsdb"
	"reflect"
	"strings"
)

const (
	DefaultDatabase = "Open_vSwitch"
	SystemTable     = "Open_vSwitch"
	ControllerTable = "Controller"
	InterfaceTable  = "Interface"
)

var quit chan bool
var update chan *libovsdb.TableUpdates
var cache map[string]map[string]libovsdb.Row

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
	Port     string
}

type Interface struct {
	uuid       string
	Name       string
	MTU        string
	Status     string
	Statistics string
}

func (c *OpenFlowController) String() string {
	return fmt.Sprintf("OpenFlowController(uuid: \"%v\", Name: \"%v\", Address: \"%v\", Protocol: \"%v\", Port: \"%v\")", c.uuid, c.Name, c.Address, c.Protocol, c.Port)
}

type OVSClient struct {
	Address    string
	Protocol   string
	Port       string
	Connection *libovsdb.OvsdbClient
	Database   string
	Config     []byte
}

func (o *OVSClient) String() string {
	return fmt.Sprintf("OVSClient(Address: \"%v\", Protocol: \"%v\", Port: \"%v\")", o.Address, o.Protocol, o.Port)
}

func NewOVSClient(address, protocol, port, privateKeyPath, publicKeyPath, caPath string) (*OVSClient, error) {
	var err error

	o := OVSClient{Address: address, Protocol: protocol, Port: port, Database: DefaultDatabase}

	o.Connection, err = libovsdb.ConnectUsingProtocolWithTLS(o.Protocol, fmt.Sprintf("%v:%v", o.Address, o.Port), privateKeyPath, publicKeyPath, caPath)
	if err != nil {
		log.Fatalf("failed to dial: %v", err)
	}

	return &o, nil
}

func (o *OVSClient) InitializeConfig() {
	system, err := o.GetSystemInformation()
	if err == nil {
		log.Info(system)
	}

	controller, err := o.GetOpenFlowControllers()
	if err == nil {
		log.Info(controller)
	}

	err = o.SetOpenFlowController(controller)

	controller, err = o.GetOpenFlowControllers()
	if err == nil {
		log.Info(controller)
	}

	err = o.MonitorAll()
	if err != nil {
		log.Error(err)
	}
}

func (o *OVSClient) GetSystemInformation() (*System, error) {
	selectOp := libovsdb.Operation{
		Op:    "select",
		Table: SystemTable,
	}

	operations := []libovsdb.Operation{selectOp}
	reply, _ := o.Connection.Transact(o.Database, operations...)

	if len(reply) < len(operations) {
		fmt.Println("Number of Replies should be at least equal to number of Operations")
	}
	ok := true
	for i, o := range reply {
		if o.Error != "" && i < len(operations) {
			fmt.Println("Transaction Failed due to an error :", o.Error, " details:", o.Details, " in ", operations[i])
			ok = false
		} else if o.Error != "" {
			fmt.Println("Transaction Failed due to an error :", o.Error)
			ok = false
		}
	}
	if ok {
		uuid := reply[0].Rows[0]["_uuid"].([]interface{})[1].(string)
		version := reply[0].Rows[0]["ovs_version"].(string)
		hostname := reply[0].Rows[0]["external_ids"].([]interface{})[1].([]interface{})[0].([]interface{})[1].(string)

		return &System{uuid: uuid, Version: version, Hostname: hostname}, nil
	}

	return nil, errors.New("unable to get system information")
}

func (o *OVSClient) GetOpenFlowControllers() (*OpenFlowController, error) {
	selectOp := libovsdb.Operation{
		Op:    "select",
		Table: ControllerTable,
	}

	operations := []libovsdb.Operation{selectOp}
	reply, _ := o.Connection.Transact(o.Database, operations...)

	if len(reply) < len(operations) {
		fmt.Println("Number of Replies should be at least equal to number of Operations")
	}
	ok := true
	for i, o := range reply {
		if o.Error != "" && i < len(operations) {
			fmt.Println("Transaction Failed due to an error :", o.Error, " details:", o.Details, " in ", operations[i])
			ok = false
		} else if o.Error != "" {
			fmt.Println("Transaction Failed due to an error :", o.Error)
			ok = false
		}
	}
	if ok {
		uuid := reply[0].Rows[0]["_uuid"].([]interface{})[1].(string)
		s := strings.Split(reply[0].Rows[0]["target"].(string), ":")
		return &OpenFlowController{uuid: uuid, Name: uuid, Protocol: s[0], Address: s[1], Port: s[2]}, nil
	}

	return nil, errors.New("unable to get openflow controller information")
}

func (o *OVSClient) GetOpenFlowController(uuid string) (*OpenFlowController, error) {
	condition := libovsdb.NewCondition("_uuid", "==", uuid)

	selectOp := libovsdb.Operation{
		Op:    "select",
		Table: ControllerTable,
		Where: condition,
	}

	operations := []libovsdb.Operation{selectOp}
	reply, _ := o.Connection.Transact(o.Database, operations...)

	if len(reply) < len(operations) {
		fmt.Println("Number of Replies should be at least equal to number of Operations")
	}
	ok := true
	for i, o := range reply {
		if o.Error != "" && i < len(operations) {
			fmt.Println("Transaction Failed due to an error :", o.Error, " details:", o.Details, " in ", operations[i])
			ok = false
		} else if o.Error != "" {
			fmt.Println("Transaction Failed due to an error :", o.Error)
			ok = false
		}
	}
	if ok {
		s := strings.Split(reply[0].Rows[0]["target"].(string), ":")
		return &OpenFlowController{Protocol: s[0], Address: s[1], Port: s[2]}, nil
	}

	return nil, errors.New("unable to get openflow controller information")
}

func (o *OVSClient) SetOpenFlowController(controller *OpenFlowController) error {
	condition := libovsdb.NewCondition("_uuid", "==", libovsdb.UUID{GoUUID: controller.uuid})

	row := make(map[string]interface{})
	row["target"] = fmt.Sprintf("%v:%v:%v", controller.Protocol, controller.Address, controller.Port)

	updateOp := libovsdb.Operation{
		Op:    "update",
		Table: ControllerTable,
		Where: []interface{}{condition},
		Row:   row,
	}

	log.Info(updateOp)

	operations := []libovsdb.Operation{updateOp}
	reply, _ := o.Connection.Transact(o.Database, operations...)

	if len(reply) < len(operations) {
		fmt.Println("Number of Replies should be at least equal to number of Operations")
	}
	ok := true
	for i, o := range reply {
		if o.Error != "" && i < len(operations) {
			fmt.Println("Transaction Failed due to an error :", o.Error, " details:", o.Details, " in ", operations[i])
			ok = false
		} else if o.Error != "" {
			fmt.Println("Transaction Failed due to an error :", o.Error)
			ok = false
		}
	}
	if ok {
		return nil
	}

	return errors.New("unable to set openflow controller address")
}

type Notifier struct {
}

func (n Notifier) Update(context interface{}, tableUpdates libovsdb.TableUpdates) {
	populateCache(tableUpdates)
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

func receivedMonitorUpdate() {
	for {
		select {
		case currUpdate := <-update:
			for table, tableUpdate := range currUpdate.Updates {
				log.Info(table)
				log.Info(tableUpdate)
			}
		}
	}

}

func populateCache(updates libovsdb.TableUpdates) {
	for table, tableUpdate := range updates.Updates {
		if _, ok := cache[table]; !ok {
			cache[table] = make(map[string]libovsdb.Row)
		}
		for uuid, row := range tableUpdate.Rows {
			empty := libovsdb.Row{}
			if !reflect.DeepEqual(row.New, empty) {
				cache[table][uuid] = row.New
			} else {
				delete(cache[table], uuid)
			}
		}
	}
}

func (o *OVSClient) MonitorAll() error {
	quit = make(chan bool)
	update = make(chan *libovsdb.TableUpdates)
	cache = make(map[string]map[string]libovsdb.Row)

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

	var notifier Notifier
	o.Connection.Register(notifier)

	initial, err := o.Connection.Monitor(DefaultDatabase, "", requests)
	if err != nil {
		log.Error(err)
		return err
	}

	log.Info(initial)
	go receivedMonitorUpdate()
	<-quit

	return nil
}
