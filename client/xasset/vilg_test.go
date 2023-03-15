package xasset

import (
	"fmt"
	"testing"

	"github.com/xuperchain/xasset-sdk-go/client/base"
)

func getTestHandler() *AssetOper {
	handle, _ := NewAssetOperCli(base.TestGetXassetConfig(), &base.TestLogger{})
	return handle
}

func TestVilgText2Img(t *testing.T) {
	h := getTestHandler()

	param := base.VilgText2ImgParam{
		Text:       "机器人大战人类",
		Style:      6,
		Resolution: 1,
	}
	resp, _, err := h.VilgText2Img(&param)
	if err != nil {
		t.Fatal(err)
	}

	fmt.Println("taskId:", resp.TaskId)
}

func TestVilgGetImg(t *testing.T) {
	h := getTestHandler()

	var taskId int64 = 14643233

	resp, _, err := h.VilgGetImg(taskId)
	if err != nil {
		t.Fatal(err)
	}

	fmt.Println("status:", resp.Task.Status)
	fmt.Println("img:", resp.Task.Img)
}

func TestGetBalance(t *testing.T) {
	h := getTestHandler()
	resp, _, err := h.VilgBalance()
	if err != nil {
		t.Fatal(err)
	}

	fmt.Println("balance:", resp.Balance)
}
