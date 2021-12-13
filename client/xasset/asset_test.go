package xasset

import (
	"errors"
	"fmt"
	"os"
	"testing"
	"time"

	"github.com/xuperchain/xasset-sdk-go/auth"
	"github.com/xuperchain/xasset-sdk-go/client/base"
)

func TestGetStoken(t *testing.T) {
	param := &base.GetStokenParam{
		Account: base.TestAccount,
	}
	handle, _ := NewAssetOperCli(base.TestGetXassetConfig(), &base.TestLogger{})
	_, _, err := handle.GetStoken(param)
	if err != nil {
		t.Errorf("get stoken error.err:%v", err)
		return
	}
}

func TestUploadFile(t *testing.T) {
	file, err := os.Open("file_path")
	if err != nil {
		return
	}
	defer file.Close()

	// FileInfo:
	stats, err := file.Stat()
	if err != nil {
		return
	}

	// []byte
	data := make([]byte, stats.Size())
	_, err = file.Read(data)
	if err != nil {
		return
	}

	param := &base.UploadFileParam{
		Account:  base.TestAccount,
		FileName: "test_bytes_go.jpg",
		DataByte: data,
		Property: "200_300",
	}

	handle, _ := NewAssetOperCli(base.TestGetXassetConfig(), &base.TestLogger{})
	resp, res, err := handle.UploadFile(param)
	if err != nil {
		t.Errorf("query asset error.err:%v", err)
		return
	}

	fmt.Println(res, resp)
}

type ProcedureFunc func(account *auth.Account, handle *AssetOper, resp interface{}) (*AssetOper, interface{}, error)

func BridgeTest(account *auth.Account, handle *AssetOper, resp interface{}, procedures ...ProcedureFunc) (*AssetOper, interface{}, error) {
	procedureA := procedures[0]
	handle, resp, err := procedureA(account, handle, resp)
	if err != nil {
		return nil, resp, err
	}
	procedureB := procedures[1:]
	for _, procedure := range procedureB {
		_, resp, err = procedure(account, handle, resp)
		if err != nil {
			return nil, resp, err
		}
	}
	return handle, resp, nil
}

func CreateTestAsset(account *auth.Account, handle *AssetOper, resp interface{}) (*AssetOper, interface{}, error) {
	param := base.CreateAssetParam{
		Amount: 100,
		AssetInfo: &base.CreateAssetInfo{
			AssetCate: base.AssetCateArt,
			Title:     "我是一个小画家",
			Thumb: []string{
				"bos_v1://bucket/object/1000_500",
			},
			ShortDesc: "我是一个小画家",
			ImgDesc: []string{
				"bos_v1://bucket/object/1000_500",
			},
			AssetUrl: []string{
				"bos_v1://bucket/object/1000_500",
			},
		},
		Account: account,
	}
	handle, _ = NewAssetOperCli(base.TestGetXassetConfig(), &base.TestLogger{})
	resp, _, err := handle.CreateAsset(&param)
	return handle, resp, err
}

func TestCreateAsset(t *testing.T) {
	_, _, err := CreateTestAsset(base.TestAccount, nil, nil)
	if err != nil {
		t.Errorf("create asset error. err: %v", err)
		return
	}
}

func TestAlterAsset(t *testing.T) {
	handle, resp, err := CreateTestAsset(base.TestAccount, nil, nil)
	if err != nil {
		t.Errorf("CreateAsset error, err: %v", err)
		return
	}
	value, ok := resp.(*base.CreateAssetResp)
	if !ok {
		t.Errorf("transfer to createResp error.")
		return
	}
	param := base.AlterAssetParam{
		Amount:  200,
		AssetId: value.AssetId,
		Account: base.TestAccount,
	}
	_, _, err = handle.AlterAsset(&param)
	if err != nil {
		t.Errorf("alter asset error, err: %v, asset_id: %d", err, value.AssetId)
		return
	}
}

func PublishTestAsset(account *auth.Account, handle *AssetOper, resp interface{}) (*AssetOper, interface{}, error) {
	value, ok := resp.(*base.CreateAssetResp)
	if !ok {
		return nil, nil, errors.New("transfer error")
	}
	param := &base.PublishAssetParam{
		AssetId: value.AssetId,
		Account: account,
	}
	_, _, err := handle.PublishAsset(param)
	// return *CreateAssetResp
	return handle, resp, err
}

func TestPublishAsset(t *testing.T) {
	if _, _, err := BridgeTest(base.TestAccount, nil, nil, CreateTestAsset, PublishTestAsset); err != nil {
		t.Errorf("publish asset error. err: %v", err)
	}
}

func TestQueryAsset(t *testing.T) {
	handle, resp, err := CreateTestAsset(base.TestAccount, nil, nil)
	if err != nil {
		t.Errorf("CreateAsset error, err: %v", err)
		return
	}
	value, ok := resp.(*base.CreateAssetResp)
	if !ok {
		t.Errorf("transfer to createResp error.")
		return
	}
	param := &base.QueryAssetParam{
		AssetId: value.AssetId,
	}
	_, _, err = handle.QueryAsset(param)
	if err != nil {
		t.Errorf("query asset error. err: %v", err)
		return
	}
}

