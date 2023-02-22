package BLC

//交易的输入

type TxInput struct {
	//不是指当前交易的hash
	TxHash []byte
	//
	Vout int

	ScriptSig string
}
