package main

import "blockchain/walletUTXO/BLC"

func main() {
	cli := new(BLC.CLI)
	cli.Run()
}
