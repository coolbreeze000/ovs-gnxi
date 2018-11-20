package main

import (
	"fmt"
	"github.com/socketplane/libovsdb"
	"strings"
)

const (
	DefaultDatabase = "Open_vSwitch"
	ControllerTable = "Controller"
)

type OpenFlowController struct {
	Address  string
	Protocol string
	Port     string
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
	controller, err := o.GetOpenFlowController()
	if err == nil {
		log.Info(controller)
	}
}

func (o *OVSClient) GetOpenFlowController() (*OpenFlowController, error) {
	selectOp := libovsdb.Operation{
		Op:      "select",
		Table:   ControllerTable,
		Columns: []string{"target"},
	}

	operations := []libovsdb.Operation{selectOp}
	reply, _ := o.Connection.Transact(o.Database, operations...)

	if len(reply) < len(operations) {
		fmt.Println("Number of Replies should be atleast equal to number of Operations")
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

	return nil, nil
}

/*
func (o *OVSClient) SetOpenFlowControllerIP() (OpenFlowController, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()

	ops := []ovsdb.TransactOp{ovsdb.Update{
		Table:   ControllerTable,
		Columns: []string{"target"},
	}}

	result, err := o.Connection.Transact(ctx, o.Database, ops)
	if err != nil {
		log.Fatalf("failed to complete transaction: %v", err)
	}

	s := strings.Split(result[0]["target"].(string), ":")

	return OpenFlowController{Protocol: s[0], Address: s[1], Port: s[2]}, nil
}
*/
