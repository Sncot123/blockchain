package BLC

import (
	"bytes"
	"crypto/elliptic"
	"encoding/gob"
	"errors"
	"io/fs"
	"io/ioutil"
	"log"
	"os"
)

//钱包集合管理的文件

const walletsFile = "wallets.dat"

//实现钱包集合的基本结构

type Wallets struct {
	//key:地址
	//value:钱包结构
	Wallets map[string]*Wallet
}

// 初始化钱包集合
func NewWallets() *Wallets {
	//从钱包文件中获取钱包信息
	//1、先判断文件是否存在
	if _, err := os.Stat(walletsFile); !errors.Is(fs.ErrExist, err) {
		wallets := &Wallets{}
		wallets.Wallets = make(map[string]*Wallet)
		return wallets
	}
	//2、文件存在，读取内容
	fileContent, err := ioutil.ReadFile(walletsFile)
	if nil != err {
		log.Panicf("read the file content failed! err:%v\n", err)
	}

	var wallets Wallets
	gob.Register(elliptic.P256())
	decoder := gob.NewDecoder(bytes.NewReader(fileContent))
	err = decoder.Decode(&wallets)
	if nil != err {
		log.Panicf("decoder the file content failed! err:%v\n", err)
	}
	return &wallets
}

func (wallets *Wallets) CreateWallets() {
	//1、创建钱包
	wallet := NewWallet()
	//2、添加
	wallets.Wallets[string(wallet.GetAddress())] = wallet
	//3、持久化钱包信息
	wallets.SaveWallets()
}

//持久化钱包信息（存储到文件中）
func (w *Wallets) SaveWallets() {
	var content bytes.Buffer //钱包内容
	//注册256椭圆，注册之后，可以在内部对curve的接口进行编码
	gob.Register(elliptic.P256())
	encoder := gob.NewEncoder(&content)
	err := encoder.Encode(&w)
	if nil != err {
		log.Panicf("encoder the struct of wallets failed %v\n", err)
	}
	err = ioutil.WriteFile(walletsFile, content.Bytes(), 0644)
	if nil != err {
		log.Panicf("write the content of wallet into file %s failed! err:%v\n ", walletsFile, err)

	}

}
