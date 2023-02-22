package main

import (
	"bytes"
	"crypto/cipher"
	"crypto/des"
	"encoding/hex"
	"fmt"
)

func Encrypt3DES(src, key []byte) []byte {
	//创建加密区块
	block, err := des.NewTripleDESCipher(key)
	if nil != err {
		panic(err)
	}
	//填充数据
	src = PaddingText(src, block.BlockSize())

	//加密
	blockmode := cipher.NewCBCEncrypter(block, key[:block.BlockSize()])
	dst := make([]byte, len(src))
	blockmode.CryptBlocks(dst, src)
	return dst
}

func Decrypt3DES(dst, key []byte) []byte {
	//创建解密块
	block, err := des.NewTripleDESCipher(key)
	if nil != err {
		panic(err)
	}
	//创建解密模式
	blockmode := cipher.NewCBCDecrypter(block, key[:block.BlockSize()])
	src := make([]byte, len(dst))
	blockmode.CryptBlocks(src, dst)
	//去除填充的数据
	src = UnPaddingText(src)
	return src
}

func main() {
	encrypt_msg := []byte("你好 XXX")
	key := []byte("123456788765432112345678")
	encrypted3DES_msg := Encrypt3DES(encrypt_msg, key)
	fmt.Println(hex.EncodeToString(encrypted3DES_msg))

	decrypt_msg := Decrypt3DES(encrypted3DES_msg, key)
	fmt.Println(string(decrypt_msg))
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
