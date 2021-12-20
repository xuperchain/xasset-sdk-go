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

func TestRetrieveAccountByMnemonic(t *testing.T) {
	acc, err := NewXchainEcdsaAccount(MnemStrgthMedium, MnemLangCN)
	if err != nil {
		t.Errorf("new account failed.err:%v", err)
		return
	}
	fmt.Println(acc)

	rc, err := RetrieveAccountByMnemonic(acc.Mnemonic, int(MnemLangCN))
	if err != nil {
		t.Errorf("retrieve account failed.err:%v", err)
		return
	}
	if rc.PrivateKey != acc.PrivateKey {
		t.Errorf("retrieve account. before:%v, after:%v", acc.PrivateKey, rc.PrivateKey)
	}
}
