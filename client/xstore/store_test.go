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
