package main

import (
	"context"
	"errors"
	"fmt"
	"github.com/digitalocean/go-openvswitch/ovsdb"
	"time"
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

type OVSClient struct {
	Address    string
	Protocol   string
	Port       string
	Connection *ovsdb.Client
	Database   string
	Config     []byte
}

func (o *OVSClient) String() string {
	return fmt.Sprintf("OVSClient(Address: \"%v\", Protocol: \"%v\", Port: \"%v\")", o.Address, o.Protocol, o.Port)
}

func NewOVSClient(address, protocol, port string) (*OVSClient, error) {
	var err error
	o := OVSClient{Address: address, Protocol: protocol, Port: port}
	o.Connection, err = ovsdb.Dial(o.Protocol, fmt.Sprintf("%v:%v", o.Address, o.Port))
	if err != nil {
		log.Fatalf("failed to dial: %v", err)
	}

	o.Database, err = o.GetDefaultDatabase()

	return &o, nil
}

func (o *OVSClient) GetDatabases() ([]string, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()

	dbs, err := o.Connection.ListDatabases(ctx)
	if err != nil {
		log.Fatalf("failed to list databases: %v", err)
	}

	return dbs, err
}

func (o *OVSClient) GetDefaultDatabase() (string, error) {
	dbs, err := o.GetDatabases()

	if err != nil {
		log.Fatalf("failed to get default database: %v", err)
	}

	for _, d := range dbs {
		if d == DefaultDatabase {
			return d, nil
		}
	}

	return "", errors.New("default database not found")
}

func (o *OVSClient) InitializeConfig() {
	controller, err := o.GetOpenFlowControllerIP()
	if err != nil {
		log.Info(controller)
	}
}

func (o *OVSClient) ClientCallback() {

}

func (o *OVSClient) GetOpenFlowControllerIP() (string, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()

	ops := []ovsdb.TransactOp{ovsdb.Select{
		Table:   ControllerTable,
		Columns: []string{"target"},
	}}

	result, err := o.Connection.Transact(ctx, o.Database, ops)
	if err != nil {
		log.Fatalf("failed to complete transaction: %v", err)
	}

	log.Info(result[0])

	/*
		s := strings.Split(c, ":")
		c.Protocol, c.Address, c.Port = s[0], s[1], s[2]
	*/

	return "UNIMPLEMENTED!", nil
}

func (o *OVSClient) SetOpenFlowControllerIP() (string, error) {
	return "UNIMPLEMENTED!", nil
}
