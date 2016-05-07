package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strconv"
	"time"
)

type CoinHeight struct {
	Height int    `json:"height"`
	Status string `json:"status"`
}

type Block struct {
	ID      string `json:"id"`
	Jsonrpc string `json:"jsonrpc"`
	Result  struct {
		Block struct {
			AlreadyGeneratedCoins        string  `json:"alreadyGeneratedCoins"`
			AlreadyGeneratedTransactions int     `json:"alreadyGeneratedTransactions"`
			BaseReward                   int     `json:"baseReward"`
			BlockSize                    int     `json:"blockSize"`
			Depth                        int     `json:"depth"`
			Difficulty                   int     `json:"difficulty"`
			EffectiveSizeMedian          int     `json:"effectiveSizeMedian"`
			Hash                         string  `json:"hash"`
			Height                       int     `json:"height"`
			MajorVersion                 int     `json:"major_version"`
			MinorVersion                 int     `json:"minor_version"`
			Nonce                        int     `json:"nonce"`
			OrphanStatus                 bool    `json:"orphan_status"`
			Penalty                      float64 `json:"penalty"`
			PrevHash                     string  `json:"prev_hash"`
			Reward                       int     `json:"reward"`
			SizeMedian                   int     `json:"sizeMedian"`
			Timestamp                    int     `json:"timestamp"`
			TotalFeeAmount               int     `json:"totalFeeAmount"`
			Transactions                 []struct {
				AmountOut int    `json:"amount_out"`
				Fee       int    `json:"fee"`
				Hash      string `json:"hash"`
				Size      int    `json:"size"`
			} `json:"transactions"`
			TransactionsCumulativeSize int `json:"transactionsCumulativeSize"`
		} `json:"block"`
		Status string `json:"status"`
	} `json:"result"`
}

func check(e error) {
	if e != nil {
		panic(e)
	}
}

func main() {

	// Get height
	resp, err := http.Get("http://localhost:8081/getheight")
	check(err)
	defer resp.Body.Close()
	body, _ := ioutil.ReadAll(resp.Body)

	var ch CoinHeight
	err = json.Unmarshal(body, &ch)
	check(err)

	printed := 0

	for i := 981539; i < ch.Height; i++ {
		// Don't kill the server!
		time.Sleep(10 * time.Millisecond)

		var height = i
		height_string := strconv.Itoa(height)
		jsonStr := []byte(`{"method": "f_block_json","params": {"hash": "` + height_string + `"}}`)
		req, err := http.NewRequest("POST", "http://localhost:8081/json_rpc", bytes.NewBuffer(jsonStr))
		req.Header.Set("Content-Type", "application/json")

		client := &http.Client{}
		resp, err := client.Do(req)
		check(err)
		defer resp.Body.Close()
		body, _ = ioutil.ReadAll(resp.Body)

		var block Block
		err = json.Unmarshal(body, &block)
		check(err)

		if printed == 0 && block.Result.Block.MajorVersion == 3 {
			fmt.Println("Height v3.0: ", height_string)
			printed = 1
		}
	}
}
