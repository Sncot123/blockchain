package BLC

// 区块链
type BlockChain struct {
	Blocks []*Block
}

func initChain() *Block {
	block := NewBlock(1, nil, []byte("创世块"))
	return block
}

// CreateBlockChain 初始化区块链
func CreateBlockChain() *BlockChain {
	block := initChain()
	chain := &BlockChain{
		[]*Block{block},
	}
	return chain
}

// AddBlock 向区块链添加区块
func (c *BlockChain) AddBlock(height int64, data []byte) {
	//给当前区块添加上一区块的hash
	block := NewBlock(height, c.Blocks[len(c.Blocks)-1].Hash, data)
	c.Blocks = append(c.Blocks, block)
}
