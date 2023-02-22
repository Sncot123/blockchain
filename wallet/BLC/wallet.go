package BLC

import (
	"bytes"
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/sha256"
	"fmt"
	"log"
	"math/big"

	"golang.org/x/crypto/ripemd160"
)

//校验和
const addressCheckSumLen = 4

type Wallet struct {
	//私钥
	PrivateKey ecdsa.PrivateKey
	//公钥
	PublicKey []byte
}
type PublicKey struct {
	elliptic.Curve
	X, Y *big.Int
}

type PrivateKey struct {
	PublicKey
	D *big.Int
}

func NewWallet() *Wallet {
	// 公钥-私钥赋值
	privateKey, PublicKey := newKeyPair()
	return &Wallet{privateKey, PublicKey}
}

// 通过钱包生成公钥-私钥对
func newKeyPair() (ecdsa.PrivateKey, []byte) {
	// 1、获取一个椭圆
	curve := elliptic.P256()
	// 2、通过椭圆相关算法生成私钥
	priv, err := ecdsa.GenerateKey(curve, rand.Reader)
	if nil != err {
		log.Panicf("general privateKry failed   err:%v\n", err)
	}
	// 3、通过私钥生成公钥
	pubKey := append(priv.PublicKey.X.Bytes(), priv.PublicKey.Y.Bytes()...)
	return *priv, pubKey
}

//实现双hash
func Ripemd160Hash(pubkey []byte) []byte {
	//1、sha256
	hash256 := sha256.New()
	hash256.Write(pubkey)
	hash := hash256.Sum(nil)

	//2、ripemd160
	rmd160 := ripemd160.New()
	rmd160.Write(hash)
	return rmd160.Sum(nil)

}

//生成校验和
func CheckSum(input []byte) []byte {
	first_hash := sha256.Sum256(input)
	second_hash := sha256.Sum256(first_hash[:])
	return second_hash[:addressCheckSumLen]
}

//通过钱包（公钥）获取地址
func (wallet *Wallet) GetAddress() []byte {
	hash := Ripemd160Hash(wallet.PublicKey)
	sum := CheckSum(hash)
	address := append(hash, sum...)
	b58Bytes := Base58Encode(address)
	fmt.Println(b58Bytes, "----")
	return b58Bytes
}

// 判断地址的有效性
func IsValidForAddress(address []byte) bool {
	//1、地址通过base58decode解码（长度为24）
	sumHash := Base58Decode(address)
	//2、拆分进行校验
	checkSumBytes := sumHash[len(sumHash)-addressCheckSumLen:]
	pubkeyBytes := sumHash[:len(sumHash)-addressCheckSumLen]

	checkS := CheckSum(pubkeyBytes)
	if bytes.Compare(checkSumBytes, checkS) == 0 {
		return true
	}
	return false
}
