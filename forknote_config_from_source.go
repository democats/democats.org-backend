// Generate Forknote config file based on the source code of a Cryptonote coin
// NOTICE: It works only for new coins
/*
Missing values:
BYTECOIN_NETWORK    Used for network packages in order not to mix two different cryptocoin networks
CHECKPOINT  Checkpoints. Format: HEIGHT:HASH
CRYPTONOTE_BLOCK_GRANTED_FULL_REWARD_ZONE_V1    The maximum size of a block not resulting into penelty. Used only by old (v1) coins
CRYPTONOTE_BLOCK_GRANTED_FULL_REWARD_ZONE_V2    The maximum size of a block not resulting into penelty. Used only by old (v2) coins
seed_nodes
*/

package main

import (
    "bufio"
    "fmt"
    "os"
    "os/exec"
    "regexp"
)

func checkRegexpInLine(line string, regex *regexp.Regexp, err error) string {
    output := ""

    if err != nil {
        panic(err)
    }
    result := regex.FindStringSubmatch(line)

    if len(result) > 1 {
        for k, v := range result {
            if k == 0 {
                continue
            }
            if k == 1 && v == "P2P_DEFAULT_PORT" {
                output += "p2p-bind-port"
            } else if k == 1 && v == "RPC_DEFAULT_PORT" {
                output += "rpc-bind-port"
            } else {
                output += v
            }
            if k == 1 {
                output += " = "
            }
        }
        output += "\n"
    }

    return output
}

func main() {
    output := ""
    isCryptonoteCoin := true

    if len(os.Args) != 2 {
        fmt.Println("Usage: ./forknote_config_from_source github_repository_url")
        os.Exit(1)
    }
    arg := os.Args[1]

    cloneCmd := exec.Command("git", "clone", arg, "tmp_repository")

    _, err := cloneCmd.Output()
    if err != nil {
        panic(err)
    }

    file, err := os.Open("tmp_repository/src/CryptoNoteConfig.h")
    if err != nil {
        panic(err)
    }
    defer file.Close()

    scanner := bufio.NewScanner(file)
    for scanner.Scan() {
        regex, err := regexp.Compile(`(CRYPTONOTE_BLOCK_GRANTED_FULL_REWARD_ZONE)\s+=\s+(.+);`)
        output += checkRegexpInLine(scanner.Text(), regex, err)
        regex, err = regexp.Compile(`(CRYPTONOTE_BLOCK_GRANTED_FULL_REWARD_ZONE_V1)\s+=\s+(.+);`)
        output += checkRegexpInLine(scanner.Text(), regex, err)
        regex, err = regexp.Compile(`(CRYPTONOTE_BLOCK_GRANTED_FULL_REWARD_ZONE_V2)\s+=\s+(.+);`)
        output += checkRegexpInLine(scanner.Text(), regex, err)
        regex, err = regexp.Compile(`(CRYPTONOTE_DISPLAY_DECIMAL_POINT)\s+=\s+(.+);`)
        output += checkRegexpInLine(scanner.Text(), regex, err)
        regex, err = regexp.Compile(`(CRYPTONOTE_MINED_MONEY_UNLOCK_WINDOW)\s+=\s+(.+);`)
        output += checkRegexpInLine(scanner.Text(), regex, err)
        regex, err = regexp.Compile(`(CRYPTONOTE_NAME)\[\]\s+=\s+(.+);`)
        output += checkRegexpInLine(scanner.Text(), regex, err)
        regex, err = regexp.Compile(`(CRYPTONOTE_PUBLIC_ADDRESS_BASE58_PREFIX)\s+=\s+(.+);`)
        output += checkRegexpInLine(scanner.Text(), regex, err)
        regex, err = regexp.Compile(`(DEFAULT_DUST_THRESHOLD)\s+=\s+(.+);`)
        output += checkRegexpInLine(scanner.Text(), regex, err)
        regex, err = regexp.Compile(`(DIFFICULTY_CUT)\s+=\s+(.+);`)
        output += checkRegexpInLine(scanner.Text(), regex, err)
        regex, err = regexp.Compile(`(DIFFICULTY_TARGET)\s+=\s+(.+);`)
        output += checkRegexpInLine(scanner.Text(), regex, err)
        regex, err = regexp.Compile(`(EMISSION_SPEED_FACTOR)\s+=\s+(.+);`)
        output += checkRegexpInLine(scanner.Text(), regex, err)
        regex, err = regexp.Compile(`(EXPECTED_NUMBER_OF_BLOCKS_PER_DAY)\s+=\s+(.+);`)
        output += checkRegexpInLine(scanner.Text(), regex, err)
        regex, err = regexp.Compile(`(GENESIS_BLOCK_REWARD)\s+=\s+(.+);`)
        output += checkRegexpInLine(scanner.Text(), regex, err)
        regex, err = regexp.Compile(`(GENESIS_COINBASE_TX_HEX)\[\]\s+=\s+(.+);`)
        output += checkRegexpInLine(scanner.Text(), regex, err)
        regex, err = regexp.Compile(`(KILL_HEIGHT)\s+=\s+(.+);`)
        output += checkRegexpInLine(scanner.Text(), regex, err)
        regex, err = regexp.Compile(`(MANDATORY_TRANSACTION)\s+=\s+(.+);`)
        output += checkRegexpInLine(scanner.Text(), regex, err)
        regex, err = regexp.Compile(`(MAX_BLOCK_SIZE_INITIAL)\s+=\s+(.+);`)
        output += checkRegexpInLine(scanner.Text(), regex, err)
        regex, err = regexp.Compile(`(MINIMUM_FEE)\s+=\s+(.+);`)
        output += checkRegexpInLine(scanner.Text(), regex, err)
        regex, err = regexp.Compile(`(MONEY_SUPPLY)\s+=\s+(.+);`)
        output += checkRegexpInLine(scanner.Text(), regex, err)
        regex, err = regexp.Compile(`(UPGRADE_HEIGHT_V2)\s+=\s+(.+);`)
        latestOutput := checkRegexpInLine(scanner.Text(), regex, err)
        if latestOutput != "" {
            output += latestOutput
            isCryptonoteCoin = false
        }
        regex, err = regexp.Compile(`(UPGRADE_HEIGHT_V3)\s+=\s+(.+);`)
        output += checkRegexpInLine(scanner.Text(), regex, err)
        regex, err = regexp.Compile(`(P2P_DEFAULT_PORT)\s+=\s+(.+);`)
        output += checkRegexpInLine(scanner.Text(), regex, err)
        regex, err = regexp.Compile(`(RPC_DEFAULT_PORT)\s+=\s+(.+);`)
        output += checkRegexpInLine(scanner.Text(), regex, err)
    }

    if err := scanner.Err(); err != nil {
        panic(err)
    }

    output += "P2P_STAT_TRUSTED_PUB_KEY=\n"

    if isCryptonoteCoin == true {
        output += "UPGRADE_HEIGHT_V2=4294967294\n"
        output += "UPGRADE_HEIGHT_V3=4294967295\n\n"
        output += "CRYPTONOTE_COIN_VERSION=1\n"
    }

    rmCmd := exec.Command("rm", "-r", "tmp_repository")

    _, err = rmCmd.Output()
    if err != nil {
        panic(err)
    }

    fmt.Printf("%s\n", output)
}
