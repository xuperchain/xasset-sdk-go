package auth

import (
	"encoding/hex"
	"fmt"
	"testing"
)

func TestXassetSignECDSA(t *testing.T) {
	acc, err := NewXchainEcdsaAccount(MnemStrgthMedium, MnemLangCN)
	if err != nil {
		t.Errorf("new account failed.err:%v", err)
		return
	}

	sign, err := XassetSignECDSA(acc.PrivateKey, []byte("hello world"))
	if err != nil {
		t.Errorf("xasset sign failed.err:%v", err)
		return
	}

	fmt.Printf("sign:%s\n", sign)
}

func TestXassetVerifyECDSA(t *testing.T) {
	acc, err := NewXchainEcdsaAccount(MnemStrgthMedium, MnemLangCN)
	if err != nil {
		t.Errorf("new account failed.err:%v", err)
		return
	}

	msg := []byte("hello world")
	sign, err := XassetSignECDSA(acc.PrivateKey, msg)
	if err != nil {
		t.Errorf("xasset sign failed.err:%v", err)
		return
	}
	fmt.Printf("sign:%s\n", sign)

	res, err := XassetVerifyECDSA(acc.PublicKey, sign, msg)
	if err != nil {
		t.Errorf("xasset sign verify failed.err:%v", err)
		return
	}

	fmt.Printf("result:%v\n", res)
}

func TestGetAddrByPubKey(t *testing.T) {
	acc, err := NewXchainEcdsaAccount(MnemStrgthMedium, MnemLangCN)
	if err != nil {
		t.Errorf("new account failed.err:%v", err)
		return
	}

	pk, _ := GetEcdsaPubKeyByJsStr(acc.PublicKey)
	addr, err := GetAddrByPubKey(pk)
	if err != nil {
		t.Errorf("get addr form public key failed.err:%v", err)
		return
	}

	fmt.Println(addr)
}

func TestVerifyAddrByPubKey(t *testing.T) {
	acc, err := NewXchainEcdsaAccount(MnemStrgthMedium, MnemLangCN)
	if err != nil {
		t.Errorf("new account failed.err:%v", err)
		return
	}

	pk, _ := GetEcdsaPubKeyByJsStr(acc.PublicKey)
	addr, err := GetAddrByPubKey(pk)
	if err != nil {
		t.Errorf("get addr form public key failed.err:%v", err)
		return
	}

	res, index := VerifyAddrByPubKey(addr, pk)
	fmt.Printf("addr:%v result:%v index:%v\n", addr, res, index)
}

func TestHashBySM3(t *testing.T) {
	msg := []byte("test")
	hash := HashBySM3(msg)
	fmt.Println(hex.EncodeToString(hash))
}
