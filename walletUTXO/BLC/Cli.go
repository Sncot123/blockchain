package BLC

import (
	"errors"
	"flag"
	"fmt"
	"io/fs"
	"log"
	"os"

	"github.com/boltdb/bolt"
)

//对blockchain进行命令行的管理

//blockchain对象
type CLI struct {
}

//用法展示
func PrintUseage() {
	fmt.Println("Useage:")
	//初始化区块链
	fmt.Printf("\tcreateblockchain -address address--创建区块链\n")
	fmt.Printf("\taddblock -data DATA --添加区块\n")
	fmt.Printf("\tprintchain --打印区块链信息\n")
	fmt.Printf("\tcreatewallets --创建钱包\n")
	//通过命令行转账
	fmt.Printf("\tsend -from FROM to TO -amount AMOUNT -- 发起转账\n")
	fmt.Printf("\t\t-from FROM --转账源地址\n")
	fmt.Printf("\t\t-to TO --转账目标地址\n")
	fmt.Printf("\t\t-amount AMOUNT  --转账金额\n")
	fmt.Printf("\tgetBalance -address FROM -- 查询指定地址的余额\n")
	fmt.Printf("\t查询余额参数说明\n")
	fmt.Printf("\t\t-address 查询的地址")
}
func BlockChainObject() *BlockChain {
	db, err := bolt.Open(dbName, 0600, nil)
	if err != nil {
		log.Panicf("get blockchain object failed err:%v\n", err)
	}
	var tip []byte
	err = db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(blockTableName))
		if b != nil {
			tip = b.Get([]byte("l"))

		}
		return nil
	})
	if err != nil {
		log.Panicf("get the blockchain object failed err:%v\n", err)
	}
	return &BlockChain{db, tip}
}

//判断数据库文件是否存在
func dbExist() bool {
	if _, err := os.Stat(dbName); !errors.Is(err, fs.ErrExist) {
		//数据库文件存在
		fmt.Println(err, "||||", fs.ErrExist)
		return true
	}
	return false
}

// 参数检验
func IsValidArgs() {
	if len(os.Args) < 2 {
		PrintUseage()
		//直接退出
		os.Exit(1)
	}
}
func (cli *CLI) Run() {
	IsValidArgs()

	addBlockCmd := flag.NewFlagSet("addblock", flag.ExitOnError)
	flagAddBlockArg := addBlockCmd.String("data", "sent 100 btc to player", "添加区块数据")

	createBlockChainCmd := flag.NewFlagSet("createblockchain", flag.ExitOnError)
	flagCreateBlockChainArgs := createBlockChainCmd.String("address", "Tom", "指定接收系统奖励的矿工地址")

	sendCmd := flag.NewFlagSet("send", flag.ExitOnError)
	flagSendFromArg := sendCmd.String("from", "", "转账源地址")
	flagSendToArg := sendCmd.String("to", "", "转账目标地址")
	flagSendAmountArg := sendCmd.String("amount", "", "转账金额")

	getBalanceCMD := flag.NewFlagSet("getBalance", flag.ExitOnError)
	flagGetBalanceArg := getBalanceCMD.String("address", "", "获取余额的地址")

	printBlockChainCmd := flag.NewFlagSet("printblockchain", flag.ExitOnError)

	createWalletsCmd := flag.NewFlagSet("createwallets", flag.ExitOnError)

	getAccountsCmd := flag.NewFlagSet("getaccounts", flag.ExitOnError)

	//判断命令
	switch os.Args[1] {
	case "getBalance":
		if err := getBalanceCMD.Parse(os.Args[2:]); err != nil {
			log.Panicf("parse getBalance failed err:%v\n", err)
		}
	case "send":
		if err := sendCmd.Parse(os.Args[2:]); err != nil {
			log.Panicf("parse sendCmd failed!  err:%v\n", err)
		}

	case "addblock":
		if err := addBlockCmd.Parse(os.Args[2:]); err != nil {
			log.Panicf("parse printchainCmd failed!  err:%v\n", err)
		}

	case "printchain":
		if err := printBlockChainCmd.Parse(os.Args[2:]); err != nil {
			log.Panicf("parse printblockchain failed err:%v\n", err)
		}

	case "createblockchain":
		if err := createBlockChainCmd.Parse(os.Args[2:]); err != nil {
			log.Panicf("parse createblockchain failed err:%v\n", err)
		}
	case "createwallets":
		if err := createWalletsCmd.Parse(os.Args[2:]); nil != err {
			log.Panicf("parse createWallets failed! err:%v\n", err)
		}
	case "getaccounts":
		if err := getAccountsCmd.Parse(os.Args[2:]); nil != err {
			log.Panicf("parse getaccounts failed! err:%v\n", err)
		}

	default:
		//没有传递任何命令或者传递的命令不在上述列表之中
		PrintUseage()
		os.Exit(1)
	}
	if printBlockChainCmd.Parsed() {
		cli.printChain()
	}
	if createBlockChainCmd.Parsed() {
		if *flagCreateBlockChainArgs == "" {
			PrintUseage()
			os.Exit(1)
		}
		cli.createBlockChain(*flagCreateBlockChainArgs)
	}
	if addBlockCmd.Parsed() {
		if *flagAddBlockArg == "" {
			PrintUseage()
			os.Exit(1)
		}
		cli.addBlock(*flagAddBlockArg)
	}
	if sendCmd.Parsed() {
		if *flagSendFromArg == "" || *flagSendAmountArg == "" || *flagSendToArg == "" {
			fmt.Println("参数不正确 请重新输入")
			PrintUseage()
			os.Exit(1)
		}
		fmt.Printf("\tfrom:%s\n", JSONToSlice(*flagSendFromArg))
		fmt.Printf("\tto:%s\n", JSONToSlice(*flagSendToArg))
		fmt.Printf("\tamount:%s\n", JSONToSlice(*flagSendAmountArg))

		cli.send(JSONToSlice(*flagSendFromArg), JSONToSlice(*flagSendToArg), JSONToSlice(*flagSendAmountArg))
	}
	if getBalanceCMD.Parsed() {
		if *flagGetBalanceArg == "" {
			log.Panicf("查询地址为空...")
			os.Exit(1)
		}
		cli.getBalance(*flagGetBalanceArg)
	}
	if createWalletsCmd.Parsed() {
		cli.CreateWallets()
	}
	if getAccountsCmd.Parsed() {
		cli.getAccounts()
	}

}
