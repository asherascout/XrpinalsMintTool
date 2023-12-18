package main

import (
	"flag"
	"fmt"
	"github.com/Xrpinals-Protocol/XrpinalsMintTool/conf"
	"github.com/Xrpinals-Protocol/XrpinalsMintTool/key"
	"github.com/Xrpinals-Protocol/XrpinalsMintTool/logger"
	"github.com/Xrpinals-Protocol/XrpinalsMintTool/mining"
	"github.com/Xrpinals-Protocol/XrpinalsMintTool/tx_builder"
	"github.com/Xrpinals-Protocol/XrpinalsMintTool/utils"
	"os"
)

func appInit() {
	// init log
	err := logger.InitAppLog("mint-tool.log")
	if err != nil {
		panic(err)
	}
}

func main() {
	appInit()

	if os.Args[1] == "mint" {
		fs := flag.NewFlagSet("mint", flag.ExitOnError)

		var addr string
		var asset string
		fs.StringVar(&addr, "addr", "", "your address")
		fs.StringVar(&asset, "asset", "", "asset name you want to mint")
		err := fs.Parse(os.Args[2:])
		if err != nil {
			panic(err)
		}
		mining.MintAssetName = asset

		ok, err := key.IsAddressExisted(addr)
		if err != nil {
			panic(err)
		}

		if ok {
			mining.PrivateKey, err = key.GetAddressKey(addr)
			if err != nil {
				panic(err)
			}
		} else {
			panic(fmt.Errorf("address %s is not in the storage", addr))
		}

		mining.StartMining()

	} else if os.Args[1] == "import_key" {
		fs := flag.NewFlagSet("import_key", flag.ExitOnError)

		fmt.Println(os.Args[1])
		var pKey string
		fs.StringVar(&pKey, "key", "", "your private key")
		err := fs.Parse(os.Args[2:])
		if err != nil {
			panic(err)
		}

		fmt.Println(pKey)
		addr, err := key.ImportPrivateKey(pKey)
		if err != nil {
			panic(err)
		}

		fmt.Printf("private key of address %s imported\n", addr)
	} else if os.Args[1] == "check_address" {
		fs := flag.NewFlagSet("check_address", flag.ExitOnError)

		var addr string
		fs.StringVar(&addr, "addr", "", "your address")
		err := fs.Parse(os.Args[2:])
		if err != nil {
			panic(err)
		}

		ok, err := key.IsAddressExisted(addr)
		if err != nil {
			panic(err)
		}

		if ok {
			fmt.Printf("address %s is in the storage\n", addr)
		} else {
			panic(fmt.Errorf("address %s is not in the storage", addr))
		}
	} else if os.Args[1] == "get_balance" {
		fs := flag.NewFlagSet("get_balance", flag.ExitOnError)

		var addr string
		var asset string
		fs.StringVar(&addr, "addr", "", "your address")
		fs.StringVar(&asset, "asset", "", "asset name you want to query")
		err := fs.Parse(os.Args[2:])
		if err != nil {
			panic(err)
		}

		resp, err := utils.GetAssetInfo(conf.GetConfig().WalletRpcUrl, asset)
		if err != nil {
			panic(err)
		}

		balance, err := utils.GetAddressBalance(conf.GetConfig().WalletRpcUrl, addr, resp.Result.Id)
		if err != nil {
			panic(err)
		}

		fmt.Println("balance: ", balance.String())
	} else if os.Args[1] == "get_deposit_address" {
		result, err := utils.GetDepositAddress(conf.GetConfig().WalletRpcUrl, "BTC")
		if err != nil {
			panic(err)
		}

		fmt.Println("BTC deposit address: ", result.BindAccountHot)
	} else if os.Args[1] == "transfer" {
		fs := flag.NewFlagSet("transfer", flag.ExitOnError)

		var from string
		var to string
		var asset string
		var amount string
		var keyWif string
		fs.StringVar(&from, "from", "", "your address")
		fs.StringVar(&to, "to", "", "receiver address")
		fs.StringVar(&asset, "asset", "", "asset name you want to transfer")
		fs.StringVar(&amount, "amount", "", "asset amount you want to transfer")
		err := fs.Parse(os.Args[2:])
		if err != nil {
			panic(err)
		}

		ok, err := key.IsAddressExisted(from)
		if err != nil {
			panic(err)
		}

		if ok {
			keyWif, err = key.GetAddressKey(from)
			if err != nil {
				panic(err)
			}
		} else {
			panic(fmt.Errorf("address %s is not in the storage", from))
		}

		txHash, err := tx_builder.Transfer(from, to, asset, amount, keyWif)
		if err != nil {
			panic(fmt.Errorf("transfer error:%v", err))
		}

		fmt.Printf("transfer success,txHash:%v", txHash)
	}
}
