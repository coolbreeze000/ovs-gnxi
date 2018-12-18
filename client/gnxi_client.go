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

package main

import (
	"flag"
	"github.com/google/gnxi/utils"
	"github.com/google/go-cmp/cmp"
	"os"
	"ovs-gnxi/client/gnmi"
	"ovs-gnxi/shared/logging"
	"time"
)

type arrayFlags []string

func (i *arrayFlags) String() string {
	return "my string representation"
}

func (i *arrayFlags) Set(value string) error {
	*i = append(*i, value)
	return nil
}

var (
	log             = logging.New("ovs-gnxi-client")
	targetAddr      = flag.String("target_addr", "target:10161", "The target address in the format of host:port")
	targetName      = flag.String("target_name", "target.gnxi.lan", "The target name use to verify the hostname returned by TLS handshake")
	encodingName    = flag.String("encoding", "JSON_IETF", "The value encoding format to be used")
	timeOut         = flag.Duration("time_out", 10*time.Second, "Timeout for the Get request, 10 seconds by default")
	method          = flag.String("method", "", "A valid gNMI specification method to execute against the target")
	subscribeMode   = flag.String("subscribe_mode", "Once", "The gNMI data subscription mode")
	getXPaths       arrayFlags
	deleteXPaths    arrayFlags
	replaceXPaths   arrayFlags
	updateXPaths    arrayFlags
	subscribeXPaths arrayFlags
)

func setGNMIClientFlags() {
	flag.Var(&getXPaths, "get_xpath", "The gNMI Get XPaths to query for")
	flag.Var(&deleteXPaths, "set_delete_xpath", "The gNMI Set Delete XPaths to query for")
	flag.Var(&replaceXPaths, "set_replace_xpath", "The gNMI Set Replace XPaths to query for")
	flag.Var(&updateXPaths, "set_update_xpath", "The gNMI Set Update XPaths to query for")
	flag.Var(&subscribeXPaths, "subscribe_xpath", "The gNMI Subscribe XPaths to query for")

	err := flag.Set("ca", "certs/ca.crt")
	if err != nil {
		log.Fatalf("Unable to set ca flag: %v", err)
	}
	err = flag.Set("cert", "certs/client.crt")
	if err != nil {
		log.Fatalf("Unable to set cert flag: %v", err)
	}
	err = flag.Set("key", "certs/client.key")
	if err != nil {
		log.Fatalf("Unable to set key flag: %v", err)
	}
	err = flag.Set("username", "admin")
	if err != nil {
		log.Fatalf("Unable to set username flag: %v", err)
	}
	err = flag.Set("password", "testpassword")
	if err != nil {
		log.Fatalf("Unable to set password flag: %v", err)
	}

	flag.Parse()
}

func main() {
	defer os.Exit(0)
	defer log.Info("Exiting Open vSwitch gNXI client tester\n")

	setGNMIClientFlags()

	log.Info("Started Open vSwitch gNXI client tester\n")

	gnmiClient := gnmi.NewGNMIClient(*targetAddr, *targetName, *encodingName, *timeOut)

	switch *method {
	case "Capabilities":
		resp, err := gnmiClient.Capabilities()
		if err != nil {
			log.Fatal(err)
		}

		log.Debug("== Response:")
		utils.PrintProto(resp)
	case "Get":
		resp, err := gnmiClient.Get(getXPaths)
		if err != nil {
			log.Fatal(err)
		}

		log.Debug("== Response:")
		utils.PrintProto(resp)
	case "Set":
		resp, err := gnmiClient.Set(deleteXPaths, replaceXPaths, updateXPaths)
		if err != nil {
			log.Fatal(err)
		}

		log.Debug("== Response:")
		utils.PrintProto(resp)
	case "Subscribe":
		log.Info(subscribeXPaths)
		resp, err := gnmiClient.Subscribe(subscribeXPaths, *subscribeMode)
		if err != nil {
			log.Fatal(err)
		}

		log.Debug("== Response:")
		utils.PrintProto(resp)
	default:
		RunGNMICapabilitiesTests(gnmiClient)
		RunGNMIGetTests(gnmiClient)
	}

	log.Info("Finished Open vSwitch gNXI client tester\n")
}

func RunGNMICapabilitiesTests(c *gnmi.Client) {
	for _, td := range gnmi.CapabilitiesTests {
		log.Infof("Testing GNMI Capabilities(%v)...", td.Desc)

		resp, err := c.Capabilities()
		if err != nil {
			log.Fatal(err)
		}

		if resp.GNMIVersion != td.ExpVersion {
			log.Errorf("Capabilities(%v): expected %v, actual %v", td.Desc, td.ExpVersion, resp.GNMIVersion)
		}

		for _, respModel := range resp.SupportedModels {
			match := false

			for _, expModel := range td.ExpModels {
				if cmp.Equal(respModel, expModel) {
					match = true
					break
				}
			}

			if match != true {
				log.Errorf("Capabilities(%v): model %v not found in %v", td.Desc, respModel, td.ExpModels)
			}
		}

		for _, respEncoding := range resp.SupportedEncodings {
			match := false

			for _, expEncoding := range td.ExpEncodings {
				if cmp.Equal(respEncoding, expEncoding) {
					match = true
					break
				}
			}

			if match != true {
				log.Errorf("Capabilities(%v): encoding %v not found in %v", td.Desc, respEncoding, td.ExpEncodings)
			}
		}

		for _, respExtension := range resp.Extension {
			match := false

			for _, expExtension := range td.ExpExtensions {
				if cmp.Equal(respExtension, expExtension) {
					match = true
					break
				}
			}

			if match != true {
				log.Errorf("Capabilities(%v): extension %v not found in %v", td.Desc, respExtension, td.ExpExtensions)
			}
		}
	}
}

func RunGNMIGetTests(c *gnmi.Client) {
	for _, td := range gnmi.GetTests {
		log.Infof("Testing GNMI Get(%v)...", td.XPaths)

		resp, err := c.Get(td.XPaths)
		if err != nil {
			log.Fatal(err)
			continue
		}

		log.Infof("Successfully verified GNMI Get(%v)", td.XPaths)

		if td.Extractor(resp.Notification) != td.ExpResp {
			log.Errorf("Get(%v): expected %v, actual %v", td.XPaths, td.ExpResp, td.Extractor(resp.Notification))
		}
	}
}
