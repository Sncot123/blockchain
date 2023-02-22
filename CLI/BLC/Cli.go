package BLC

import (
	"flag"
	"fmt"
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
	fmt.Printf("\tcreateblockchain --创建区块链\n")
	fmt.Printf("\taddblock -data DATA --添加区块\n")
	fmt.Printf("\tprintchain --打印区块链信息\n")
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
	if _, err := os.Stat(dbName); os.IsExist(err) {
		//数据库文件不存在
		return false
	}
	return true
}

//创建区块链
func (cli *CLI) createBlockChain() {
	//fmt.Println(dbExist())
	//if dbExist() {
	//	fmt.Println("创世区块已存在...")
	//	os.Exit(1)
	//}
	CreateBlockChain()
}

//添加区块
func (cli *CLI) addBlock(data string) {
	if !dbExist() {
		fmt.Println("数据库不存在...")
		os.Exit(1)
	}
	blockchain := BlockChainObject()
	blockchain.AddBlock(data)
}

//打印区块链信息
func (cli *CLI) printChain() {
	if !dbExist() {
		fmt.Println("数据库不存在...")
		os.Exit(1)
	}
	//获取区块链对象
	blockChain := BlockChainObject()
	//打印区块链信息
	blockChain.PrintChain()
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

	printBlockChainCmd := flag.NewFlagSet("printblockchain", flag.ExitOnError)

	createBlockChainCmd := flag.NewFlagSet("createblockchain", flag.ExitOnError)

	flagAddBlockArg := addBlockCmd.String("data", "sent 100 btc to player", "添加区块数据")

	//判断命令
	switch os.Args[1] {
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
			fmt.Println("=======")
			log.Panicf("parse createblockchain failed err:%v\n", err)
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
		cli.createBlockChain()
	}
	if addBlockCmd.Parsed() {
		if *flagAddBlockArg == "" {
			PrintUseage()
			os.Exit(1)
		}
		cli.addBlock(*flagAddBlockArg)
	}
}
