package ovs

import (
	"errors"
	"fmt"
	"github.com/socketplane/libovsdb"
	"ovs-gnxi/target/logging"
)

const (
	DefaultDatabase = "Open_vSwitch"
	SystemTable     = "Open_vSwitch"
	ControllerTable = "Controller"
	InterfaceTable  = "Interface"
)

var log = logging.New("ovs-gnxi")

var quit chan bool
var update chan *libovsdb.TableUpdates

type Client struct {
	Address    string
	Protocol   string
	Port       string
	Connection *libovsdb.OvsdbClient
	Database   string
	Notifier   *Notifier
	Config     *Config
}

func (o *Client) String() string {
	return fmt.Sprintf("OVSClient(Address: \"%v\", Protocol: \"%v\", Port: \"%v\")", o.Address, o.Protocol, o.Port)
}

func NewClient(address, protocol, port, privateKeyPath, publicKeyPath, caPath string) (*Client, error) {
	var err error

	o := Client{Address: address, Protocol: protocol, Port: port, Database: DefaultDatabase, Config: NewConfig()}

	o.Connection, err = libovsdb.ConnectUsingProtocolWithTLS(o.Protocol, fmt.Sprintf("%v:%v", o.Address, o.Port), privateKeyPath, publicKeyPath, caPath)
	if err != nil {
		log.Fatalf("failed to dial: %v", err)
	}

	return &o, nil
}

func (o *Client) StartMonitorAll() {
	err := o.MonitorAll()
	if err != nil {
		log.Error(err)
	}

	o.Config.DumpRawCache()
}

func (o *Client) SetOpenFlowController(controller *OpenFlowController) error {
	condition := libovsdb.NewCondition("_uuid", "==", libovsdb.UUID{GoUUID: controller.uuid})

	row := make(map[string]interface{})
	row["target"] = fmt.Sprintf("%v:%v:%v", controller.Target.Protocol, controller.Target.Address, controller.Target.Port)

	updateOp := libovsdb.Operation{
		Op:    "update",
		Table: ControllerTable,
		Where: []interface{}{condition},
		Row:   row,
	}

	log.Info(updateOp)

	operations := []libovsdb.Operation{updateOp}
	reply, _ := o.Connection.Transact(o.Database, operations...)

	if len(reply) < len(operations) {
		fmt.Println("Number of Replies should be at least equal to number of Operations")
	}
	ok := true
	for i, o := range reply {
		if o.Error != "" && i < len(operations) {
			fmt.Println("Transaction Failed due to an error :", o.Error, " details:", o.Details, " in ", operations[i])
			ok = false
		} else if o.Error != "" {
			fmt.Println("Transaction Failed due to an error :", o.Error)
			ok = false
		}
	}
	if ok {
		return nil
	}

	return errors.New("unable to set openflow controller address")
}

type Notifier struct {
	config *Config
}

func (n Notifier) Update(context interface{}, tableUpdates libovsdb.TableUpdates) {
	n.config.SyncCache(&tableUpdates)
	update <- &tableUpdates
}
func (n Notifier) Locked([]interface{}) {
}
func (n Notifier) Stolen([]interface{}) {
}
func (n Notifier) Echo([]interface{}) {
}
func (n Notifier) Disconnected(client *libovsdb.OvsdbClient) {
}

func receivedMonitorUpdate() {
	for {
		select {
		case currUpdate := <-update:
			for tableName, tableUpdate := range currUpdate.Updates {
				log.Debugf("Received Table update for \"%v\" with content: %v", tableName, tableUpdate)
			}
		}
	}
}

func (o *Client) MonitorAll() error {
	quit = make(chan bool)
	update = make(chan *libovsdb.TableUpdates)

	request := libovsdb.MonitorRequest{
		Select: libovsdb.MonitorSelect{
			Initial: true,
			Insert:  true,
			Delete:  true,
			Modify:  true,
		},
	}

	requests := make(map[string]libovsdb.MonitorRequest)
	requests[SystemTable] = request
	requests[ControllerTable] = request
	requests[InterfaceTable] = request

	o.Notifier = &Notifier{config: o.Config}
	o.Connection.Register(o.Notifier)

	initial, err := o.Connection.Monitor(DefaultDatabase, "", requests)
	if err != nil {
		log.Error(err)
		return err
	}

	o.Config.InitializeCache(initial)

	go receivedMonitorUpdate()
	<-quit

	return nil
}
