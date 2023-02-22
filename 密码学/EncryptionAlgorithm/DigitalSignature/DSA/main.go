package main

import (
	"crypto/dsa"
	"crypto/rand"
	"fmt"
	"math/big"
)

//对消息进行加密生成数字签名
func GenDSA() (*dsa.PrivateKey, *dsa.PublicKey) {
	//生成私钥所需要的参数
	var para dsa.Parameters
	//GenerateParameters的第三个参数决定L、N的长度，长度越长加密程度越高
	dsa.GenerateParameters(&para, rand.Reader, dsa.L1024N160)
	//将初始化好的参数传入给私钥
	var priKey dsa.PrivateKey
	priKey.Parameters = para
	dsa.GenerateKey(&priKey, rand.Reader)
	return &priKey, &priKey.PublicKey
}

//验证数字签名
func EncryptDSA(priKey *dsa.PrivateKey, src []byte) (r, s *big.Int) {
	r, s, err := dsa.Sign(rand.Reader, priKey, src)
	if nil != err {
		panic(err)
	}
	return r, s
}

//验证数字签名的真实性
func VerifyDSA(pubKey *dsa.PublicKey, src []byte, r, s *big.Int) bool {
	return dsa.Verify(pubKey, src, r, s)
}

//对消息本身进行签名
func main() {
	priKey, publicKey := GenDSA()
	msg := []byte("你好！ XXX")
	r, s := EncryptDSA(priKey, msg)
	verifyDSA := VerifyDSA(publicKey, msg, r, s)
	if verifyDSA {
		fmt.Println("验证通过！")
	} else {
		fmt.Println("验证失败！")
	}
}
