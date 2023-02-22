package main

import (
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/x509"
	"crypto/x509/pkix"
	"encoding/pem"
	"flag"
	"io"
	"io/ioutil"
	"log"
	"math/big"
	"net"
	"net/mail"
	"os"
	"path/filepath"
	"strings"
	"time"
)

var serialNumberLimit = new(big.Int).Lsh(big.NewInt(1), 128)
var errorLog = log.New(os.Stderr, "ERROR:", log.Ldate)

//创建服务器端客户端证书
//args 用户输入的参数
//path 保存路径
//defaultDN 证书主题
//defaultSan 备用主题 相对于服务器来说可以是ip 联系邮箱
//defaultValidity证书有效时间
//keyUsage 证书用途（加密，数字签名）
//extKeyUsage 额外用途（服务端认证，客户端认证）
func createCertificate(args []string, path, defaultDN, defaultSan string, defaultValidity int, keyUsage x509.KeyUsage, extKeyUsage []x509.ExtKeyUsage) {
	flagSet := flag.NewFlagSet("ca", flag.PanicOnError)
	dn := flagSet.String("dn", defaultDN, "证书主题")
	san := flagSet.String("san", defaultSan, "证书备用主题")
	validity := flagSet.Int("validity", defaultValidity, "证书有效期")
	//解析命令行参数
	err := flagSet.Parse(args)
	if nil != err {
		errorLog.Fatalf("命令行参数解析失败：%s", err)
	}
	//ca/ica/client
	//找到证书保存路径的父级目录
	ca := filepath.Dir(path)
	var caCert *x509.Certificate
	var caKey *ecdsa.PrivateKey
	//解析父级证书
	caCert = parseCert(ca)
	if !caCert.IsCA {
		errorLog.Fatalf("%s不是证书颁发机构", filepath.Dir(path))
	}
	//解析父级密钥
	caKey = parseKey(ca)

	//生成密钥对
	priKey, derKey := GenKey()

	//生成证书序列号
	serialNumber, err := rand.Int(rand.Reader, serialNumberLimit)
	if nil != err {
		errorLog.Fatalf("证书序列号生成失败：%s", err)
	}

	notBefore := time.Now().UTC()
	notAfter := notBefore.AddDate(0, 0, *validity)
	template := x509.Certificate{
		//序列号
		SerialNumber: serialNumber,
		//主题
		Subject: *parseDN(caCert, *dn),
		//生效时间
		NotBefore: notBefore,
		//失效时间
		NotAfter: notAfter,
		//是否CA
		IsCA: false,

		//证书的用途：用于数字签名和证书签署
		KeyUsage: keyUsage,
		//额外用途
		ExtKeyUsage: extKeyUsage,
		//邮箱地址
		EmailAddresses: []string{},
		//ip地址
		IPAddresses: []net.IP{},
	}
	//从备用主题中解析邮箱和ip
	parseIpAndEmail(*san, &template)

	certificate, err := x509.CreateCertificate(rand.Reader, &template, caCert, &priKey.PublicKey, caKey)
	if nil != err {
		errorLog.Fatalf("证书生成失败：%s", err)
	}
	//保存证书
	SaveCert(path, certificate)
	//保存私钥
	savePriKey(path, derKey)

}

//从备用主题中解析邮箱和ip到模板中
func parseIpAndEmail(san string, template *x509.Certificate) {
	if san != "" {
		for _, h := range strings.Split(san, ",") {
			if ip := net.ParseIP(h); ip != nil {
				template.IPAddresses = append(template.IPAddresses, ip)
			} else if email, err := mail.ParseAddress(h); err != nil || email != nil {
				template.EmailAddresses = append(template.EmailAddresses, email.Address)
			}
		}
	}
}

