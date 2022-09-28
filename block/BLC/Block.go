package BLC

import (
	"bytes"
	"crypto/sha256"
	"time"
)

//区块基本结构与功能文件

// 实现一个基本的区块结构
type Block struct {
	TimeStamp    int64  //区块时间
	Hash         []byte //当前区块hash
	PreBlockHash []byte //上一个区块的hash
	Height       int64  //作为区块的高度
	Data         []byte //数据
	Nonce        int64  //在运行pow时生成的哈希变化值，也代表pow运行时动态修改的数据
}

func NewBlock(height int64, preBlockHash, data []byte) *Block {
	block := &Block{
		TimeStamp:    time.Now().Unix(),
		PreBlockHash: preBlockHash,
		Height:       height,
		Data:         data,
	}
	//block.setHash()
	pow := NewProofOfWork(block)
	hash, nonce := pow.Run()
	block.Nonce = nonce
	block.Hash = hash
	return block
}

// setHash 调用生成当前区块的Hash并放入区块
func (b *Block) setHash() {
	heightBytes := IntToHex(b.Height)
	timeBytes := IntToHex(b.TimeStamp)
	blockBytes := bytes.Join([][]byte{
		timeBytes,
		b.PreBlockHash,
		heightBytes,
		b.Data,
	}, []byte{})
	currentHash := sha256.Sum256(blockBytes)
	b.Hash = currentHash[:] //int32转换成int64
}
