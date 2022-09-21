package xasset

import (
	"fmt"
	"os"
	"testing"
	"time"

	"github.com/xuperchain/xasset-sdk-go/auth"
	"github.com/xuperchain/xasset-sdk-go/client/base"
	"github.com/xuperchain/xasset-sdk-go/utils"
)

var (
	handle   *AssetOper
	sleepT   = time.Second * 30
	AccountA = base.TestAccount
	AccountB = base.TestTransAccount
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

func ChainReady(checkFunc func(...interface{}) error) <-chan interface{} {
	heartbeat := make(chan interface{}, 1)
	go func() {
		defer close(heartbeat)
		for {
			if err := checkFunc(); err == nil {
				heartbeat <- struct{}{}
				return
			}
			time.Sleep(sleepT)
		}
	}()

	return heartbeat
}

func CreatetAnAsset(account *auth.Account) (int64, error) {
	param := base.CreateAssetParam{
		Price:  10010,
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
	h, _ := NewAssetOperCli(base.TestGetXassetConfig(), &base.TestLogger{})
	if handle == nil {
		handle = h
	}

	resp, _, err := h.CreateAsset(&param)
	if err != nil {
		return 0, err
	}
	return resp.AssetId, nil
}

func checkAssetDone(assetId int64, status int) error {
	qResp, _, err := handle.QueryAsset(&base.QueryAssetParam{
		AssetId: assetId,
	})
	if err != nil {
		return err
	}
	if qResp.Meta.Status != status {
		return fmt.Errorf("oops, asset status: %d, want: %d", qResp.Meta.Status, status)
	}
	return nil
}

func checkShardDone(assetId int64, shardId int64, status int) error {
	qResp, _, err := handle.QueryShard(&base.QueryShardParam{
		AssetId: assetId,
		ShardId: shardId,
	})
	if err != nil {
		return err
	}
	if qResp.Meta.Status != status {
		return fmt.Errorf("oops, shard status: %d, want: %d", qResp.Meta.Status, status)
	}
	return nil
}

func TestAlterAsset(t *testing.T) {
	assetId, err := CreatetAnAsset(AccountA)
	if err != nil {
		t.Errorf("CreateAsset error, err: %v", err)
		return
	}
	param := base.AlterAssetParam{
		Price:   10011,
		Amount:  200,
		AssetId: assetId,
		Account: AccountA,
	}
	_, _, err = handle.AlterAsset(&param)
	if err != nil {
		t.Errorf("alter asset error, err: %v, asset_id: %d", err, assetId)
		return
	}
}

func TestScene(t *testing.T) {
	handle, _ := NewAssetOperCli(base.TestGetXassetConfig(), &base.TestLogger{})

	param1 := base.SceneListShardByAddrParam{}
	resp1, _, err := handle.SceneListShardByAddr(&param1)
	if err != nil {
		t.Errorf("scene listshardbyaddr error, err: %v, addr: %s", err, param1.Addr)
		return
	}

	fmt.Println(resp1.List[0].Title)

	param2 := base.SceneQueryShardParam{
		AssetId: 963581463228618626,
		ShardId: 9091098622548975525,
	}

	resp2, _, err := handle.SceneQueryShard(&param2)
	if err != nil {
		t.Errorf("scene queryshard error, err: %v, addr: %s", err, param2.Addr)
		return
	}

	fmt.Println(resp2.Meta.TxId)
}

func TestXasset(t *testing.T) {
	// Create Asset
	assetId, err := CreatetAnAsset(AccountA)
	if err != nil {
		t.Errorf("CreateAsset error, err: %v", err)
		return
	}
	param := &base.PublishAssetParam{
		AssetId:    assetId,
		Account:    AccountA,
		IsEvidence: 1,
	}

	// Publish Asset
	_, _, err = handle.PublishAsset(param)
	if err != nil {
		t.Errorf("PublishAsset error, err: %v", err)
		return
	}

	// check onChain status
	checkPublishFunc := func(...interface{}) error {
		return checkAssetDone(assetId, 4)
	}
	done := ChainReady(checkPublishFunc)
	<-done

	// Grant Shard
	grantParam := &base.GrantAssetParam{
		AssetId: assetId,
		Price:   10012,
		Account: AccountA,
		Addr:    AccountA.Address,
		ToAddr:  AccountB.Address,
	}
	grantResp, _, err := handle.GrantAsset(grantParam)
	if err != nil {
		t.Errorf("GrantAsset error, err: %v", err)
		return
	}
	shardId := grantResp.ShardId

	// check shard onChain status
	checkShardOnChain := func(...interface{}) error {
		return checkShardDone(assetId, shardId, 0)
	}
	grantDone := ChainReady(checkShardOnChain)
	<-grantDone

	// Transfer Shard
	transParam := &base.TransferAssetParam{
		AssetId: assetId,
		ShardId: shardId,
		Price:   10013,
		Account: AccountB,
		Addr:    AccountB.Address,
		ToAddr:  AccountA.Address,
	}
	_, _, err = handle.TransferAsset(transParam)
	if err != nil {
		t.Errorf("TransferAsset error, err: %v", err)
		return
	}

	// check shard onChain status
	transferDone := ChainReady(checkShardOnChain)
	<-transferDone

	// Consume Shard
	nonce := utils.GenNonce()
	signMsg := fmt.Sprintf("%d%d", assetId, nonce)
	sign, _ := auth.XassetSignECDSA(AccountA.PrivateKey, []byte(signMsg))
	consumeParam := &base.ConsumeShardParam{
		AssetId:  assetId,
		ShardId:  shardId,
		Nonce:    nonce,
		UAddr:    AccountA.Address,
		USign:    sign,
		UPKey:    AccountA.PublicKey,
		CAccount: AccountA,
	}
	_, _, err = handle.ConsumeShard(consumeParam)
	if err != nil {
		t.Errorf("consume shard error. err: %v", err)
		return
	}

	// check shard consume status
	checkShardConsume := func(...interface{}) error {
		return checkShardDone(assetId, shardId, 6)
	}
	consumeDone := ChainReady(checkShardConsume)
	<-consumeDone

	// Freeze Asset
	freezeParam := &base.FreezeAssetParam{
		AssetId: assetId,
		Account: AccountA,
	}
	_, _, err = handle.FreezeAsset(freezeParam)
	if err != nil {
		t.Errorf("freeze shard error. err: %v", err)
		return
	}

	// GetEvidenceInfo
	evidenceParam := &base.GetEvidenceInfoParam{
		AssetId: assetId,
	}
	_, _, err = handle.GetEvidenceInfo(evidenceParam)
	if err != nil {
		t.Errorf("GetEvidenceInfo error.err:%v", err)
		return
	}

	// ListAssetHistory
	listAstHisParam := &base.ListAssetHisParam{
		AssetId: assetId,
		Page:    1,
		Limit:   5,
	}
	_, _, err = handle.ListAssetHistory(listAstHisParam)
	if err != nil {
		t.Errorf("ListAssetHistory error.err:%v", err)
		return
	}

	listSdsByAddrParam := &base.ListShardsByAddrParam{
		Addr:  AccountA.Address,
		Page:  1,
		Limit: 20,
	}
	lResp, _, err := handle.ListShardsByAddr(listSdsByAddrParam)
	if err != nil {
		t.Errorf("list asset error. err: %v", err)
		return
	}
	t.Logf("Query srds, param: %+v, resp: %+v", param, lResp)

	listAstByAddrParam := &base.ListAssetsByAddrParam{
		Addr:   AccountA.Address,
		Status: 1,
		Page:   1,
		Limit:  20,
	}
	rResp, _, err := handle.ListAssetsByAddr(listAstByAddrParam)
	if err != nil {
		t.Errorf("list asset error. err: %v", err)
		return
	}
	t.Logf("Query ast, param: %+v, resp: %+v", param, rResp)

	listSdsByAstParam := &base.ListShardsByAssetParam{
		AssetId: assetId,
		Limit:   20,
	}

	qResp, _, err := handle.ListShardsByAsset(listSdsByAstParam)
	if err != nil {
		t.Errorf("list asset error. err: %v", err)
		return
	}
	t.Logf("Query srds, param: %+v, resp: %+v", param, qResp)
}

func TestSceneHasAsset(t *testing.T) {
	handle, _ := NewAssetOperCli(base.TestGetXassetConfig(), &base.TestLogger{})

	resp1, _, err := handle.SceneListAddr(base.UnionId)
	if err != nil {
		t.Errorf("scene list addr error, err: %v, union_id: %s", err, base.UnionId)
		return
	}

	fmt.Println(len(resp1.List))
	if len(resp1.List) == 0 {
		return
	}

	param1 := &base.SceneHasAssetByAddrParam{}
	param1.Addr = resp1.List[0].Addr
	param1.Token = resp1.List[0].Token
	param1.AssetIds = "[623062333210793858, 623062333210793859, 623062333210793856]"

	resp2, _, err := handle.SceneHasAssetByAddr(param1)
	if err != nil {
		t.Errorf("scene list addr error, err: %v, param: %+v", err, param1)
		return
	}
	fmt.Println(resp2)
}

func TestXassetDid(t *testing.T) {
	handle, _ := NewAssetOperCli(base.TestGetXassetConfig(), &base.TestLogger{})

	param1 := &base.BdBoxRegisterParam{
		OpenId: base.OpenId,
		AppKey: base.AppKey,
	}

	resp1, _, err := handle.BdBoxRegister(param1)
	if err != nil {
		t.Errorf("bdbox register error, err: %v, param: %+v", err, param1)
		return
	}

	fmt.Println(resp1.Address, resp1.IsNew)

	param2 := &base.BdBoxBindParam{
		OpenId:   base.OpenId,
		AppKey:   base.AppKey,
		Mnemonic: resp1.Mnemonic,
	}
	_, _, err = handle.BdBoxBind(param2)
	if err != nil {
		t.Errorf("bdbox bind error, err: %v, param: %+v", err, param2)
		return
	}

	param3 := &base.BindByUnionIdParam{
		UnionId:  base.UnionId,
		Mnemonic: resp1.Mnemonic,
	}
	_, _, err = handle.BindByUnionId(param3)
	if err != nil {
		t.Errorf("bdbox bind by union id error, err: %v, param: %+v", err, param3)
		return
	}

	resp2, _, err := handle.GetAddrByUnionId(base.UnionId)
	if err != nil {
		t.Errorf("get addr by union id error, err: %v, union id: %s", err, base.UnionId)
		return
	}
	fmt.Println(resp2.Address)
}
