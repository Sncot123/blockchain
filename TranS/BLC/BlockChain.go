package BLC

import (
	"log"
	"strconv"

	"github.com/boltdb/bolt"
)

// 区块链

const (
	dbName         = "block.db"
	blockTableName = "blocks"
)

type BlockChain struct {
	//Blocks []*Block
	DB  *bolt.DB //数据库对象
	Tip []byte   //保存最新区块的hash
}

func initChain(tx *Transaction) *Block {
	block := NewBlock(1, nil, []*Transaction{tx})
	return block
}

// CreateBlockChain 初始化区块链
func CreateBlockChain(address string) *BlockChain {
	//保存最新区块的hash
	var latestBlockHash []byte
	//1、打开数据库
	db, err := bolt.Open(dbName, 0600, nil)
	if err != nil {
		log.Panicf("create db [%s] failed  err:%v\n ", db, err)
	}
	//2、创建桶
	err = db.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(blockTableName))
		if b == nil {
			//没找到桶
			b, err := tx.CreateBucket([]byte(blockTableName))
			if err != nil {
				log.Panicf("create bucket [%s] failed  err:%v\n", blockTableName, err)
			}
			//生成一个coinbase交易
			tx := NewCoinbaseTransaction(address)
			//生成创世区块
			block := initChain(tx)
			//存储
			//1、key、value分别以什么样的数据存入
			//2、如何把block结构存入数据库中
			err = b.Put(block.Hash, block.Serialize())
			if err != nil {
				log.Panicf("insert the genesis block failed  err:%v\n ", err)
			}
			latestBlockHash = block.Hash
			//存储最新区块的hash
			//l:latestHash
			err = b.Put([]byte("l"), block.Hash)
			if err != nil {
				log.Panicf(" save the hash of genesis block failed err:%v\n", err)
			}
		}
		if err != nil {
			log.Panicf("init bolckChain failed err:%v\n", err)
		}
		return nil
	})
	//3、把创世区块存入数据库中
	return &BlockChain{DB: db, Tip: latestBlockHash}
}

// AddBlock 向区块链添加区块
func (c *BlockChain) AddBlock(tr *Transaction) {
	//更新区块数据（insert）
	//给当前区块添加上一区块的hash
	err := c.DB.Update(func(tx *bolt.Tx) error {
		//1、获取数据库桶
		b := tx.Bucket([]byte(blockTableName))
		if b != nil {
			//2、获取最后插入区块的hash
			lastBlockBytes := b.Get(c.Tip)
			//3、反序列化获取对应区块
			lastBlock := DeSerialize(lastBlockBytes)
			//4、创建新的区块
			newBlock := NewBlock(lastBlock.Height+1, lastBlock.Hash, []*Transaction{tr})
			//5、更新Tip
			c.Tip = newBlock.Hash
			//6、序列化新区块并存入数据库，并更新数据库中的最新区块的hash
			newBlockBytes := newBlock.Serialize()
			err := b.Put([]byte("l"), newBlock.Hash)
			if err != nil {
				log.Panicf("update the lastest hash of block failed  err:%v\n", err)
			}
			err = b.Put(newBlock.Hash, newBlockBytes)
			if err != nil {
				log.Panicf("insert the newBlock failed  err:%v\n", err)
			}
		}
		return nil
	})
	if err != nil {
		log.Panicf("insert newBlock failederr:%v\n", err)
	}
}

//实现挖矿功能
//实现通过接受交易，生成区块
func (bc *BlockChain) MineNewBlock(from, to, amount []string) {
	//搁置交易生成步骤
	var txs []*Transaction
	var block *Block
	value, _ := strconv.Atoi(amount[0])
	//生成新的交易
	tx := NewSimpleTransaction(from[0], to[0], value)
	txs = append(txs, tx)
	bc.DB.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(blockTableName))
		if b != nil {
			lHash := b.Get([]byte("l"))
			blockBytes := b.Get(lHash)
			block = DeSerialize(blockBytes)
		}
		return nil
	})
	//通过数据库中最新的区块去生成新的区块
	block = NewBlock(block.Height+1, block.Hash, txs)
	//持久化新生成的区块的到数据库中
	bc.DB.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(blockTableName))
		if b != nil {
			err := b.Put(block.Hash, block.Serialize())
			if err != nil {
				log.Panicf("update the latest block hash to db failed err:%v\n", err)
			}
			err = b.Put([]byte("l"), block.Hash)
			if err != nil {
				log.Panicf("update the latest hash to db failed err:%v\n", err)
			}
		}
		return nil
	})
	bc.Tip = block.Hash
}
