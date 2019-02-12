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

// Package gnxi implements a gnxi server.
package service

import (
	"bytes"
	"compress/gzip"
	"crypto/tls"
	"crypto/x509"
	"encoding/json"
	"fmt"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/reflection"
	"io"
	"io/ioutil"
	"net"
	"ovs-gnxi/shared"
	"ovs-gnxi/shared/logging"
	"reflect"
	"strconv"
	"sync"
	"time"

	"golang.org/x/net/context"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"ovs-gnxi/target/gnxi/service/gnmi"

	"github.com/golang/protobuf/proto"
	"github.com/openconfig/gnmi/value"
	"github.com/openconfig/ygot/experimental/ygotutils"
	"github.com/openconfig/ygot/ygot"

	dpb "github.com/golang/protobuf/protoc-gen-go/descriptor"
	pbg "github.com/openconfig/gnmi/proto/gnmi"
	cpb "google.golang.org/genproto/googleapis/rpc/code"
	pbc "ovs-gnxi/shared/gnoi/modeldata/generated/cert"
	pbs "ovs-gnxi/shared/gnoi/modeldata/generated/system"
)

var log = logging.New("ovs-gnxi")

var (
	pbRootPath         = &pbg.Path{}
	supportedEncodings = []pbg.Encoding{pbg.Encoding_JSON, pbg.Encoding_JSON_IETF}
	gnxiProtocol       = "tcp"
	gnxiPort           = "10161"
)

type ConfigSetupCallback func(ygot.ValidatedGoStruct) error

// ConfigChangeCallback is the signature of the function to apply a validated config to the physical device.
type ConfigChangeCallback func(ygot.ValidatedGoStruct) error

type RebootCallback func() error
type RotateCertificatesCallback func(certificates *shared.TargetCertificates) error

type CallbackHandler struct {
	CallbackSetup       ConfigSetupCallback
	CallbackChange      ConfigChangeCallback
	CallbackReboot      RebootCallback
	CallbackRotateCerts RotateCertificatesCallback
}

// Service struct maintains the data structure for device config and implements the gnxi interface. It supports Capabilities, Get, Set and Subscribe APIs.
type Service struct {
	g            *grpc.Server
	socket       net.Listener
	certs        *shared.TargetCertificates
	auth         *shared.Authenticator
	model        *gnmi.Model
	config       ygot.ValidatedGoStruct
	ch           *CallbackHandler
	mu           sync.RWMutex // mu is the RW lock to protect the access to config
	ConfigUpdate chan bool

	timeout time.Duration
}

// NewService creates an instance of Service with given json config.
func NewService(auth *shared.Authenticator, model *gnmi.Model, certs *shared.TargetCertificates, config []byte,
	callbackSetup ConfigSetupCallback, callbackChange ConfigChangeCallback, callbackReboot RebootCallback, callbackRotateCerts RotateCertificatesCallback) (*Service, error) {
	rootStruct, err := model.NewConfigStruct(config)

	if err != nil {
		return nil, err
	}
	s := &Service{
		certs:        certs,
		auth:         auth,
		model:        model,
		config:       rootStruct,
		ConfigUpdate: make(chan bool),
		ch: &CallbackHandler{
			CallbackSetup:       callbackSetup,
			CallbackChange:      callbackChange,
			CallbackReboot:      callbackReboot,
			CallbackRotateCerts: callbackRotateCerts,
		},
	}

	if config != nil && s.ch.CallbackSetup != nil {
		if err := s.ch.CallbackSetup(rootStruct); err != nil {
			return nil, err
		}
	}

	return s, nil
}

func (s *Service) LockService() {
	s.mu.Lock()
}

func (s *Service) UnlockService() {
	s.mu.Unlock()
}

// checkEncodingAndModel checks whether encoding and models are supported by the server. Return error if anything is unsupported.
func (s *Service) checkEncodingAndModel(encoding pbg.Encoding, models []*pbg.ModelData) error {
	hasSupportedEncoding := false
	for _, supportedEncoding := range supportedEncodings {
		if encoding == supportedEncoding {
			hasSupportedEncoding = true
			break
		}
	}
	if !hasSupportedEncoding {
		return fmt.Errorf("unsupported encoding: %s", pbg.Encoding_name[int32(encoding)])
	}
	for _, m := range models {
		isSupported := false
		for _, supportedModel := range s.model.ModelData {
			if reflect.DeepEqual(m, supportedModel) {
				isSupported = true
				break
			}
		}
		if !isSupported {
			return fmt.Errorf("unsupported model: %v", m)
		}
	}
	return nil
}

