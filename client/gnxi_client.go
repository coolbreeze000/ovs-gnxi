package main

import (
	"context"
	"github.com/digitalocean/go-openvswitch/ovsdb"
	"github.com/op/go-logging"
	"os"
	"time"
)

var (
	logModule = "gnxi-client"
	log       = logging.MustGetLogger(logModule)
)

func main() {
	defer os.Exit(0)
	defer log.Info("Exiting Open vSwitch gNXI client tester\n")

	log.Info("Starting Open vSwitch gNXI client tester\n")

	// Dial an OVSDB connection and create a *ovsdb.Client.
	c, err := ovsdb.Dial("tcp", "ovs.gnxi.lan:6640")
	if err != nil {
		log.Fatalf("failed to dial: %v", err)
	}
	// Be sure to close the connection!
	defer c.Close()

	// Ask ovsdb-server for all of its databases, but only allow the RPC
	// a limited amount of time to complete before timing out.
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()

	dbs, err := c.ListDatabases(ctx)
	if err != nil {
		log.Fatalf("failed to list databases: %v", err)
	}

	for _, d := range dbs {
		log.Info(d)
	}

	db_name := dbs[0]

	if db_name == "Open_vSwitch" {
		ops := []ovsdb.TransactOp{ovsdb.Select{
			Table:   "Controller",
			Columns: []string{"target"},
		}}

		controller, err := c.Transact(ctx, db_name, ops)
		if err != nil {
			log.Fatalf("failed to complete transaction: %v", err)
		}

		log.Info(controller)
	}
}
