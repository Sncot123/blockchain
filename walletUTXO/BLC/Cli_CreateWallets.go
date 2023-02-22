package BLC

//创建钱包
func (cli *CLI) CreateWallets() {
	wallets := NewWallets()
	wallets.CreateWallets()
}
