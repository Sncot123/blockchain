package BLC

import (
	"fmt"
	"os"
)

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