// doDelete deletes the path from the json tree if the path exists. If success,
// it calls the callback function to apply the change to the device hardware.
func (s *Service) doDelete(jsonTree map[string]interface{}, prefix, path *pbg.Path) (*pbg.UpdateResult, error) {
	// Update json tree of the device config
	var curNode interface{} = jsonTree
	pathDeleted := false
	fullPath := gnmiFullPath(prefix, path)
	schema := s.model.SchemaTreeRoot
	for i, elem := range fullPath.Elem { // Delete sub-tree or leaf node.
		node, ok := curNode.(map[string]interface{})
		if !ok {
			break
		}

		// Delete node
		if i == len(fullPath.Elem)-1 {
			if elem.GetKey() == nil {
				delete(node, elem.Name)
				pathDeleted = true
				break
			}
			pathDeleted = deleteKeyedListEntry(node, elem)
			break
		}

		if curNode, schema = gnmi.GetChildNode(node, schema, elem, false); curNode == nil {
			break
		}
	}
	if reflect.DeepEqual(fullPath, pbRootPath) { // Delete root
		for k := range jsonTree {
			delete(jsonTree, k)
		}
	}

	// Apply the validated operation to the config tree and device.
	if pathDeleted {
		newConfig, err := s.toGoStruct(jsonTree)
		if err != nil {
			return nil, status.Error(codes.Internal, err.Error())
		}
		if s.ch.CallbackChange != nil {
			if applyErr := s.ch.CallbackChange(newConfig); applyErr != nil {
				if rollbackErr := s.ch.CallbackChange(s.config); rollbackErr != nil {
					return nil, status.Errorf(codes.Internal, "error in rollback the failed operation (%v): %v", applyErr, rollbackErr)
				}
				return nil, status.Errorf(codes.Aborted, "error in applying operation to device: %v", applyErr)
			}
		}
	}
	return &pbg.UpdateResult{
		Path: path,
		Op:   pbg.UpdateResult_DELETE,
	}, nil
}

// doReplaceOrUpdate validates the replace or update operation to be applied to
// the device, modifies the json tree of the config struct, then calls the
// callback function to apply the operation to the device hardware.
func (s *Service) doReplaceOrUpdate(jsonTree map[string]interface{}, op pbg.UpdateResult_Operation, prefix, path *pbg.Path, val *pbg.TypedValue) (*pbg.UpdateResult, error) {
	// Validate the operation.
	fullPath := gnmiFullPath(prefix, path)
	emptyNode, stat := ygotutils.NewNode(s.model.StructRootType, fullPath)
	if stat.GetCode() != int32(cpb.Code_OK) {
		return nil, status.Errorf(codes.NotFound, "path %v is not found in the config structure: %v", fullPath, stat)
	}
	var nodeVal interface{}
	nodeStruct, ok := emptyNode.(ygot.ValidatedGoStruct)
	if ok {
		if err := s.model.JSONUnmarshaler(val.GetJsonIetfVal(), nodeStruct); err != nil {
			return nil, status.Errorf(codes.InvalidArgument, "unmarshaling json data to config struct fails: %v", err)
		}
		if err := nodeStruct.Validate(); err != nil {
			return nil, status.Errorf(codes.InvalidArgument, "config data validation fails: %v", err)
		}
		var err error
		if nodeVal, err = ygot.ConstructIETFJSON(nodeStruct, &ygot.RFC7951JSONConfig{}); err != nil {
			msg := fmt.Sprintf("error in constructing IETF JSON tree from config struct: %v", err)
			log.Error(msg)
			return nil, status.Error(codes.Internal, msg)
		}
	} else {
		var err error
		if nodeVal, err = value.ToScalar(val); err != nil {
			return nil, status.Errorf(codes.Internal, "cannot convert leaf node to scalar type: %v", err)
		}
	}

	// Update json tree of the device config.
	var curNode interface{} = jsonTree
	schema := s.model.SchemaTreeRoot
	for i, elem := range fullPath.Elem {
		switch node := curNode.(type) {
		case map[string]interface{}:
			// Set node value.
			if i == len(fullPath.Elem)-1 {
				if elem.GetKey() == nil {
					if grpcStatusError := setPathWithoutAttribute(op, node, elem, nodeVal); grpcStatusError != nil {
						return nil, grpcStatusError
					}
					break
				}
				if grpcStatusError := setPathWithAttribute(op, node, elem, nodeVal); grpcStatusError != nil {
					return nil, grpcStatusError
				}
				break
			}

			if curNode, schema = gnmi.GetChildNode(node, schema, elem, true); curNode == nil {
				return nil, status.Errorf(codes.NotFound, "path elem not found: %v", elem)
			}
		case []interface{}:
			return nil, status.Errorf(codes.NotFound, "uncompatible path elem: %v", elem)
		default:
			return nil, status.Errorf(codes.Internal, "wrong node type: %T", curNode)
		}
	}
	if reflect.DeepEqual(fullPath, pbRootPath) { // Replace/Update root.
		if op == pbg.UpdateResult_UPDATE {
			return nil, status.Error(codes.Unimplemented, "update the root of config tree is unsupported")
		}
		nodeValAsTree, ok := nodeVal.(map[string]interface{})
		if !ok {
			return nil, status.Errorf(codes.InvalidArgument, "expect a tree to replace the root, got a scalar value: %T", nodeVal)
		}
		for k := range jsonTree {
			delete(jsonTree, k)
		}
		for k, v := range nodeValAsTree {
			jsonTree[k] = v
		}
	}
	newConfig, err := s.toGoStruct(jsonTree)
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	// Apply the validated operation to the device.
	if s.ch.CallbackChange != nil {
		if applyErr := s.ch.CallbackChange(newConfig); applyErr != nil {
			if rollbackErr := s.ch.CallbackChange(s.config); rollbackErr != nil {
				return nil, status.Errorf(codes.Internal, "error in rollback the failed operation (%v): %v", applyErr, rollbackErr)
			}
			return nil, status.Errorf(codes.Aborted, "error in applying operation to device: %v", applyErr)
		}
	}
	return &pbg.UpdateResult{
		Path: path,
		Op:   op,
	}, nil
}

