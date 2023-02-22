package BLC

import (
	"bytes"
	"crypto/ecdsa"
	"encoding/hex"
	"fmt"
	"log"
	"math/big"
	"os"
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
	// 遍历交易的参与者
	for index, address := range from {
		value, _ := strconv.Atoi(amount[index])
		//生成新的交易
		tx := NewSimpleTransaction(address, to[index], value, bc, txs)
		//追加到txs的交易列表中去
		txs = append(txs, tx)
		tx = NewCoinbaseTransaction(address)
		txs = append(txs, tx)
	}

	bc.DB.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(blockTableName))
		if b != nil {
			lHash := b.Get([]byte("l"))
			blockBytes := b.Get(lHash)
			block = DeSerialize(blockBytes)
		}
		return nil
	})
	//在此处进行交易签名的验证
	//对txs中的每一笔交易都进行验证
	for _, tx := range txs {
		//验证签名，只要有一笔签名验证失败就panic
		bc.VerifyTransaction(tx)
	}
	//通过数据库中最新的区块去生成新的区块(交易的打包)
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

//获取所有已花费输出
func (bc *BlockChain) SpentOutputs(address string) map[string][]int {
	//已花费输出缓存
	spentTXoutputs := make(map[string][]int)
	//获取迭代器对象
	bcit := bc.Iterate()
	for {
		block := bcit.Next()
		for _, tx := range block.Txh {
			if !tx.IsCoinbaseTransaction() {
				for _, in := range tx.Vins {
					if in.UnLockRipemd160Hash(StringToHash160(address)) {
						key := hex.EncodeToString(in.TxHash)
						//添加到已花费输出的缓存中
						spentTXoutputs[key] = append(spentTXoutputs[key], in.Vout)
					}
					//if in.CheckPubKeyWithAddress(address) {
					//	key := hex.EncodeToString(in.TxHash)
					//	//添加到已花费输出的缓存中
					//	spentTXoutputs[key] = append(spentTXoutputs[key], in.Vout)
					//}
				}
			}
		}
		//退出循环条件
		var hashInt big.Int
		hashInt.SetBytes(block.PreBlockHash)
		if hashInt.Cmp(big.NewInt(0)) == 0 {
			break
		}
	}
	return spentTXoutputs
}

//查询未花费的utxo
func (bc *BlockChain) UnUTXOS(from string, txs []*Transaction) []*UTXO {
	//1、遍历数据库，查找与address相关的交易
	//获取迭代器
	bcit := bc.Iterate()
	//
	var unTxoutput []*UTXO
	//获取指定地址的已花费输出
	spentTXOutputs := bc.SpentOutputs(from)
	//缓存迭代
	//查找缓存中的已花费输出
	for _, tx := range txs {
		//添加一个缓存输出跳转
	WorkCacheTx:
		for index, vout := range tx.Vouts {
			if vout.UnLockScriptPubkeyWithAddress(from) {
				//if vout.CheckPubKeyWithAddress(from) {
				if len(spentTXOutputs) != 0 {
					var isUtxoTx bool //判断交易是否被其他交易所引用
					for txHash, indexArray := range spentTXOutputs {
						txHashStr := hex.EncodeToString(tx.TxHash)
						if txHash == txHashStr {
							//当前遍历到的交易已经有输出被其他交易的输入所引用
							isUtxoTx = true
							//添加状态变量，判断指定的output是否被引用
							var isSpentUTXO bool
							for _, voutIndex := range indexArray {
								if index == voutIndex {
									//该输出别引用了
									isSpentUTXO = true
									//跳出当前vout判断逻辑，进行下一个输出判断
									continue WorkCacheTx
								}
							}
							if isSpentUTXO == false {
								utxo := &UTXO{tx.TxHash, index, vout}
								unTxoutput = append(unTxoutput, utxo)
							}
						}
					}
					if isUtxoTx == false {
						//说并当前交易中所有与address有关的outputs都是utxo
						utxo := &UTXO{tx.TxHash, index, vout}
						unTxoutput = append(unTxoutput, utxo)
					}
				} else {
					utxo := &UTXO{tx.TxHash, index, vout}
					unTxoutput = append(unTxoutput, utxo)
				}
			}
		}

	}
	//优先遍历缓存中的utxo，如果余额足够，直接返回，否则在遍历数据库中的utxo
	//
	//数据库迭代，不断获取下一个区块
	for {
		block := bcit.Next()
		//遍历区块链中的每笔交易
		for _, tx := range block.Txh {
			//跳转
		work:
			for index, vout := range tx.Vouts {
				//index 输出当前交易中的索引位置
				//vout 当前输出
				if vout.UnLockScriptPubkeyWithAddress(from) {
					//if vout.CheckPubKeyWithAddress(from) {
					//当前vout属于传入地址
					if len(spentTXOutputs) != 0 {
						var isSpentOutput bool
						for txHash, indexArray := range spentTXOutputs {
							for _, i := range indexArray {
								//txHash当前输出所引用的交易hash
								//indexArray hash关联的vout索引列表
								if txHash == hex.EncodeToString(tx.TxHash) && index == i {
									//txHash==hex.EncodeToString(tx.TxHash)说明当前的交易tx已经有输出被其他交易所引用
									//index==i 说明正好是当前的输出被其他交易引用
									isSpentOutput = true
									continue work
								}
							}
						}
						/*
							type UTXO struct {
								//UTXO对应的交易hash
								TxHash []byte
								//UTXO在所属交易的输出列表中的索引
								Index int
								//Output本身
								Output *TxOutput
							}



						*/
						if isSpentOutput == false {
							utxo := &UTXO{tx.TxHash, index, vout}
							unTxoutput = append(unTxoutput, utxo)
						}
					} else {
						//将当前所有输出都添加到未花费输出中
						utxo := &UTXO{tx.TxHash, index, vout}
						unTxoutput = append(unTxoutput, utxo)
					}
				}
			}
		}
		//退出循环条件
		var hashInt big.Int
		hashInt.SetBytes(block.PreBlockHash)
		if hashInt.Cmp(big.NewInt(0)) == 0 {
			break
		}
	}
	return nil
}

