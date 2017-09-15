package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"math/rand"
	"os"
	"sync"
	"time"
)

type Host struct {
	IP     string `json:"IP"`
	Status int    `json:"status"`
}

type Group struct {
	Min   int    `json:"min"`
	Max   int    `json:"max"`
	Hosts []Host `json:"Hosts"`
}

type Root struct {
	Listen        string  `json:"Listen"`
	Mode          string  `json:"Mode"`
	CheckInterval int64   `json:"CheckInterval"`
	AllHost       []Group `json:"AllHost"`
}

var ALLHOST Root
var lock sync.RWMutex

func LoadConfig() {
	rand.Seed(time.Now().UnixNano())
	file, err := ioutil.ReadFile("config.json")
	if err != nil {
		fmt.Println("ReadFile Err: ", err)
		os.Exit(-1)
		return
	}
	err = json.Unmarshal(file, &ALLHOST)
	if err != nil {
		fmt.Println("json.Unmarshal Err: ", err)
		os.Exit(-1)
		return
	}

	fmt.Printf("\nAllHost=%+v\n", ALLHOST)
}
