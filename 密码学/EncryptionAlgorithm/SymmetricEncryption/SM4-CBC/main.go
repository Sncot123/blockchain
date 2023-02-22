package main

import (
	"bytes"
	"crypto/cipher"
	"encoding/hex"
	"fmt"

	"github.com/tjfoc/gmsm/sm4"
)

func EncryptSM4(src, pwd []byte) []byte {
	//创建加密块
	block, err := sm4.NewCipher([]byte(pwd))
	if nil != err {
		panic(err)
	}
	//填充数据
	src = PaddingText(src, block.BlockSize())
	//初始化向量
	iv := []byte("1234567887654321")
	//设置加密模式
	blockmode := cipher.NewCBCEncrypter(block, iv)
	dst := make([]byte, len(src))
	blockmode.CryptBlocks(dst, src)
	return dst
}

func DecryptSM4(src, pwd []byte) []byte {
	//创建解密块
	block, err := sm4.NewCipher(pwd)
	if nil != err {
		panic(err)
	}
	//初始化向量
	iv := []byte("1234567887654321")
	//创建解密模式
	blockmode := cipher.NewCBCDecrypter(block, iv)
	dst := make([]byte, len(src))
	//解密
	blockmode.CryptBlocks(dst, src)

	//去除填充
	dst = UnPaddingText(dst)
	return dst
}

func main() {
	msg := []byte("潘丽萍")
	pwd := []byte("1234567887654321")
	encryptedMsg := EncryptSM4(msg, pwd)
	decryptMsg := DecryptSM4(encryptedMsg, pwd)
	fmt.Printf("%v\n%v", hex.EncodeToString(encryptedMsg), string(decryptMsg))
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
