package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
)

type PoolsConfig struct {
	Pools []struct {
		Name       string `json:"name"`
		Poolmining string `json:"poolmining"`
		Poolrpc    string `json:"poolrpc"`
		Ticker     string `json:"ticker"`
		Daemonrpc  string `json:"daemonrpc"`
	} `json:"pools"`
}

func main() {

	dir, err := filepath.Abs(filepath.Dir(os.Args[0]))
	if err != nil {
		log.Fatal(err)
	}

	b, err := ioutil.ReadFile(dir + "/pools.json")
	if err != nil {
		log.Fatal(err)
	}

	var pc PoolsConfig
	err = json.Unmarshal(b, &pc)
	if err != nil {
		log.Fatal(err)
	}

	out, err := json.Marshal(pc)

	fmt.Printf("%s", out)
}