func (s *Service) toGoStruct(jsonTree map[string]interface{}) (ygot.ValidatedGoStruct, error) {
	jsonDump, err := json.Marshal(jsonTree)
	if err != nil {
		return nil, fmt.Errorf("error in marshaling IETF JSON tree to bytes: %v", err)
	}
	goStruct, err := s.model.NewConfigStruct(jsonDump)
	if err != nil {
		return nil, fmt.Errorf("error in creating config struct from IETF JSON data: %v", err)
	}
	return goStruct, nil
}

// getGNMIServiceVersion returns a pointer to the gNMI service version string.
// The method is non-trivial because of the way it is defined in the proto file.
func getGNMIServiceVersion() (*string, error) {
	gzB, _ := (&pbg.Update{}).Descriptor()
	r, err := gzip.NewReader(bytes.NewReader(gzB))
	if err != nil {
		return nil, fmt.Errorf("error in initializing gzip reader: %v", err)
	}
	defer r.Close()
	b, err := ioutil.ReadAll(r)
	if err != nil {
		return nil, fmt.Errorf("error in reading gzip data: %v", err)
	}
	desc := &dpb.FileDescriptorProto{}
	if err := proto.Unmarshal(b, desc); err != nil {
		return nil, fmt.Errorf("error in unmarshaling proto: %v", err)
	}
	ver, err := proto.GetExtension(desc.Options, pbg.E_GnmiService)
	if err != nil {
		return nil, fmt.Errorf("error in getting version from proto extension: %v", err)
	}
	return ver.(*string), nil
}

// deleteKeyedListEntry deletes the keyed list entry from node that matches the
// path elem. If the entry is the only one in keyed list, deletes the entire
// list. If the entry is found and deleted, the function returns true. If it is
// not found, the function returns false.
func deleteKeyedListEntry(node map[string]interface{}, elem *pbg.PathElem) bool {
	curNode, ok := node[elem.Name]
	if !ok {
		return false
	}

	keyedList, ok := curNode.([]interface{})
	if !ok {
		return false
	}
	for i, n := range keyedList {
		m, ok := n.(map[string]interface{})
		if !ok {
			log.Errorf("expect map[string]interface{} for a keyed list entry, got %T", n)
			return false
		}
		keyMatching := true
		for k, v := range elem.Key {
			attrVal, ok := m[k]
			if !ok {
				return false
			}
			if v != fmt.Sprintf("%v", attrVal) {
				keyMatching = false
				break
			}
		}
		if keyMatching {
			listLen := len(keyedList)
			if listLen == 1 {
				delete(node, elem.Name)
				return true
			}
			keyedList[i] = keyedList[listLen-1]
			node[elem.Name] = keyedList[0 : listLen-1]
			return true
		}
	}
	return false
}

// gnmiFullPath builds the full path from the prefix and path.
func gnmiFullPath(prefix, path *pbg.Path) *pbg.Path {
	fullPath := &pbg.Path{Origin: path.Origin}
	if path.GetElement() != nil {
		fullPath.Element = append(prefix.GetElement(), path.GetElement()...)
	}
	if path.GetElem() != nil {
		fullPath.Elem = append(prefix.GetElem(), path.GetElem()...)
	}
	return fullPath
}

// isNIl checks if an interface is nil or its value is nil.
func isNil(i interface{}) bool {
	if i == nil {
		return true
	}
	switch kind := reflect.ValueOf(i).Kind(); kind {
	case reflect.Chan, reflect.Func, reflect.Interface, reflect.Map, reflect.Ptr, reflect.Slice:
		return reflect.ValueOf(i).IsNil()
	default:
		return false
	}
}

