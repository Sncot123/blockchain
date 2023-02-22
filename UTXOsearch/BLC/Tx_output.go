package BLC

//交易输出
type TxOutput struct {

	//金额
	value        int
	ScriptPubKey string
}

//验证当前utxo是否属于当前指定的地址

func (txot *TxOutput) CheckPubKeyWithAddress(address string) bool {
	return txot.ScriptPubKey == address
}
