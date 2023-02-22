package main

import (
	"bytes"
	"crypto/aes"
	"crypto/cipher"
	"encoding/hex"
	"fmt"
)

func EncryptAES(src, key []byte) []byte {
	block, err := aes.NewCipher(key)
	if nil != err {
		panic(err)
	}
	//填充
	src = PaddingText(src, block.BlockSize())
	//初始化向量
	iv := []byte("12345678abcdefgh")
	//创建加密模式
	blockmode := cipher.NewCBCEncrypter(block, iv)
	dst := make([]byte, len(src))
	//加密
	blockmode.CryptBlocks(dst, src)
	return dst
}

func DecryptAES(src, key []byte) []byte {
	block, err := aes.NewCipher(key)
	if nil != err {
		panic(err)
	}
	//初始化向量
	iv := []byte("12345678abcdefgh")
	blockmode := cipher.NewCBCDecrypter(block, iv)
	//解密
	dst := make([]byte, len(src))
	blockmode.CryptBlocks(dst, src)
	//去除填充
	dst = UnPaddingText(dst)
	return dst
}

func main() {
	src := []byte("你好！潘丽萍")
	key := []byte("8765432112345678")
	encryptedMsg := EncryptAES(src, key)
	fmt.Println("加密后:", hex.EncodeToString(encryptedMsg))
	decryptedMsg := DecryptAES(encryptedMsg, key)
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
