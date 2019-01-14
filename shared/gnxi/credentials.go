/* Copyright 2017 Google Inc.

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
package gnxi

import (
	"errors"
	"fmt"
	"golang.org/x/net/context"
	"google.golang.org/grpc/metadata"
)

type Authenticator struct {
	Users map[string]User
}

func (a *Authenticator) AddUser(user *User) {
	a.Users[user.username] = *user
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

	if _, ok := a.Users[username]; ok {
		if password == a.Users[username].password && username == a.Users[username].username {
			return true, nil
		}
	}

	return false, errors.New(fmt.Sprintf("not authorized with \"%s:%s\"", username, password))
}

type User struct {
	username string
	password string
}

func NewUser(username, password string) *User {
	return &User{
		username: username,
		password: password,
	}
}
