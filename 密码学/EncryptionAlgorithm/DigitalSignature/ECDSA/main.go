package main

import (
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/sha256"
	"fmt"
	"math/big"
)

//椭圆加密算法

//生成密钥对
func GenECDSA() (*ecdsa.PrivateKey, *ecdsa.PublicKey) {
	privateKey, err := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
	if nil != err {
		panic(err)
	}
	return privateKey, &privateKey.PublicKey
}
func EncryptECDSA(priKey *ecdsa.PrivateKey, src string) (r, s *big.Int) {
	//对信息进行hash
	bytes := sha256.Sum256([]byte(src))
	//对散列值进行加密
	r, s, err := ecdsa.Sign(rand.Reader, priKey, bytes[:])
	if nil != err {
		panic(err)
	}
	return r, s
}

func VerifyECDSA(pubKey *ecdsa.PublicKey, src string, r, s *big.Int) bool {
	//对消息进行hash
	bytes := sha256.Sum256([]byte(src))
	//验证签名
	return ecdsa.Verify(pubKey, bytes[:], r, s)

}
func main() {
	msg := "你好！XXX"
	priKey, publicKey := GenECDSA()
	r, s := EncryptECDSA(priKey, msg)
	verifyECDSA := VerifyECDSA(publicKey, msg, r, s)
	if verifyECDSA {
		fmt.Println("验证通过")
	} else {
		fmt.Println("验证失败")
	}
}
