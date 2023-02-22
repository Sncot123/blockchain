package main

import (
	"crypto/cipher"
	"crypto/des"
	"encoding/hex"
	"fmt"
)

func EncryptDES(src, key []byte) []byte {
	//创建加密模式
	block, err := des.NewCipher(key)
	if nil != err {
		panic(err)
	}

	//初始化向量
	iv := []byte("12345678")
	blockmode := cipher.NewCFBEncrypter(block, iv)

	//加密
	dst := make([]byte, len(src))
	blockmode.XORKeyStream(dst, src)
	return dst
}

func DecryptDES(src, key []byte) []byte {
	//创建解密模式
	block, err := des.NewCipher(key)
	if nil != err {
		panic(err)
	}
	//初始化向量
	iv := []byte("12345678")
	blockmode := cipher.NewCFBDecrypter(block, iv)
	dst := make([]byte, len(src))
	blockmode.XORKeyStream(dst, src)
	return dst
}

func main() {
	msg := []byte("1223454")
	key := []byte("12345678")
	encryptDES := EncryptDES(msg, key)
	fmt.Println(hex.EncodeToString(encryptDES))

	decryptDES := DecryptDES(encryptDES, key)
	fmt.Println(string(decryptDES))
}
