package BLC

import (
	"fmt"
	"log"
	"math/big"

	"github.com/boltdb/bolt"
)

// Print 借助for循环遍历区块链
func (bc *BlockChain) Print() {
	err := bc.DB.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(blockTableName))
		var block *Block
		var blockBytes []byte
		var lastHash = bc.Tip
		var hashInt big.Int
		if b != nil {
			//获取最新hash
			for {
				//循环通过区块的preBlockHash进行遍历，打印当前区块信息
				blockBytes = b.Get(lastHash)
				block = DeSerialize(blockBytes)
				lastHash = block.PreBlockHash
				fmt.Println("----------------------------")
				fmt.Printf("\tTimeStamp:%v\n", block.TimeStamp)
				fmt.Printf("\tHash:%x\n", block.Hash)
				fmt.Printf("\tPreBlockHash:%x\n", block.PreBlockHash)
				fmt.Printf("\t:Height%d\n", block.Height)
				fmt.Printf("\tData:%v\n", block.Data)
				fmt.Printf("\tNonce:%d\n", block.Nonce)
				fmt.Println("-----------------------------")
				//判断创世块，退出循环
				if big.NewInt(0).Cmp(hashInt.SetBytes(block.PreBlockHash)) == 0 {
					//遍历到创世块，退出循环
					break
				}
			}

		}
		return nil
	})
	if err != nil {
		log.Panicf("blockchain print failed err:%v", err)
	}
}

//  区块链的遍历方式由此分开 上方是for循环的普通遍历  下方是迭代器遍历

type Iterater struct {
	DB          *bolt.DB
	CurrentHash []byte
}

func (bc *BlockChain) Iterate() *Iterater {
	return &Iterater{bc.DB, bc.Tip}
}

func (i *Iterater) Next() *Block {
	var block *Block
	err := i.DB.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(blockTableName))
		if b != nil {
			blockBytes := b.Get(i.CurrentHash)
			block = DeSerialize(blockBytes)
			i.CurrentHash = block.PreBlockHash
			fmt.Println("----------------------------")
			fmt.Printf("\tTimeStamp:%v\n", block.TimeStamp)
			fmt.Printf("\tHash:%x\n", block.Hash)
			fmt.Printf("\tPreBlockHash:%x\n", block.PreBlockHash)
			fmt.Printf("\tHeight:%d\n", block.Height)
			fmt.Printf("\tData:%s\n", block.Data)
			fmt.Printf("\tNonce:%d\n", block.Nonce)
			fmt.Println("-----------------------------")
		}
		return nil
	})
	if err != nil {
		log.Panicf("next failed")
	}
	return block
}
func (b *BlockChain) PrintChain() {
	iter := b.Iterate()
	var hashInt big.Int
	var block *Block
	for {
		block = iter.Next()
		if big.NewInt(0).Cmp(hashInt.SetBytes(block.PreBlockHash)) == 0 {
			//遍历到创世块，退出循环
			break
		}
	}
}
