package main

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/hex"
	"encoding/pem"
	"fmt"
	"os"
)

//生成密钥对
func Gen(bits int) {
	//1、生成密钥对
	//rand.Reader是一个全局的、共享的密码随机生成器
	privKey, err := rsa.GenerateKey(rand.Reader, bits)
	if nil != err {
		panic(err)
	}
	//2、将私钥序列化，
	//x509是通用的证书格式，包括序列号、签名算法、颁发者、有效时间、持有者、公钥
	//PKCS是rsa实验室联合一些其他安全系统开发商制定的一系列标准
	priStream := x509.MarshalPKCS1PrivateKey(privKey)
	block := pem.Block{
		Type:  "RSA PRIVATE KEY",
		Bytes: priStream,
	}
	//3、将密钥写入文件
	//pem是密钥格式
	priFile, err := os.Create("./privateKey.pem")
	if nil != err {
		panic(err)
	}
	defer priFile.Close()
	err = pem.Encode(priFile, &block)
	if nil != err {
		panic("私钥写入失败")
	}
	//从私钥获取公钥
	pubKey := privKey.PublicKey
	//序列化公钥
	pubStream := x509.MarshalPKCS1PublicKey(&pubKey)
	block = pem.Block{
		Type:  "RSA PUBLIC KEY",
		Bytes: pubStream,
	}
	//创建公钥文件
	pubFile, err := os.Create("./publicKey.pem")
	if nil != err {
		panic(err)
	}
	defer pubFile.Close()
	err = pem.Encode(pubFile, &block)
	if nil != err {
		panic("公钥写入文件失败")
	}
	return
}

//从公钥文件获取公钥，对明文进行加密
func EncryptRSA(src []byte, pubPath string) []byte {
	//1、从公钥文件获取公钥
	pubFile, err := os.Open(pubPath)
	if nil != err {
		panic(err)
	}
	pubFileInfo, err := pubFile.Stat()
	if nil != err {
		panic(err)
	}
	buf := make([]byte, pubFileInfo.Size())
	pubFile.Read(buf)
	//2、将公钥转换为块

	pubBlock, _ := pem.Decode(buf) //第一个参数为存储公钥的切片，第二个参数为未解码的数据
	//3、将块中的公钥反序列化

	pubKey, err := x509.ParsePKCS1PublicKey(pubBlock.Bytes)
	if nil != err {
		panic(err)
	}
	//4、用公钥对明文进行加密
	dst, err := rsa.EncryptPKCS1v15(rand.Reader, pubKey, src)
	if nil != err {
		panic(err)
	}
	return dst
}

func DecryptRSA(src []byte, priPath string) []byte {
	//1、从私钥文件读取私钥
	priFile, err := os.Open(priPath)
	if nil != err {
		panic(err)
	}
	priFileInfo, err := priFile.Stat()
	if nil != err {
		panic(err)
	}
	buf := make([]byte, priFileInfo.Size())
	priFile.Read(buf)
	//2、将私钥转换为块
	priBlock, _ := pem.Decode(buf)
	//3、将私钥块反序列化
	priKey, err := x509.ParsePKCS1PrivateKey(priBlock.Bytes)
	if nil != err {
		panic(err)
	}
	dst, err := rsa.DecryptPKCS1v15(rand.Reader, priKey, src)
	if nil != err {
		panic(err)
	}
	return dst
}

func main() {
	Gen(2048)
	encryptRSA := EncryptRSA([]byte("你好！潘丽萍"), "publicKey.pem")
	fmt.Println(hex.EncodeToString(encryptRSA))
	decryptRSA := DecryptRSA(encryptRSA, "privateKey.pem")
	fmt.Println(string(decryptRSA))
}
