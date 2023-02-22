package BLC

import (
	"bytes"
	"crypto/sha256"
	"encoding/gob"
	"encoding/hex"
	"fmt"
	"log"
)

//交易管理文件

//定义一个交易结构体

type Transaction struct {
	//交易hash
	TxHash []byte
	//输入列表
	Vins []*TxInput
	//输出列表
	Vouts []*TxOutput
}

func NewCoinbaseTransaction(address string) *Transaction {

	//输入
	//coinbase特点
	//txHash:nil
	//vout:-1(用于判断当前交易是否为coinbase交易)
	//ScriptSig:系统奖励
	txInput := &TxInput{[]byte{}, -1, "system reward"}
	//输出
	//value
	//address
	txOutput := &TxOutput{10, address}
	//输入输出组装交易
	txCoinbase := &Transaction{
		nil,
		[]*TxInput{txInput},
		[]*TxOutput{txOutput},
	}
	txCoinbase.HashTransaction()
	return txCoinbase
}

//生成交易hash
func (tx *Transaction) HashTransaction() {
	var result bytes.Buffer
	//设置编码对象
	encoder := gob.NewEncoder(&result)
	if err := encoder.Encode(tx); err != nil {
		log.Panicf("tx hash encoder failed  err:%v\n", err)
	}
	//生成hash值
	hash := sha256.Sum256(result.Bytes())
	tx.TxHash = hash[:]
}

//生成普通的交易
func NewSimpleTransaction(from, to string, amount int, bc *BlockChain, txs []*Transaction) *Transaction {
	var txInputs []*TxInput   //输入列表
	var txOutputs []*TxOutput //输出列表
	//带哦用可花费UTXO函数
	money, spendableUTXODic := bc.FindSpendableUTXO(from, amount, txs)
	fmt.Printf("money:%v\n", money)
	//输入
	for txHash, indexArray := range spendableUTXODic {
		txHashesBytes, err := hex.DecodeString(txHash)
		if err != nil {
			log.Panicf("decode string to []byte failed err:%v\n", err)
		}
		//遍历索引列表
		for _, index := range indexArray {
			txInput := &TxInput{txHashesBytes, index, from}
			txInputs = append(txInputs, txInput)
		}
	}
	//生成新的交易

	//输入
	//txInput := &TxInput{[]byte("00007443c8e0ef0b2ace632c7c7a7a56d1974f0d08f3fe3d15f091923dc3e46a"), 0, from}
	//txInputs = append(txInputs, txInput)
	//输出
	txOutput := &TxOutput{amount, to}
	txOutputs = append(txOutputs, txOutput)
	//输出找零
	if money < amount {
		txOutput = &TxOutput{amount - amount, from}
		txOutputs = append(txOutputs, txOutput)
	} else {
		log.Panicf("余额不足。。。\n")
	}

	tx := Transaction{nil, txInputs, txOutputs}
	tx.HashTransaction()
	return &tx

}

//判断指定交易是否是coinbase交易
func (tx *Transaction) IsCoinbaseTransaction() bool {
	return tx.Vins[0].Vout == -1 && len(tx.Vins[0].TxHash) == 0
}
