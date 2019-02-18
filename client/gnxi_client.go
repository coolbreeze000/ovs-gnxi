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
	"bytes"
	"crypto/x509"
	"encoding/pem"
	"flag"
	"github.com/google/gnxi/utils"
	"github.com/google/gnxi/utils/entity"
	"github.com/google/go-cmp/cmp"
	pb "github.com/openconfig/gnmi/proto/gnmi"
	"io/ioutil"
	"net"
	"os"
	"ovs-gnxi/client/gnmi"
	"ovs-gnxi/client/gnoi"
	"ovs-gnxi/shared/logging"
	"ovs-gnxi/target/ovs"
	"strings"
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

const (
	ovsProtocol = "tcp"
)

var (
	log             = logging.New("ovs-gnxi-client")
	targetAddr      = flag.String("target_addr", "target", "The target address")
	targetGNMIPort  = flag.String("target_gnmi_port", "10161", "The target gNMI port")
	targetGNOIPort  = flag.String("target_gnoi_port", "10161", "The target gNOI port")
	targetName      = flag.String("target_name", "target.gnxi.lan", "The target name use to verify the hostname returned by TLS handshake")
	encodingName    = flag.String("encoding", "JSON_IETF", "The value encoding format to be used")
	timeOut         = flag.Duration("time_out", 10*time.Second, "Timeout for the Get request, 10 seconds by default")
	method          = flag.String("method", "", "A valid gNMI specification method to execute against the target")
	subscribeMode   = flag.String("subscribe_mode", "ONCE", "The gNMI data subscription mode")
	ovsAddr         = flag.String("ovs_address", "target.gnxi.lan:6640", "The ovs address in the format of host:port")
	caCert          = "certs/ca.crt"
	caKey           = "certs/ca.key"
	clientCert      = "certs/client.crt"
	clientKey       = "certs/client.key"
	getXPaths       arrayFlags
	deleteXPaths    arrayFlags
	replaceXPaths   arrayFlags
	updateXPaths    arrayFlags
	subscribeXPaths arrayFlags
)

func loadCAEntity(cert, key string) *entity.Entity {
	caEntity, err := entity.FromFile(cert, key)
	if err != nil {
		log.Fatalf("Failed to load certificate and key from file: %v", err)
	}

	return caEntity
}