//查询余额
func (bc *BlockChain) getbalance(address string) int {
	var amount int //余额
	utxos := bc.UnUTXOS(address, []*Transaction{})
	for _, utxo := range utxos {
		amount += utxo.Output.value
	}
	return amount
}

//查找指定地址的可用utxo 超过amount就中断查找
// 更新当前数据库中指定地址的utxo数量
// txs:缓存中的交易列表（用于多笔交易处理）
func (bc *BlockChain) FindSpendableUTXO(from string, amount int, txs []*Transaction) (int, map[string][]int) {
	//可用的UTXO
	spendableUTXO := make(map[string][]int)

	var value int
	utxos := bc.UnUTXOS(from, txs)
	//遍历utxo
	for _, utxo := range utxos {
		//计算交易hash
		value += utxo.Output.value
		hash := hex.EncodeToString(utxo.TxHash)
		spendableUTXO[hash] = append(spendableUTXO[hash], utxo.Index)
		if value >= amount {
			break
		}
	}
	//所有都遍历完，仍然小于amount
	//资金不足
	if value < amount {
		fmt.Printf("地址[%s]余额不足，当前余额[%d]，转账金额[%d]\n", from, value, amount)
		os.Exit(1)
	}
	return value, spendableUTXO
}

//通过指定的交易hash查找交易
func (bc *BlockChain) FindTransaction(ID []byte) Transaction {
	bcit := bc.Iterate()
	for {
		block := bcit.Next()
		for _, tx := range block.Txh {
			if bytes.Compare(ID, tx.TxHash) == 0 {
				//没找到该交易
				return *tx
			}
		}
		//退出
		var hashInt big.Int
		hashInt.SetBytes(block.PreBlockHash)
		if big.NewInt(0).Cmp(&hashInt) == 0 {
			break
		}
	}
	fmt.Printf("没找到交易[%x]! err:%v\n", ID)
	return Transaction{}
}

//交易签名
func (bc *BlockChain) SignTransaction(tx *Transaction, privateKey ecdsa.PrivateKey) {
	//coinbase交易不需签名
	if tx.IsCoinbaseTransaction() {
		return
	}
	//处理交易的input，查找tx中input所引用的vout所属交易（查找发送者）
	//对我们所花费的每一笔utxo进行签名
	//存储引用的交易
	prevTxs := make(map[string]Transaction)
	for _, vin := range tx.Vins {
		//查找当前交易输入所引用的交易
		//vin.TxHash
		tx := bc.FindTransaction(vin.TxHash)
		prevTxs[hex.EncodeToString(tx.TxHash)] = tx

	}
	//签名
	tx.Sign(privateKey, prevTxs)

}

//验证签名
func (bc *BlockChain) VerifyTransaction(tx *Transaction) bool {
	prevTxs := make(map[string]Transaction)
	//查找输入所引用的交易
	for _, vin := range tx.Vins {
		tx := bc.FindTransaction(vin.TxHash)
		prevTxs[hex.EncodeToString(tx.TxHash)] = tx
	}
	return tx.Verify(prevTxs)
}
