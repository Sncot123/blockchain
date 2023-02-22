package BLC

import "bytes"

//交易的输入

type TxInput struct {
	//不是指当前交易的hash
	TxHash []byte
	//引用的上一笔交易的输出索引号
	Vout int
	//数字签名
	Signature []byte
	//公钥
	PublicKey []byte
}

//验证引用的地址是否符合条件
//func (txInput *TxInput) CheckPubKeyWithAddress(address string) bool {
//	return txInput.ScriptSig == address
//}

// 传递哈希160进行判断
func (in *TxInput) UnLockRipemd160Hash(ripemd160Hash []byte) bool {
	//获取input的ripemd160哈希值
	inputRipemd160Hash := Ripemd160Hash(in.PublicKey)
	return bytes.Compare(inputRipemd160Hash, ripemd160Hash) == 0
}
