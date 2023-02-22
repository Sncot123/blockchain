package main

import (
	"crypto/tls"
	"fmt"
	"io/ioutil"
	"net/http"
)

func main() {
	tr := &http.Transport{
		TLSClientConfig: &tls.Config{
			//跳过客户端对服务端证书的认证
			InsecureSkipVerify: true,
		},
	}
	//创建一个客户端
	client := http.Client{
		Transport: tr,
	}
	resp, _ := client.Get("https://127.0.0.1:8080")
	defer resp.Body.Close()
	all, _ := ioutil.ReadAll(resp.Body)
	fmt.Println(string(all))

}