func setGNMIClientFlags() {
	flag.Var(&getXPaths, "get_xpath", "The gNMI Get XPaths to query for")
	flag.Var(&deleteXPaths, "set_delete_xpath", "The gNMI Set Delete XPaths to query for")
	flag.Var(&replaceXPaths, "set_replace_xpath", "The gNMI Set Replace XPaths to query for")
	flag.Var(&updateXPaths, "set_update_xpath", "The gNMI Set Update XPaths to query for")
	flag.Var(&subscribeXPaths, "subscribe_xpath", "The gNMI Subscribe XPaths to query for")

	err := flag.Set("ca", caCert)
	if err != nil {
		log.Fatalf("Unable to set ca flag: %v", err)
	}
	err = flag.Set("cert", clientCert)
	if err != nil {
		log.Fatalf("Unable to set cert flag: %v", err)
	}
	err = flag.Set("key", clientKey)
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

	gnmiTarget := net.JoinHostPort(*targetAddr, *targetGNMIPort)
	gnoiTarget := net.JoinHostPort(*targetAddr, *targetGNOIPort)

	gnmiClient := gnmi.NewClient(gnmiTarget, *targetName, *encodingName)
	gnoiClient := gnoi.NewClient(gnoiTarget, *targetName, *encodingName, loadCAEntity(caCert, caKey))

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
		RunGNOIRebootTests(gnoiClient)
		RunGNMIGetTests(gnmiClient)
		RunGNOIGetCertificatesTests(gnoiClient)
		//RunGNOIRotateCertificatesTests(gnoiClient)
		RunGNMIGetTests(gnmiClient)
		RunGNMISetTests(gnmiClient)
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

	con := strings.Split(*ovsAddr, ":")

	ovsClient, err := NewClient(con[0], ovsProtocol, con[1], clientKey, clientCert, caCert)
	if err != nil {
		log.Errorf("Unable to initialize OVS Client: %v", err)
		os.Exit(1)
	}

	for _, td := range gnmi.SetTests {
		log.Infof("Testing GNMI Set(%v, %v, %v)...", td.DeleteXPaths, td.ReplaceXPaths, td.UpdateXPaths)

		result, err := ovsClient.Get(td.OVSDataBefore, ovs.ControllerTable)
		if val, ok := result[td.OVSResultKey]; !ok {
			log.Errorf("Set(%v) OVS Data: expected %v, actual %v", td.UpdateXPaths, "", val)
		}

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

		if len(td.DeleteXPaths) > 0 {
			if verifiedDeletePaths != len(td.DeleteXPaths) {
				log.Errorf("Set(%v): expected %v deletes, actual %v deletes", len(td.DeleteXPaths), verifiedDeletePaths)
			}
		}

		if len(td.ReplaceXPaths) > 0 {
			if verifiedReplacePaths != len(td.ReplaceXPaths) {
				log.Errorf("Set(%v): expected %v replaces, actual %v replaces", len(td.ReplaceXPaths), verifiedReplacePaths)
			}
		}

		if len(td.ReplaceXPaths) > 0 {
			if verifiedUpdatePaths != len(td.UpdateXPaths) {
				log.Errorf("Set(%v): expected %v updates, actual %v updates", len(td.UpdateXPaths), verifiedUpdatePaths)
			}
		}

		var getXPathsAfter []string

		for _, getXPath := range td.UpdateXPaths {
			path := strings.Split(getXPath, ":")

			getXPathsAfter = append(getXPathsAfter, path[0])
		}

		respGet, err := c.Get(ctx, getXPathsAfter)
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

		result, err = ovsClient.Get(td.OVSDataBefore, ovs.ControllerTable)
		if val, ok := result[td.OVSResultKey]; ok {
			log.Errorf("Set(%v) OVS Data: expected %v, actual %v", td.UpdateXPaths, "", val)
		}

		result, err = ovsClient.Get(td.OVSDataAfter, ovs.ControllerTable)
		if val, ok := result[td.OVSResultKey]; !ok {
			log.Errorf("Set(%v) OVS Data: expected %v, actual %v", td.UpdateXPaths, td.OVSDataAfter, val)
		}

		respSet, err = c.Set(ctx, td.RollbackDeleteXPaths, td.RollbackReplaceXPaths, td.RollbackUpdateXPaths)
		if err != nil {
			log.Fatal(err)
			continue
		}

		var getXPathsBefore []string

		for _, getXPath := range td.RollbackUpdateXPaths {
			path := strings.Split(getXPath, ":")

			getXPathsBefore = append(getXPathsBefore, path[0])
		}

		respGet, err = c.Get(ctx, getXPathsBefore)
		if err != nil {
			log.Fatal(err)
			continue
		}

		if td.ExtractorString != nil {
			actResp := td.ExtractorString(respGet.Notification)

			if actResp != td.RollbackExpResp {
				log.Errorf("Rollback Set(%v) Update: expected %v, actual %v", td.UpdateXPaths, td.RollbackExpResp, actResp)
			} else {
				log.Infof("Successfully verified Rollback GNMI Set(%v) Update with response value %v", td.UpdateXPaths, actResp)
			}
		} else if td.ExtractorUInt != nil {
			actResp := td.ExtractorUInt(respGet.Notification)

			if td.ExpResp != nil {
				if actResp != td.RollbackExpResp.(uint64) {
					log.Errorf("Rollback Set(%v) Update: expected %v, actual %v", td.UpdateXPaths, td.RollbackExpResp, actResp)
				} else {
					log.Infof("Successfully verified Rollback GNMI Set(%v) Update with response value %v", td.UpdateXPaths, actResp)
				}
			}
		}

		result, err = ovsClient.Get(td.OVSDataBefore, ovs.ControllerTable)
		if val, ok := result[td.OVSResultKey]; !ok {
			log.Errorf("Set(%v) OVS Data: expected %v, actual %v", td.UpdateXPaths, "", val)
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

func RunGNOIRebootTests(c *gnoi.Client) {
	ctx, cancel := context.WithTimeout(context.Background(), *timeOut)
	defer cancel()

	for _, td := range gnoi.RebootTests {
		log.Infof("Testing GNOI Reboot(%v)...", td.Desc)

		resp, err := c.Reboot(ctx, td.Message)
		if err != nil {
			log.Fatal(err)
		}

		if resp == nil {
			log.Errorf("Unexpected GNOI Reboot response")
		}

		for {
			conn, _ := net.DialTimeout("tcp", net.JoinHostPort(*targetAddr, *targetGNOIPort), *timeOut)
			if conn != nil {
				conn.Close()
				log.Info("Successfully verified that target device is back up")
				time.Sleep(10 * time.Second)
				break
			}
		}

		log.Infof("Successfully verified GNOI Reboot(%v)", td.Desc)
	}
}

func RunGNOIGetCertificatesTests(c *gnoi.Client) {
	ctx, cancel := context.WithTimeout(context.Background(), *timeOut)
	defer cancel()

	for _, td := range gnoi.GetCertificatesTests {
		log.Infof("Testing GNOI GetCertificates(%v)...", td.Desc)

		resp, err := c.GetCertificates(ctx)
		if err != nil {
			log.Fatal(err)
		}

		certFile, err := ioutil.ReadFile(td.ExpCertPath)
		if err != nil {
			log.Fatal(err)
		}

		certBlock, _ := pem.Decode(certFile)
		if certBlock == nil {
			log.Fatal(err)
		}

		certs, err := x509.ParseCertificates(certBlock.Bytes)
		if err != nil {
			log.Fatal(err)
		}

		if bytes.Compare(resp[td.ExpCertID].Signature, certs[0].Signature) == 0 {
			log.Infof("Successfully verified GNOI GetCertificates(%v)", td.Desc)
		} else {
			log.Errorf("GetCertificates(%v): expected %v, actual %v", td.Desc, string(certs[0].Signature), string(resp[td.ExpCertID].Signature))
		}
	}
}

func RunGNOIRotateCertificatesTests(c *gnoi.Client) {
	ctx, cancel := context.WithTimeout(context.Background(), *timeOut)
	defer cancel()

	for _, td := range gnoi.RotateCertificatesTests {
		log.Infof("Testing GNOI Reboot(%v)...", td.Desc)

		err := c.RotateCertificates(ctx, td.CertID)
		if err != nil {
			log.Fatal(err)
		} else {
			log.Infof("Successfully verified GNOI Rotate Certificates(%v)", td.Desc)
		}
	}
}
