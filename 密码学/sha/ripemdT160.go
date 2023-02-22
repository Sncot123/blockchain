package main

import (
	"fmt"

	"golang.org/x/crypto/ripemd160"
)

func main() {
	r160 := ripemd160.New()
	r160.Write([]byte("eth1804"))
	bytes := r160.Sum(nil)
	fmt.Printf("ripemd160:%x \n", bytes)
}
