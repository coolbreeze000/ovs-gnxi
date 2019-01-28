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

package gnmi

import (
	"bytes"
	"fmt"
	"github.com/go-errors/errors"
	"github.com/golang/protobuf/proto"
	"github.com/google/gnxi/utils/credentials"
	"github.com/google/gnxi/utils/xpath"
	"google.golang.org/grpc"
	"io"
	"io/ioutil"
	"ovs-gnxi/shared/logging"
	"strconv"
	"strings"

	pb "github.com/openconfig/gnmi/proto/gnmi"
	"golang.org/x/net/context"
)

var (
	log = logging.New("ovs-gnxi-client")
)

type Client struct {
	targetAddress string
	targetName    string
	encodingName  string
}

// NewGNMIClient returns an instance of GNMIClient struct.
func NewGNMIClient(targetAddress, targetName, encodingName string) *Client {
	return &Client{
		targetAddress: targetAddress,
		targetName:    targetName,
		encodingName:  encodingName,
	}
}

func (c *Client) Capabilities(ctx context.Context) (*pb.CapabilityResponse, error) {
	opts := credentials.ClientCredentials(c.targetName)
	conn, err := grpc.Dial(c.targetAddress, opts...)
	if err != nil {
		return nil, err
	}
	defer conn.Close()

	cli := pb.NewGNMIClient(conn)

	request := &pb.CapabilityRequest{}

	log.Debug("== Request:")
	log.Debug(proto.MarshalTextString(request))

	response, err := cli.Capabilities(ctx, request)
	if err != nil {
		return nil, errors.New(fmt.Sprintf("error in getting capabilities: %v", err))
	}

	return response, nil
}

func (c *Client) Get(ctx context.Context, getXPaths []string) (*pb.GetResponse, error) {
	opts := credentials.ClientCredentials(c.targetName)
	conn, err := grpc.Dial(c.targetAddress, opts...)
	if err != nil {
		return nil, err
	}
	defer conn.Close()

	cli := pb.NewGNMIClient(conn)

	encoding, ok := pb.Encoding_value[c.encodingName]
	if !ok {
		var gnmiEncodingList []string
		for _, name := range pb.Encoding_name {
			gnmiEncodingList = append(gnmiEncodingList, name)
		}
		return nil, fmt.Errorf("supported encodings: %s", strings.Join(gnmiEncodingList, ", "))
	}

	var pbPathList []*pb.Path
	var pbModelDataList []*pb.ModelData
	for _, xPath := range getXPaths {
		pbPath, err := xpath.ToGNMIPath(xPath)
		if err != nil {
			return nil, errors.New(fmt.Sprintf("error in parsing xpath %q to gnxi path", xPath))
		}
		pbPathList = append(pbPathList, pbPath)
	}

	request := &pb.GetRequest{
		Encoding:  pb.Encoding(encoding),
		Path:      pbPathList,
		UseModels: pbModelDataList,
	}

	log.Debug("== Request:")
	log.Debug(proto.MarshalTextString(request))

	response, err := cli.Get(ctx, request)
	if err != nil {
		return nil, errors.New(fmt.Sprintf("Get failed: %v", err))
	}

	return response, nil
}

func buildPbUpdateList(pathValuePairs []string) ([]*pb.Update, error) {
	var pbUpdateList []*pb.Update
	for _, item := range pathValuePairs {
		splitIndex := strings.LastIndexAny(item, ":")
		if splitIndex < 1 {
			return nil, fmt.Errorf("invalid path-value pair: %v", item)
		}
		pathValuePair := []string{item[:splitIndex], item[(splitIndex + 1):]}
		pbPath, err := xpath.ToGNMIPath(pathValuePair[0])
		if err != nil {
			return nil, fmt.Errorf("error in parsing xpath %q to gnxi path", pathValuePair[0])
		}
		var pbVal *pb.TypedValue
		if pathValuePair[1][0] == '@' {
			jsonFile := pathValuePair[1][1:]
			jsonConfig, err := ioutil.ReadFile(jsonFile)
			if err != nil {
				return nil, fmt.Errorf("cannot read data from file %v", jsonFile)
			}
			jsonConfig = bytes.Trim(jsonConfig, " \r\n\t")
			pbVal = &pb.TypedValue{
				Value: &pb.TypedValue_JsonIetfVal{
					JsonIetfVal: jsonConfig,
				},
			}
		} else {
			if strVal, err := strconv.Unquote(pathValuePair[1]); err == nil {
				pbVal = &pb.TypedValue{
					Value: &pb.TypedValue_StringVal{
						StringVal: strVal,
					},
				}
			} else {
				if intVal, err := strconv.ParseInt(pathValuePair[1], 10, 64); err == nil {
					pbVal = &pb.TypedValue{
						Value: &pb.TypedValue_IntVal{
							IntVal: intVal,
						},
					}
				} else if floatVal, err := strconv.ParseFloat(pathValuePair[1], 32); err == nil {
					pbVal = &pb.TypedValue{
						Value: &pb.TypedValue_FloatVal{
							FloatVal: float32(floatVal),
						},
					}
				} else if boolVal, err := strconv.ParseBool(pathValuePair[1]); err == nil {
					pbVal = &pb.TypedValue{
						Value: &pb.TypedValue_BoolVal{
							BoolVal: boolVal,
						},
					}
				} else {
					pbVal = &pb.TypedValue{
						Value: &pb.TypedValue_StringVal{
							StringVal: pathValuePair[1],
						},
					}
				}
			}
		}
		pbUpdateList = append(pbUpdateList, &pb.Update{Path: pbPath, Val: pbVal})
	}

	return pbUpdateList, nil
}

