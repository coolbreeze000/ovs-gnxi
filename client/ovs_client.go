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

func (o *Client) Get(param, table string) (map[string]interface{}, error) {
	// TODO(dherkel@google.com): This needs to be more generic if anything other than Controller needs to be tested.
	switch table {
	case ControllerTable:
		return o.GetOpenFlowControllerTarget(param)
	}

	return nil, fmt.Errorf("unimplemented OVS Get Type")
}

func (o *Client) GetOpenFlowControllerTarget(target string) (map[string]interface{}, error) {
	condition := libovsdb.NewCondition("target", "==", target)

	selectOp := libovsdb.Operation{
		Op:      "select",
		Table:   ControllerTable,
		Where:   []interface{}{condition},
		Columns: []string{"target"},
	}

	operations := []libovsdb.Operation{selectOp}
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
		if len(reply[0].Rows) > 0 {
			return reply[0].Rows[0], nil
		}

		return map[string]interface{}{}, nil
	}

	return nil, fmt.Errorf("unable to set system information")
}
