package main

import (
	"crypto/md5"
	"encoding/hex"
	"fmt"
	"io"
	"os"
)

//对文件进行hash处理

func HashFileMD5(path string) string {
	//1、获取文件句柄
	file, err := os.Open(path)
	if nil != err {
		panic(err)
	}
	//2、将文件内容写入到hash
	hash := md5.New()
	_, err = io.Copy(hash, file)
	if nil != err {
		panic(err)
	}
	//3、计算散列值
	bytes := hash.Sum(nil)
	return hex.EncodeToString(bytes)
}
func main() {

	//散列对象创建
	hash := md5.New()
	//将内容写入对象
	hash.Write([]byte("区块链"))
	//计算出散列值
	sum := hash.Sum(nil)
	fmt.Println(hex.EncodeToString(sum))
	//75d8fafb0706c9381d4c91e3b184f19d

	hash.Reset()
	hash.Write([]byte("区块链"))
	bytes := hash.Sum([]byte("123"))
	fmt.Println(hex.EncodeToString(bytes))
	//313233  75d8fafb0706c9381d4c91e3b184f19d
	i := md5.Sum([]byte("区块链"))
	fmt.Println(hex.EncodeToString(i[:]))
	fileMD5 := HashFileMD5("privateKey.pem")
	fmt.Println(fileMD5, "-------3234")
}
