package main

import (
	"encoding/json"
	"fmt"
	"github.com/callum-ramage/jsonconfig"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"time"
)

// your JSON structure as a byte slice
//var j = []byte(`{"foo":1,"bar":2,"baz":[3,4]}`)

type PoolStats struct {
	Currency            string  `json:"currency"`
	Workers             float64 `json:"workers"`
	MinPaymentThreshold float64 `json:"min_payment_threshold"`
	Fee                 float64 `json:"fee"`
	Donation            float64 `json:"donation"`
	Height              float64 `json:"height"`
	LastBlock           string  `json:"last_block"`
	Difficulty          float64 `json:"difficulty"`
	CoinUnits           float64 `json:"coin_units"`
	Depth               float64 `json:"depth"`
	TotalPayments       float64 `json:"total_payments"`
	TotalMinersPaid     float64 `json:"total_miners_paid"`
	PoolHashrate        float64 `json:"pool_hashrate"`
	WorldHashrate       float64 `json:"world_hashrate"`
}

func (a *PoolStats) UnmarshalJSON(b []byte) error {
	var f interface{}
	json.Unmarshal(b, &f)

	m := f.(map[string]interface{})

	configmap := m["config"]
	v1 := configmap.(map[string]interface{})

	a.Currency = v1["coin"].(string)
	a.MinPaymentThreshold = v1["minPaymentThreshold"].(float64)
	a.Fee = v1["fee"].(float64)
	a.CoinUnits = v1["coinUnits"].(float64)
	a.Depth = v1["depth"].(float64)

	poolmap := m["pool"]
	v2 := poolmap.(map[string]interface{})

	a.PoolHashrate = v2["hashrate"].(float64)
	a.Workers = v2["miners"].(float64)
	a.TotalPayments = v2["totalPayments"].(float64)
	a.TotalMinersPaid = v2["totalMinersPaid"].(float64)

	networkmap := m["network"]
	v3 := networkmap.(map[string]interface{})

	a.Difficulty = v3["difficulty"].(float64)
	a.Height = v3["height"].(float64)
	a.LastBlock = v3["hash"].(string)
	a.WorldHashrate = v3["difficulty"].(float64) / v1["coinDifficultyTarget"].(float64)

	return nil
}

func main() {
	timeout := time.Duration(3 * time.Second)
	httpClient := http.Client{
		Timeout: timeout,
	}

	stats := []PoolStats{}

	dir, err := filepath.Abs(filepath.Dir(os.Args[0]))
	if err != nil {
		log.Fatal(err)
	}

	pools, err := jsonconfig.LoadAbstract(dir+"/pools.json", "")

	for _, element := range pools["pools"].Arr[:] {
		resp, err := httpClient.Get(element.Obj["poolrpc"].Str + "stats")
		if err != nil {
			// Next if the pool is down
			continue
		}
		defer resp.Body.Close()
		body, err := ioutil.ReadAll(resp.Body)

		if len(body) == 0 {
			continue
		}

		var s PoolStats
		err = json.Unmarshal(body, &s)
		if err != nil {
			log.Fatal(err)
			return
		}

		stats = append(stats, s)
	}

	b, err := json.Marshal(stats)
	if err != nil {
		log.Fatal(err)
		return
	}

	fmt.Printf("%s", string(b))
}
