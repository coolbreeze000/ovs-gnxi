package main

var RunEnv = NewRuntimeEnvironment()

type RuntimeEnvironment struct {
}

func NewRuntimeEnvironment() *RuntimeEnvironment {
	r := RuntimeEnvironment{}
	return &r
}
