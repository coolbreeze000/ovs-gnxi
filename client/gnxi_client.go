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
	pb "github.com/openconfig/gnmi/proto/gnmi"
	"os"
	"ovs-gnxi/client/gnmi"
	"ovs-gnxi/shared/logging"
	"time"

	"golang.org/x/net/context"
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
	subscribeMode   = flag.String("subscribe_mode", "ONCE", "The gNMI data subscription mode")
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

	gnmiClient := gnmi.NewGNMIClient(*targetAddr, *targetName, *encodingName)

	switch *method {
	case "Capabilities":
		ctx, cancel := context.WithTimeout(context.Background(), *timeOut)
		defer cancel()

		resp, err := gnmiClient.Capabilities(ctx)
		if err != nil {
			log.Fatal(err)
		}

		log.Debug("== Response:")
		utils.PrintProto(resp)
	case "Get":
		ctx, cancel := context.WithTimeout(context.Background(), *timeOut)
		defer cancel()

		resp, err := gnmiClient.Get(ctx, getXPaths)
		if err != nil {
			log.Fatal(err)
		}

		log.Debug("== Response:")
		utils.PrintProto(resp)
	case "Set":
		ctx, cancel := context.WithTimeout(context.Background(), *timeOut)
		defer cancel()

		resp, err := gnmiClient.Set(ctx, deleteXPaths, replaceXPaths, updateXPaths)
		if err != nil {
			log.Fatal(err)
		}

		log.Debug("== Response:")
		utils.PrintProto(resp)
	case "Subscribe":
		log.Info(subscribeXPaths)

		switch *subscribeMode {
		case "STREAM":
			ctx, cancel := context.WithCancel(context.Background())
			defer cancel()

			respChan := make(chan *pb.SubscribeResponse)
			errChan := make(chan error)

			go gnmiClient.SubscribeStream(ctx, subscribeXPaths, respChan, errChan)

			for {
				select {
				case resp := <-respChan:
					log.Debug("== Response:")
					utils.PrintProto(resp)
				case err := <-errChan:
					log.Fatal(err)
				}
			}
		case "POLL":
			// TODO(dherkel@google.com): pb.SubscriptionList_POLL
		default:
			ctx, cancel := context.WithTimeout(context.Background(), *timeOut)
			defer cancel()

			resp, err := gnmiClient.SubscribeOnce(ctx, subscribeXPaths)
			if err != nil {
				log.Fatal(err)
			}

			log.Debug("== Response:")
			utils.PrintProto(resp)
		}
	default:
		RunGNMICapabilitiesTests(gnmiClient)
		RunGNMIGetTests(gnmiClient)
		// TODO(dherkel@google.com): Activate gNMI Set Test Suite - RunGNMISetTests(gnmi.Client)
		RunGNMISubscribeOnceTests(gnmiClient)
		RunGNMISubscribeStreamTests(gnmiClient)
	}

	log.Info("Finished Open vSwitch gNXI client tester\n")
}

