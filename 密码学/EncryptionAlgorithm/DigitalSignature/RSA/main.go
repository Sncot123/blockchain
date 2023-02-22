package main

import (
	"crypto"
	"crypto/md5"
	"crypto/rand"
	"crypto/rsa"
	"fmt"
)

//生成密钥对
func GenKey() (*rsa.PrivateKey, *rsa.PublicKey) {
	PriKey, err := rsa.GenerateKey(rand.Reader, 2048)
	if nil != err {
		panic(err)
	}
	return PriKey, &PriKey.PublicKey
}

//用私钥和指定的hash函数对消息的散列值进行加密形成数字签名
func EncryptRSA(priKey *rsa.PrivateKey, src string) []byte {
	//对消息进行hash
	hash := md5.New()
	hash.Write([]byte(src))
	bytes := hash.Sum(nil)
	//对散列值进行加密
	opt := &rsa.PSSOptions{SaltLength: rsa.PSSSaltLengthAuto, Hash: crypto.MD5}
	sig, err := rsa.SignPSS(rand.Reader, priKey, crypto.MD5, bytes, opt)
	if nil != err {
		panic(err)
	}
	return sig
}

//验证数字签名的真实性
func VerifyRSA(pubKey *rsa.PublicKey, src string, sig []byte) bool {
	//对消息进行散列处理
	hash := md5.New()
	hash.Write([]byte(src))
	bytes := hash.Sum(nil)
	//对数字签名进行验证
	opt := &rsa.PSSOptions{SaltLength: rsa.PSSSaltLengthAuto, Hash: crypto.MD5}
	err := rsa.VerifyPSS(pubKey, crypto.MD5, bytes, sig, opt)
	if nil != err {
		fmt.Println(err)
		return false
	} else {
		return true
	}
}

//rsa对消息的散列值进行数字签名和验证
func main() {
	priKey, publicKey := GenKey()
	sig := EncryptRSA(priKey, "你好！XXX")
	verifyRSA := VerifyRSA(publicKey, "你好！XXX", sig)
	if verifyRSA {
		fmt.Println("验证通过！")
	} else {
		fmt.Println("验证失败！")
	}
}
