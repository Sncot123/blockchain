package BLC

//交易的输入

type TxInput struct {
	//不是指当前交易的hash
	TxHash []byte
	//
	Vout int

	ScriptSig string
}

//验证引用的地址是否符合条件
func (txInput *TxInput) CheckPubKeyWithAddress(address string) bool {
	return txInput.ScriptSig == address
}
