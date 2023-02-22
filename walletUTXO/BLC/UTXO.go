package BLC

type UTXO struct {
	//UTXO对应的交易hash
	TxHash []byte
	//UTXO在所属交易的输出列表中的索引
	Index int
	//Output本身
	Output *TxOutput
}