//创建ca及ica证书
func createCA(args []string, path, defaultDN string, defaultValidity int) {
	flagSet := flag.NewFlagSet("ca", flag.PanicOnError)
	dn := flagSet.String("dn", defaultDN, "证书主题")
	validity := flagSet.Int("validity", defaultValidity, "证书有效期")
	//解析命令行参数
	err := flagSet.Parse(args)
	if nil != err {
		errorLog.Fatalf("命令行参数解析失败：%s", err)
	}
	//ca/ica/client
	//找到证书保存路径的父级目录
	ca := filepath.Dir(path)
	var caCert *x509.Certificate
	var maxPathLength int = 5
	var caKey *ecdsa.PrivateKey
	//如果父级目录不等于.，说明该证书是有父级证书
	if ca != "." {
		//解析父级证书
		caCert = parseCert(ca)
		if !caCert.IsCA {
			errorLog.Fatalf("%s不是证书颁发机构", ca)
		} else if !(caCert.MaxPathLen > 0) {
			errorLog.Fatalf("%s不能签发证书", ca)
		}
		maxPathLength = caCert.MaxPathLen - 1
		caKey = parseKey(ca)
	}

	//生成证书序列号
	serialNumber, err := rand.Int(rand.Reader, serialNumberLimit)
	if nil != err {
		errorLog.Fatalf("证书序列号生成失败：%s", err)
	}

	notBefore := time.Now().UTC()
	notAfter := notBefore.AddDate(0, 0, *validity)
	template := x509.Certificate{
		//序列号
		SerialNumber: serialNumber,
		//主题
		Subject: *parseDN(caCert, *dn),
		//生效时间
		NotBefore: notBefore,
		//失效时间
		NotAfter: notAfter,
		//是否CA
		IsCA: true,
		//表示是否是ca，MaxPathLen，MaxPathLenZero是否合法
		BasicConstraintsValid: true,
		//可以颁发的证书数量
		MaxPathLen:     maxPathLength,
		MaxPathLenZero: maxPathLength == 0,
		//证书的用途：用于数字签名和证书签署
		KeyUsage: x509.KeyUsageDigitalSignature | x509.KeyUsageCertSign,
	}
	//生成密钥对
	priKey, derKey := GenKey()

	if caCert == nil {
		caCert = &template
		caKey = priKey
	}
	certificate, err := x509.CreateCertificate(rand.Reader, &template, caCert, &caKey.PublicKey, caKey)
	if nil != err {
		errorLog.Fatalf("证书生成失败：%s", err)
	}
	//保存证书
	SaveCert(path, certificate)
	//保存私钥
	savePriKey(path, derKey)

}

//保存证书
func SaveCert(dir string, crt []byte) {
	//创建文件夹
	createDir(dir)
	//拼接文件名
	//ca/ica->ca/ica/ica.crt
	filename := filepath.Join(dir, filepath.Base(dir)+".crt")
	//打开文件
	certFile, err := os.OpenFile(filename, os.O_WRONLY|os.O_CREATE, 0777)
	if nil != err {
		errorLog.Fatalf("文件打开失败：%s", err)
	}
	defer certFile.Close()
	//将证书封装到块中
	block := pem.Block{
		Type:  "CERTIFICATE",
		Bytes: crt,
	}
	//将块写入文件
	if err = pem.Encode(certFile, &block); err != nil {
		errorLog.Fatalf("证书写入文件失败：%s", err)
	}

	//ca/ica->ca
	//ca-->.
	if filepath.Dir(dir) != "." {

		//ca/ca.crt
		//获取父级证书的路径
		caFile, err := os.Open(filepath.Join(filepath.Dir(dir), filepath.Base(filepath.Dir(dir))) + ".crt")
		if nil != err {
			errorLog.Fatalf("文件打开失败：%s", err)
		}
		defer caFile.Close()
		//将父级证书拷贝到自己的证书后面
		_, err = io.Copy(certFile, caFile)
		if nil != err {
			errorLog.Fatalf("证书拷贝失败：%s", err)
		}
		//将数据刷新到磁盘
		err = certFile.Sync()
		if nil != err {
			errorLog.Fatalf("失败：%s", err)
		}

	}
}

//保存私钥
func savePriKey(dir string, priKey []byte) {
	filename := filepath.Join(dir, filepath.Base(dir)+".key")
	//打开文件
	certFile, err := os.OpenFile(filename, os.O_WRONLY|os.O_CREATE, 0777)
	if nil != err {
		errorLog.Fatalf("文件打开失败：%s", err)
	}
	defer certFile.Close()
	//将证书封装到块中
	block := pem.Block{
		Type:  "EC PRIVATE KEY",
		Bytes: priKey,
	}
	//将块写入文件
	if err = pem.Encode(certFile, &block); err != nil {
		errorLog.Fatalf("证书写入文件失败：%s", err)
	}

}

//rsa加密生成密钥对
func GenKey() (*ecdsa.PrivateKey, []byte) {
	priKey, err := ecdsa.GenerateKey(elliptic.P384(), rand.Reader)

	if nil != err {
		errorLog.Fatalf("私钥生成失败：%s", err)
	}
	key, err := x509.MarshalECPrivateKey(priKey)
	if nil != err {
		errorLog.Fatalf("私钥序列化失败：%s", err)
	}
	return priKey, key
}

