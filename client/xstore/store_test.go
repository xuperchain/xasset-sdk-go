package xstore

import (
	"fmt"
	"testing"
	"time"

	"github.com/xuperchain/xasset-sdk-go/client/base"
)

var (
	handle   *StoreOper
	sleepT   = time.Second * 30
	AccountA = base.TestAccount
)

func TestListAct(t *testing.T) {
	param := &base.ListActParam{
		StoreId: 0,
	}
	handle, _ := NewXstoreOper(base.TestGetXassetConfig(), &base.TestLogger{})
	resp, _, err := handle.ListAct(param)
	if err != nil {
		t.Errorf("list act error.err:%v", err)
		return
	}

	fmt.Println(resp.List[0])
}

func TestStoreOper_CheckRefund(t *testing.T) {
	param := &base.CheckRefundParam{
		Oid: 1,
	}
	handle, _ := NewXstoreOper(base.TestGetXassetConfig(), &base.TestLogger{})
	resp, _, err := handle.CheckRefund(param)
	if err != nil {
		t.Errorf("check order refundable error.err:%v", err)
		return
	}
	fmt.Println(resp)
}

func TestStoreOper_CreateRefund(t *testing.T) {
	param := &base.CreateRefundParam{
		Oid:     1,
		Address: "xx",
		Reason:  "refund",
	}
	handle, _ := NewXstoreOper(base.TestGetXassetConfig(), &base.TestLogger{})
	resp, _, err := handle.CreateRefund(param)
	if err != nil {
		t.Errorf("create refund error.err:%v", err)
		return
	}
	fmt.Println(resp)
}

func TestStoreOper_CancelRefund(t *testing.T) {
	param := &base.CancelRefundParam{
		Rid:     1,
		Address: "xx",
	}
	handle, _ := NewXstoreOper(base.TestGetXassetConfig(), &base.TestLogger{})
	resp, _, err := handle.CancelRefund(param)
	if err != nil {
		t.Errorf("cancel refund error.err:%v", err)
		return
	}
	fmt.Println(resp)
}

func TestStoreOper_ConfirmRefund(t *testing.T) {
	param := &base.ConfirmRefundParam{
		Rid: 1,
	}
	handle, _ := NewXstoreOper(base.TestGetXassetConfig(), &base.TestLogger{})
	resp, _, err := handle.ConfirmRefund(param)
	if err != nil {
		t.Errorf("confirm refund error.err:%v", err)
		return
	}
	fmt.Println(resp)
}

func TestStoreOper_RefuseRefund(t *testing.T) {
	param := &base.RefuseRefundParam{
		Rid:     1,
		Message: "rejected",
	}
	handle, _ := NewXstoreOper(base.TestGetXassetConfig(), &base.TestLogger{})
	resp, _, err := handle.RefuseRefund(param)
	if err != nil {
		t.Errorf("refuse refund error.err:%v", err)
		return
	}
	fmt.Println(resp)
}
