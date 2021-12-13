package auth

import (
	"fmt"
	"testing"
)

func TestNewXchainEcdsaAccount(t *testing.T) {
	acc, err := NewXchainEcdsaAccount(MnemStrgthMedium, MnemLangCN)
	if err != nil {
		t.Errorf("new account failed.err:%v", err)
		return
	}
	fmt.Println(acc)

	acc, err = NewXchainEcdsaAccount(MnemStrgthStrong, MnemLangEN)
	if err != nil {
		t.Errorf("new account failed.err:%v", err)
		return
	}
	fmt.Println(acc)
}
