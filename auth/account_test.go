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

func TestInferLanguage(t *testing.T) {
	mn1 := "used doll unit rifle learn extend drop brass pipe fancy daughter access expose link unveil fox sight enrich bargain amateur chunk large rough broken"
	mn2 := "伯 错 很 企 响 奉 党 府 惠 校 阅 科 富 沙 怀 数 纳 予"
	mn3 := "abc def fff"
	mn4 := "梨 我 他"

	if InferLanguage(mn1) != int(MnemLangEN) {
		t.Fatal(mn1)
	}
	if InferLanguage(mn2) != int(MnemLangCN) {
		t.Fatal(mn2)
	}
	if InferLanguage(mn3) != 0 {
		t.Fatal(mn3)
	}
	if InferLanguage(mn4) != 0 {
		t.Fatal(mn4)
	}
}
