package BLC

import "fmt"

func (cli *CLI) getAccounts() {
	wallets := NewWallets()
	fmt.Println("\t账号列表")
	for key, _ := range wallets.Wallets {
		fmt.Printf("\t\t[%s]\n", key)
	}
}
