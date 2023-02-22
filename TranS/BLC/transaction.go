package BLC

import (
	"bytes"
	"crypto/sha256"
	"encoding/gob"
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
func NewSimpleTransaction(from, to string, amount int) *Transaction {
	var txInputs []*TxInput   //输入列表
	var txOutputs []*TxOutput //输出列表
	//生成新的交易

	//输入
	txInput := &TxInput{[]byte("00007443c8e0ef0b2ace632c7c7a7a56d1974f0d08f3fe3d15f091923dc3e46a"), 0, from}
	txInputs = append(txInputs, txInput)
	//输出
	txOutput := &TxOutput{amount, to}
	txOutputs = append(txOutputs, txOutput)
	//输出找零
	if amount < 10 {
		txOutput = &TxOutput{10 - amount, from}
		txOutputs = append(txOutputs, txOutput)
	}

	tx := Transaction{nil, txInputs, txOutputs}
	tx.HashTransaction()
	return &tx

}
