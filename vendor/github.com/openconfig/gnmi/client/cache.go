/*
Copyright 2017 Google Inc.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package client

import (
	"fmt"

	"context"
	log "github.com/golang/glog"
	"github.com/openconfig/gnmi/ctree"
)

// Client adds a caching layer on top of a simple query client.
//
// It works similarly to BaseClient and adds the Leaves method to return
// current tree state.
type CacheClient struct {
	*BaseClient
	*ctree.Tree
	synced        chan struct{}
	clientHandler NotificationHandler
}

var _ Client = &CacheClient{}

// New returns an initialized caching client.
func New() *CacheClient {
	c := &CacheClient{
		BaseClient: &BaseClient{},
		Tree:       &ctree.Tree{},
		synced:     make(chan struct{}),
	}
	return c
}

// Subscribe implements the Client interface.
func (c *CacheClient) Subscribe(ctx context.Context, q Query, clientType ...string) error {
	q.ProtoHandler = nil
	if q.NotificationHandler != nil {
		c.clientHandler = q.NotificationHandler
	}
	q.NotificationHandler = c.defaultHandler
	return c.BaseClient.Subscribe(ctx, q, clientType...)
}

// defaultHandler is passed into the client specific implementations. It will
// be called for each leaf notification generated by the client.
func (c *CacheClient) defaultHandler(n Notification) error {
	switch v := n.(type) {
	default:
		return fmt.Errorf("invalid type %#v", v)
	case Connected: // Ignore.
	case Error:
		return fmt.Errorf("received error: %v", v)
	case Update:
		c.Add(v.Path, TreeVal{TS: v.TS, Val: v.Val})
	case Delete:
		c.Delete(v.Path)
	case Sync:
		select {
		default:
			close(c.synced)
		case <-c.synced:
		}
	}
	if c.clientHandler != nil {
		return c.clientHandler(n)
	}
	return nil
}

// Poll implements the Client interface.
// Poll also closes the channel returned by Synced and resets it.
func (c *CacheClient) Poll() error {
	select {
	default:
		close(c.synced)
	case <-c.synced:
	}
	return c.BaseClient.Poll()
}

// Synced will close when a sync is recieved from the query.
func (c *CacheClient) Synced() <-chan struct{} {
	return c.synced
}

// Leaves returns the current state of the received tree. It's safe to call at
// any point after New.
func (c *CacheClient) Leaves() Leaves {
	// Convert node items into Leaf (expand TreeVal leaves).
	var pvs Leaves
	c.WalkSorted(func(path []string, _ *ctree.Leaf, value interface{}) {
		tv, ok := value.(TreeVal)
		if !ok {
			log.Errorf("Invalid value in tree: %s=%#v", path, value)
			return
		}
		pvs = append(pvs, Leaf{Path: path, Val: tv.Val, TS: tv.TS})
	})
	return pvs
}