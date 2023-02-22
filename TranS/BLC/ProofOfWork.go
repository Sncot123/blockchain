package BLC

import (
	"bytes"
	"crypto/sha256"
	"fmt"
	"math/big"
)

//目标难度值
const targetBit = 16

type ProofOfWork struct {
	//需要共识验证的区块
	Block *Block
	//目标难度的哈希（大数据存储）
	target *big.Int
}

func NewProofOfWork(block *Block) *ProofOfWork {
	target := big.NewInt(1)
	//数据长度为8位
	//需求：需要满足前两位为0，才能解决问题
	//256=32*8
	target = target.Lsh(target, 256-targetBit)
	return &ProofOfWork{Block: block, target: target}
}

//执行pow，返回哈希
//返回哈希值，以及碰撞次数
func (pow *ProofOfWork) Run() ([]byte, int64) {
	//碰撞次数
	var nonce = int64(0)
	var hashInt big.Int
	var hash [32]byte
	//无限循环，生成符合条件的哈希值
	for {
		//生成准备数据
		dataBytes := pow.prepareData(int64(nonce))
		hash = sha256.Sum256(dataBytes)
		hashInt.SetBytes(hash[:])
		//检测生成的哈希值是否符合条件
		if pow.target.Cmp(&hashInt) == 1 {
			//找到了符合条件的哈希，中断循环
			break
		}
		nonce++
	}
	fmt.Printf("碰撞次数:%v\n", nonce)
	return hash[:], nonce
}

//生成准备数据
func (pow *ProofOfWork) prepareData(nonce int64) []byte {
	var data []byte
	heightBytes := IntToHex(pow.Block.Height)
	timeBytes := IntToHex(pow.Block.TimeStamp)
	data = bytes.Join([][]byte{
		timeBytes,
		pow.Block.PreBlockHash,
		heightBytes,
		pow.Block.HashTransaction(),
		IntToHex(nonce),
		IntToHex(targetBit),
	}, []byte{})

	return data
}
