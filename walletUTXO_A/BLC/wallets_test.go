package BLC

import (
	"fmt"
	"testing"
)

func TestWallets_CreateWallets(t *testing.T) {
	wallets := NewWallets()
	wallets.CreateWallets()
	fmt.Printf("wallets:%v\n", wallets.Wallets)
}
