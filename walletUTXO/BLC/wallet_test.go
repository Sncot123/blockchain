package BLC

import (
	"fmt"
	"testing"
)

func TestNewWallet(t *testing.T) {
	wallet := NewWallet()
	fmt.Println("privateKey:", wallet.PrivateKey)
	fmt.Println("publicKey:", wallet.PublicKey)
	fmt.Println("wallet:", wallet)
}

func TestWallet_GetAddress(t *testing.T) {
	wallet := NewWallet()
	address := wallet.GetAddress()
	fmt.Printf("the address of wallet is [%s]\n", address)
	fmt.Printf("the validation of current address is %v\n", IsValidForAddress(address))
}
