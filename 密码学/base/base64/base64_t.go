package main

import (
	"encoding/base64"
	"fmt"
)

func main() {
	//base64编码
	msg := "Man"
	encoded := base64.StdEncoding.EncodeToString([]byte(msg))
	fmt.Printf("base64:%v\n", encoded)

	//base64 解码
	decoded, _ := base64.StdEncoding.DecodeString("TWFu")
	fmt.Printf("base64:%x\n", decoded)
}
