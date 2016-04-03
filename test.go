package main

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"net/http"
)

func main() {
	url := "http://pool.democats.org:7611/json_rpc"
	fmt.Println("URL:>", url)

	//    var jsonStr = []byte(`{method":"f_get_blockchain_settings"}`)
	var jsonStr = []byte(`{"method": "f_blocks_list_json",
    "params": {
        "hash": 10
    }}`)
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonStr))
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()

	fmt.Println("response Status:", resp.Status)
	fmt.Println("response Headers:", resp.Header)
	body, _ := ioutil.ReadAll(resp.Body)
	fmt.Println("response Body:", string(body))
}
