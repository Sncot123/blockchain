package main

import (
	"blockchain/block/BLC"
	"fmt"
)

func main() {
	chain := BLC.CreateBlockChain()
	chain.AddBlock(2, []byte("第二区块"))
	chain.AddBlock(2, []byte("第三区块"))
	for _, v := range chain.Blocks {
		fmt.Printf("Hash:%x\n", v.Hash)
	}

}
