package BLC

//创建区块链
func (cli *CLI) createBlockChain(address string) {
	//fmt.Println(dbExist())
	//if dbExist() {
	//	fmt.Println("创世区块已存在...")
	//	os.Exit(1)
	//}
	CreateBlockChain(address)
}
