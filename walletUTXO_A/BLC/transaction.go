package BLC

import (
	"bytes"
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/sha256"
	"encoding/gob"
	"encoding/hex"
	"fmt"
	"log"
	"math/big"
	"time"
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
	txInput := &TxInput{[]byte{}, -1, nil, nil}
	//输出
	//value
	//address
	txOutput := NewTxOutput(10, address)
	//输入输出组装交易
	txCoinbase := &Transaction{
		nil,
		[]*TxInput{txInput},
		[]*TxOutput{txOutput},
	}
	txCoinbase.HashTransaction()
	return txCoinbase
}

//生成交易hash(交易序列化)
//不同时间生成的交易hash值不同
func (tx *Transaction) HashTransaction() {
	var result bytes.Buffer
	//设置编码对象
	encoder := gob.NewEncoder(&result)
	if err := encoder.Encode(tx); err != nil {
		log.Panicf("tx hash encoder failed  err:%v\n", err)
	}
	//添加时间戳标识，不添加会导致所有的coinbase交易的hash相同
	tm := time.Now().UnixNano()
	//生成hash值
	txHashBytes := bytes.Join([][]byte{result.Bytes(), IntToHex(tm)}, []byte{})
	hash := sha256.Sum256(txHashBytes)
	tx.TxHash = hash[:]
}

//生成普通的交易
func NewSimpleTransaction(from, to string, amount int, bc *BlockChain, txs []*Transaction) *Transaction {
	var txInputs []*TxInput   //输入列表
	var txOutputs []*TxOutput //输出列表
	//带哦用可花费UTXO函数
	money, spendableUTXODic := bc.FindSpendableUTXO(from, amount, txs)
	fmt.Printf("money:%v\n", money)
	//获取钱包集合对象
	wallets := NewWallets()
	//查找对应钱包结构
	wallet := wallets.Wallets[from]
	//输入
	for txHash, indexArray := range spendableUTXODic {
		txHashesBytes, err := hex.DecodeString(txHash)
		if err != nil {
			log.Panicf("decode string to []byte failed err:%v\n", err)
		}
		//遍历索引列表
		for _, index := range indexArray {
			txInput := &TxInput{txHashesBytes, index, nil, wallet.PublicKey}
			txInputs = append(txInputs, txInput)
		}
	}
	//生成新的交易

	//输入
	//txInput := &TxInput{[]byte("00007443c8e0ef0b2ace632c7c7a7a56d1974f0d08f3fe3d15f091923dc3e46a"), 0, from}
	//txInputs = append(txInputs, txInput)
	//输出
	//txOutput := &TxOutput{amount, to}
	txOutput := NewTxOutput(amount, to)
	txOutputs = append(txOutputs, txOutput)
	//输出找零
	if money < amount {
		//txOutput = &TxOutput{money - amount, from}
		txOutput = NewTxOutput(money-amount, from)
		txOutputs = append(txOutputs, txOutput)
	} else {
		log.Panicf("余额不足。。。\n")
	}

	tx := Transaction{nil, txInputs, txOutputs}
	tx.HashTransaction() //生成一笔完整的交易

	//对交易进行签名
	bc.SignTransaction(&tx, wallet.PrivateKey)
	return &tx

}

//判断指定交易是否是coinbase交易
func (tx *Transaction) IsCoinbaseTransaction() bool {
	return tx.Vins[0].Vout == -1 && len(tx.Vins[0].TxHash) == 0
}

