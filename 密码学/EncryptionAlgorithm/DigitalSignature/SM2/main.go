package main

import (
	"crypto/rand"
	"fmt"

	"github.com/tjfoc/gmsm/sm2"
)

//生成密钥对
func GenKeyPair() (*sm2.PrivateKey, *sm2.PublicKey) {
	privateKey, err := sm2.GenerateKey(rand.Reader)
	if nil != err {
		panic(err)
	}
	return privateKey, &privateKey.PublicKey
}

//对消息进行加密
func EncryptSM2(pubKey *sm2.PublicKey, src string) []byte {
	sign, err := pubKey.EncryptAsn1([]byte(src), rand.Reader)
	if nil != err {
		panic(err)
	}
	return sign
}

//对消息进行数字签名验证
func VerifySM2(priKey *sm2.PrivateKey, msg string, sign []byte) bool {
	//验证数字签名的正确性
	asn1, err := priKey.DecryptAsn1([]byte(msg))
	if nil != err {
		fmt.Println(err)
		return false
	} else {
		fmt.Println(string(asn1))
		return true
	}

}

func main() {
	msg := "你好！xxx"
	priKey, publicKey := GenKeyPair()
	sign := EncryptSM2(publicKey, msg)
	verifySM2 := VerifySM2(priKey, msg, sign)
	if verifySM2 {
		fmt.Println("验证通过")
	} else {
		fmt.Println("验证失败")
	}
}

func GenSM2KeyPairFile(priKeyPath, PubKeyPath, pwd string) {
	//privateKey, err := sm2.GenerateKey(rand.Reader)
	//if nil != err {
	//	panic(err)
	//}
	//sm2.W

}
