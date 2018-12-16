/* Copyright 2018 Google Inc.

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