func (c *Client) Set(ctx context.Context, deleteXPaths, replaceXPaths, updateXPaths []string) (*pb.SetResponse, error) {
	opts := credentials.ClientCredentials(c.targetName)
	conn, err := grpc.Dial(c.targetAddress, opts...)
	if err != nil {
		return nil, err
	}
	defer conn.Close()

	cli := pb.NewGNMIClient(conn)

	var deleteList []*pb.Path
	for _, xPath := range deleteXPaths {
		pbPath, err := xpath.ToGNMIPath(xPath)
		if err != nil {
			return nil, fmt.Errorf("error in parsing xpath %q to gnxi path", xPath)
		}
		deleteList = append(deleteList, pbPath)
	}
	replaceList, err := buildPbUpdateList(replaceXPaths)
	if err != nil {
		return nil, err
	}
	updateList, err := buildPbUpdateList(updateXPaths)
	if err != nil {
		return nil, err
	}

	request := &pb.SetRequest{
		Delete:  deleteList,
		Replace: replaceList,
		Update:  updateList,
	}

	log.Debug("== Request:")
	log.Debug(proto.MarshalTextString(request))

	response, err := cli.Set(ctx, request)
	if err != nil {
		return nil, fmt.Errorf("set failed: %v", err)
	}

	return response, nil
}

func (c *Client) SubscribeStream(ctx context.Context, subscribeXPaths []string, respChan chan<- *pb.SubscribeResponse, errChan chan<- error) {
	opts := credentials.ClientCredentials(c.targetName)
	conn, err := grpc.Dial(c.targetAddress, opts...)
	if err != nil {
		errChan <- err
		return
	}
	defer conn.Close()

	cli := pb.NewGNMIClient(conn)

	defer close(respChan)

	encoding, ok := pb.Encoding_value[c.encodingName]
	if !ok {
		var gnmiEncodingList []string
		for _, name := range pb.Encoding_name {
			gnmiEncodingList = append(gnmiEncodingList, name)
		}
		errChan <- fmt.Errorf("supported encodings: %s", strings.Join(gnmiEncodingList, ", "))
		return
	}

	var pbModelDataList []*pb.ModelData
	var subscriptions []*pb.Subscription

	for _, xPath := range subscribeXPaths {
		pbPath, err := xpath.ToGNMIPath(xPath)
		if err != nil {
			errChan <- err
			return
		}
		subscriptions = append(subscriptions, &pb.Subscription{Path: pbPath})
	}

	request := &pb.SubscribeRequest{
		Request: &pb.SubscribeRequest_Subscribe{
			Subscribe: &pb.SubscriptionList{
				Prefix:       &pb.Path{Target: c.targetName},
				Mode:         pb.SubscriptionList_STREAM,
				UseModels:    pbModelDataList,
				Subscription: subscriptions,
				Encoding:     pb.Encoding(encoding),
			},
		},
	}

	log.Debug("== Request:")
	log.Debug(proto.MarshalTextString(request))

	subClient, err := cli.Subscribe(ctx)
	if err != nil {
		errChan <- err
		return
	}

	err = subClient.Send(request)
	if err != nil {
		errChan <- err
		return
	}

	for {
		resp, err := subClient.Recv()
		if err != nil {
			if err == io.EOF {
				return
			}
			errChan <- err
			return
		}

		select {
		case <-ctx.Done():
			log.Info("Stopped gNMI Subscribe Stream Client")
			return
		case respChan <- resp:
		}
	}
}

func (c *Client) SubscribeOnce(ctx context.Context, subscribeXPaths []string) (*pb.SubscribeResponse, error) {
	opts := credentials.ClientCredentials(c.targetName)
	conn, err := grpc.Dial(c.targetAddress, opts...)
	if err != nil {
		return nil, err
	}
	defer conn.Close()

	cli := pb.NewGNMIClient(conn)

	encoding, ok := pb.Encoding_value[c.encodingName]
	if !ok {
		var gnmiEncodingList []string
		for _, name := range pb.Encoding_name {
			gnmiEncodingList = append(gnmiEncodingList, name)
		}
		return nil, fmt.Errorf("supported encodings: %s", strings.Join(gnmiEncodingList, ", "))
	}

	var pbModelDataList []*pb.ModelData
	var subscriptions []*pb.Subscription

	for _, xPath := range subscribeXPaths {
		pbPath, err := xpath.ToGNMIPath(xPath)
		if err != nil {
			return nil, fmt.Errorf("error in parsing xpath %q to gnxi path", xPath)
		}
		subscriptions = append(subscriptions, &pb.Subscription{Path: pbPath})
	}

	request := &pb.SubscribeRequest{
		Request: &pb.SubscribeRequest_Subscribe{
			Subscribe: &pb.SubscriptionList{
				Prefix:       &pb.Path{Target: c.targetName},
				Mode:         pb.SubscriptionList_ONCE,
				UseModels:    pbModelDataList,
				Subscription: subscriptions,
				Encoding:     pb.Encoding(encoding),
			},
		},
	}

	log.Debug("== Request:")
	log.Debug(proto.MarshalTextString(request))

	subClient, err := cli.Subscribe(ctx)
	if err != nil {
		return nil, fmt.Errorf("subscribe failed: %v", err)
	}

	err = subClient.Send(request)
	if err != nil {
		return nil, fmt.Errorf("subscribe send failed: %v", err)
	}
	response, err := subClient.Recv()
	if err != nil {
		return nil, fmt.Errorf("subscribe recv failed: %v", err)
	}

	return response, nil
}
