package main

import (
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/sha256"
	"fmt"
	"log"
)

func main() {

	//签名
	key, err := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
	if nil != err {
		log.Panicf("get privaeKey failed! err:%v\n", err)
	}
	msg := "this is a message"
	hash := sha256.Sum256([]byte(msg))
	sign, b, err := ecdsa.Sign(rand.Reader, key, hash[:])
	if nil != err {
		panic(err)
	}
	fmt.Printf("signature:(0x%x,0x%x\n", sign, b)

	//验证
	valid := ecdsa.Verify(&key.PublicKey, hash[:], sign, b)
	fmt.Printf("验证:%v\n", valid)
}
