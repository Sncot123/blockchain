package main

import (
	"encoding/hex"
	"fmt"

	"github.com/tjfoc/gmsm/sm3"
)

func EncryptSM3(msg string) string {
	hash := sm3.New()
	hash.Write([]byte(msg))
	bytes := hash.Sum(nil)
	return hex.EncodeToString(bytes)
}
func main() {
	encryptSM3 := EncryptSM3("你好！XXX")
	fmt.Println(encryptSM3)
	sm3Sum := sm3.Sm3Sum([]byte("你好！XXX"))
	fmt.Println(hex.EncodeToString(sm3Sum))
}
