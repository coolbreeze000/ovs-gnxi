package main

import (
	"errors"
	"fmt"
	"github.com/socketplane/libovsdb"
	"strings"
)

const (
	DefaultDatabase = "Open_vSwitch"
	SystemTable     = "Open_vSwitch"
	ControllerTable = "Controller"
	InterfaceTable  = "Interface"
)

type System struct {
	uuid     string
	Platform string
	Version  string
	Hostname string
}

func (s *System) String() string {
	return fmt.Sprintf("System(uuid: \"%v\", Platform: \"%v\", Version: \"%v\", Hostname: \"%v\")", s.uuid, s.Platform, s.Version, s.Hostname)
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
	return fmt.Sprintf("OpenFlowController(Address: \"%v\", Protocol: \"%v\", Port: \"%v\")", c.Address, c.Protocol, c.Port)
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
	log.Info("TEST")

	system, err := o.GetSystemInformation()
	log.Info(err)
	if err == nil {
		log.Info(system)
	}

	log.Info("TEST2")

	/*
		controller, err := o.GetOpenFlowControllers()
		if err == nil {
			log.Info(controller)
		}

		err = o.SetOpenFlowController(controller)

		controller, err = o.GetOpenFlowControllers()
		if err == nil {
			log.Info(controller)
		}*/
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
		uuid := strings.Replace(strings.Split(fmt.Sprint(reply[0].Rows[0]["_uuid"]), " ")[1], "]", "", -1)
		log.Info(uuid)
		version := reply[0].Rows[0]["ovs_version"].(string)
		log.Info(version)
		hostname := reply[0].Rows[0]["external_ids"].(map[string]string)
		log.Info(hostname)
		return &System{uuid: uuid, Version: version, Hostname: fmt.Sprint(hostname)}, nil
	}

	return nil, errors.New("unable to get openflow controller information")
}

func (o *OVSClient) GetOpenFlowControllers() (*OpenFlowController, error) {
	selectOp := libovsdb.Operation{
		Op:      "select",
		Table:   ControllerTable,
		Columns: []string{"target"},
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
		log.Info(reply[0])
		s := strings.Split(reply[0].Rows[0]["target"].(string), ":")
		return &OpenFlowController{Protocol: s[0], Address: s[1], Port: s[2]}, nil
	}

	return nil, errors.New("unable to get openflow controller information")
}

func (o *OVSClient) GetFirstOpenFlowController() (*OpenFlowController, error) {
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
		log.Info(reply[0])
		s := strings.Split(reply[0].Rows[0]["target"].(string), ":")
		return &OpenFlowController{Protocol: s[0], Address: s[1], Port: s[2]}, nil
	}

	return nil, errors.New("unable to get openflow controller information")
}

func (o *OVSClient) GetOpenFlowController(uuid string) (*OpenFlowController, error) {
	condition := libovsdb.NewCondition("_uuid", "==", uuid)

	selectOp := libovsdb.Operation{
		Op:      "select",
		Table:   ControllerTable,
		Where:   condition,
		Columns: []string{"target"},
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
		log.Info(reply[0])
		s := strings.Split(reply[0].Rows[0]["target"].(string), ":")
		return &OpenFlowController{Protocol: s[0], Address: s[1], Port: s[2]}, nil
	}

	return nil, errors.New("unable to get openflow controller information")
}

func (o *OVSClient) SetOpenFlowController(controller *OpenFlowController) error {
	row := make(map[string]interface{})
	row["target"] = fmt.Sprintf("%v:%v:%v", controller.Protocol, controller.Address, controller.Port)

	updateOp := libovsdb.Operation{
		Op:    "update",
		Table: ControllerTable,
		Row:   row,
	}

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
		log.Info(reply[0])
		return nil
	}

	return errors.New("unable to set openflow controller address")
}