// setPathWithAttribute replaces or updates a child node of curNode in the IETF
// JSON config tree, where the child node is indexed by pathElem with attribute.
// The function returns grpc status error if unsuccessful.
func setPathWithAttribute(op pbg.UpdateResult_Operation, curNode map[string]interface{}, pathElem *pbg.PathElem, nodeVal interface{}) error {
	nodeValAsTree, ok := nodeVal.(map[string]interface{})
	if !ok {
		return status.Errorf(codes.InvalidArgument, "expect nodeVal is a json node of map[string]interface{}, received %T", nodeVal)
	}
	m := gnmi.GetKeyedListEntry(curNode, pathElem, true)
	if m == nil {
		return status.Errorf(codes.NotFound, "path elem not found: %v", pathElem)
	}
	if op == pbg.UpdateResult_REPLACE {
		for k := range m {
			delete(m, k)
		}
	}
	for attrKey, attrVal := range pathElem.GetKey() {
		m[attrKey] = attrVal
		if asNum, err := strconv.ParseFloat(attrVal, 64); err == nil {
			m[attrKey] = asNum
		}
		for k, v := range nodeValAsTree {
			if k == attrKey && fmt.Sprintf("%v", v) != attrVal {
				return status.Errorf(codes.InvalidArgument, "invalid config data: %v is a path attribute", k)
			}
		}
	}
	for k, v := range nodeValAsTree {
		m[k] = v
	}
	return nil
}

// setPathWithoutAttribute replaces or updates a child node of curNode in the
// IETF config tree, where the child node is indexed by pathElem without
// attribute. The function returns grpc status error if unsuccessful.
func setPathWithoutAttribute(op pbg.UpdateResult_Operation, curNode map[string]interface{}, pathElem *pbg.PathElem, nodeVal interface{}) error {
	target, hasElem := curNode[pathElem.Name]
	nodeValAsTree, nodeValIsTree := nodeVal.(map[string]interface{})
	if op == pbg.UpdateResult_REPLACE || !hasElem || !nodeValIsTree {
		curNode[pathElem.Name] = nodeVal
		return nil
	}
	targetAsTree, ok := target.(map[string]interface{})
	if !ok {
		return status.Errorf(codes.Internal, "error in setting path: expect map[string]interface{} to update, got %T", target)
	}
	for k, v := range nodeValAsTree {
		targetAsTree[k] = v
	}
	return nil
}

// Capabilities returns supported encodings and supported models.
func (s *Service) Capabilities(ctx context.Context, req *pbg.CapabilityRequest) (*pbg.CapabilityResponse, error) {
	authorized, err := s.auth.AuthorizeUser(ctx)
	if !authorized {
		log.Infof("denied a Capabilities request: %v", err)
		return nil, status.Error(codes.PermissionDenied, fmt.Sprint(err))
	}
	log.Infof("allowed a Capabilities request")

	ver, err := getGNMIServiceVersion()
	if err != nil {
		return nil, status.Errorf(codes.Internal, "error in getting gnxi service version: %v", err)
	}

	resp := &pbg.CapabilityResponse{
		SupportedModels:    s.model.ModelData,
		SupportedEncodings: supportedEncodings,
		GNMIVersion:        *ver,
	}

	log.Infof("Send Capability response to client: %v", resp)

	return resp, nil
}

