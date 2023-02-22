package BLC

import "bytes"

//交易输出
type TxOutput struct {

	//金额
	value int
	//用户名（utxo的所有者）
	//ScriptPubKey string
	//用户名（utxo的所有者）
	Ripemd160Hash []byte
}

//验证当前utxo是否属于当前指定的地址

//func (txot *TxOutput) CheckPubKeyWithAddress(address string) bool {
//	return txot.ScriptPubKey == address
//}

// output身份验证
func (txout *TxOutput) UnLockScriptPubkeyWithAddress(address string) bool {
	hash160 := StringToHash160(address)
	return bytes.Compare(hash160, txout.Ripemd160Hash) == 0
}

func NewTxOutput(value int, address string) *TxOutput {
	txOutput := &TxOutput{}
	hash160 := StringToHash160(address)
	txOutput.value = value
	txOutput.Ripemd160Hash = hash160
	return txOutput
}
