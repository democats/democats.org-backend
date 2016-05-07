package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/callum-ramage/jsonconfig"
	"io/ioutil"
	"net/http"
	"os"
	"os/exec"
	"os/user"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"
)

// Change to 0 to connect to remote host
var localhost_for_daemon_rpc = 0
var parent_dir = "/home/fork"

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
			GENESISBLOCKREWARD                     int      `json:"GENESIS_BLOCK_REWARD"`
			GENESISCOINBASETXHEX                   string   `json:"GENESIS_COINBASE_TX_HEX"`
			MAXBLOCKSIZEINITIAL                    int      `json:"MAX_BLOCK_SIZE_INITIAL"`
			MINIMUMFEE                             int      `json:"MINIMUM_FEE"`
			MONEYSUPPLY                            string   `json:"MONEY_SUPPLY"`
			P2PDEFAULTPORT                         int      `json:"P2P_DEFAULT_PORT"`
			RPCDEFAULTPORT                         int      `json:"RPC_DEFAULT_PORT"`
			SEEDNODES                              []string `json:"SEED_NODES"`
			UPGRADEHEIGHT                          int      `json:"UPGRADE_HEIGHT"`
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

func createDir(path string) {
	os.MkdirAll(path, 0775)
}

func writeFile(path string, content []byte) {
	os.MkdirAll(filepath.Dir(path), 0775)
	err := ioutil.WriteFile(path, content, 0644)
	check(err)
}

func printCommand(cmd *exec.Cmd) {
	fmt.Printf("==> Executing: %s\n", strings.Join(cmd.Args, " "))
}

func printError(err error) {
	if err != nil {
		os.Stderr.WriteString(fmt.Sprintf("==> Error: %s\n", err.Error()))
	}
}

func printOutput(outs []byte) {
	if len(outs) > 0 {
		fmt.Printf("==> Output: %s\n", string(outs))
	}
}

func main() {
	dir, err := filepath.Abs(filepath.Dir(os.Args[0]))
	check(err)
	pools, err := jsonconfig.LoadAbstract(dir+"/pools.json", "")

	if parent_dir == "" {
		usr, err := user.Current()
		check(err)
		parent_dir = usr.HomeDir
	}

	for _, element := range pools["pools"].Arr[:] {
		var daemon_rpc = element.Obj["daemonrpc"].Str
		if localhost_for_daemon_rpc == 1 {
			r, err := regexp.Compile(`.+:`)
			check(err)
			daemon_rpc = r.ReplaceAllString(daemon_rpc, "http://localhost:")
		}

		// Get height
		resp, err := http.Get(element.Obj["daemonrpc"].Str + "getheight")
		check(err)
		defer resp.Body.Close()
		body, _ := ioutil.ReadAll(resp.Body)

		var ch CoinHeight
		err = json.Unmarshal(body, &ch)
		check(err)

		// Get settings
		var jsonStr = []byte(`{"method":"f_get_blockchain_settings"}`)
		req, err := http.NewRequest("POST", daemon_rpc+"json_rpc", bytes.NewBuffer(jsonStr))
		req.Header.Set("Content-Type", "application/json")

		client := &http.Client{}
		resp, err = client.Do(req)
		check(err)
		defer resp.Body.Close()
		body, _ = ioutil.ReadAll(resp.Body)

		var cc CoinConfig
		err = json.Unmarshal(body, &cc)
		check(err)

		if string(cc.Result.Core.CRYPTONOTENAME) == "bytecoin" {

		}

		blockchain_dir := parent_dir + "/." + string(cc.Result.Core.CRYPTONOTENAME)
		main_output_dir := dir + "/blockchains_output/"
		blockchain_output_dir := main_output_dir + string(cc.Result.Core.CRYPTONOTENAME) + "-blockchain/"
		createDir(blockchain_output_dir)
		fmt.Println("blockchain_dir: ", blockchain_dir)
		fmt.Println("blockchain_output_dir: ", blockchain_output_dir)

		// Copy blockchain into folder
		// blockindexes.dat
		cmd := exec.Command("cp", "-r", blockchain_dir+"/blockindexes.dat", blockchain_output_dir)
		printCommand(cmd)
		output, err := cmd.CombinedOutput()
		printError(err)
		printOutput(output)

		// blocks.dat
		cmd = exec.Command("cp", "-r", blockchain_dir+"/blocks.dat", blockchain_output_dir)
		printCommand(cmd)
		output, err = cmd.CombinedOutput()
		printError(err)
		printOutput(output)

		// lasthash
		height_string := strconv.Itoa(ch.Height - 1)
		jsonStr = []byte(`{"method": "f_block_json","params": {"hash": "` + height_string + `"}}`)
		req, err = http.NewRequest("POST", daemon_rpc+"json_rpc", bytes.NewBuffer(jsonStr))
		req.Header.Set("Content-Type", "application/json")

		client = &http.Client{}
		resp, err = client.Do(req)
		check(err)
		defer resp.Body.Close()
		body, _ = ioutil.ReadAll(resp.Body)

		var block Block
		err = json.Unmarshal(body, &block)
		check(err)

		writeFile(blockchain_output_dir+"lasthash", []byte(block.Result.Block.Hash))
		fmt.Println("Network: ", string(cc.Result.Core.CRYPTONOTENAME), "   Height: ", height_string, "   Hash: ", string(block.Result.Block.Hash))

		// Zip blockchain
		cmd = exec.Command("zip", "-r", main_output_dir+string(cc.Result.Core.CRYPTONOTENAME)+"-blockchain.zip", blockchain_output_dir)
		printCommand(cmd)
		output, err = cmd.CombinedOutput()
		printError(err)
		printOutput(output)

		// Delete folders
		cmd = exec.Command("rm", "-rf", blockchain_output_dir)
		printCommand(cmd)
		output, err = cmd.CombinedOutput()
		printError(err)
		printOutput(output)
	}
}
