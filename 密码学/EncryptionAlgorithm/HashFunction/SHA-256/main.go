package main

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
)

func main() {
	hash := sha256.New()
	hash.Write([]byte("你好！XXX"))
	bytes := hash.Sum(nil)
	fmt.Println(hex.EncodeToString(bytes))
	//36eda9afdaae39041fcbd101618717330a07b188373ab21f28e87e8793a7014d
}