// 交易签名
//prevTxs:代表当前交易的输入所引用的所有output所属的交易
func (tx *Transaction) Sign(privateKey ecdsa.PrivateKey, prevTxs map[string]Transaction) {
	//处理输入，保证交易的正确性
	//检查每个tx中每一个输入所引用的交易哈希是否包含在prevTxs中
	//如果没有包含在里面，说明该交易被人修改了
	for _, vin := range tx.Vins {
		if prevTxs[hex.EncodeToString(vin.TxHash)].TxHash == nil {
			log.Panicf("error:prev transaction is not correct!\n")
		}
	}
	//提取需要签名的属性
	txCopy := tx.TrimmedCopy()
	//出阿里交易副本的输入
	for vin_id, vin := range txCopy.Vins {
		//获取关联交易
		prevTX := prevTxs[hex.EncodeToString(vin.TxHash)]
		//找到发送者（当前输入引用的hash-输出的hash）
		txCopy.Vins[vin_id].PublicKey = prevTX.Vouts[vin_id].Ripemd160Hash
		//生成交易副本的hash
		txCopy.TxHash = txCopy.Hash()
		//调用核心签名函数
		r, s, err := ecdsa.Sign(rand.Reader, &privateKey, txCopy.TxHash)
		if nil != err {
			log.Panicf("sign to transaction failed! err:%v\n ", err)
		}
		//组成交易签名
		signature := append(r.Bytes(), s.Bytes()...)
		tx.Vins[vin_id].Signature = signature
	}
}

//交易copy，生成一个专门用于交易的副本
func (tx *Transaction) TrimmedCopy() Transaction {
	//重新组装生成一个新的交易
	var inputs []*TxInput
	var outputs []*TxOutput
	//组装input
	for _, vin := range tx.Vins {
		inputs = append(inputs, &TxInput{vin.TxHash, vin.Vout, nil, nil})

	}
	//组装outputs
	for _, vout := range tx.Vouts {

		outputs = append(outputs, &TxOutput{vout.value, vout.Ripemd160Hash})

	}
	txCopy := Transaction{tx.TxHash, inputs, outputs}
	return txCopy
}

// 设置用于签名交易的hash
func (tx *Transaction) Hash() []byte {
	txCopy := tx
	txCopy.TxHash = []byte{}
	hash := sha256.Sum256(tx.Serialize())
	return hash[:]
}

//交易序列化
func (tx *Transaction) Serialize() []byte {
	var buffer bytes.Buffer
	//新建编码对象
	encoder := gob.NewEncoder(&buffer)
	//编码(序列化)
	if err := encoder.Encode(tx); nil != err {
		log.Panicf("serialize the tx to []byte failed! err:%v\n", err)

	}
	return buffer.Bytes()
}

//验证签名
func (tx *Transaction) Verify(prevTxs map[string]Transaction) bool {
	//检查能否找到交易hash
	for _, vin := range tx.Vins {
		if prevTxs[hex.EncodeToString(vin.TxHash)].TxHash == nil {
			log.Panicf("verify error: transaction verify failed!\n")
		}
	}
	//提取相同的交易签名属性
	txCopy := tx.TrimmedCopy()
	//使用相同的椭圆
	curve := elliptic.P256()
	//遍历tx输入，对每一笔输入所引用的输出进行验证
	for vinId, vin := range tx.Vins {
		prevTx := prevTxs[hex.EncodeToString(vin.TxHash)]
		//找到发送者（当前输入引用的hash-输出的hash）
		txCopy.Vins[vinId].PublicKey = prevTx.Vouts[vin.Vout].Ripemd160Hash
		//有需要验证的数据生成交易hash，必须要逾千名是的数据完全一致
		txCopy.TxHash = txCopy.Hash()
		//在比特币中，签名是一个数值对，r,s代表签名
		//从要输入的signature中获取
		//获取r,s   注意：r，s的长度是相等的
		r := big.Int{}
		s := big.Int{}
		sigLen := len(vin.Signature)
		r.SetBytes(vin.Signature[:sigLen/2])
		s.SetBytes(vin.Signature[sigLen/2:])
		// 获取公钥
		//公钥是由X,Y的坐标组成的
		x := big.Int{}
		y := big.Int{}
		pubKeyLen := len(vin.PublicKey)
		x.SetBytes(vin.PublicKey[:pubKeyLen/2])
		y.SetBytes(vin.PublicKey[pubKeyLen/2:])
		rawPublicKey := ecdsa.PublicKey{curve, &x, &y}
		if !ecdsa.Verify(&rawPublicKey, txCopy.TxHash, &r, &s) {
			return false
		}

	}
	//验证签名核心函数
	return true
}
