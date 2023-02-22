package main

import (
	"bytes"
	"crypto/des"
	"encoding/hex"
	"fmt"
)

func EncryptDES(src, key []byte) []byte {
	//创建一个加密块
	block, err := des.NewCipher(key)
	if nil != err {
		panic(err)
	}

	//填充
	data := PaddingText(src, block.BlockSize())
	dst := make([]byte, 0, len(data))
	length := block.BlockSize()
	temp := make([]byte, length)
	//加密
	for len(data) > 0 {
		block.Encrypt(temp, data[:length])
		data = data[length:]
		dst = append(dst, temp...)
	}
	return dst
}

func DecryptDES(src, key []byte) []byte {
	//创建一个解密块
	block, err := des.NewCipher(key)
	if nil != err {
		panic(err)
	}
	//解密
	dst := make([]byte, 0, len(src))
	length := block.BlockSize()
	temp := make([]byte, length)
	for len(src) > 0 {
		block.Decrypt(temp, src[:length])
		src = src[length:]
		dst = append(dst, temp...)
	}
	//去除补充的
	dst = UnPaddingText(dst)
	return dst
}

func main() {
	//使用ECB模式的DES加密
	msg := []byte("你好 你好 ")
	key := []byte("12345678")
	encryptDES := EncryptDES(msg, key)
	fmt.Println(hex.EncodeToString(encryptDES))
	decryptDES := DecryptDES(encryptDES, key)
	fmt.Println(string(decryptDES))
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
