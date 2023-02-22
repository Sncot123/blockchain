package BLC

import (
	"fmt"
	"os"
)

//发起交易函数
func (cli *CLI) send(from, to, amount []string) {
	if !dbExist() {
		fmt.Println("数据库不存在...")
		os.Exit(1)
	}
	//获取区块链对象
	blockChain := BlockChainObject()
	defer blockChain.DB.Close()
	if len(from) != len(to) || len(from) != len(amount) {
		fmt.Println("参数有误！请检查参数的一致性")
	}
	blockChain.MineNewBlock(from, to, amount)
}