// Get implements the Get RPC in gNMI spec and provides user auth.
func (s *Service) Get(ctx context.Context, req *pbg.GetRequest) (*pbg.GetResponse, error) {
	authorized, err := s.auth.AuthorizeUser(ctx)
	if !authorized {
		log.Infof("denied a Get request: %v", err)
		return nil, status.Error(codes.PermissionDenied, fmt.Sprint(err))
	}
	log.Infof("allowed a Get request")

	if req.GetType() != pbg.GetRequest_ALL {
		return nil, status.Errorf(codes.Unimplemented, "unsupported request type: %s", pbg.GetRequest_DataType_name[int32(req.GetType())])
	}
	if err := s.checkEncodingAndModel(req.GetEncoding(), req.GetUseModels()); err != nil {
		return nil, status.Error(codes.Unimplemented, err.Error())
	}

	prefix := req.GetPrefix()
	paths := req.GetPath()
	notifications := make([]*pbg.Notification, len(paths))

	s.mu.RLock()
	defer s.mu.RUnlock()

	for i, path := range paths {
		// Get schema node for path from config struct.
		fullPath := path
		if prefix != nil {
			fullPath = gnmiFullPath(prefix, path)
		}
		if fullPath.GetElem() == nil && fullPath.GetElement() != nil {
			return nil, status.Error(codes.Unimplemented, "deprecated path element type is unsupported")
		}
		node, stat := ygotutils.GetNode(s.model.SchemaTreeRoot, s.config, fullPath)
		if isNil(node) || stat.GetCode() != int32(cpb.Code_OK) {
			return nil, status.Errorf(codes.NotFound, "path %v not found", fullPath)
		}

		ts := time.Now().UnixNano()

		nodeStruct, ok := node.(ygot.GoStruct)
		// Return leaf node.
		if !ok {
			var val *pbg.TypedValue
			switch kind := reflect.ValueOf(node).Kind(); kind {
			case reflect.Ptr, reflect.Interface:
				var err error
				val, err = value.FromScalar(reflect.ValueOf(node).Elem().Interface())
				if err != nil {
					msg := fmt.Sprintf("leaf node %v does not contain a scalar type value: %v", path, err)
					log.Error(msg)
					return nil, status.Error(codes.Internal, msg)
				}
			case reflect.Int64:
				enumMap, ok := s.model.EnumData[reflect.TypeOf(node).Name()]
				if !ok {
					return nil, status.Error(codes.Internal, "not a GoStruct enumeration type")
				}
				val = &pbg.TypedValue{
					Value: &pbg.TypedValue_StringVal{
						StringVal: enumMap[reflect.ValueOf(node).Int()].Name,
					},
				}
			default:
				return nil, status.Errorf(codes.Internal, "unexpected kind of leaf node type: %v %v", node, kind)
			}

			update := &pbg.Update{Path: path, Val: val}
			notifications[i] = &pbg.Notification{
				Timestamp: ts,
				Prefix:    prefix,
				Update:    []*pbg.Update{update},
			}
			continue
		}

		// Return all leaf nodes of the sub-tree.
		if len(req.GetUseModels()) != len(s.model.ModelData) && req.GetEncoding() != pbg.Encoding_JSON_IETF {
			results, err := ygot.TogNMINotifications(nodeStruct, ts, ygot.GNMINotificationsConfig{UsePathElem: true, PathElemPrefix: fullPath.Elem})
			if err != nil {
				return nil, status.Errorf(codes.Internal, "error in serializing GoStruct to notifications: %v", err)
			}
			if len(results) != 1 {
				return nil, status.Errorf(codes.Internal, "ygot.TogNMINotifications() return %d notifications instead of one", len(results))
			}
			notifications[i] = results[0]
			continue
		}

		// Return IETF JSON for the sub-tree.
		jsonTree, err := ygot.ConstructIETFJSON(nodeStruct, &ygot.RFC7951JSONConfig{AppendModuleName: true})
		if err != nil {
			msg := fmt.Sprintf("error in constructing IETF JSON tree from requested node: %v", err)
			log.Error(msg)
			return nil, status.Error(codes.Internal, msg)
		}
		jsonDump, err := json.Marshal(jsonTree)
		if err != nil {
			msg := fmt.Sprintf("error in marshaling IETF JSON tree to bytes: %v", err)
			log.Error(msg)
			return nil, status.Error(codes.Internal, msg)
		}
		update := &pbg.Update{
			Path: path,
			Val: &pbg.TypedValue{
				Value: &pbg.TypedValue_JsonIetfVal{
					JsonIetfVal: jsonDump,
				},
			},
		}
		notifications[i] = &pbg.Notification{
			Timestamp: ts,
			Prefix:    prefix,
			Update:    []*pbg.Update{update},
		}
	}

	resp := &pbg.GetResponse{Notification: notifications}

	log.Infof("Send Get response to client: %v", resp)

	return resp, nil
}

