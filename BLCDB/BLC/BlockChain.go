package BLC

import (
	"log"

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

func initChain() *Block {
	block := NewBlock(1, nil, []byte("创世块"))
	return block
}

// CreateBlockChain 初始化区块链
func CreateBlockChain() *BlockChain {
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
			//生成创世区块
			block := initChain()
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
func (c *BlockChain) AddBlock(data string) {
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
			newBlock := NewBlock(lastBlock.Height+1, lastBlock.Hash, []byte(data))
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