func RunGNMICapabilitiesTests(c *gnmi.Client) {
	ctx, cancel := context.WithTimeout(context.Background(), *timeOut)
	defer cancel()

	for _, td := range gnmi.CapabilitiesTests {
		log.Infof("Testing GNMI Capabilities(%v)...", td.Desc)

		resp, err := c.Capabilities(ctx)
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
	ctx, cancel := context.WithTimeout(context.Background(), *timeOut)
	defer cancel()

	for _, td := range gnmi.GetTests {
		log.Infof("Testing GNMI Get(%v)...", td.XPaths)

		resp, err := c.Get(ctx, td.XPaths)
		if err != nil {
			log.Fatal(err)
			continue
		}

		if td.ExtractorString != nil {
			actResp := td.ExtractorString(resp.Notification)

			if actResp != td.ExpResp {
				log.Errorf("Get(%v): expected %v, actual %v", td.XPaths, td.ExpResp, actResp)
			} else {
				log.Infof("Successfully verified GNMI Get(%v) with response value %v", td.XPaths, actResp)
			}
		} else if td.ExtractorUInt != nil {
			actResp := td.ExtractorUInt(resp.Notification)

			if td.ExpResp != nil {
				if actResp != td.ExpResp.(uint64) {
					log.Errorf("Get(%v): expected %v, actual %v", td.XPaths, td.ExpResp, actResp)
				} else {
					log.Infof("Successfully verified GNMI Get(%v) with response value %v", td.XPaths, actResp)
				}
			} else if td.MinResp != nil {
				if actResp <= td.MinResp.(uint64) {
					log.Errorf("Get(%v): expected higher than %v, actual %v", td.XPaths, td.MinResp, actResp)
				} else {
					log.Infof("Successfully verified GNMI Subscribe(%v) with response value %v", td.XPaths, actResp)
				}
			}
		}

	}
}

func RunGNMISetTests(c *gnmi.Client) {
	ctx, cancel := context.WithTimeout(context.Background(), *timeOut)
	defer cancel()

	for _, td := range gnmi.SetTests {
		log.Infof("Testing GNMI Set(%v, %v, %v)...", td.DeleteXPaths, td.ReplaceXPaths, td.UpdateXPaths)

		respSet, err := c.Set(ctx, td.DeleteXPaths, td.ReplaceXPaths, td.UpdateXPaths)
		if err != nil {
			log.Fatal(err)
			continue
		}

		verifiedDeletePaths, verifiedReplacePaths, verifiedUpdatePaths := 0, 0, 0

		for _, delPath := range td.DeleteXPaths {
			for _, resp := range respSet.Response {
				if delPath == resp.Path.Target && resp.Op == pb.UpdateResult_DELETE {
					verifiedDeletePaths++
				}
			}
		}

		if verifiedDeletePaths != len(td.DeleteXPaths) {
			log.Errorf("Set(%v): expected %v deletes, actual %v deletes", len(td.DeleteXPaths), verifiedDeletePaths)
		} else {
			log.Infof("Successfully verified GNMI Set(%v) with %v deletes", verifiedDeletePaths)
		}

		if verifiedReplacePaths != len(td.ReplaceXPaths) {
			log.Errorf("Set(%v): expected %v replaces, actual %v replaces", len(td.ReplaceXPaths), verifiedReplacePaths)
		} else {
			log.Infof("Successfully verified GNMI Set(%v) with %v replaces", verifiedReplacePaths)
		}

		if verifiedUpdatePaths != len(td.UpdateXPaths) {
			log.Errorf("Set(%v): expected %v updates, actual %v updates", len(td.UpdateXPaths), verifiedUpdatePaths)
		} else {
			log.Infof("Successfully verified GNMI Set(%v) with %v updates", verifiedUpdatePaths)
		}

		respGet, err := c.Get(ctx, td.UpdateXPaths)
		if err != nil {
			log.Fatal(err)
			continue
		}

		if td.ExtractorString != nil {
			actResp := td.ExtractorString(respGet.Notification)

			if actResp != td.ExpResp {
				log.Errorf("Set(%v) Update: expected %v, actual %v", td.UpdateXPaths, td.ExpResp, actResp)
			} else {
				log.Infof("Successfully verified GNMI Set(%v) Update with response value %v", td.UpdateXPaths, actResp)
			}
		} else if td.ExtractorUInt != nil {
			actResp := td.ExtractorUInt(respGet.Notification)

			if td.ExpResp != nil {
				if actResp != td.ExpResp.(uint64) {
					log.Errorf("Set(%v) Update: expected %v, actual %v", td.UpdateXPaths, td.ExpResp, actResp)
				} else {
					log.Infof("Successfully verified GNMI Set(%v) Update with response value %v", td.UpdateXPaths, actResp)
				}
			}
		}

	}
}

func RunGNMISubscribeOnceTests(c *gnmi.Client) {
	ctx, cancel := context.WithTimeout(context.Background(), *timeOut)
	defer cancel()

	for _, td := range gnmi.SubscribeOnceTests {
		log.Infof("Testing GNMI Subscribe ONCE(%v)...", td.XPaths)

		resp, err := c.SubscribeOnce(ctx, td.XPaths)
		if err != nil {
			log.Fatal(err)
			continue
		}

		update, ok := resp.GetResponse().(*pb.SubscribeResponse_Update)
		if !ok {
			log.Errorf("Invalid subscribe ONCE(%v) response update: %v", td.XPaths, update)
			continue
		}

		actResp := td.ExtractorUInt([]*pb.Notification{update.Update})

		if actResp <= td.MinResp.(uint64) {
			log.Errorf("Subscribe ONCE(%v): expected higher than %v, actual %v", td.XPaths, td.MinResp, td.ExtractorUInt([]*pb.Notification{update.Update}))
		} else {
			log.Infof("Successfully verified GNMI Subscribe ONCE(%v) with response value %v", td.XPaths, actResp)
		}
	}
}

func RunGNMISubscribeStreamTests(c *gnmi.Client) {
	for _, td := range gnmi.SubscribeStreamTests {
		ctx, cancel := context.WithCancel(context.Background())

		log.Infof("Testing GNMI Subscribe STREAM(%v)...", td.XPaths)

		respChan := make(chan *pb.SubscribeResponse)
		errChan := make(chan error)

		go c.SubscribeStream(ctx, td.XPaths, respChan, errChan)

		currentStream := 1

	l:
		for {
			select {
			case resp := <-respChan:
				update, ok := resp.GetResponse().(*pb.SubscribeResponse_Update)
				if !ok {
					log.Errorf("Invalid subscribe STREAM(%v) %v/%v response update: %v", td.XPaths, currentStream, td.MaxStreamResp, update)
					continue
				}

				actResp := td.ExtractorUInt([]*pb.Notification{update.Update})

				if actResp <= td.MinResp.(uint64) {
					log.Errorf("Subscribe STREAM(%v) %v/%v: expected higher than %v, actual %v", td.XPaths, currentStream, td.MaxStreamResp, td.MinResp, td.ExtractorUInt([]*pb.Notification{update.Update}))
				} else {
					log.Infof("Successfully verified GNMI Subscribe STREAM(%v) %v/%v with response value %v", td.XPaths, currentStream, td.MaxStreamResp, actResp)
				}

				if currentStream >= td.MaxStreamResp {
					cancel()
					break l
				}

				currentStream++
			case err := <-errChan:
				log.Fatal(err)
				continue
			}
		}
	}
}
