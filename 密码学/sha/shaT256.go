package main

import (
	"crypto/sha256"
	"fmt"
)

func main() {
	hash := sha256.New()
	hash.Write([]byte("eth1804"))
	byte := hash.Sum(nil)
	fmt.Printf("sha:%x \n", byte)
}
