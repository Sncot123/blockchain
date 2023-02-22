package BLC

//交易的输入

type TxInput struct {
	//不是指当前交易的hash
	TxHash []byte
	//引用的上一笔交易的输出索引号
	Vout int
	//签名
	ScriptSig string
}

//验证引用的地址是否符合条件
func (txInput *TxInput) CheckPubKeyWithAddress(address string) bool {
	return txInput.ScriptSig == address
}
