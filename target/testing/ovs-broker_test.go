package testing

import (
	"context"
	"github.com/digitalocean/go-openvswitch/ovsdb"
	"log"
	"testing"
	"time"
)

func TestOVSDBConnection(t *testing.T) {
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
		log.Println(d)
	}
}
