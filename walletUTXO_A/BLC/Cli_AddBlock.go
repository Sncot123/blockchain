package BLC

import (
	"fmt"
	"os"
)

//添加区块
func (cli *CLI) addBlock(address string) {
	if !dbExist() {
		fmt.Println("数据库不存在...")
		os.Exit(1)
	}
	//blockchain := BlockChainObject()
	//blockchain.AddBlock(address)
}
