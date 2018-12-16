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
package gnmi

import (
	"bytes"
	"compress/gzip"
	"encoding/json"
	"fmt"
	"github.com/google/gnxi/utils/credentials"
	"github.com/openconfig/gnmi/cache"
	"github.com/openconfig/gnmi/client"
	"github.com/openconfig/gnmi/coalesce"
	"github.com/openconfig/gnmi/ctree"
	"github.com/openconfig/gnmi/path"
	"google.golang.org/grpc/peer"
	"io"
	"io/ioutil"
	"ovs-gnxi/shared/logging"
	"reflect"
	"strconv"
	"sync"
	"time"

	"github.com/openconfig/gnmi/match"
	"golang.org/x/net/context"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/golang/protobuf/proto"
	"github.com/openconfig/gnmi/value"
	"github.com/openconfig/ygot/experimental/ygotutils"
	"github.com/openconfig/ygot/ygot"

	dpb "github.com/golang/protobuf/protoc-gen-go/descriptor"
	pb "github.com/openconfig/gnmi/proto/gnmi"
	cpb "google.golang.org/genproto/googleapis/rpc/code"
)

var (
	pbRootPath         = &pb.Path{}
	supportedEncodings = []pb.Encoding{pb.Encoding_JSON, pb.Encoding_JSON_IETF}
)

var log = logging.New("ovs-gnxi")

// ConfigCallback is the signature of the function to apply a validated config to the physical device.
type ConfigCallback func(ygot.ValidatedGoStruct) error

// Service struct maintains the data structure for device config and implements the gnmi interface. It supports Capabilities, Get, Set and Subscribe APIs.
type Service struct {
	model    *Model
	callback ConfigCallback
	config   ygot.ValidatedGoStruct
	mu       sync.RWMutex // mu is the RW lock to protect the access to config

	cache *cache.Cache // The cache queries are performed against.
	m     *match.Match // Structure to match updates against active subscriptions.
	// subscribeSlots is a channel of size SubscriptionLimit to restrict how many
	// queries are in flight.
	subscribeSlots chan struct{}
	timeout        time.Duration
}

// NewService creates an instance of Service with given json config.
func NewService(model *Model, config []byte, callback ConfigCallback, cache *cache.Cache, subscriptionLimit int) (*Service, error) {
	rootStruct, err := model.NewConfigStruct(config)

	if err != nil {
		return nil, err
	}
	s := &Service{
		model:    model,
		config:   rootStruct,
		callback: callback,
		cache:    cache,
	}

	s.cache.SetClient(func(l *ctree.Leaf) {
		log.Error(l.Value())
	})

	if subscriptionLimit > 0 {
		s.subscribeSlots = make(chan struct{}, subscriptionLimit)
	}

	if config != nil && s.callback != nil {
		if err := s.callback(rootStruct); err != nil {
			return nil, err
		}
	}
	return s, nil
}

