package main

import (
	"blockchain/BLCDB/BLC"
)

func main() {
	chain := BLC.CreateBlockChain()
	chain.AddBlock("a")
	chain.AddBlock("b")
	chain.AddBlock("c")
	chain.AddBlock("d")
	chain.PrintChain()
}