// Set implements the Set RPC in gNMI spec and provides user auth.
func (s *Service) Set(ctx context.Context, req *pbg.SetRequest) (*pbg.SetResponse, error) {
	authorized, err := s.auth.AuthorizeUser(ctx)
	if !authorized {
		log.Infof("denied a Set request: %v", err)
		return nil, status.Error(codes.PermissionDenied, fmt.Sprint(err))
	}
	log.Infof("allowed a Set request")

	s.mu.Lock()
	defer s.mu.Unlock()

	jsonTree, err := ygot.ConstructIETFJSON(s.config, &ygot.RFC7951JSONConfig{})
	if err != nil {
		msg := fmt.Sprintf("error in constructing IETF JSON tree from config struct: %v", err)
		log.Error(msg)
		return nil, status.Error(codes.Internal, msg)
	}

	prefix := req.GetPrefix()
	var results []*pbg.UpdateResult

	for _, path := range req.GetDelete() {
		res, grpcStatusError := s.doDelete(jsonTree, prefix, path)
		if grpcStatusError != nil {
			return nil, grpcStatusError
		}
		results = append(results, res)
	}
	for _, upd := range req.GetReplace() {
		res, grpcStatusError := s.doReplaceOrUpdate(jsonTree, pbg.UpdateResult_REPLACE, prefix, upd.GetPath(), upd.GetVal())
		if grpcStatusError != nil {
			return nil, grpcStatusError
		}
		results = append(results, res)
	}
	for _, upd := range req.GetUpdate() {
		res, grpcStatusError := s.doReplaceOrUpdate(jsonTree, pbg.UpdateResult_UPDATE, prefix, upd.GetPath(), upd.GetVal())
		if grpcStatusError != nil {
			return nil, grpcStatusError
		}
		results = append(results, res)
	}

	jsonDump, err := json.Marshal(jsonTree)
	if err != nil {
		msg := fmt.Sprintf("error in marshaling IETF JSON tree to bytes: %v", err)
		log.Error(msg)
		return nil, status.Error(codes.Internal, msg)
	}
	rootStruct, err := s.model.NewConfigStruct(jsonDump)
	if err != nil {
		msg := fmt.Sprintf("error in creating config struct from IETF JSON data: %v", err)
		log.Error(msg)
	}
	s.config = rootStruct

	resp := &pbg.SetResponse{
		Prefix:   req.GetPrefix(),
		Response: results,
	}

	log.Infof("Send Set response to client: %v", resp)

	return resp, nil
}

// Overwrites the internal gNMI config.
func (s *Service) OverwriteConfig(jsonConfig []byte) {
	s.mu.Lock()
	defer s.mu.Unlock()

	rootStruct, err := s.model.NewConfigStruct(jsonConfig)
	if err != nil {
		msg := fmt.Sprintf("error in creating config struct from IETF JSON data: %v", err)
		log.Error(msg)
	}
	s.config = rootStruct
}

func (s *Service) Subscribe(stream pbg.GNMI_SubscribeServer) error {
	authorized, err := s.auth.AuthorizeUser(stream.Context())
	if !authorized {
		log.Infof("denied a Subscribe request: %v", err)
		return status.Error(codes.PermissionDenied, fmt.Sprint(err))
	}

	req, err := stream.Recv()

	log.Infof("allowed Subscribe request: %v", req)

	switch {
	case err == io.EOF:
		return nil
	case err != nil:
		return err
	case req.GetSubscribe() == nil:
		return status.Errorf(codes.InvalidArgument, "request must contain a subscription %#v", req)
	}

	if err := s.checkEncodingAndModel(req.GetSubscribe().GetEncoding(), req.GetSubscribe().UseModels); err != nil {
		return status.Error(codes.Unimplemented, err.Error())
	}

	errChan := make(chan error)

	switch req.GetSubscribe().Mode {
	case pbg.SubscriptionList_STREAM:
		go s.subscribeStream(stream, req, errChan)
	case pbg.SubscriptionList_POLL:
		// TODO(dherkel@google.com): Subscribe POLL currently not implemented.
		return status.Errorf(codes.Unimplemented, "unsupported subscribe mode: %s", req.GetSubscribe().Mode)
	case pbg.SubscriptionList_ONCE:
		s.subscribeOnce(stream, req, errChan)
	default:
		return status.Errorf(codes.Unimplemented, "unsupported subscribe mode: %s", req.GetSubscribe().Mode)
	}

	return <-errChan
}

