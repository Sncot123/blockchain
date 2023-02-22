package main

import (
	"bytes"
	"crypto/cipher"
	"crypto/des"
	"encoding/hex"
	"fmt"
)

//使用DES加密
//src 待加密明文 ， key 密钥

//加密
func EncryptDES(src, key []byte) []byte {
	//创建cipher.Block接口 其对应的就是一个加密块
	block, err := des.NewCipher(key)
	if nil != err {
		panic(err)
	}
	length := block.BlockSize()
	//填充最后一组数据
	src = PaddingText(src, length)
	//初始化向量
	iv := []byte("12345678")
	//创建cbc加密模式
	blockMode := cipher.NewCBCEncrypter(block, iv)
	//创建切片，用于存储加密后的数据
	dst := make([]byte, len(src))
	blockMode.CryptBlocks(dst, src)
	return dst
}

//解密

func DecryptDES(src, key []byte) []byte {
	//创建解密块
	block, err := des.NewCipher(key)
	if nil != err {
		panic(err)
	}
	iv := []byte("12345678")
	//创建cbc解密模式
	blockMode := cipher.NewCBCDecrypter(block, iv)
	dst := make([]byte, len(src))
	blockMode.CryptBlocks(dst, src)
	return UnPaddingText(dst)
}

func main() {
	src := []byte("你好！XXX")
	key := []byte("87654321")
	encryptedMsg := EncryptDES(src, key)
	fmt.Println("加密后:", hex.EncodeToString(encryptedMsg))
	decryptedMsg := DecryptDES(encryptedMsg, key)
	fmt.Println("解密后的明文:", string(decryptedMsg))

}

//给最后一组数据填充至64字节
func PaddingText(src []byte, blockSize int) []byte {
	//求出最后一个分组需要填充的字节数
	padding := blockSize - len(src)%blockSize
	//创建新的切片，切片字节数为padding
	padText := bytes.Repeat([]byte{byte(padding)}, padding)
	//将新创建的切片和带填充的数据进行拼接
	nextText := append(src, padText...)
	return nextText

}

//取出数据尾部填充的赘余字符

func UnPaddingText(src []byte) []byte {
	//获取待处理数据长度
	len := len(src)
	//取出最后一个字符
	num := int(src[len-1])
	newText := src[:len-num]
	return newText
}