func GrantTestAsset(account *auth.Account, handle *AssetOper, resp interface{}) (*AssetOper, interface{}, error) {
	value, ok := resp.(*base.CreateAssetResp)
	if !ok {
		return nil, nil, errors.New("transfer error")
	}
	param := &base.GrantAssetParam{
		AssetId: value.AssetId,
		Account: account,
		Addr:    account.Address,
		ToAddr:  base.TestTransAccount.Address,
	}
	nResp, _, err := handle.GrantAsset(param)
	if err != nil {
		q := &base.QueryAssetParam{
			AssetId: value.AssetId,
		}
		qRes, _, _ := handle.QueryAsset(q)
		return nil, qRes.Meta.Status, err
	}
	return handle, nResp, err
}

// TestGrantNTransAsset use -timeout=500s when executing go test
func TestGrantNTransAsset(t *testing.T) {
	// do CreateAsset & PublishAsset
	handle, resp, err := BridgeTest(base.TestAccount, nil, nil, CreateTestAsset, PublishTestAsset)
	if err != nil {
		t.Errorf("publish asset error. err: %v", err)
		return
	}

	// waiting for data be send to the chain
	waitT := time.Duration(60)
	time.Sleep(waitT * time.Second)
	// do GrantAsset
	handle, resp, err = BridgeTest(base.TestAccount, handle, resp, GrantTestAsset)
	if err != nil {
		t.Errorf("grant asset error. err: %v, status: %d", err, resp.(int))
		return
	}
	// get asset_id
	value := resp.(*base.GrantAssetResp)
	assetId := value.AssetId
	shardId := value.ShardId

	time.Sleep(waitT * time.Second)
	// do TransferAsset
	handle, resp, err = BridgeTest(base.TestTransAccount, handle, resp, TransferTestAsset)
	if err != nil {
		t.Errorf("transfer asset error. err: %v", err)
		return
	}

	time.Sleep(waitT * time.Second)
	// do QueryShardAsset
	srdP := &base.QueryShardParam{
		AssetId: assetId,
		ShardId: shardId,
	}
	nResp, _, err := handle.QueryShard(srdP)
	if err != nil {
		t.Errorf("query shard error. err: %v", err)
		return
	}
	if nResp.Meta.OwnerAddr != base.TestAccount.Address {
		t.Errorf("query shard error. owner: %s", nResp.Meta.OwnerAddr)
	}
}

func TestListShardsByAddr(t *testing.T) {
	handle, _ := NewAssetOperCli(base.TestGetXassetConfig(), &base.TestLogger{})
	param := &base.ListShardsByAddrParam{
		Addr:  base.TestTransAccount.Address,
		Page:  1,
		Limit: 20,
	}
	lResp, _, err := handle.ListShardsByAddr(param)
	if err != nil {
		t.Errorf("list asset error. err: %v", err)
		return
	}
	if lResp.TotalCnt <= 0 {
		t.Error("read asset error")
		return
	}
}

func EvidenceTestAsset(account *auth.Account, handle *AssetOper, resp interface{}) (*AssetOper, interface{}, error) {
	value, ok := resp.(*base.CreateAssetResp)
	if !ok {
		return nil, nil, errors.New("transfer error")
	}
	param := &base.PublishAssetParam{
		AssetId:    value.AssetId,
		Account:    account,
		IsEvidence: 1,
	}
	_, _, err := handle.PublishAsset(param)
	// return *CreateAssetResp
	return handle, resp, err
}

func TestGetEvidence(t *testing.T) {
	// do CreateAsset & PublishAsset
	handle, resp, err := BridgeTest(base.TestAccount, nil, nil, CreateTestAsset, EvidenceTestAsset)
	if err != nil {
		t.Errorf("publish asset error. err: %v", err)
		return
	}

	// waiting for data be send to the chain
	waitT := time.Duration(60)
	time.Sleep(waitT * time.Second)
	// do GrantAsset
	handle, resp, err = BridgeTest(base.TestAccount, handle, resp, GrantTestAsset)
	if err != nil {
		t.Errorf("grant asset error. err: %v, status: %d", err, resp.(int))
		return
	}

	time.Sleep(waitT * time.Second)
	// get asset_id
	value := resp.(*base.GrantAssetResp)
	assetId := value.AssetId
	shardId := value.ShardId
	eParam := &base.GetEvidenceInfoParam{
		AssetId: assetId,
		ShardId: shardId,
	}
	_, _, err = handle.GetEvidenceInfo(eParam)
	if err != nil {
		t.Errorf("create asset error.err:%v", err)
		return
	}
}

func TransferTestAsset(account *auth.Account, handle *AssetOper, resp interface{}) (*AssetOper, interface{}, error) {
	value, ok := resp.(*base.GrantAssetResp)
	if !ok {
		return nil, nil, errors.New("transfer error")
	}
	param := &base.TransferAssetParam{
		AssetId: value.AssetId,
		ShardId: value.ShardId,
		Account: account,
		Addr:    account.Address,
		ToAddr:  base.TestAccount.Address,
	}
	nResp, _, err := handle.TransferAsset(param)
	return handle, nResp, err
}
