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

func TestListActAst(t *testing.T) {
	param := &base.BaseActParam{
		ActId: 5407,
	}
	handle, _ := NewXstoreOper(base.TestGetXassetConfig(), &base.TestLogger{})
	resp, _, err := handle.ListActAst(param)
	if err != nil {
		t.Errorf("list act error.err:%v", err)
		return
	}

	fmt.Println(resp.List[0].ActId)
}
