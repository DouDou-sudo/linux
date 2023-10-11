package main

import (
	"log"
	"os"
	"sync"

	"gopkg.in/yaml.v2"
	// "gopkg.in/yaml.v3"
)

var ipAddrs chan string = make(chan string)
var wg sync.WaitGroup

type HostInfo struct {
	Group  string   `yaml:"group"`
	IpList []string `yaml:"ip"`
}

type Groups struct {
	List []*HostInfo `yaml:"groups"`
}

func OpenCfgFile() {
	databyte, err := os.ReadFile("go-ansible.yml")
	if err != nil {
		log.Fatalf("read cfg file fail.", err)
	}

	g := new(Groups)
	err = yaml.Unmarshal(databyte, &g)
	if err != nil {
		log.Fatalf("read cfg file fail.", err)
	}
	for _,hostinfo :=range g.List 
}



func main() {

}
