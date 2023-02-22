package BLC

import (
	"bytes"
	"crypto/sha256"
	"encoding/gob"
	"log"
	"time"
)

//区块基本结构与功能文件

// 实现一个基本的区块结构
type Block struct {
	TimeStamp    int64  //区块时间
	Hash         []byte //当前区块hash
	PreBlockHash []byte //上一个区块的hash
	Height       int64  //作为区块的高度
	//Data         []byte //数据
	Txh   []*Transaction //切片存储交易
	Nonce int64          //在运行pow时生成的哈希变化值，也代表pow运行时动态修改的数据
}

func NewBlock(height int64, preBlockHash []byte, tr []*Transaction) *Block {
	block := &Block{
		TimeStamp:    time.Now().Unix(),
		PreBlockHash: preBlockHash,
		Height:       height,
		Txh:          tr,
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
		b.HashTransaction(),
	}, []byte{})
	currentHash := sha256.Sum256(blockBytes)
	b.Hash = currentHash[:] //int32转换成int64
}

// 把区块序列化和反序列化  使用内置的gob包
func (b *Block) Serialize() []byte {
	var buffer = new(bytes.Buffer)
	//新建编码对象
	encoder := gob.NewEncoder(buffer)
	//序列化
	if err := encoder.Encode(b); err != nil {
		log.Panicf("Serialize the block to []byte failed  err:%v\n", err)
	}
	return buffer.Bytes()
}
func DeSerialize(data []byte) *Block {
	var block Block
	//新建decoder对象
	decoder := gob.NewDecoder(bytes.NewReader(data))
	if err := decoder.Decode(&block); err != nil {
		log.Panicf("DeSerialize the []byte to block failed err:%v\n", err)
	}
	return &block
}

func (block *Block) HashTransaction() []byte {
	var txHashs [][]byte
	//将指定区块中的所有交易哈希进行拼接
	for _, tx := range block.Txh {
		txHashs = append(txHashs, tx.TxHash)
	}
	txHash := sha256.Sum256(bytes.Join(txHashs, []byte{}))
	return txHash[:]
}
