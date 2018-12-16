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

// Package gnmi implements a gnmi server.
package gnxi

import (
	"ovs-gnxi/shared/logging"
	"ovs-gnxi/target/gnxi/gnmi"
	"ovs-gnxi/target/gnxi/gnoi"
)

var log = logging.New("ovs-gnxi")

// Server struct maintains the data structure for device config and implements the interface of gnmi server. It supports Capabilities, Get, and Set APIs.
type Server struct {
	ServiceGNMI *gnmi.Service
	ServiceGNOI *gnoi.Service
}

// NewServer creates an instance of Server.
func NewServer(serviceGNMI *gnmi.Service, serviceGNOI *gnoi.Service) *Server {
	log.Info("Initializing gNXI Server...")
	return &Server{ServiceGNMI: serviceGNMI, ServiceGNOI: serviceGNOI}
}