// checkEncodingAndModel checks whether encoding and models are supported by the server. Return error if anything is unsupported.
func (s *Service) checkEncodingAndModel(encoding pb.Encoding, models []*pb.ModelData) error {
	hasSupportedEncoding := false
	for _, supportedEncoding := range supportedEncodings {
		if encoding == supportedEncoding {
			hasSupportedEncoding = true
			break
		}
	}
	if !hasSupportedEncoding {
		return fmt.Errorf("unsupported encoding: %s", pb.Encoding_name[int32(encoding)])
	}
	for _, m := range models {
		isSupported := false
		for _, supportedModel := range s.model.modelData {
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
func (s *Service) doDelete(jsonTree map[string]interface{}, prefix, path *pb.Path) (*pb.UpdateResult, error) {
	// Update json tree of the device config
	var curNode interface{} = jsonTree
	pathDeleted := false
	fullPath := gnmiFullPath(prefix, path)
	schema := s.model.schemaTreeRoot
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

		if curNode, schema = getChildNode(node, schema, elem, false); curNode == nil {
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
		if s.callback != nil {
			if applyErr := s.callback(newConfig); applyErr != nil {
				if rollbackErr := s.callback(s.config); rollbackErr != nil {
					return nil, status.Errorf(codes.Internal, "error in rollback the failed operation (%v): %v", applyErr, rollbackErr)
				}
				return nil, status.Errorf(codes.Aborted, "error in applying operation to device: %v", applyErr)
			}
		}
	}
	return &pb.UpdateResult{
		Path: path,
		Op:   pb.UpdateResult_DELETE,
	}, nil
}

// doReplaceOrUpdate validates the replace or update operation to be applied to
// the device, modifies the json tree of the config struct, then calls the
// callback function to apply the operation to the device hardware.
func (s *Service) doReplaceOrUpdate(jsonTree map[string]interface{}, op pb.UpdateResult_Operation, prefix, path *pb.Path, val *pb.TypedValue) (*pb.UpdateResult, error) {
	// Validate the operation.
	fullPath := gnmiFullPath(prefix, path)
	emptyNode, stat := ygotutils.NewNode(s.model.structRootType, fullPath)
	if stat.GetCode() != int32(cpb.Code_OK) {
		return nil, status.Errorf(codes.NotFound, "path %v is not found in the config structure: %v", fullPath, stat)
	}
	var nodeVal interface{}
	nodeStruct, ok := emptyNode.(ygot.ValidatedGoStruct)
	if ok {
		if err := s.model.jsonUnmarshaler(val.GetJsonIetfVal(), nodeStruct); err != nil {
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
	schema := s.model.schemaTreeRoot
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

			if curNode, schema = getChildNode(node, schema, elem, true); curNode == nil {
				return nil, status.Errorf(codes.NotFound, "path elem not found: %v", elem)
			}
		case []interface{}:
			return nil, status.Errorf(codes.NotFound, "uncompatible path elem: %v", elem)
		default:
			return nil, status.Errorf(codes.Internal, "wrong node type: %T", curNode)
		}
	}
	if reflect.DeepEqual(fullPath, pbRootPath) { // Replace/Update root.
		if op == pb.UpdateResult_UPDATE {
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
	if s.callback != nil {
		if applyErr := s.callback(newConfig); applyErr != nil {
			if rollbackErr := s.callback(s.config); rollbackErr != nil {
				return nil, status.Errorf(codes.Internal, "error in rollback the failed operation (%v): %v", applyErr, rollbackErr)
			}
			return nil, status.Errorf(codes.Aborted, "error in applying operation to device: %v", applyErr)
		}
	}
	return &pb.UpdateResult{
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
	gzB, _ := (&pb.Update{}).Descriptor()
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
	ver, err := proto.GetExtension(desc.Options, pb.E_GnmiService)
	if err != nil {
		return nil, fmt.Errorf("error in getting version from proto extension: %v", err)
	}
	return ver.(*string), nil
}

// deleteKeyedListEntry deletes the keyed list entry from node that matches the
// path elem. If the entry is the only one in keyed list, deletes the entire
// list. If the entry is found and deleted, the function returns true. If it is
// not found, the function returns false.
func deleteKeyedListEntry(node map[string]interface{}, elem *pb.PathElem) bool {
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
func gnmiFullPath(prefix, path *pb.Path) *pb.Path {
	fullPath := &pb.Path{Origin: path.Origin}
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
func setPathWithAttribute(op pb.UpdateResult_Operation, curNode map[string]interface{}, pathElem *pb.PathElem, nodeVal interface{}) error {
	nodeValAsTree, ok := nodeVal.(map[string]interface{})
	if !ok {
		return status.Errorf(codes.InvalidArgument, "expect nodeVal is a json node of map[string]interface{}, received %T", nodeVal)
	}
	m := getKeyedListEntry(curNode, pathElem, true)
	if m == nil {
		return status.Errorf(codes.NotFound, "path elem not found: %v", pathElem)
	}
	if op == pb.UpdateResult_REPLACE {
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
func setPathWithoutAttribute(op pb.UpdateResult_Operation, curNode map[string]interface{}, pathElem *pb.PathElem, nodeVal interface{}) error {
	target, hasElem := curNode[pathElem.Name]
	nodeValAsTree, nodeValIsTree := nodeVal.(map[string]interface{})
	if op == pb.UpdateResult_REPLACE || !hasElem || !nodeValIsTree {
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
func (s *Service) Capabilities(ctx context.Context, req *pb.CapabilityRequest) (*pb.CapabilityResponse, error) {
	ver, err := getGNMIServiceVersion()
	if err != nil {
		return nil, status.Errorf(codes.Internal, "error in getting gnmi service version: %v", err)
	}
	return &pb.CapabilityResponse{
		SupportedModels:    s.model.modelData,
		SupportedEncodings: supportedEncodings,
		GNMIVersion:        *ver,
	}, nil
}

// Get implements the Get RPC in gNMI spec and provides user auth.
func (s *Service) Get(ctx context.Context, req *pb.GetRequest) (*pb.GetResponse, error) {
	msg, ok := credentials.AuthorizeUser(ctx)
	if !ok {
		log.Infof("denied a Get request: %v", msg)
		return nil, status.Error(codes.PermissionDenied, msg)
	}
	log.Infof("allowed a Get request: %v", msg)

	if req.GetType() != pb.GetRequest_ALL {
		return nil, status.Errorf(codes.Unimplemented, "unsupported request type: %s", pb.GetRequest_DataType_name[int32(req.GetType())])
	}
	if err := s.checkEncodingAndModel(req.GetEncoding(), req.GetUseModels()); err != nil {
		return nil, status.Error(codes.Unimplemented, err.Error())
	}

	prefix := req.GetPrefix()
	paths := req.GetPath()
	notifications := make([]*pb.Notification, len(paths))

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
		node, stat := ygotutils.GetNode(s.model.schemaTreeRoot, s.config, fullPath)
		if isNil(node) || stat.GetCode() != int32(cpb.Code_OK) {
			return nil, status.Errorf(codes.NotFound, "path %v not found", fullPath)
		}

		ts := time.Now().UnixNano()

		nodeStruct, ok := node.(ygot.GoStruct)
		// Return leaf node.
		if !ok {
			var val *pb.TypedValue
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
				enumMap, ok := s.model.enumData[reflect.TypeOf(node).Name()]
				if !ok {
					return nil, status.Error(codes.Internal, "not a GoStruct enumeration type")
				}
				val = &pb.TypedValue{
					Value: &pb.TypedValue_StringVal{
						StringVal: enumMap[reflect.ValueOf(node).Int()].Name,
					},
				}
			default:
				return nil, status.Errorf(codes.Internal, "unexpected kind of leaf node type: %v %v", node, kind)
			}

			update := &pb.Update{Path: path, Val: val}
			notifications[i] = &pb.Notification{
				Timestamp: ts,
				Prefix:    prefix,
				Update:    []*pb.Update{update},
			}
			continue
		}

		// Return all leaf nodes of the sub-tree.
		if len(req.GetUseModels()) != len(s.model.modelData) && req.GetEncoding() != pb.Encoding_JSON_IETF {
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
		update := &pb.Update{
			Path: path,
			Val: &pb.TypedValue{
				Value: &pb.TypedValue_JsonIetfVal{
					JsonIetfVal: jsonDump,
				},
			},
		}
		notifications[i] = &pb.Notification{
			Timestamp: ts,
			Prefix:    prefix,
			Update:    []*pb.Update{update},
		}
	}

	return &pb.GetResponse{Notification: notifications}, nil
}

// Set implements the Set RPC in gNMI spec and provides user auth.
func (s *Service) Set(ctx context.Context, req *pb.SetRequest) (*pb.SetResponse, error) {
	msg, ok := credentials.AuthorizeUser(ctx)
	if !ok {
		log.Infof("denied a Set request: %v", msg)
		return nil, status.Error(codes.PermissionDenied, msg)
	}
	log.Infof("allowed a Set request: %v", msg)

	s.mu.Lock()
	defer s.mu.Unlock()

	jsonTree, err := ygot.ConstructIETFJSON(s.config, &ygot.RFC7951JSONConfig{})
	if err != nil {
		msg := fmt.Sprintf("error in constructing IETF JSON tree from config struct: %v", err)
		log.Error(msg)
		return nil, status.Error(codes.Internal, msg)
	}

	prefix := req.GetPrefix()
	var results []*pb.UpdateResult

	for _, path := range req.GetDelete() {
		res, grpcStatusError := s.doDelete(jsonTree, prefix, path)
		if grpcStatusError != nil {
			return nil, grpcStatusError
		}
		results = append(results, res)
	}
	for _, upd := range req.GetReplace() {
		res, grpcStatusError := s.doReplaceOrUpdate(jsonTree, pb.UpdateResult_REPLACE, prefix, upd.GetPath(), upd.GetVal())
		if grpcStatusError != nil {
			return nil, grpcStatusError
		}
		results = append(results, res)
	}
	for _, upd := range req.GetUpdate() {
		res, grpcStatusError := s.doReplaceOrUpdate(jsonTree, pb.UpdateResult_UPDATE, prefix, upd.GetPath(), upd.GetVal())
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
	return &pb.SetResponse{
		Prefix:   req.GetPrefix(),
		Response: results,
	}, nil
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

func (s *Service) OverwriteCallback(callback ConfigCallback) {
	s.mu.Lock()
	defer s.mu.Unlock()

	s.callback = callback
}

// Update passes a streaming update to registered clients.
func (s *Service) Update(n *ctree.Leaf) {
	switch v := n.Value().(type) {
	case client.Delete:
		s.m.Update(n, v.Path)
	case client.Update:
		s.m.Update(n, v.Path)
	case *pb.Notification:
		p := path.ToStrings(v.Prefix, true)
		if len(v.Update) > 0 {
			p = append(p, path.ToStrings(v.Update[0].Path, false)...)
		} else if len(v.Delete) > 0 {
			p = append(p, path.ToStrings(v.Delete[0], false)...)
		}
		// If neither update nor delete notification exists,
		// just go with the path in the prefix
		s.m.Update(n, p)
	default:
		log.Errorf("update is not a known type; type is %T", v)
	}
}

type streamClient struct {
	target string
	sr     *pb.SubscribeRequest
	queue  *coalesce.Queue
	stream pb.GNMI_SubscribeServer
	errC   chan<- error
}

// Subscribe implements the Subscribe RPC in gNMI spec.
func (s *Service) Subscribe(stream pb.GNMI_SubscribeServer) error {
	c := streamClient{stream: stream}
	var err error
	c.sr, err = stream.Recv()

	switch {
	case err == io.EOF:
		return nil
	case err != nil:
		return err
	case c.sr.GetSubscribe() == nil:
		return status.Errorf(codes.InvalidArgument, "request must contain a subscription %#v", c.sr)
	case c.sr.GetSubscribe().GetPrefix() == nil:
		return status.Errorf(codes.InvalidArgument, "request must contain a prefix %#v", c.sr)
	case c.sr.GetSubscribe().GetPrefix().GetTarget() == "":
		return status.Error(codes.InvalidArgument, "missing target")
	}

	c.target = c.sr.GetSubscribe().GetPrefix().GetTarget()
	if !s.cache.HasTarget(c.target) {
		return status.Errorf(codes.NotFound, "no such target: %q", c.target)
	}
	peer, _ := peer.FromContext(stream.Context())
	mode := c.sr.GetSubscribe().Mode

	log.Infof("peer: %v target: %q subscription: %s", peer.Addr, c.target, c.sr)
	defer log.Infof("peer: %v target %q subscription: end: %q", peer.Addr, c.target, c.sr)

	c.queue = coalesce.NewQueue()
	defer c.queue.Close()

	// This error channel is buffered to accept errors from all goroutines spawned
	// for this RPC. Only the first is ever read and returned causing the RPC to
	// terminate.
	errC := make(chan error, 3)
	c.errC = errC

	switch mode {
	case pb.SubscriptionList_ONCE:
		go func() {
			s.processSubscription(&c)
			c.queue.Close()
		}()
	case pb.SubscriptionList_POLL:
		go func() {
			log.Info("SUBSCRIPTION POLL")
		}()
	case pb.SubscriptionList_STREAM:
		go func() {
			log.Info("SUBSCRIPTION STREAM")
		}()
	default:
		return status.Errorf(codes.InvalidArgument, "Subscription mode %v not recognized", mode)
	}

	return <-errC
}

type syncMarker struct{}

// processSubscription walks the cache tree and inserts all of the matching
// nodes into the coalesce queue followed by a subscriptionSync response.
func (s *Service) processSubscription(c *streamClient) {
	var err error
	// Close the cache client queue on error.
	defer func() {
		if err != nil {
			log.Error(err)
			c.queue.Close()
			c.errC <- err
		}
	}()
	if s.subscribeSlots != nil {
		select {
		// Register a subscription in the channel, which will block if SubscriptionLimit queries
		// are already in flight.
		case s.subscribeSlots <- struct{}{}:
		default:
			log.Infof("subscription %s delayed waiting for 1 of %d subscription slots.", c.sr, len(s.subscribeSlots))
			s.subscribeSlots <- struct{}{}
			log.Infof("subscription %s resumed", c.sr)
		}
		// Remove subscription from the channel upon completion.
		defer func() {
			<-s.subscribeSlots
		}()
	}
	if !c.sr.GetSubscribe().GetUpdatesOnly() {
		// remove the target name from the index string
		prefix := path.ToStrings(c.sr.GetSubscribe().Prefix, true)[1:]
		for _, subscription := range c.sr.GetSubscribe().Subscription {
			path := append(prefix, path.ToStrings(subscription.Path, false)...)
			s.cache.Query(c.target, path, func(_ []string, l *ctree.Leaf, _ interface{}) {
				// Stop processing query results on error.
				if err != nil {
					return
				}
				_, err = c.queue.Insert(l)
			})
			if err != nil {
				return
			}
		}
	}

	_, err = c.queue.Insert(syncMarker{})
}

/*
// addSubscription registers all subscriptions for this client for update matching.
func addSubscription(m *match.Match, s *pb.SubscriptionList, c *matchClient) (remove func()) {
	var removes []func()
	prefix := path.ToStrings(s.Prefix, true)
	for _, p := range s.Subscription {
		if p.Path == nil {
			continue
		}
		// TODO(yusufsn) : Origin field in the Path may need to be included
		path := append(prefix, path.ToStrings(p.Path, false)...)
		removes = append(removes, m.AddQuery(path, c))
	}
	return func() {
		for _, remove := range removes {
			remove()
		}
	}
}*/

/*
// Subscribe implements the Subscribe RPC in gNMI spec.
func (s *Service) Subscribe(stream pb.GNMI_SubscribeServer) error {

	c := streamClient{stream: stream}
	var err error
	c.sr, err = stream.Recv()

	switch {
	case err == io.EOF:
		return nil
	case err != nil:
		return err
	case c.sr.GetSubscribe() == nil:
		return status.Errorf(codes.InvalidArgument, "request must contain a subscription %#v", c.sr)
	case c.sr.GetSubscribe().GetPrefix() == nil:
		return status.Errorf(codes.InvalidArgument, "request must contain a prefix %#v", c.sr)
	case c.sr.GetSubscribe().GetPrefix().GetTarget() == "":
		return status.Error(codes.InvalidArgument, "missing target")
	}

	c.target = c.sr.GetSubscribe().GetPrefix().GetTarget()
	log.Error(c.sr.GetSubscribe())
	if !s.cache.HasTarget(c.target) {
		return status.Errorf(codes.NotFound, "no such target: %q", c.target)
	}
	peer, _ := peer.FromContext(stream.Context())
	mode := c.sr.GetSubscribe().Mode

	log.Infof("peer: %v target: %q subscription: %s", peer.Addr, c.target, c.sr)
	defer log.Infof("peer: %v target %q subscription: end: %q", peer.Addr, c.target, c.sr)

	c.queue = coalesce.NewQueue()
	defer c.queue.Close()

	// This error channel is buffered to accept errors from all goroutines spawned
	// for this RPC. Only the first is ever read and returned causing the RPC to
	// terminate.
	errC := make(chan error, 3)
	c.errC = errC

	switch mode {
	case pb.SubscriptionList_ONCE:
		go func() {
			s.processSubscription(&c)
			c.queue.Close()
		}()
	case pb.SubscriptionList_POLL:
		go s.processPollingSubscription(&c)
	case pb.SubscriptionList_STREAM:
		if c.sr.GetSubscribe().GetUpdatesOnly() {
			result, err := c.queue.Insert(syncMarker{})
			if err != nil {
				return status.Errorf(codes.Unknown, err.Error())
			}
			log.Debug(result)
		}
		remove := addSubscription(s.m, c.sr.GetSubscribe(), &matchClient{q: c.queue})
		defer remove()
		if !c.sr.GetSubscribe().GetUpdatesOnly() {
			go s.processSubscription(&c)
		}
	default:
		return status.Errorf(codes.InvalidArgument, "Subscription mode %v not recognized", mode)
	}

	go s.sendStreamingResults(&c)

	return <-errC
}

type resp struct {
	stream pb.GNMI_SubscribeServer
	n      *ctree.Leaf
	dup    uint32
	t      *time.Timer // Timer used to timout the subscription.
}

// sendSubscribeResponse populates and sends a single response returned on
// the Subscription RPC output stream. Streaming queries send responses for the
// initial walk of the results as well as streamed updates and use a queue to
// ensure order.
func (s *Service) sendSubscribeResponse(r *resp) error {
	log.Error("TEEEEEEEST222222222!")

	notification, err := MakeSubscribeResponse(r.n.Value(), r.dup)
	if err != nil {
		return status.Errorf(codes.Unknown, err.Error())
	}
	// Start the timeout before attempting to send.
	r.t.Reset(s.timeout)
	// Clear the timeout upon sending.
	defer r.t.Stop()
	return r.stream.Send(notification)
}

// subscribeSync is a response indicating that a Subscribe RPC has successfully
// returned all matching nodes once for ONCE and POLLING queries and at least
// once for STREAMING queries.
var subscribeSync = &pb.SubscribeResponse{Response: &pb.SubscribeResponse_SyncResponse{true}}

type syncMarker struct{}

// cacheClient implements match.Client interface.
type matchClient struct {
	q   *coalesce.Queue
	err error
}

// Update implements the match.Client Update interface for coalesce.Queue.
func (c matchClient) Update(n interface{}) {
	// Stop processing updates on error.
	if c.err != nil {
		return
	}
	_, c.err = c.q.Insert(n)
}

type streamClient struct {
	target string
	sr     *pb.SubscribeRequest
	queue  *coalesce.Queue
	stream pb.GNMI_SubscribeServer
	errC   chan<- error
}

// processSubscription walks the cache tree and inserts all of the matching
// nodes into the coalesce queue followed by a subscriptionSync response.
func (s *Service) processSubscription(c *streamClient) {
	var err error
	// Close the cache client queue on error.
	defer func() {
		if err != nil {
			log.Error(err)
			c.queue.Close()
			c.errC <- err
		}
	}()
	if s.subscribeSlots != nil {
		select {
		// Register a subscription in the channel, which will block if SubscriptionLimit queries
		// are already in flight.
		case s.subscribeSlots <- struct{}{}:
		default:
			log.Infof("subscription %s delayed waiting for 1 of %d subscription slots.", c.sr, len(s.subscribeSlots))
			s.subscribeSlots <- struct{}{}
			log.Infof("subscription %s resumed", c.sr)
		}
		// Remove subscription from the channel upon completion.
		defer func() {
			<-s.subscribeSlots
		}()
	}
	if !c.sr.GetSubscribe().GetUpdatesOnly() {
		// remove the target name from the index string
		prefix := path.ToStrings(c.sr.GetSubscribe().Prefix, true)[1:]
		for _, subscription := range c.sr.GetSubscribe().Subscription {
			path := append(prefix, path.ToStrings(subscription.Path, false)...)
			s.cache.Query(c.target, path, func(_ []string, l *ctree.Leaf, _ interface{}) {
				// Stop processing query results on error.
				if err != nil {
					return
				}
				_, err = c.queue.Insert(l)
			})
			if err != nil {
				return
			}
		}
	}

	_, err = c.queue.Insert(syncMarker{})
}

// processPollingSubscription handles the POLL mode Subscription RPC.
func (s *Service) processPollingSubscription(c *streamClient) {
	s.processSubscription(c)
	log.Infof("polling subscription: first complete response: %q", c.sr)
	for {
		if c.queue.IsClosed() {
			return
		}
		// Subsequent receives are only triggers to poll again. The contents of the
		// request are completely ignored.
		_, err := c.stream.Recv()
		if err == io.EOF {
			return
		}
		if err != nil {
			log.Error(err)
			c.errC <- err
			return
		}
		log.Infof("polling subscription: repoll: %q", c.sr)
		s.processSubscription(c)
		log.Infof("polling subscription: repoll complete: %q", c.sr)
	}
}

// sendStreamingResults forwards all streaming updates to a given streaming
// Subscription RPC client.
func (s *Service) sendStreamingResults(c *streamClient) {
	ctx := c.stream.Context()
	peer, _ := peer.FromContext(ctx)
	t := time.NewTimer(s.timeout)

	log.Error("TEEEEEEEST1111111111!")

	// Make sure the timer doesn't expire waiting for a value to send, only
	// waiting to send.
	t.Stop()
	done := make(chan struct{})
	defer close(done)
	// If a send doesn't occur within the timeout, close the stream.
	go func() {
		select {
		case <-t.C:
			err := errors.New("subscription timed out while sending")
			c.errC <- err
			log.Errorf("%v : %v", peer, err)
		case <-done:
		}
	}()
	for {
		log.Error("TEEEEEEEST999999999!")

		item, dup, err := c.queue.Next(ctx)
		if coalesce.IsClosedQueue(err) {
			c.errC <- nil
			return
		}
		if err != nil {
			c.errC <- err
			return
		}

		// s.processSubscription will send a sync marker, handle it separately.
		if _, ok := item.(syncMarker); ok {
			if err = c.stream.Send(subscribeSync); err != nil {
				break
			}
			continue
		}

		log.Error("TEEEEEEEST8888888888!")

		n, ok := item.(*ctree.Leaf)
		if !ok || n == nil {
			c.errC <- status.Errorf(codes.Internal, "invalid cache node: %#v", item)
			return
		}
		if err = s.sendSubscribeResponse(&resp{
			stream: c.stream,
			n:      n,
			dup:    dup,
			t:      t,
		}); err != nil {
			c.errC <- err
			return
		}
		// If the only target being subscribed was deleted, stop streaming.
		if cache.IsTargetDelete(n) && c.target != "*" {
			log.Infof("Target %q was deleted. Closing stream.", c.target)
			c.errC <- nil
			return
		}
	}
}

// MakeSubscribeResponse produces a gnmi_proto.SubscribeResponse from either
// client.Notification or gnmi_proto.Notification
//
// This function modifies the message to set the duplicate count if it is
// greater than 0. The funciton clones the gnmi notification if the duplicate count needs to be set.
// You have to be working on a cloned message if you need to modify the message in any way.
func MakeSubscribeResponse(n interface{}, dup uint32) (*pb.SubscribeResponse, error) {
	var notification *pb.Notification
	switch cache.Type {
	case cache.GnmiNoti:
		var ok bool
		notification, ok = n.(*pb.Notification)
		if !ok {
			return nil, status.Errorf(codes.Internal, "invalid notification type: %#v", n)
		}

		// There may be multiple updates in a notification. Since duplicate count is just
		// an indicator that coalescion is happening, not a critical data, just the first
		// update is set with duplicate count to be on the side of efficiency.
		// Only attempt to set the duplicate count if it is greater than 0. The default
		// value in the message is already 0.
		if dup > 0 && len(notification.Update) > 0 {
			// We need a copy of the cached notification before writing a client specific
			// duplicate count as the notification is shared across all clients.
			notification = proto.Clone(notification).(*pb.Notification)
			notification.Update[0].Duplicates = dup
		}
	case cache.ClientLeaf:
		notification = &pb.Notification{}
		switch v := n.(type) {
		case client.Delete:
			notification.Delete = []*pb.Path{{Element: v.Path}}
			notification.Timestamp = v.TS.UnixNano()
		case client.Update:
			typedVal, err := value.FromScalar(v.Val)
			if err != nil {
				return nil, err
			}
			notification.Update = []*pb.Update{{Path: &pb.Path{Element: v.Path}, Val: typedVal, Duplicates: dup}}
			notification.Timestamp = v.TS.UnixNano()
		}
	}
	response := &pb.SubscribeResponse{
		Response: &pb.SubscribeResponse_Update{
			Update: notification,
		},
	}

	log.Error(response)

	return response, nil
}
*/
