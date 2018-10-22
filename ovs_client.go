package main

import "fmt"

type OVSClient struct {
	Address string
	Port    string
	Config  []byte
}

func (o *OVSClient) String() string {
	return fmt.Sprintf("OVSClient(Address: \"%v\", Port: \"%v\")", o.Address, o.Port)
}

func NewOVSClient(address, port string) (*OVSClient, error) {
	o := OVSClient{Address: address, Port: port}
	return &o, nil
}

func (o *OVSClient) ClientCallback() {

}

func (o *OVSClient) GetOpenFlowControllerIP() (string, error) {
	return "UNIMPLEMENTED!", nil
}

func (o *OVSClient) SetOpenFlowControllerIP() (string, error) {
	return "UNIMPLEMENTED!", nil
}