func (s *Service) processSubscribe(req *pbg.SubscribeRequest, respChan chan<- *pbg.SubscribeResponse, errChan chan<- error) {
	log.Debug("process subscribe")

	prefix := req.GetSubscribe().GetPrefix()
	paths := req.GetSubscribe().GetSubscription()
	var notification *pbg.Notification

	s.mu.RLock()
	defer s.mu.RUnlock()

	for _, path := range paths {
		// Get schema node for path from config struct.
		fullPath := path.Path
		if prefix != nil {
			fullPath = gnmiFullPath(prefix, path.Path)
		}
		if fullPath.GetElem() == nil && fullPath.GetElement() != nil {
			errChan <- status.Error(codes.Unimplemented, "deprecated path element type is unsupported")
			return
		}
		node, stat := ygotutils.GetNode(s.model.SchemaTreeRoot, s.config, fullPath)
		if isNil(node) || stat.GetCode() != int32(cpb.Code_OK) {
			errChan <- status.Errorf(codes.NotFound, "path %v not found", fullPath)
			return
		}

		ts := time.Now().UnixNano()

		nodeStruct, ok := node.(ygot.GoStruct)
		// Return leaf node.
		if !ok {
			var val *pbg.TypedValue
			switch kind := reflect.ValueOf(node).Kind(); kind {
			case reflect.Ptr, reflect.Interface:
				var err error
				val, err = value.FromScalar(reflect.ValueOf(node).Elem().Interface())
				if err != nil {
					msg := fmt.Sprintf("leaf node %v does not contain a scalar type value: %v", path, err)
					log.Error(msg)
					errChan <- status.Error(codes.Internal, msg)
					return
				}
			case reflect.Int64:
				enumMap, ok := s.model.EnumData[reflect.TypeOf(node).Name()]
				if !ok {
					errChan <- status.Error(codes.Internal, "not a GoStruct enumeration type")
					return
				}
				val = &pbg.TypedValue{
					Value: &pbg.TypedValue_StringVal{
						StringVal: enumMap[reflect.ValueOf(node).Int()].Name,
					},
				}
			default:
				errChan <- status.Errorf(codes.Internal, "unexpected kind of leaf node type: %v %v", node, kind)
				return
			}

			update := &pbg.Update{Path: path.Path, Val: val}
			notification = &pbg.Notification{
				Timestamp: ts,
				Prefix:    prefix,
				Update:    []*pbg.Update{update},
			}
			continue
		}

		// Return all leaf nodes of the sub-tree.
		if len(req.GetSubscribe().UseModels) != len(s.model.ModelData) && req.GetSubscribe().GetEncoding() != pbg.Encoding_JSON_IETF {
			results, err := ygot.TogNMINotifications(nodeStruct, ts, ygot.GNMINotificationsConfig{UsePathElem: true, PathElemPrefix: fullPath.Elem})
			if err != nil {
				errChan <- status.Errorf(codes.Internal, "error in serializing GoStruct to notifications: %v", err)
				return
			}
			if len(results) != 1 {
				errChan <- status.Errorf(codes.Internal, "ygot.TogNMINotifications() return %d notifications instead of one", len(results))
				return
			}
			notification = results[0]
			continue
		}

		// Return IETF JSON for the sub-tree.
		jsonTree, err := ygot.ConstructIETFJSON(nodeStruct, &ygot.RFC7951JSONConfig{AppendModuleName: true})
		if err != nil {
			msg := fmt.Sprintf("error in constructing IETF JSON tree from requested node: %v", err)
			log.Error(msg)
			errChan <- status.Error(codes.Internal, msg)
			return
		}
		jsonDump, err := json.Marshal(jsonTree)
		if err != nil {
			msg := fmt.Sprintf("error in marshaling IETF JSON tree to bytes: %v", err)
			log.Error(msg)
			errChan <- status.Error(codes.Internal, msg)
			return
		}
		update := &pbg.Update{
			Path: path.Path,
			Val: &pbg.TypedValue{
				Value: &pbg.TypedValue_JsonIetfVal{
					JsonIetfVal: jsonDump,
				},
			},
		}
		notification = &pbg.Notification{
			Timestamp: ts,
			Prefix:    prefix,
			Update:    []*pbg.Update{update},
		}
	}

	resp := &pbg.SubscribeResponse{
		Response: &pbg.SubscribeResponse_Update{
			Update: notification,
		},
	}

	log.Debugf("prepared subscribe response: %v", resp)

	respChan <- resp
}

func (s *Service) subscribeOnce(stream pbg.GNMI_SubscribeServer, req *pbg.SubscribeRequest, errChan chan<- error) {
	log.Infof("serving subscribe ONCE")

	respChan := make(chan *pbg.SubscribeResponse)
	go s.processSubscribe(req, respChan, errChan)

	for {
		select {
		case resp := <-respChan:
			log.Infof("Send Subscribe ONCE response to client: %v", resp)

			err := stream.Send(resp)
			if err != nil {
				errChan <- status.Error(codes.Unimplemented, err.Error())
			}

			return
		}
	}
}

func (s *Service) subscribeStream(stream pbg.GNMI_SubscribeServer, req *pbg.SubscribeRequest, errChan chan<- error) {
	log.Infof("serving subscribe STREAM")

	respChan := make(chan *pbg.SubscribeResponse)

	for {
		select {
		default:
			<-s.ConfigUpdate
			go s.processSubscribe(req, respChan, errChan)

			resp := <-respChan
			log.Infof("Send Subscribe STREAM response to client: %v", resp)

			err := stream.Send(resp)
			if err != nil {
				errChan <- status.Error(codes.Unimplemented, err.Error())
				return
			}
		}
	}
}

func (s *Service) Reboot(ctx context.Context, req *pbs.RebootRequest) (*pbs.RebootResponse, error) {
	authorized, err := s.auth.AuthorizeUser(ctx)
	if !authorized {
		log.Infof("denied a Reboot request: %v", err)
		return nil, status.Error(codes.PermissionDenied, fmt.Sprint(err))
	}
	log.Infof("allowed a Reboot request")

	s.mu.Lock()

	go s.ch.CallbackReboot()

	defer s.mu.Unlock()

	resp := &pbs.RebootResponse{}

	log.Infof("Send Reboot response to client: %v", resp)

	return resp, nil
}

