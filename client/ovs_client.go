package main

import (
	"fmt"
	"github.com/socketplane/libovsdb"
)

const (
	DefaultDatabase = "Open_vSwitch"
	ControllerTable = "Controller"
)

type OpenFlowControllerTarget struct {
	Address  string
	Protocol string
	Port     uint16
}

type OpenFlowController struct {
	uuid      string
	Name      string
	Connected bool
	Target    *OpenFlowControllerTarget
}

type Client struct {
	Address    string
	Protocol   string
	Port       string
	Connection *libovsdb.OvsdbClient
	Database   string
}

func NewClient(address, protocol, port, privateKeyPath, publicKeyPath, caPath string) (*Client, error) {
	var err error

	o := Client{Address: address, Protocol: protocol, Port: port, Database: DefaultDatabase}

	o.Connection, err = libovsdb.ConnectUsingProtocolWithTLS(o.Protocol, fmt.Sprintf("%v:%v", o.Address, o.Port), privateKeyPath, publicKeyPath, caPath)
	if err != nil {
		log.Errorf("failed to dial: %v", err)
		return nil, err
	}

	return &o, nil
}

func GetOpenFlowController(target string) error {
	condition := libovsdb.NewCondition("target", "==", target)

	updateOp := libovsdb.Operation{
		Op:      "select",
		Table:   ControllerTable,
		Where:   []interface{}{condition},
		Columns: []string{"target"},
	}

	log.Info(updateOp)

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