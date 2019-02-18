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

// Package credentials loads certificates and validates user credentials.
package shared

import (
	"errors"
	"fmt"
	"golang.org/x/net/context"
	"google.golang.org/grpc/metadata"
)

type User struct {
	username string
	password string
}

type Authenticator struct {
	users map[string]User
}

func NewAuthenticator(adminUsername, adminPassword string) *Authenticator {
	a := &Authenticator{users: make(map[string]User)}
	a.AddUser(adminUsername, adminPassword)
	return a
}

func (a *Authenticator) AddUser(username, password string) {
	a.users[username] = User{
		username: username,
		password: password,
	}
}

// AuthorizeUser checks for valid credentials in the context Metadata.
func (a *Authenticator) AuthorizeUser(ctx context.Context) (bool, error) {
	headers, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		return false, errors.New("no Metadata found")
	}

	u, ok := headers["username"]
	if !ok || len(u) == 0 {
		return false, errors.New("no username in Metadata")
	}

	username := u[0]

	p, ok := headers["password"]
	if !ok || len(p) == 0 {
		return false, errors.New(fmt.Sprintf("found username \"%s\" but no password in Metadata", username))
	}

	password := p[0]

	if _, ok := a.users[username]; ok {
		if password == a.users[username].password && username == a.users[username].username {
			return true, nil
		}
	}

	return false, errors.New(fmt.Sprintf("not authorized with \"%s:%s\"", username, password))
}