func (s *Service) RebootStatus(ctx context.Context, req *pbs.RebootStatusRequest) (*pbs.RebootStatusResponse, error) {
	return nil, status.Error(codes.Unimplemented, "RebootStatus is not implemented.")
}

func (s *Service) CancelReboot(ctx context.Context, req *pbs.CancelRebootRequest) (*pbs.CancelRebootResponse, error) {
	return nil, status.Error(codes.Unimplemented, "CancelReboot is not implemented.")
}

func (s *Service) Ping(req *pbs.PingRequest, stream pbs.System_PingServer) error {
	return status.Error(codes.Unimplemented, "Ping is not implemented.")
}

func (s *Service) Traceroute(req *pbs.TracerouteRequest, stream pbs.System_TracerouteServer) error {
	return status.Error(codes.Unimplemented, "Traceroute is not implemented.")
}

func (s *Service) Time(ctx context.Context, req *pbs.TimeRequest) (*pbs.TimeResponse, error) {
	return nil, status.Error(codes.Unimplemented, "Time is not implemented.")
}

func (s *Service) SetPackage(stream pbs.System_SetPackageServer) error {
	return status.Error(codes.Unimplemented, "SetPackage is not implemented.")
}

func (s *Service) SwitchControlProcessor(ctx context.Context, req *pbs.SwitchControlProcessorRequest) (*pbs.SwitchControlProcessorResponse, error) {
	return nil, status.Error(codes.Unimplemented, "SwitchControlProcessor is not implemented.")
}

func (s *Service) Rotate(stream pbc.CertificateManagement_RotateServer) error {
	authorized, err := s.auth.AuthorizeUser(stream.Context())
	if !authorized {
		log.Infof("denied a Rotate request: %v", err)
		return status.Error(codes.PermissionDenied, fmt.Sprint(err))
	}

	req, err := stream.Recv()

	log.Infof("allowed a Rotate request: %v", req)

	switch {
	case err == io.EOF:
		return nil
	case err != nil:
		return err
	}

	errChan := make(chan error)

	return <-errChan
}

func (s *Service) GetCertificates(ctx context.Context, req *pbc.GetCertificatesRequest) (*pbc.GetCertificatesResponse, error) {
	authorized, err := s.auth.AuthorizeUser(ctx)
	if !authorized {
		log.Infof("denied a GetCertificates request: %v", err)
		return nil, status.Error(codes.PermissionDenied, fmt.Sprint(err))
	}
	log.Infof("allowed a GetCertificates request")

	resp := &pbc.GetCertificatesResponse{
		CertificateInfo: s.certs.CertInfo,
	}

	log.Infof("Send GetCertificates response to client: %v", resp)

	return resp, nil
}

func (s *Service) Install(stream pbc.CertificateManagement_InstallServer) error {
	return status.Error(codes.Unimplemented, "Install is not implemented.")
}

func (s *Service) RevokeCertificates(ctx context.Context, req *pbc.RevokeCertificatesRequest) (*pbc.RevokeCertificatesResponse, error) {
	return nil, status.Error(codes.Unimplemented, "RevokeCertificates is not implemented.")
}

func (s *Service) CanGenerateCSR(ctx context.Context, req *pbc.CanGenerateCSRRequest) (*pbc.CanGenerateCSRResponse, error) {
	return nil, status.Error(codes.Unimplemented, "CanGenerateCSR is not implemented.")
}

func (s *Service) prepareService(certificates []tls.Certificate, certPool *x509.CertPool) {
	opts := []grpc.ServerOption{grpc.Creds(credentials.NewTLS(&tls.Config{
		ClientAuth:   tls.RequireAndVerifyClientCert,
		Certificates: certificates,
		ClientCAs:    certPool,
	}))}

	s.g = grpc.NewServer(opts...)

	s.registerService()
}

func (s *Service) registerService() {
	pbg.RegisterGNMIServer(s.g, s)
	pbs.RegisterSystemServer(s.g, s)
	pbc.RegisterCertificateManagementServer(s.g, s)
	reflection.Register(s.g)
}

func (s *Service) StartService() {
	log.Info("Start gNXI Service")

	s.prepareService(s.certs.TLSCertificates, s.certs.CertPool)

	var err error

	log.Infof("Starting to listen")
	s.socket, err = net.Listen(gnxiProtocol, fmt.Sprintf(":%s", gnxiPort))
	if err != nil {
		log.Fatalf("Failed to listen: %v", err)
	}

	log.Info("Starting to serve gNXI")
	if err := s.g.Serve(s.socket); err != nil {
		log.Fatalf("Failed to serve: %v", err)
	}
}

func (s *Service) StopService() {
	log.Info("Stop gNXI Service")
	s.g.Stop()
}
