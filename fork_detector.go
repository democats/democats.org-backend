// Known issues: hashrate is overestimated for 1h

package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/callum-ramage/jsonconfig"
	"io/ioutil"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"time"
)

// Change to 0 to connect to remote host
var localhost_for_daemon_rpc = 0
var start_height = 290000
var end_height = 300000

type CoinHeight struct {
	Height int    `json:"height"`
	Status string `json:"status"`
}

type CoinConfig struct {
	Jsonrpc string `json:"jsonrpc"`
	Result  struct {
		BaseCoin struct {
			Git  string `json:"git"`
			Name string `json:"name"`
		} `json:"base_coin"`
		Core struct {
			BYTECOINNETWORK                        string   `json:"BYTECOIN_NETWORK"`
			CHECKPOINTS                            []string `json:"CHECKPOINTS"`
			CRYPTONOTEBLOCKGRANTEDFULLREWARDZONE   int      `json:"CRYPTONOTE_BLOCK_GRANTED_FULL_REWARD_ZONE"`
			CRYPTONOTEBLOCKGRANTEDFULLREWARDZONEV1 int      `json:"CRYPTONOTE_BLOCK_GRANTED_FULL_REWARD_ZONE_V1"`
			CRYPTONOTEBLOCKGRANTEDFULLREWARDZONEV2 int      `json:"CRYPTONOTE_BLOCK_GRANTED_FULL_REWARD_ZONE_V2"`
			CRYPTONOTECOINVERSION                  int      `json:"CRYPTONOTE_COIN_VERSION"`
			CRYPTONOTEDISPLAYDECIMALPOINT          int      `json:"CRYPTONOTE_DISPLAY_DECIMAL_POINT"`
			CRYPTONOTEMINEDMONEYUNLOCKWINDOW       int      `json:"CRYPTONOTE_MINED_MONEY_UNLOCK_WINDOW"`
			CRYPTONOTENAME                         string   `json:"CRYPTONOTE_NAME"`
			CRYPTONOTEPUBLICADDRESSBASE58PREFIX    int      `json:"CRYPTONOTE_PUBLIC_ADDRESS_BASE58_PREFIX"`
			DEFAULTDUSTTHRESHOLD                   int      `json:"DEFAULT_DUST_THRESHOLD"`
			DIFFICULTYCUT                          int      `json:"DIFFICULTY_CUT"`
			DIFFICULTYLAG                          int      `json:"DIFFICULTY_LAG"`
			DIFFICULTYTARGET                       int      `json:"DIFFICULTY_TARGET"`
			EMISSIONSPEEDFACTOR                    int      `json:"EMISSION_SPEED_FACTOR"`
			EXPECTEDNUMBEROFBLOCKSPERDAY           int      `json:"EXPECTED_NUMBER_OF_BLOCKS_PER_DAY"`
			GENESISBLOCKREWARD                     string   `json:"GENESIS_BLOCK_REWARD"`
			GENESISCOINBASETXHEX                   string   `json:"GENESIS_COINBASE_TX_HEX"`
			KILLHEIGHT                             int      `json:"KILL_HEIGHT"`
			MANDATORYTRANSACTION                   int      `json:"MANDATORY_TRANSACTION"`
			MAXBLOCKSIZEINITIAL                    int      `json:"MAX_BLOCK_SIZE_INITIAL"`
			MINIMUMFEE                             int      `json:"MINIMUM_FEE"`
			MONEYSUPPLY                            string   `json:"MONEY_SUPPLY"`
			P2PDEFAULTPORT                         int      `json:"P2P_DEFAULT_PORT"`
			RPCDEFAULTPORT                         int      `json:"RPC_DEFAULT_PORT"`
			SEEDNODES                              []string `json:"SEED_NODES"`
			UPGRADEHEIGHTV2                        int      `json:"UPGRADE_HEIGHT_V2"`
			UPGRADEHEIGHTV3                        int      `json:"UPGRADE_HEIGHT_V3"`
		} `json:"core"`
		Extensions []string `json:"extensions"`
		Status     string   `json:"status"`
	} `json:"result"`
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

	dir, err := filepath.Abs(filepath.Dir(os.Args[0]))
	check(err)

	nodes, err := jsonconfig.LoadAbstract(dir+"/nodes.json", "")

	for i := start_height + 1; i < end_height; i++ {
		prevHash := ""
		for _, node_address := range nodes["nodes"].Arr[:] {
			// Don't kill the server!
			time.Sleep(10 * time.Millisecond)

			var height = i
			height_string := strconv.Itoa(height)
			jsonStr := []byte(`{"method": "f_block_json","params": {"hash": "` + height_string + `"}}`)
			req, err := http.NewRequest("POST", node_address.Str+"json_rpc", bytes.NewBuffer(jsonStr))
			req.Header.Set("Content-Type", "application/json")

			client := &http.Client{}
			resp, err := client.Do(req)
			check(err)
			defer resp.Body.Close()
			body, _ := ioutil.ReadAll(resp.Body)

			var block Block
			err = json.Unmarshal(body, &block)
			check(err)

			if height%100 == 0 {
				fmt.Println("Height: " + height_string)
			}

			if prevHash == "" {
				prevHash = block.Result.Block.Hash
			} else {
				if prevHash != block.Result.Block.Hash {
					fmt.Println(prevHash)
					os.Exit(1)
				}
			}
		}
	}
}