//创建文件夹
func createDir(dir string) {
	if _, err := os.Stat(dir); os.IsNotExist(err) {
		if err := os.MkdirAll(dir, 0777); nil != err {
			errorLog.Fatalf("文件创建失败：%s", err)
		}

	}

}

//读取父级私钥
func parseKey(path string) *ecdsa.PrivateKey {
	file, err := ioutil.ReadFile(filepath.Join(path, filepath.Base(path)+".key"))
	if nil != err {
		errorLog.Fatalf("私钥文件读取失败：%s", err)

	}
	//将字节证书转换为块
	block, _ := pem.Decode(file)
	if block == nil || block.Type != "EC PRIVATE KEY" {
		errorLog.Fatalf("证书文件解码失败！")
	}
	//将证书从块中取出
	priKey, err := x509.ParseECPrivateKey(block.Bytes)
	if nil != err {
		errorLog.Fatalf("私钥文件转码失败：%s", err)
	}
	return priKey
}

//从文件中读取证书
func parseCert(path string) *x509.Certificate {
	//path=ca/ica
	//读取父级证书
	file, err := ioutil.ReadFile(filepath.Join(path, filepath.Base(path)+".crt"))
	if nil != err {
		errorLog.Fatalf("证书文件读取失败：%s", err)
	}
	//将字节整数转换为块
	block, _ := pem.Decode(file)
	if block == nil || block.Type != "CERTIFICATE" {
		errorLog.Fatalf("证书文件解码失败：%v", err)
	}
	//将证书从块中取出
	crt, err := x509.ParseCertificate(block.Bytes)
	if nil != err {
		errorLog.Fatalf("证书文件转码失败：%v", err)
	}
	return crt
}

//解析用户输入的证书主题
//在解析主题时我们需要判断证书中是否有主题，如果有需要将解析出来的主题拼接到ca主题的后面，然后返回该主题
//如果ca没有主题，直接返回新解析出来的主题
func parseDN(ca *x509.Certificate, dn string) *pkix.Name {
	newName := &pkix.Name{}
	//-dn="/CN=China/O=xdl/OU=qkl"
	for _, element := range strings.Split(strings.Trim(dn, "/"), "/") {
		value := strings.Split(element, "=")
		if len(value) != 2 {
			errorLog.Fatalf("解析失败：%s", element)
		}
		switch strings.ToUpper(value[0]) {
		//域名
		case "CN":
			if value[1] != "" {
				newName.CommonName = value[1]
			}
			//国家
		case "C":
			if value[1] != "" {
				newName.Country = append(newName.Country, value[1])
			}
			//城市
		case "L":
			if value[1] != "" {
				newName.Locality = append(newName.Locality, value[1])
			}
			//州或省
		case "ST":
			if value[1] != "" {
				newName.Province = append(newName.Province, value[1])
			}
			//公司
		case "O":
			if value[1] != "" {
				newName.Organization = append(newName.Organization, value[1])
			}
			//部门
		case "OU":
			if value[1] != "" {
				newName.OrganizationalUnit = append(newName.OrganizationalUnit)
			}
		default:
			errorLog.Fatalf("解析失败：%s\n", element)
			return nil
		}
	}
	if ca != nil {
		newName.Country = append(ca.Subject.Country, newName.Country...)
		newName.Locality = append(ca.Subject.Locality, newName.Locality...)
		newName.Province = append(ca.Subject.Province, newName.Province...)
		newName.Organization = append(ca.Subject.Organization, newName.Organization...)
		newName.OrganizationalUnit = append(ca.Subject.OrganizationalUnit, newName.OrganizationalUnit...)
	}
	return newName
}

//certshop.exe ca
//certshop.exe ica
func main() {
	var command string
	if len(os.Args) < 2 {
		command = ""
	} else {
		command = os.Args[1]
	}
	switch command {
	case "ca":
		createCA(os.Args[2:], "ca", "/CN=localhost/O=xdl/OU=qkl/L=beijing", 365)
	case "ica":
		createCA(os.Args[2:], "ca/ica", "/CN=localhost/O=xdl/OU=qkl/L=beijing", 365)
	case "server":
		createCertificate(os.Args[2:], "ca/server", "/CN=server", "sncot@aliyun.com,127.0.0.1",
			365, x509.KeyUsageDigitalSignature|x509.KeyUsageKeyEncipherment,
			[]x509.ExtKeyUsage{x509.ExtKeyUsageServerAuth})
	case "client":
		createCertificate(os.Args[2:], "ca/client", "/CN=client", "",
			365, x509.KeyUsageDigitalSignature,
			[]x509.ExtKeyUsage{x509.ExtKeyUsageClientAuth})
	}
}
