package BLC

import "fmt"

func (cli *CLI) getBalance(from string) {
	blockchain := BlockChainObject()
	defer blockchain.DB.Close()
	amount := blockchain.getbalance(from)
	fmt.Printf("\t地址 [%s]的余额：[%d]\n", from, amount)
}
