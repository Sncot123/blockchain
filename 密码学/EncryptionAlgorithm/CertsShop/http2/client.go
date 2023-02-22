package main

import (
	"crypto/tls"
	"crypto/x509"
	"fmt"
	"io/ioutil"
	"net/http"
)

func main() {

	caCertPath := "D:\\JetBrains\\GoLand_workspace\\go-study\\blockchain\\EncryptionAlgorithm\\CertsShop\\ca\\ca.crt"
	file, _ := ioutil.ReadFile(caCertPath)
	//创建证书池
	pool := x509.NewCertPool()
	pool.AppendCertsFromPEM(file)

	tr := &http.Transport{
		//忽略客户端对服务端整数的校验
		TLSClientConfig: &tls.Config{
			//忽略客户端对服务端证书的校验
			//InsecureSkipVerify: true,
			////证书池
			RootCAs: pool,
		},
	}

	client := http.Client{Transport: tr}
	//发起请求
	resp, _ := client.Get("https://127.0.0.1:8080")

	defer resp.Body.Close()
	all, _ := ioutil.ReadAll(resp.Body)
	fmt.Println(string(all))
}
