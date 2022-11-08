package xasset

import (
	"encoding/json"
	"fmt"
	"net/url"

	auth2 "github.com/baidubce/bce-sdk-go/auth"
	"github.com/baidubce/bce-sdk-go/services/bos"

	"github.com/xuperchain/xasset-sdk-go/auth"
	xbase "github.com/xuperchain/xasset-sdk-go/client/base"
	"github.com/xuperchain/xasset-sdk-go/common/config"
	"github.com/xuperchain/xasset-sdk-go/common/logs"
	"github.com/xuperchain/xasset-sdk-go/utils"
)

type AssetOper struct {
	xbase.XassetBaseClient
}

func NewAssetOperCli(cfg *config.XassetCliConfig, logger logs.LogDriver) (*AssetOper, error) {
	obj := &AssetOper{}
	err := obj.InitClient(cfg, logger)
	if err != nil {
		return nil, err
	}

	return obj, nil
}

// genGetStokenBody Grant uses the general parameter as follows,
//
//	   {
//			   Addr     string `json:"addr"`
//			   Sign     string `json:"sign"`
//			   PKey     string `json:"pkey"`
//			   Nonce    int64  `json:"nonce"`
//		  }
func (t *AssetOper) genGetStokenBody(param *xbase.GetStokenParam) (string, error) {
	nonce := utils.GenNonce()
	signMsg := fmt.Sprintf("%d", nonce)
	sign, err := auth.XassetSignECDSA(param.Account.PrivateKey, []byte(signMsg))
	if err != nil {
		return "", xbase.ComErrAccountSignFailed
	}

	v := url.Values{}
	v.Set("addr", param.Account.Address)
	v.Set("sign", sign)
	v.Set("pkey", param.Account.PublicKey)
	v.Set("nonce", fmt.Sprintf("%d", nonce))
	body := v.Encode()
	return body, nil
}

func (t *AssetOper) GetStoken(param *xbase.GetStokenParam) (*xbase.GetStokenResp, *xbase.RequestRes, error) {
	if err := param.Valid(); err != nil {
		return nil, nil, err
	}

	body, err := t.genGetStokenBody(param)
	if err != nil {
		t.Logger.Warn("fail to generate value for getting stoken, err: %v, param: %+v", err, *param)
		return nil, nil, err
	}
	res, err := t.Post(xbase.FileApiGetStoken, body)
	if err != nil {
		t.Logger.Warn("post request xasset failed.err:%v", err)
		return nil, nil, xbase.ComErrRequsetFailed
	}
	if res.HttpCode != 200 {
		t.Logger.Warn("post request response is not 200. [http_code: %d] [url: %s] [body: %s] [trace_id: %s]",
			res.HttpCode, res.ReqUrl, res.Body, t.GetTarceId(res.Header))
		return nil, nil, xbase.ComErrRespCodeErr
	}

	var resp xbase.GetStokenResp
	err = json.Unmarshal([]byte(res.Body), &resp)
	if err != nil {
		t.Logger.Warn("unmarshal body failed. [http_code: %d] [url: %s] [body: %s] [trace_id: %s]",
			res.HttpCode, res.ReqUrl, res.Body, t.GetTarceId(res.Header))
		return nil, res, xbase.ComErrUnmarshalBodyFailed
	}
	if resp.Errno != xbase.XassetErrNoSucc {
		t.Logger.Warn("get resp failed. [url: %s] [request_id: %s] [err_no: %d] [trace_id: %s]",
			res.ReqUrl, resp.RequestId, resp.Errno, t.GetTarceId(res.Header))
		return nil, res, xbase.ComErrServRespErrnoErr
	}

	t.Logger.Trace("operate succ. [accessInfo: %v] [url: %s] [request_id: %s] [trace_id: %s]",
		resp.AccessInfo, res.ReqUrl, resp.RequestId, t.GetTarceId(res.Header))
	return &resp, res, nil
}

func (t *AssetOper) UploadFile(param *xbase.UploadFileParam) (*xbase.UploadFileResp, *xbase.RequestRes, error) {
	if err := param.Valid(); err != nil {
		return nil, nil, err
	}

	resp, res, err := t.GetStoken(&xbase.GetStokenParam{Account: param.Account})
	if err != nil {
		t.Logger.Warn("get stoken failed.[url:%s] [request_id:%s] [err_no:%d] [trace_id:%s]",
			res.ReqUrl, resp.RequestId, resp.Errno, t.GetTarceId(res.Header))
		return nil, nil, err
	}

	bosClient, err := bos.NewClient(resp.AccessInfo.AK, resp.AccessInfo.SK, resp.AccessInfo.EndPoint)
	if err != nil {
		t.Logger.Warn("create bos client failed.err:%v", err)
		return nil, nil, err
	}
	stsCredential, err := auth2.NewSessionBceCredentials(resp.AccessInfo.AK, resp.AccessInfo.SK, resp.AccessInfo.SessionToken)
	if err != nil {
		t.Logger.Warn("create sts credential object failed.err:%v", err)
		return nil, nil, err
	}
	bosClient.Config.Credentials = stsCredential

	key := fmt.Sprintf("/%s%s", resp.AccessInfo.ObjectPath, param.FileName)

	if param.FilePath != "" {
		_, err = bosClient.PutObjectFromFile(resp.AccessInfo.Bucket, key, param.FilePath, nil)
		if err != nil {
			t.Logger.Warn("upload file through local file failed.err:%v", err)
			return nil, nil, err
		}
	} else if param.DataByte != nil {
		_, err = bosClient.PutObjectFromBytes(resp.AccessInfo.Bucket, key, param.DataByte, nil)
		if err != nil {
			t.Logger.Warn("upload file through bytes failed.err:%v", err)
			return nil, nil, err
		}
	} else {
		t.Logger.Warn("unsupported upload file method")
		return nil, nil, fmt.Errorf("wrong upload file method")
	}

	link := fmt.Sprintf("bos_v1://%s%s/%s", resp.AccessInfo.Bucket, key, param.Property)
	t.Logger.Trace("upload file succ.[link:%s] [url:%s] [request_id:%s] [trace_id:%s]",
		link, res.ReqUrl, resp.RequestId, t.GetTarceId(res.Header))

	return &xbase.UploadFileResp{
		Link:       link,
		AccessInfo: resp.AccessInfo,
	}, res, nil
}

// GenCreateAssetBody uses the parameter as follows,
//
//	{
//			AssetId   int64  `json:"asset_id"`
//			Price     int64  `json:"price"`
//			Amount    int    `json:"amount"`
//			AssetInfo string `json:"asset_info"`
//			Addr      string `json:"addr"`
//			Sign      string `json:"sign"`
//			PKey      string `json:"pkey"`
//			Nonce     int64  `json:"nonce"`
//			UserId    int64  `json:"user_id,omitempty"`
//			FileHash  string `json:"file_hash,omitempty"`
//	}
func (t *AssetOper) genCreateAssetBody(appid int64, param *xbase.CreateAssetParam) (string, error) {
	nonce := utils.GenNonce()
	assetId := utils.GenAssetId(appid)
	signMsg := fmt.Sprintf("%d%d", assetId, nonce)
	sign, err := auth.XassetSignECDSA(param.Account.PrivateKey, []byte(signMsg))
	if err != nil {
		return "", xbase.ComErrAccountSignFailed
	}
	assetInfo, err := json.Marshal(param.AssetInfo)
	if err != nil {
		return "", xbase.ComErrJsonMarFailed
	}

	v := url.Values{}
	v.Set("asset_id", fmt.Sprintf("%d", assetId))
	v.Set("price", fmt.Sprintf("%d", param.Price))
	v.Set("amount", fmt.Sprintf("%d", param.Amount))
	v.Set("asset_info", string(assetInfo))
	v.Set("addr", param.Account.Address)
	v.Set("sign", sign)
	v.Set("pkey", param.Account.PublicKey)
	v.Set("nonce", fmt.Sprintf("%d", nonce))
	if err := xbase.IdValid(param.UserId); err == nil {
		v.Set("user_id", fmt.Sprintf("%d", param.UserId))
	}
	if param.FileHash != "" {
		v.Set("file_hash", param.FileHash)
	}
	body := v.Encode()
	return body, nil
}

func (t *AssetOper) CreateAsset(param *xbase.CreateAssetParam) (*xbase.CreateAssetResp, *xbase.RequestRes, error) {
	if err := param.Valid(); err != nil {
		return nil, nil, err
	}

	body, err := t.genCreateAssetBody(t.GetConfig().Credentials.AppId, param)
	if err != nil {
		t.Logger.Warn("fail to generate value for creating, err: %v, param: %+v", err, *param)
		return nil, nil, err
	}
	res, err := t.Post(xbase.AssetApiCreate, body)
	if err != nil {
		t.Logger.Warn("post request xasset failed, uri: %s, err: %v", xbase.AssetApiCreate, err)
		return nil, nil, xbase.ComErrRequsetFailed
	}
	if res.HttpCode != 200 {
		t.Logger.Warn("post request response is not 200. [http_code: %d] [url: %s] [body: %s] [trace_id: %s]",
			res.HttpCode, res.ReqUrl, res.Body, t.GetTarceId(res.Header))
		return nil, nil, xbase.ComErrRespCodeErr
	}

	var resp xbase.CreateAssetResp
	err = json.Unmarshal([]byte(res.Body), &resp)
	if err != nil {
		t.Logger.Warn("unmarshal body failed. [http_code: %d] [url: %s] [body: %s] [trace_id: %s]",
			res.HttpCode, res.ReqUrl, res.Body, t.GetTarceId(res.Header))
		return nil, res, xbase.ComErrUnmarshalBodyFailed
	}
	if resp.Errno != xbase.XassetErrNoSucc {
		t.Logger.Warn("get resp failed. [url: %s] [request_id: %s] [err_no: %d] [trace_id: %s]",
			res.ReqUrl, resp.RequestId, resp.Errno, t.GetTarceId(res.Header))
		return nil, res, xbase.ComErrServRespErrnoErr
	}

	t.Logger.Trace("operate succ. [asset_id: %d] [url: %s] [request_id: %s] [trace_id: %s]",
		resp.AssetId, res.ReqUrl, resp.RequestId, t.GetTarceId(res.Header))
	return &resp, res, nil
}

// GenAlterAssetBody uses the parameter as follows,
//
//	{
//			AssetId   int64  `json:"asset_id"`
//			Addr      string `json:"addr"`
//			Sign      string `json:"sign"`
//			PKey      string `json:"pkey"`
//			Nonce     int64  `json:"nonce"`
//			Amount    int    `json:"amount"`
//			AssetInfo string `json:"asset_info"`
//			FileHash  string `json:"file_hash"`
//	}
func (t *AssetOper) genAlterAssetBody(param *xbase.AlterAssetParam) (string, error) {
	nonce := utils.GenNonce()
	signMsg := fmt.Sprintf("%d%d", param.AssetId, nonce)
	sign, err := auth.XassetSignECDSA(param.Account.PrivateKey, []byte(signMsg))
	if err != nil {
		return "", xbase.ComErrAccountSignFailed
	}

	v := url.Values{}
	v.Set("asset_id", fmt.Sprintf("%d", param.AssetId))
	v.Set("addr", param.Account.Address)
	v.Set("sign", sign)
	v.Set("pkey", param.Account.PublicKey)
	v.Set("nonce", fmt.Sprintf("%d", nonce))

	if err := xbase.PriceInvalid(param.Price); err == nil {
		v.Set("price", fmt.Sprintf("%d", param.Price))
	}
	if err := xbase.AmountInvalid(param.Amount); err == nil {
		v.Set("amount", fmt.Sprintf("%d", param.Amount))
	}
	if err := xbase.AlterAssetInfoValid(param.AssetInfo); err == nil {
		assetInfo, err := json.Marshal(param.AssetInfo)
		if err != nil {
			return "", xbase.ComErrJsonMarFailed
		}
		v.Set("asset_info", string(assetInfo))
	}
	if param.FileHash != "" {
		v.Set("file_hash", param.FileHash)
	}
	body := v.Encode()
	return body, nil
}

// AlterAsset Empty price makes the asset with a zero price value. If you don't want to alter the price parameter, set price to -1.
// Empty amount makes the asset with an endless supply of shards. If you don't want to alter the amount parameter, set amount to -1.
func (t *AssetOper) AlterAsset(param *xbase.AlterAssetParam) (*xbase.BaseResp, *xbase.RequestRes, error) {
	if err := param.Valid(); err != nil {
		return nil, nil, err
	}

	body, err := t.genAlterAssetBody(param)
	if err != nil {
		t.Logger.Warn("fail to generate value for altering, err: %v, param: %+v", err, *param)
		return nil, nil, err
	}
	res, err := t.Post(xbase.AssetApiAlter, body)
	if err != nil {
		t.Logger.Warn("post request xasset failed. err: %v", err)
		return nil, nil, xbase.ComErrRequsetFailed
	}
	if res.HttpCode != 200 {
		t.Logger.Warn("post request response is not 200. [http_code: %d] [url: %s] [body: %s] [trace_id: %s]",
			res.HttpCode, res.ReqUrl, res.Body, t.GetTarceId(res.Header))
		return nil, nil, xbase.ComErrRespCodeErr
	}

	var resp xbase.BaseResp
	err = json.Unmarshal([]byte(res.Body), &resp)
	if err != nil {
		t.Logger.Warn("unmarshal body failed. [http_code: %d] [url: %s] [body: %s] [trace_id: %s]",
			res.HttpCode, res.ReqUrl, res.Body, t.GetTarceId(res.Header))
		return nil, res, xbase.ComErrUnmarshalBodyFailed
	}
	if resp.Errno != xbase.XassetErrNoSucc {
		t.Logger.Warn("get resp failed. [url: %s] [request_id: %s] [err_no: %d] [trace_id: %s]",
			res.ReqUrl, resp.RequestId, resp.Errno, t.GetTarceId(res.Header))
		return nil, res, xbase.ComErrServRespErrnoErr
	}

	t.Logger.Trace("operate succ. [asset_id: %d] [url: %s] [request_id: %s] [trace_id: %s]",
		param.AssetId, res.ReqUrl, resp.RequestId, t.GetTarceId(res.Header))
	return &resp, res, nil
}

// GenPublishAssetBody uses the parameter as follows,
//
//	{
//			AssetId    int64  `json:"asset_id"`
//			Addr       string `json:"addr"`
//			Sign       string `json:"sign"`
//			PKey       string `json:"pkey"`
//			Nonce      int64  `json:"nonce"`
//		    IsEvidence int    `json:"is_evidence,omitempty"`
//	}
func (t *AssetOper) genPublishAssetBody(param *xbase.PublishAssetParam) (string, error) {
	nonce := utils.GenNonce()
	signMsg := fmt.Sprintf("%d%d", param.AssetId, nonce)
	sign, err := auth.XassetSignECDSA(param.Account.PrivateKey, []byte(signMsg))
	if err != nil {
		return "", xbase.ComErrAccountSignFailed
	}

	v := url.Values{}
	v.Set("asset_id", fmt.Sprintf("%d", param.AssetId))
	v.Set("addr", param.Account.Address)
	v.Set("sign", sign)
	v.Set("pkey", param.Account.PublicKey)
	v.Set("nonce", fmt.Sprintf("%d", nonce))
	v.Set("is_evidence", fmt.Sprintf("%d", param.IsEvidence))
	body := v.Encode()
	return body, nil
}

func (t *AssetOper) PublishAsset(param *xbase.PublishAssetParam) (*xbase.BaseResp, *xbase.RequestRes, error) {
	if err := param.Valid(); err != nil {
		return nil, nil, err
	}
	body, err := t.genPublishAssetBody(param)
	if err != nil {
		t.Logger.Warn("fail to generate value for publishing, err: %v, param: %+v", err, *param)
		return nil, nil, err
	}
	res, err := t.Post(xbase.AssetApiPublish, body)
	if err != nil {
		t.Logger.Warn("post request xasset failed. err: %v", err)
		return nil, nil, xbase.ComErrRequsetFailed
	}
	if res.HttpCode != 200 {
		t.Logger.Warn("post request response is not 200. [http_code: %d] [url: %s] [body: %s] [trace_id: %s]",
			res.HttpCode, res.ReqUrl, res.Body, t.GetTarceId(res.Header))
		return nil, nil, xbase.ComErrRespCodeErr
	}

	var resp xbase.BaseResp
	err = json.Unmarshal([]byte(res.Body), &resp)
	if err != nil {
		t.Logger.Warn("unmarshal body failed. [http_code: %d] [url: %s] [body: %s] [trace_id: %s]",
			res.HttpCode, res.ReqUrl, res.Body, t.GetTarceId(res.Header))
		return nil, res, xbase.ComErrUnmarshalBodyFailed
	}
	if resp.Errno != xbase.XassetErrNoSucc {
		t.Logger.Warn("get resp failed. [url: %s] [request_id: %s] [err_no: %d] [trace_id: %s]",
			res.ReqUrl, resp.RequestId, resp.Errno, t.GetTarceId(res.Header))
		return nil, res, xbase.ComErrServRespErrnoErr
	}

	t.Logger.Trace("operate succ. [asset_id: %d] [url: %s] [request_id: %s] [trace_id: %s]",
		param.AssetId, res.ReqUrl, resp.RequestId, t.GetTarceId(res.Header))
	return &resp, res, nil
}

// GenQueryAssetBody uses the parameter as follows,
//
//	{
//			AssetId int64 `json:"asset_id"`
//	}
func (t *AssetOper) genQueryAssetBody(param *xbase.QueryAssetParam) (string, error) {
	v := url.Values{}
	v.Set("asset_id", fmt.Sprintf("%d", param.AssetId))
	body := v.Encode()
	return body, nil
}

func (t *AssetOper) QueryAsset(param *xbase.QueryAssetParam) (*xbase.QueryAssetResp, *xbase.RequestRes, error) {
	if err := param.Valid(); err != nil {
		return nil, nil, err
	}
	body, _ := t.genQueryAssetBody(param)

	res, err := t.Post(xbase.AssetApiQueryAsset, body)
	if err != nil {
		t.Logger.Warn("post request xasset failed. err: %v", err)
		return nil, nil, xbase.ComErrRequsetFailed
	}
	if res.HttpCode != 200 {
		t.Logger.Warn("post request response is not 200. [http_code: %d] [url: %s] [body: %s] [trace_id: %s]",
			res.HttpCode, res.ReqUrl, res.Body, t.GetTarceId(res.Header))
		return nil, nil, xbase.ComErrRespCodeErr
	}

	var resp xbase.QueryAssetResp
	err = json.Unmarshal([]byte(res.Body), &resp)
	if err != nil {
		t.Logger.Warn("unmarshal body failed. [http_code: %d] [url: %s] [body: %s] [trace_id: %s]",
			res.HttpCode, res.ReqUrl, res.Body, t.GetTarceId(res.Header))
		return nil, res, xbase.ComErrUnmarshalBodyFailed
	}
	if resp.Errno != xbase.XassetErrNoSucc {
		t.Logger.Warn("get resp failed. [url: %s] [request_id: %s] [err_no: %d] [trace_id: %s]",
			res.ReqUrl, resp.RequestId, resp.Errno, t.GetTarceId(res.Header))
		return nil, res, xbase.ComErrServRespErrnoErr
	}

	t.Logger.Trace("operate succ. [meta: %+v] [url: %s] [request_id: %s] [trace_id: %s]",
		resp.Meta, res.ReqUrl, resp.RequestId, t.GetTarceId(res.Header))
	return &resp, res, nil
}

// GenListAssetsByAddrBody uses the general parameter as follows,
//
//	   {
//			   Addr   string `json:"addr"`
//			   Page   int    `json:"page"`
//			   Limit  int    `json:"limit"`
//			   Status int	 `json:"status"`
//		  }
func (t *AssetOper) genListAssetByAddrBody(param *xbase.ListAssetsByAddrParam) (string, error) {
	v := url.Values{}
	v.Set("addr", param.Addr)
	v.Set("status", fmt.Sprintf("%d", param.Status))
	v.Set("page", fmt.Sprintf("%d", param.Page))
	v.Set("limit", fmt.Sprintf("%d", param.Limit))
	body := v.Encode()
	return body, nil
}

func (t *AssetOper) ListAssetsByAddr(param *xbase.ListAssetsByAddrParam) (*xbase.ListAssetsByAddrResp, *xbase.RequestRes, error) {
	if err := param.Valid(); err != nil {
		return nil, nil, err
	}
	body, _ := t.genListAssetByAddrBody(param)

	res, err := t.Post(xbase.AssetApiListAssetByAddr, body)
	if err != nil {
		t.Logger.Warn("post request xasset failed. err: %v", err)
		return nil, nil, xbase.ComErrRequsetFailed
	}
	if res.HttpCode != 200 {
		t.Logger.Warn("post request response is not 200. [http_code: %d] [url: %s] [body: %s] [trace_id: %s]",
			res.HttpCode, res.ReqUrl, res.Body, t.GetTarceId(res.Header))
		return nil, nil, xbase.ComErrRespCodeErr
	}

	var resp xbase.ListAssetsByAddrResp
	err = json.Unmarshal([]byte(res.Body), &resp)
	if err != nil {
		t.Logger.Warn("unmarshal body failed. [http_code: %d] [url: %s] [body: %s] [trace_id: %s]",
			res.HttpCode, res.ReqUrl, res.Body, t.GetTarceId(res.Header))
		return nil, res, xbase.ComErrUnmarshalBodyFailed
	}
	if resp.Errno != xbase.XassetErrNoSucc {
		t.Logger.Warn("get resp failed. [url: %s] [request_id: %s] [err_no: %d] [trace_id: %s]",
			res.ReqUrl, resp.RequestId, resp.Errno, t.GetTarceId(res.Header))
		return nil, res, xbase.ComErrServRespErrnoErr
	}

	t.Logger.Trace("operate succ. [total_cnt: %d] [url: %s] [request_id: %s] [trace_id: %s]", resp.TotalCnt,
		res.ReqUrl, resp.RequestId, t.GetTarceId(res.Header))
	return &resp, res, nil
}

// GenListDiffByAddrBody uses the general parameter as follows,
//
//	   {
//		    Addr   string `json:"addr"`
//	   	Limit  int    `json:"limit"`
//	  	Cursor string `json:"cursor"`
//	   	OpTyps string `json:"op_types"`
//		  }
func (t *AssetOper) genListDiffByAddrBody(param *xbase.ListDiffByAddrParam) (string, error) {
	v := url.Values{}
	v.Set("addr", param.Addr)
	if param.Limit > 0 {
		v.Set("limit", fmt.Sprintf("%d", param.Limit))
	}
	if param.Cursor != "" {
		v.Set("cursor", param.Cursor)
	}
	if param.OpTyps != "" {
		v.Set("op_types", param.OpTyps)
	}
	body := v.Encode()
	return body, nil
}

func (t *AssetOper) ListDiffByAddr(param *xbase.ListDiffByAddrParam) (*xbase.ListDiffByAddrResp, *xbase.RequestRes, error) {
	if err := param.Valid(); err != nil {
		return nil, nil, err
	}
	body, _ := t.genListDiffByAddrBody(param)

	res, err := t.Post(xbase.AssetApiListDiffByAddr, body)
	if err != nil {
		t.Logger.Warn("post request xasset failed. err: %v", err)
		return nil, nil, xbase.ComErrRequsetFailed
	}
	if res.HttpCode != 200 {
		t.Logger.Warn("post req resp not 200.[http_code: %d] [url: %s] [body: %s] [trace_id: %s]",
			res.HttpCode, res.ReqUrl, res.Body, t.GetTarceId(res.Header))
		return nil, nil, xbase.ComErrRespCodeErr
	}

	var resp xbase.ListDiffByAddrResp
	err = json.Unmarshal([]byte(res.Body), &resp)
	if err != nil {
		t.Logger.Warn("unmarshal body failed.err:%v [http_code: %d] [url: %s] [body: %s] [trace_id: %s]",
			err, res.HttpCode, res.ReqUrl, res.Body, t.GetTarceId(res.Header))
		return nil, res, xbase.ComErrUnmarshalBodyFailed
	}
	if resp.Errno != xbase.XassetErrNoSucc {
		t.Logger.Warn("get resp failed. [url: %s] [request_id: %s] [err_no: %d] [trace_id: %s]",
			res.ReqUrl, resp.RequestId, resp.Errno, t.GetTarceId(res.Header))
		return nil, res, xbase.ComErrServRespErrnoErr
	}

	t.Logger.Trace("operate succ. [url: %s] [request_id: %s] [trace_id: %s]",
		res.ReqUrl, resp.RequestId, t.GetTarceId(res.Header))
	return &resp, res, nil
}

// GenGrantAssetBody Grant uses the general parameter as follows,
//
//	   {
//			   AssetId  int64  `json:"asset_id"`
//			   ShardId  int64  `json:"shard_id"`
//			   Addr     string `json:"addr"`
//			   Sign     string `json:"sign"`
//			   PKey     string `json:"pkey"`
//			   Nonce    int64  `json:"nonce"`
//			   ToAddr   string `json:"to_addr"`
//			   ToUserId int64  `json:"to_userid,omitempty"`
//	 	   Price 	int64  `json:"price",omitempty`
//		  }
func (t *AssetOper) genGrantAssetBody(appid int64, param *xbase.GrantAssetParam) (string, error) {
	nonce := utils.GenNonce()
	signMsg := fmt.Sprintf("%d%d", param.AssetId, nonce)
	sign, err := auth.XassetSignECDSA(param.Account.PrivateKey, []byte(signMsg))
	if err != nil {
		return "", xbase.ComErrAccountSignFailed
	}

	// 未指定shard_id，生成一个唯一值
	shardId := param.ShardId
	if shardId < 1 {
		shardId = utils.GenNonce()
	}

	v := url.Values{}
	v.Set("asset_id", fmt.Sprintf("%d", param.AssetId))
	v.Set("shard_id", fmt.Sprintf("%d", shardId))
	v.Set("price", fmt.Sprintf("%d", param.Price))
	v.Set("addr", param.Addr)
	v.Set("sign", sign)
	v.Set("pkey", param.Account.PublicKey)
	v.Set("nonce", fmt.Sprintf("%d", nonce))
	v.Set("to_addr", param.ToAddr)
	if err := xbase.IdValid(param.ToUserId); err == nil {
		v.Set("to_userid", fmt.Sprintf("%d", param.ToUserId))
	}
	return v.Encode(), nil
}

// GrantAsset grants a random shard to the specific address for the very first time after the maker publishes its asset.
func (t *AssetOper) GrantAsset(param *xbase.GrantAssetParam) (*xbase.GrantAssetResp, *xbase.RequestRes, error) {
	if err := param.Valid(); err != nil {
		return nil, nil, err
	}

	body, err := t.genGrantAssetBody(t.GetConfig().Credentials.AppId, param)
	if err != nil {
		t.Logger.Warn("fail to generate value for granting, err: %v, param: %+v", err, *param)
		return nil, nil, err
	}
	res, err := t.Post(xbase.AssetApiGrant, body)
	if err != nil {
		t.Logger.Warn("post request xasset failed, uri: %s, asset_id: %d, err: %v", xbase.AssetApiGrant, param.AssetId, err)
		return nil, nil, xbase.ComErrRequsetFailed
	}
	if res.HttpCode != 200 {
		t.Logger.Warn("post request response is not 200. [http_code: %d] [url: %s] [asset_id: %d] [body: %s] [trace_id: %s]",
			res.HttpCode, res.ReqUrl, param.AssetId, res.Body, t.GetTarceId(res.Header))
		return nil, nil, xbase.ComErrRespCodeErr
	}

	var resp xbase.GrantAssetResp
	err = json.Unmarshal([]byte(res.Body), &resp)
	if err != nil {
		t.Logger.Warn("unmarshal body failed. [http_code: %d] [url: %s] [asset_id: %d] [body: %s] [trace_id: %s]",
			res.HttpCode, res.ReqUrl, param.AssetId, res.Body, t.GetTarceId(res.Header))
		return nil, res, xbase.ComErrUnmarshalBodyFailed
	}
	if resp.Errno != xbase.XassetErrNoSucc {
		t.Logger.Warn("get resp failed. [url: %s] [asset_id: %d] [request_id: %s] [err_no: %d] [trace_id: %s]",
			res.ReqUrl, param.AssetId, resp.RequestId, resp.Errno, t.GetTarceId(res.Header))
		return nil, res, xbase.ComErrServRespErrnoErr
	}

	t.Logger.Trace("operate succ. [asset_id: %v] [shard_id: %v] [from: %s] [to: %s] [url: %s] [request_id: %s] [trace_id: %s]",
		resp.AssetId, resp.ShardId, param.Addr, param.ToAddr, res.ReqUrl, resp.RequestId, t.GetTarceId(res.Header))
	return &resp, res, nil
}

// GenTransferAssetBody uses the general parameter as follows,
//
//	   {
//			   AssetId  int64  `json:"asset_id"`
//			   ShardId  int64  `json:"shard_id"`
//			   Price	int64  `json:"price"`
//			   Addr     string `json:"addr"`
//			   Sign     string `json:"sign"`
//			   PKey     string `json:"pkey"`
//			   Nonce    int64  `json:"nonce"`
//			   ToAddr   string `json:"to_addr"`
//			   ToUserId int64  `json:"to_userid,omitempty"`
//		  }
func (t *AssetOper) genTransferAssetBody(param *xbase.TransferAssetParam) (string, error) {
	nonce := utils.GenNonce()
	signMsg := fmt.Sprintf("%d%d", param.AssetId, nonce)
	sign, err := auth.XassetSignECDSA(param.Account.PrivateKey, []byte(signMsg))
	if err != nil {
		return "", xbase.ComErrAccountSignFailed
	}

	v := url.Values{}
	v.Set("asset_id", fmt.Sprintf("%d", param.AssetId))
	v.Set("shard_id", fmt.Sprintf("%d", param.ShardId))
	v.Set("price", fmt.Sprintf("%d", param.Price))
	v.Set("addr", param.Addr)
	v.Set("sign", sign)
	v.Set("pkey", param.Account.PublicKey)
	v.Set("nonce", fmt.Sprintf("%d", nonce))
	v.Set("to_addr", param.ToAddr)
	if err := xbase.IdValid(param.ToUserId); err == nil {
		v.Set("to_userid", fmt.Sprintf("%d", param.ToUserId))
	}
	return v.Encode(), nil
}

// GrantAsset transfer th specific shard from address A to address B.
func (t *AssetOper) TransferAsset(param *xbase.TransferAssetParam) (*xbase.BaseResp, *xbase.RequestRes, error) {
	if err := param.Valid(); err != nil {
		return nil, nil, err
	}

	body, err := t.genTransferAssetBody(param)
	if err != nil {
		t.Logger.Warn("fail to generate value for transferring, err: %v, param: %+v", err, *param)
		return nil, nil, err
	}
	res, err := t.Post(xbase.AssetApiTransfer, body)
	if err != nil {
		t.Logger.Warn("post request xasset failed, uri: %s, err: %v", xbase.AssetApiTransfer, err)
		return nil, nil, xbase.ComErrRequsetFailed
	}
	if res.HttpCode != 200 {
		t.Logger.Warn("post request response is not 200. [http_code: %d] [url: %s] [body: %s] [trace_id: %s]",
			res.HttpCode, res.ReqUrl, res.Body, t.GetTarceId(res.Header))
		return nil, nil, xbase.ComErrRespCodeErr
	}

	var resp xbase.BaseResp
	err = json.Unmarshal([]byte(res.Body), &resp)
	if err != nil {
		t.Logger.Warn("unmarshal body failed. [http_code: %d] [url: %s] [body: %s] [trace_id: %s]",
			res.HttpCode, res.ReqUrl, res.Body, t.GetTarceId(res.Header))
		return nil, res, xbase.ComErrUnmarshalBodyFailed
	}
	if resp.Errno != xbase.XassetErrNoSucc {
		t.Logger.Warn("get resp failed. [url: %s] [request_id: %s] [err_no: %d] [trace_id: %s]",
			res.ReqUrl, resp.RequestId, resp.Errno, t.GetTarceId(res.Header))
		return nil, res, xbase.ComErrServRespErrnoErr
	}

	t.Logger.Trace("operate succ. [asset_id: %d] [shard_id: %d] [from: %s] [to: %s] [url: %s] [request_id: %s] [trace_id: %s]",
		param.AssetId, param.ShardId, param.Addr, param.ToAddr, res.ReqUrl, resp.RequestId, t.GetTarceId(res.Header))
	return &resp, res, nil
}

// GenQueryShardsBody uses the general parameter as follows,
//
//	   {
//			   AssetId  int64  `json:"asset_id"`
//			   ShardId  int64  `json:"shard_id"`
//		  }
func (t *AssetOper) genQueryShardsBody(param *xbase.QueryShardParam) (string, error) {
	v := url.Values{}
	v.Set("asset_id", fmt.Sprintf("%d", param.AssetId))
	v.Set("shard_id", fmt.Sprintf("%d", param.ShardId))
	body := v.Encode()
	return body, nil
}

func (t *AssetOper) QueryShard(param *xbase.QueryShardParam) (*xbase.QueryShardResp, *xbase.RequestRes, error) {
	if err := param.Valid(); err != nil {
		return nil, nil, err
	}
	body, _ := t.genQueryShardsBody(param)

	res, err := t.Post(xbase.AssetApiQueryShard, body)
	if err != nil {
		t.Logger.Warn("post request xasset failed. err: %v", err)
		return nil, nil, xbase.ComErrRequsetFailed
	}
	if res.HttpCode != 200 {
		t.Logger.Warn("post request response is not 200. [http_code: %d] [url: %s] [body: %s] [trace_id: %s]",
			res.HttpCode, res.ReqUrl, res.Body, t.GetTarceId(res.Header))
		return nil, nil, xbase.ComErrRespCodeErr
	}

	var resp xbase.QueryShardResp
	err = json.Unmarshal([]byte(res.Body), &resp)
	if err != nil {
		t.Logger.Warn("unmarshal body failed. [http_code: %d] [url: %s] [body: %s] [trace_id: %s]",
			res.HttpCode, res.ReqUrl, res.Body, t.GetTarceId(res.Header))
		return nil, res, xbase.ComErrUnmarshalBodyFailed
	}
	if resp.Errno != xbase.XassetErrNoSucc {
		t.Logger.Warn("get resp failed. [url: %s] [request_id: %s] [err_no: %d] [trace_id: %s]",
			res.ReqUrl, resp.RequestId, resp.Errno, t.GetTarceId(res.Header))
		return nil, res, xbase.ComErrServRespErrnoErr
	}

	t.Logger.Trace("operate succ. [meta:%+v] [url:%s] [request_id:%s] [trace_id:%s]",
		resp.Meta, res.ReqUrl, resp.RequestId, t.GetTarceId(res.Header))
	return &resp, res, nil
}

// GenListShardsByAddrBody uses the general parameter as follows,
//
//	   {
//			   Addr  string `json:"addr"`
//			   Page  int    `json:"page"`
//			   Limit int    `json:"limit"`
//		  }
func (t *AssetOper) genListShardsByAddrBody(param *xbase.ListShardsByAddrParam) (string, error) {
	v := url.Values{}
	v.Set("addr", param.Addr)
	v.Set("page", fmt.Sprintf("%d", param.Page))
	v.Set("limit", fmt.Sprintf("%d", param.Limit))
	if param.AssetId > 0 {
		v.Set("asset_id", fmt.Sprintf("%d", param.AssetId))
	}
	body := v.Encode()
	return body, nil
}

func (t *AssetOper) ListShardsByAddr(param *xbase.ListShardsByAddrParam) (*xbase.ListShardsByAddrResp, *xbase.RequestRes, error) {
	if err := param.Valid(); err != nil {
		return nil, nil, err
	}
	body, _ := t.genListShardsByAddrBody(param)

	res, err := t.Post(xbase.AssetApiListShardsByAddr, body)
	if err != nil {
		t.Logger.Warn("post request xasset failed. err: %v", err)
		return nil, nil, xbase.ComErrRequsetFailed
	}
	if res.HttpCode != 200 {
		t.Logger.Warn("post request response is not 200. [http_code: %d] [url: %s] [body: %s] [trace_id: %s]",
			res.HttpCode, res.ReqUrl, res.Body, t.GetTarceId(res.Header))
		return nil, nil, xbase.ComErrRespCodeErr
	}

	var resp xbase.ListShardsByAddrResp
	err = json.Unmarshal([]byte(res.Body), &resp)
	if err != nil {
		t.Logger.Warn("unmarshal body failed. [http_code: %d] [url: %s] [body: %s] [trace_id: %s]",
			res.HttpCode, res.ReqUrl, res.Body, t.GetTarceId(res.Header))
		return nil, res, xbase.ComErrUnmarshalBodyFailed
	}
	if resp.Errno != xbase.XassetErrNoSucc {
		t.Logger.Warn("get resp failed. [url: %s] [request_id: %s] [err_no: %d] [trace_id: %s]",
			res.ReqUrl, resp.RequestId, resp.Errno, t.GetTarceId(res.Header))
		return nil, res, xbase.ComErrServRespErrnoErr
	}

	t.Logger.Trace("operate succ. [total_cnt: %d] [url: %s] [request_id: %s] [trace_id: %s]", resp.TotalCnt,
		res.ReqUrl, resp.RequestId, t.GetTarceId(res.Header))
	return &resp, res, nil
}

// GenListShardsByAssetBody uses the general parameter as follows,
//
//	   {
//			   AssetId  int64 	`json:"asset_id"`
//			   Cursor   string 	`json:"cursor"`
//			   Limit  	int    	`json:"limit"`
//		  }
func (t *AssetOper) genListShardsByAssetBody(param *xbase.ListShardsByAssetParam) (string, error) {
	v := url.Values{}
	v.Set("asset_id", fmt.Sprintf("%d", param.AssetId))
	v.Set("cursor", param.Cursor)
	v.Set("limit", fmt.Sprintf("%d", param.Limit))
	body := v.Encode()
	return body, nil
}

func (t *AssetOper) ListShardsByAsset(param *xbase.ListShardsByAssetParam) (*xbase.ListShardsByAssetResp, *xbase.RequestRes, error) {
	if err := param.Valid(); err != nil {
		return nil, nil, err
	}
	body, _ := t.genListShardsByAssetBody(param)

	res, err := t.Post(xbase.AssetListShardsByAsset, body)
	if err != nil {
		t.Logger.Warn("post request xasset failed. err: %v", err)
		return nil, nil, xbase.ComErrRequsetFailed
	}
	if res.HttpCode != 200 {
		t.Logger.Warn("post request response is not 200. [http_code: %d] [url: %s] [body: %s] [trace_id: %s]",
			res.HttpCode, res.ReqUrl, res.Body, t.GetTarceId(res.Header))
		return nil, nil, xbase.ComErrRespCodeErr
	}

	var resp xbase.ListShardsByAssetResp
	err = json.Unmarshal([]byte(res.Body), &resp)
	if err != nil {
		t.Logger.Warn("unmarshal body failed. [http_code: %d] [url: %s] [body: %s] [trace_id: %s]",
			res.HttpCode, res.ReqUrl, res.Body, t.GetTarceId(res.Header))
		return nil, res, xbase.ComErrUnmarshalBodyFailed
	}
	if resp.Errno != xbase.XassetErrNoSucc {
		t.Logger.Warn("get resp failed. [url: %s] [request_id: %s] [err_no: %d] [trace_id: %s]",
			res.ReqUrl, resp.RequestId, resp.Errno, t.GetTarceId(res.Header))
		return nil, res, xbase.ComErrServRespErrnoErr
	}

	t.Logger.Trace("operate succ. [cursor: %s] [has_more: %d] [url: %s] [request_id: %s] [trace_id: %s]", resp.Cursor, resp.HasMore,
		res.ReqUrl, resp.RequestId, t.GetTarceId(res.Header))
	return &resp, res, nil
}

func (t *AssetOper) ListAssetHistory(param *xbase.ListAssetHisParam) (*xbase.ListAssetHistoryResp, *xbase.RequestRes, error) {
	if err := param.Valid(); err != nil {
		return nil, nil, err
	}

	v := url.Values{}
	v.Set("asset_id", fmt.Sprintf("%d", param.AssetId))
	v.Set("page", fmt.Sprintf("%d", param.Page))
	if param.Limit > 0 {
		v.Set("limit", fmt.Sprintf("%d", param.Limit))
	}
	if param.ShardId > 0 {
		v.Set("shard_id", fmt.Sprintf("%d", param.ShardId))
	}
	body := v.Encode()

	res, err := t.Post(xbase.ListAssetHistory, body)
	if err != nil {
		t.Logger.Warn("post request xasset failed, uri: %s, err: %v", xbase.ListAssetHistory, err)
		return nil, nil, xbase.ComErrRequsetFailed
	}
	if res.HttpCode != 200 {
		t.Logger.Warn("post request response is not 200. [http_code: %d] [url: %s] [body: %s] [trace_id: %s]",
			res.HttpCode, res.ReqUrl, res.Body, t.GetTarceId(res.Header))
		return nil, nil, xbase.ComErrRespCodeErr
	}

	var resp xbase.ListAssetHistoryResp
	err = json.Unmarshal([]byte(res.Body), &resp)
	if err != nil {
		t.Logger.Warn("unmarshal body failed. [http_code: %d] [url: %s] [body: %s] [trace_id: %s]",
			res.HttpCode, res.ReqUrl, res.Body, t.GetTarceId(res.Header))
		return nil, nil, xbase.ComErrUnmarshalBodyFailed
	}
	t.Logger.Trace("operate succ. [asset_id: %d] [url: %s] [request_id: %s] [trace_id: %s] [resp: %+v]",
		param.AssetId, res.ReqUrl, resp.RequestId, t.GetTarceId(res.Header), resp)
	return &resp, res, nil
}

// GenEvidenceBody uses the general parameter as follows,
//
//	   {
//			   AssetId  int64  `json:"asset_id"`
//		  }
func (t *AssetOper) genEvidenceBody(param *xbase.GetEvidenceInfoParam) (string, error) {
	v := url.Values{}
	v.Set("asset_id", fmt.Sprintf("%d", param.AssetId))
	body := v.Encode()
	return body, nil
}

func (t *AssetOper) GetEvidenceInfo(param *xbase.GetEvidenceInfoParam) (*xbase.GetEvidenceInfoResp, *xbase.RequestRes, error) {
	if err := param.Valid(); err != nil {
		return nil, nil, err
	}
	body, _ := t.genEvidenceBody(param)

	res, err := t.Post(xbase.AssetApiGetEvidenceInfo, body)
	if err != nil {
		t.Logger.Warn("post request xasset failed. err: %v", err)
		return nil, nil, xbase.ComErrRequsetFailed
	}
	if res.HttpCode != 200 {
		t.Logger.Warn("post request response is not 200. [http_code: %d] [url: %s] [body: %s] [trace_id: %s]",
			res.HttpCode, res.ReqUrl, res.Body, t.GetTarceId(res.Header))
		return nil, nil, xbase.ComErrRespCodeErr
	}

	var resp xbase.GetEvidenceInfoResp
	err = json.Unmarshal([]byte(res.Body), &resp)
	if err != nil {
		t.Logger.Warn("unmarshal body failed. [http_code: %d] [url: %s] [body: %s] [trace_id: %s]",
			res.HttpCode, res.ReqUrl, res.Body, t.GetTarceId(res.Header))
		return nil, res, xbase.ComErrUnmarshalBodyFailed
	}
	if resp.Errno != xbase.XassetErrNoSucc {
		t.Logger.Warn("get resp failed. [url: %s] [request_id: %s] [err_no: %d] [trace_id: %s]",
			res.ReqUrl, resp.RequestId, resp.Errno, t.GetTarceId(res.Header))
		return nil, res, xbase.ComErrServRespErrnoErr
	}

	t.Logger.Trace("operate succ. [create_addr: %s] [tx_id: %s] [asset_info: %v] [ctime: %d] [url: %s] [request_id: %s] [trace_id: %s]",
		resp.CreateAddr, resp.TxId, resp.AssetInfo, resp.Ctime, res.ReqUrl, resp.RequestId, t.GetTarceId(res.Header))
	return &resp, res, nil
}

// GenFreezeBody uses the general parameter as follows,
//
//	   {
//			   AssetId  int64  			`json:"asset_id"`
//			   Account  *auth.Account	`json:"account"`
//		  }
func (t *AssetOper) genFreezeAssetBody(param *xbase.FreezeAssetParam) (string, error) {
	nonce := utils.GenNonce()
	signMsg := fmt.Sprintf("%d%d", param.AssetId, nonce)
	sign, err := auth.XassetSignECDSA(param.Account.PrivateKey, []byte(signMsg))
	if err != nil {
		return "", xbase.ComErrAccountSignFailed
	}

	v := url.Values{}
	v.Set("asset_id", fmt.Sprintf("%d", param.AssetId))
	v.Set("addr", param.Account.Address)
	v.Set("sign", sign)
	v.Set("pkey", param.Account.PublicKey)
	v.Set("nonce", fmt.Sprintf("%d", nonce))
	return v.Encode(), nil
}

// FreezeAsset freeze assets where granting action is forbidden.
func (t *AssetOper) FreezeAsset(param *xbase.FreezeAssetParam) (*xbase.BaseResp, *xbase.RequestRes, error) {
	if err := param.Valid(); err != nil {
		return nil, nil, err
	}

	body, err := t.genFreezeAssetBody(param)
	if err != nil {
		t.Logger.Warn("fail to generate value for freeze, err: %v, param: %+v", err, *param)
		return nil, nil, err
	}
	res, err := t.Post(xbase.AssetApiFreeze, body)
	if err != nil {
		t.Logger.Warn("post request xasset failed, uri: %s, err: %v", xbase.AssetApiFreeze, err)
		return nil, nil, xbase.ComErrRequsetFailed
	}
	if res.HttpCode != 200 {
		t.Logger.Warn("post request response is not 200. [http_code: %d] [url: %s] [body: %s] [trace_id: %s]",
			res.HttpCode, res.ReqUrl, res.Body, t.GetTarceId(res.Header))
		return nil, nil, xbase.ComErrRespCodeErr
	}

	var resp xbase.BaseResp
	err = json.Unmarshal([]byte(res.Body), &resp)
	if err != nil {
		t.Logger.Warn("unmarshal body failed. [http_code: %d] [url: %s] [body: %s] [trace_id: %s]",
			res.HttpCode, res.ReqUrl, res.Body, t.GetTarceId(res.Header))
		return nil, res, xbase.ComErrUnmarshalBodyFailed
	}
	if resp.Errno != xbase.XassetErrNoSucc {
		t.Logger.Warn("get resp failed. [url: %s] [request_id: %s] [err_no: %d] [trace_id: %s]",
			res.ReqUrl, resp.RequestId, resp.Errno, t.GetTarceId(res.Header))
		return nil, res, xbase.ComErrServRespErrnoErr
	}

	t.Logger.Trace("operate succ. [asset_id: %d] [url: %s] [request_id: %s] [trace_id: %s]",
		param.AssetId, res.ReqUrl, resp.RequestId, t.GetTarceId(res.Header))
	return &resp, res, nil
}

// GenConsumeBody uses the general parameter as follows,
//
//	   {
//				AssetId  int64         `json:"asset_id"`
//				ShardId  int64         `json:"shard_id"`
//				Nonce    int64         `json:"nonce"`
//				UAddr    string        `json:"user_addr"`
//				USign    string        `json:"user_sign"`
//				UPKey    string        `json:"user_pkey"`
//				CAccount *auth.Account `json:"create_account"`
//		  }
func (t *AssetOper) genConsumeShardBody(param *xbase.ConsumeShardParam) (string, error) {
	signMsg := fmt.Sprintf("%d%d", param.AssetId, param.Nonce)
	sign, err := auth.XassetSignECDSA(param.CAccount.PrivateKey, []byte(signMsg))
	if err != nil {
		return "", xbase.ComErrAccountSignFailed
	}

	v := url.Values{}
	v.Set("asset_id", fmt.Sprintf("%d", param.AssetId))
	v.Set("shard_id", fmt.Sprintf("%d", param.ShardId))
	v.Set("nonce", fmt.Sprintf("%d", param.Nonce))
	v.Set("addr", param.CAccount.Address)
	v.Set("sign", sign)
	v.Set("pkey", param.CAccount.PublicKey)
	v.Set("user_addr", param.UAddr)
	v.Set("user_sign", param.USign)
	v.Set("user_pkey", param.UPKey)

	return v.Encode(), nil
}

// ConsumeShard consumes shards where any other action is forbidden.
func (t *AssetOper) ConsumeShard(param *xbase.ConsumeShardParam) (*xbase.BaseResp, *xbase.RequestRes, error) {
	if err := param.Valid(); err != nil {
		return nil, nil, err
	}

	body, err := t.genConsumeShardBody(param)
	if err != nil {
		t.Logger.Warn("fail to generate value for consume, err: %v, param: %+v", err, *param)
		return nil, nil, err
	}
	res, err := t.Post(xbase.AssetApiConsume, body)
	if err != nil {
		t.Logger.Warn("post request xasset failed, uri: %s, err: %v", xbase.AssetApiConsume, err)
		return nil, nil, xbase.ComErrRequsetFailed
	}
	if res.HttpCode != 200 {
		t.Logger.Warn("post request response is not 200. [http_code: %d] [url: %s] [body: %s] [trace_id: %s]",
			res.HttpCode, res.ReqUrl, res.Body, t.GetTarceId(res.Header))
		return nil, nil, xbase.ComErrRespCodeErr
	}

	var resp xbase.BaseResp
	err = json.Unmarshal([]byte(res.Body), &resp)
	if err != nil {
		t.Logger.Warn("unmarshal body failed. [http_code: %d] [url: %s] [body: %s] [trace_id: %s]",
			res.HttpCode, res.ReqUrl, res.Body, t.GetTarceId(res.Header))
		return nil, res, xbase.ComErrUnmarshalBodyFailed
	}
	if resp.Errno != xbase.XassetErrNoSucc {
		t.Logger.Warn("get resp failed. [url: %s] [request_id: %s] [err_no: %d] [trace_id: %s]",
			res.ReqUrl, resp.RequestId, resp.Errno, t.GetTarceId(res.Header))
		return nil, res, xbase.ComErrServRespErrnoErr
	}

	t.Logger.Trace("operate succ. [asset_id: %d] [shard_id: %d] [url: %s] [request_id: %s] [trace_id: %s]",
		param.AssetId, param.ShardId, res.ReqUrl, resp.RequestId, t.GetTarceId(res.Header))
	return &resp, res, nil
}

// GenSceneListShardByAddrBody uses the general parameter as follows,
//
//	   {
//				Addr  	string     `json:"addr"`
//				Token 	string     `json:"token"`
//				Limit 	int        `json:"limit"`
//				Cursor  string     `json:"cursor"`
//		  }
func (t *AssetOper) genSceneListShardByAddrBody(param *xbase.SceneListShardByAddrParam) (string, error) {
	v := url.Values{}
	v.Set("addr", param.Addr)
	v.Set("token", param.Token)
	v.Set("cursor", param.Cursor)
	if param.Limit > 0 {
		v.Set("limit", fmt.Sprintf("%d", param.Limit))
	}
	return v.Encode(), nil
}

// SceneListShardByAddr list shards under scene authorization.
func (t *AssetOper) SceneListShardByAddr(param *xbase.SceneListShardByAddrParam) (*xbase.SceneListShardByAddrResp, *xbase.RequestRes, error) {
	if err := param.Valid(); err != nil {
		return nil, nil, err
	}

	body, err := t.genSceneListShardByAddrBody(param)
	if err != nil {
		t.Logger.Warn("fail to generate value for scene listshardbyaddr, err: %v, param: %+v", err, *param)
		return nil, nil, err
	}
	res, err := t.Post(xbase.SceneListShardByAddr, body)
	if err != nil {
		t.Logger.Warn("post request xasset failed, uri: %s, err: %v", xbase.SceneListShardByAddr, err)
		return nil, nil, xbase.ComErrRequsetFailed
	}
	if res.HttpCode != 200 {
		t.Logger.Warn("post request response is not 200. [http_code: %d] [url: %s] [body: %s] [trace_id: %s]",
			res.HttpCode, res.ReqUrl, res.Body, t.GetTarceId(res.Header))
		return nil, nil, xbase.ComErrRespCodeErr
	}

	var resp xbase.SceneListShardByAddrResp
	err = json.Unmarshal([]byte(res.Body), &resp)
	if err != nil {
		t.Logger.Warn("unmarshal body failed. [http_code: %d] [url: %s] [body: %s] [trace_id: %s]",
			res.HttpCode, res.ReqUrl, res.Body, t.GetTarceId(res.Header))
		return nil, res, xbase.ComErrUnmarshalBodyFailed
	}
	if resp.Errno != xbase.XassetErrNoSucc {
		t.Logger.Warn("get resp failed. [url: %s] [request_id: %s] [err_no: %d] [trace_id: %s]",
			res.ReqUrl, resp.RequestId, resp.Errno, t.GetTarceId(res.Header))
		return nil, res, xbase.ComErrServRespErrnoErr
	}

	t.Logger.Trace("operate succ. [addr: %s] [url: %s] [request_id: %s] [trace_id: %s]",
		param.Addr, res.ReqUrl, resp.RequestId, t.GetTarceId(res.Header))
	return &resp, res, nil
}

// GenSceneQueryShardBody uses the general parameter as follows,
//
//	   {
//				Addr  		string     	`json:"addr"`
//				Token 		string     	`json:"token"`
//				AssetId 	int64 		`json: "asset_id"`
//				ShardId 	int64 		`json: "shard_id"`
//		  }
func (t *AssetOper) genSceneQueryShardBody(param *xbase.SceneQueryShardParam) (string, error) {
	v := url.Values{}
	v.Set("addr", param.Addr)
	v.Set("token", param.Token)
	v.Set("asset_id", fmt.Sprintf("%d", param.AssetId))
	v.Set("shard_id", fmt.Sprintf("%d", param.ShardId))

	return v.Encode(), nil
}

// SceneQueryShard query shard under scene authorization.
func (t *AssetOper) SceneQueryShard(param *xbase.SceneQueryShardParam) (*xbase.SceneQueryShardResp, *xbase.RequestRes, error) {
	if err := param.Valid(); err != nil {
		return nil, nil, err
	}

	body, err := t.genSceneQueryShardBody(param)
	if err != nil {
		t.Logger.Warn("fail to generate value for scene queryshard, err: %v, param: %+v", err, *param)
		return nil, nil, err
	}
	res, err := t.Post(xbase.SceneQueryShard, body)
	if err != nil {
		t.Logger.Warn("post request xasset failed, uri: %s, err: %v", xbase.SceneQueryShard, err)
		return nil, nil, xbase.ComErrRequsetFailed
	}
	if res.HttpCode != 200 {
		t.Logger.Warn("post request response is not 200. [http_code: %d] [url: %s] [body: %s] [trace_id: %s]",
			res.HttpCode, res.ReqUrl, res.Body, t.GetTarceId(res.Header))
		return nil, nil, xbase.ComErrRespCodeErr
	}

	var resp xbase.SceneQueryShardResp
	err = json.Unmarshal([]byte(res.Body), &resp)
	if err != nil {
		t.Logger.Warn("unmarshal body failed. [http_code: %d] [url: %s] [body: %s] [trace_id: %s]",
			res.HttpCode, res.ReqUrl, res.Body, t.GetTarceId(res.Header))
		return nil, res, xbase.ComErrUnmarshalBodyFailed
	}
	if resp.Errno != xbase.XassetErrNoSucc {
		t.Logger.Warn("get resp failed. [url: %s] [request_id: %s] [err_no: %d] [trace_id: %s]",
			res.ReqUrl, resp.RequestId, resp.Errno, t.GetTarceId(res.Header))
		return nil, res, xbase.ComErrServRespErrnoErr
	}

	t.Logger.Trace("operate succ. [addr: %s] [asset_id: %d] [shard_id: %d] [url: %s] [request_id: %s] [trace_id: %s]",
		param.Addr, param.AssetId, param.ShardId, res.ReqUrl, resp.RequestId, t.GetTarceId(res.Header))
	return &resp, res, nil
}

// GenSceneListDiffByAddrBody uses the general parameter as follows,
//
//	   {
//		    Addr   string `json:"addr"`
//			Token  string `json:"token"`
//	   	Limit  int    `json:"limit"`
//	  	Cursor string `json:"cursor"`
//	   	OpTyps string `json:"op_types"`
//		  }
func (t *AssetOper) genSceneListDiffByAddrBody(param *xbase.SceneListDiffByAddrParam) (string, error) {
	v := url.Values{}
	v.Set("addr", param.Addr)
	v.Set("token", param.Token)
	if param.Limit > 0 {
		v.Set("limit", fmt.Sprintf("%d", param.Limit))
	}
	if param.Cursor != "" {
		v.Set("cursor", param.Cursor)
	}
	if param.OpTyps != "" {
		v.Set("op_types", param.OpTyps)
	}
	body := v.Encode()
	return body, nil
}

func (t *AssetOper) SceneListDiffByAddr(param *xbase.SceneListDiffByAddrParam) (*xbase.ListDiffByAddrResp, *xbase.RequestRes, error) {
	if err := param.Valid(); err != nil {
		return nil, nil, err
	}
	body, _ := t.genSceneListDiffByAddrBody(param)

	res, err := t.Post(xbase.SceneListDiffByAddr, body)
	if err != nil {
		t.Logger.Warn("post request xasset failed. err: %v", err)
		return nil, nil, xbase.ComErrRequsetFailed
	}
	if res.HttpCode != 200 {
		t.Logger.Warn("post req resp not 200.[http_code: %d] [url: %s] [body: %s] [trace_id: %s]",
			res.HttpCode, res.ReqUrl, res.Body, t.GetTarceId(res.Header))
		return nil, nil, xbase.ComErrRespCodeErr
	}

	var resp xbase.ListDiffByAddrResp
	err = json.Unmarshal([]byte(res.Body), &resp)
	if err != nil {
		t.Logger.Warn("unmarshal body failed.err:%v [http_code: %d] [url: %s] [body: %s] [trace_id: %s]",
			err, res.HttpCode, res.ReqUrl, res.Body, t.GetTarceId(res.Header))
		return nil, res, xbase.ComErrUnmarshalBodyFailed
	}
	if resp.Errno != xbase.XassetErrNoSucc {
		t.Logger.Warn("get resp failed. [url: %s] [request_id: %s] [err_no: %d] [trace_id: %s]",
			res.ReqUrl, resp.RequestId, resp.Errno, t.GetTarceId(res.Header))
		return nil, res, xbase.ComErrServRespErrnoErr
	}

	t.Logger.Trace("operate succ. [url: %s] [request_id: %s] [trace_id: %s]",
		res.ReqUrl, resp.RequestId, t.GetTarceId(res.Header))
	return &resp, res, nil
}

// SceneHasAssetByAddrBody
//
//	type SceneHasAssetByAddrParam struct {
//		Addr     string `json:"addr"`
//		Token    string `json:"token"`
//		AssetIds string `json:"asset_ids"`
//	}
func (t *AssetOper) genSceneHasAssetByAddrBody(param *xbase.SceneHasAssetByAddrParam) (string, error) {
	v := url.Values{}
	v.Set("addr", param.Addr)
	v.Set("token", param.Token)
	v.Set("asset_ids", param.AssetIds)
	body := v.Encode()
	return body, nil
}

func (t *AssetOper) SceneHasAssetByAddr(param *xbase.SceneHasAssetByAddrParam) (*xbase.SceneHasAssetByAddrResp, *xbase.RequestRes, error) {
	if err := param.Valid(); err != nil {
		return nil, nil, err
	}
	body, _ := t.genSceneHasAssetByAddrBody(param)

	res, err := t.Post(xbase.SceneHasAstByAddr, body)
	if err != nil {
		t.Logger.Warn("post request xasset failed. url: %s, err: %v", xbase.SceneHasAstByAddr, err)
		return nil, nil, xbase.ComErrRequsetFailed
	}
	if res.HttpCode != 200 {
		t.Logger.Warn("post req resp not 200.[http_code: %d] [url: %s] [body: %s] [trace_id: %s]",
			res.HttpCode, res.ReqUrl, res.Body, t.GetTarceId(res.Header))
		return nil, nil, xbase.ComErrRespCodeErr
	}

	var resp xbase.SceneHasAssetByAddrResp
	err = json.Unmarshal([]byte(res.Body), &resp)
	if err != nil {
		t.Logger.Warn("unmarshal body failed.err:%v [http_code: %d] [url: %s] [body: %s] [trace_id: %s]",
			err, res.HttpCode, res.ReqUrl, res.Body, t.GetTarceId(res.Header))
		return nil, res, xbase.ComErrUnmarshalBodyFailed
	}
	if resp.Errno != xbase.XassetErrNoSucc {
		t.Logger.Warn("get resp failed. [url: %s] [request_id: %s] [err_no: %d] [trace_id: %s]",
			res.ReqUrl, resp.RequestId, resp.Errno, t.GetTarceId(res.Header))
		return nil, res, xbase.ComErrServRespErrnoErr
	}

	t.Logger.Trace("operate succ. [addr: %s] [token: %s] [url: %s] [request_id: %s] [trace_id: %s]",
		param.Addr, param.Token, res.ReqUrl, resp.RequestId, t.GetTarceId(res.Header))
	return &resp, res, nil
}

func (t *AssetOper) SceneListAddr(uid string) (*xbase.SceneListAddrResp, *xbase.RequestRes, error) {
	if err := xbase.UnionIdValid(uid); err != nil {
		return nil, nil, err
	}
	signedUnionId, err := t.aesEncodeStr(uid)
	if err != nil {
		t.Logger.Warn("encode union id fail, union id: %s", uid)
		return nil, nil, err
	}
	v := url.Values{}
	v.Set("union_id", signedUnionId)
	body := v.Encode()

	res, err := t.Post(xbase.SceneListAddr, body)
	if err != nil {
		t.Logger.Warn("post request xasset failed. url: %s, err: %v", xbase.SceneListAddr, err)
		return nil, nil, xbase.ComErrRequsetFailed
	}
	if res.HttpCode != 200 {
		t.Logger.Warn("post req resp not 200.[http_code: %d] [url: %s] [body: %s] [trace_id: %s]",
			res.HttpCode, res.ReqUrl, res.Body, t.GetTarceId(res.Header))
		return nil, nil, xbase.ComErrRespCodeErr
	}

	var resp xbase.SceneListAddrResp
	err = json.Unmarshal([]byte(res.Body), &resp)
	if err != nil {
		t.Logger.Warn("unmarshal body failed.err:%v [http_code: %d] [url: %s] [body: %s] [trace_id: %s]",
			err, res.HttpCode, res.ReqUrl, res.Body, t.GetTarceId(res.Header))
		return nil, res, xbase.ComErrUnmarshalBodyFailed
	}
	if resp.Errno != xbase.XassetErrNoSucc {
		t.Logger.Warn("get resp failed. [url: %s] [request_id: %s] [err_no: %d] [trace_id: %s]",
			res.ReqUrl, resp.RequestId, resp.Errno, t.GetTarceId(res.Header))
		return nil, res, xbase.ComErrServRespErrnoErr
	}

	t.Logger.Trace("operate succ. [union_id: %s] [url: %s] [request_id: %s] [trace_id: %s]",
		uid, res.ReqUrl, resp.RequestId, t.GetTarceId(res.Header))
	return &resp, res, nil
}

func (t *AssetOper) BdBoxRegister(param *xbase.BdBoxRegisterParam) (*xbase.BdBoxRegisterResp, *xbase.RequestRes, error) {
	if err := param.Valid(); err != nil {
		return nil, nil, err
	}
	v := url.Values{}
	signedOpenId, err := t.aesEncodeStr(param.OpenId)
	if err != nil {
		t.Logger.Warn("encode open id fail, open id: %s", param.OpenId)
		return nil, nil, err
	}
	signedAppKey, err := t.aesEncodeStr(param.AppKey)
	if err != nil {
		t.Logger.Warn("encode app key fail, app key: %s", param.AppKey)
		return nil, nil, err
	}
	v.Set("open_id", signedOpenId)
	v.Set("app_key", signedAppKey)
	body := v.Encode()

	res, err := t.Post(xbase.DidApiRegister, body)
	if err != nil {
		t.Logger.Warn("post request xasset failed. url: %s, err: %v", xbase.DidApiRegister, err)
		return nil, nil, xbase.ComErrRequsetFailed
	}
	if res.HttpCode != 200 {
		t.Logger.Warn("post req resp not 200.[http_code: %d] [url: %s] [body: %s] [trace_id: %s]",
			res.HttpCode, res.ReqUrl, res.Body, t.GetTarceId(res.Header))
		return nil, nil, xbase.ComErrRespCodeErr
	}

	var resp xbase.BdBoxRegisterResp
	err = json.Unmarshal([]byte(res.Body), &resp)
	if err != nil {
		t.Logger.Warn("unmarshal body failed.err:%v [http_code: %d] [url: %s] [body: %s] [trace_id: %s]",
			err, res.HttpCode, res.ReqUrl, res.Body, t.GetTarceId(res.Header))
		return nil, res, xbase.ComErrUnmarshalBodyFailed
	}

	if resp.Errno != xbase.XassetErrNoSucc {
		t.Logger.Warn("get resp failed. [url: %s] [request_id: %s] [err_no: %d] [trace_id: %s]",
			res.ReqUrl, resp.RequestId, resp.Errno, t.GetTarceId(res.Header))
		return nil, res, xbase.ComErrServRespErrnoErr
	}

	decodeMnem, err := t.aesDecodeStr(resp.Mnemonic)
	if err != nil {
		t.Logger.Warn("get resp succ but cannot decode mnemonic. [url: %s] [request_id: %s] [trace_id: %s]",
			res.ReqUrl, resp.RequestId, t.GetTarceId(res.Header))
		return &resp, res, err
	}
	resp.Mnemonic = decodeMnem

	t.Logger.Trace("operate succ. [open_id: %s] [app_key: %s] [url: %s] [request_id: %s] [trace_id: %s]",
		param.OpenId, param.AppKey, res.ReqUrl, resp.RequestId, t.GetTarceId(res.Header))
	return &resp, res, nil
}

func (t *AssetOper) BdBoxBind(param *xbase.BdBoxBindParam) (*xbase.BaseResp, *xbase.RequestRes, error) {
	if err := param.Valid(); err != nil {
		return nil, nil, err
	}
	v := url.Values{}
	signedOpenId, err := t.aesEncodeStr(param.OpenId)
	if err != nil {
		t.Logger.Warn("encode open id fail, open id: %s", param.OpenId)
		return nil, nil, err
	}
	signedAppKey, err := t.aesEncodeStr(param.AppKey)
	if err != nil {
		t.Logger.Warn("encode app key fail, app key: %s", param.AppKey)
		return nil, nil, err
	}
	signedMnem, err := t.aesEncodeStr(param.Mnemonic)
	if err != nil {
		t.Logger.Warn("encode mnemonic fail, mnemonic: %s", param.Mnemonic)
		return nil, nil, err
	}
	v.Set("open_id", signedOpenId)
	v.Set("app_key", signedAppKey)
	v.Set("mnemonic", signedMnem)
	body := v.Encode()

	res, err := t.Post(xbase.DidApiBind, body)
	if err != nil {
		t.Logger.Warn("post request xasset failed. url: %s, err: %v", xbase.DidApiBind, err)
		return nil, nil, xbase.ComErrRequsetFailed
	}
	if res.HttpCode != 200 {
		t.Logger.Warn("post req resp not 200.[http_code: %d] [url: %s] [body: %s] [trace_id: %s]",
			res.HttpCode, res.ReqUrl, res.Body, t.GetTarceId(res.Header))
		return nil, nil, xbase.ComErrRespCodeErr
	}

	var resp xbase.BaseResp
	err = json.Unmarshal([]byte(res.Body), &resp)
	if err != nil {
		t.Logger.Warn("unmarshal body failed.err:%v [http_code: %d] [url: %s] [body: %s] [trace_id: %s]",
			err, res.HttpCode, res.ReqUrl, res.Body, t.GetTarceId(res.Header))
		return nil, res, xbase.ComErrUnmarshalBodyFailed
	}
	if resp.Errno != xbase.XassetErrNoSucc {
		t.Logger.Warn("get resp failed. [url: %s] [request_id: %s] [err_no: %d] [trace_id: %s]",
			res.ReqUrl, resp.RequestId, resp.Errno, t.GetTarceId(res.Header))
		return nil, res, xbase.ComErrServRespErrnoErr
	}

	t.Logger.Trace("operate succ. [open_id: %s] [app_key: %s] [url: %s] [request_id: %s] [trace_id: %s]",
		param.OpenId, param.AppKey, res.ReqUrl, resp.RequestId, t.GetTarceId(res.Header))
	return &resp, res, nil
}

func (t *AssetOper) BindByUnionId(param *xbase.BindByUnionIdParam) (*xbase.BaseResp, *xbase.RequestRes, error) {
	if err := param.Valid(); err != nil {
		return nil, nil, err
	}

	signedUnionId, err := t.aesEncodeStr(param.UnionId)
	if err != nil {
		t.Logger.Warn("encode union id fail, union id: %s", param.UnionId)
		return nil, nil, err
	}
	signedMnem, err := t.aesEncodeStr(param.Mnemonic)
	if err != nil {
		t.Logger.Warn("encode mnemonic fail, mnemonic: %s", param.Mnemonic)
		return nil, nil, err
	}
	v := url.Values{}
	v.Set("union_id", signedUnionId)
	v.Set("mnemonic", signedMnem)
	body := v.Encode()

	res, err := t.Post(xbase.DidApiBindByUid, body)
	if err != nil {
		t.Logger.Warn("post request xasset failed. url: %s, err: %v", xbase.DidApiBindByUid, err)
		return nil, nil, xbase.ComErrRequsetFailed
	}
	if res.HttpCode != 200 {
		t.Logger.Warn("post req resp not 200.[http_code: %d] [url: %s] [body: %s] [trace_id: %s]",
			res.HttpCode, res.ReqUrl, res.Body, t.GetTarceId(res.Header))
		return nil, nil, xbase.ComErrRespCodeErr
	}

	var resp xbase.BaseResp
	err = json.Unmarshal([]byte(res.Body), &resp)
	if err != nil {
		t.Logger.Warn("unmarshal body failed.err:%v [http_code: %d] [url: %s] [body: %s] [trace_id: %s]",
			err, res.HttpCode, res.ReqUrl, res.Body, t.GetTarceId(res.Header))
		return nil, res, xbase.ComErrUnmarshalBodyFailed
	}
	if resp.Errno != xbase.XassetErrNoSucc {
		t.Logger.Warn("get resp failed. [url: %s] [request_id: %s] [err_no: %d] [trace_id: %s]",
			res.ReqUrl, resp.RequestId, resp.Errno, t.GetTarceId(res.Header))
		return nil, res, xbase.ComErrServRespErrnoErr
	}

	t.Logger.Trace("operate succ. [union_id: %s] [url: %s] [request_id: %s] [trace_id: %s]",
		param.UnionId, res.ReqUrl, resp.RequestId, t.GetTarceId(res.Header))
	return &resp, res, nil
}

func (t *AssetOper) GetAddrByUnionId(uid string) (*xbase.GetAddrByUnionIdResp, *xbase.RequestRes, error) {
	if err := xbase.UnionIdValid(uid); err != nil {
		return nil, nil, err
	}
	signedUnionId, err := t.aesEncodeStr(uid)
	if err != nil {
		t.Logger.Warn("encode union id fail, union id: %s", uid)
		return nil, nil, err
	}
	v := url.Values{}
	v.Set("union_id", signedUnionId)
	body := v.Encode()

	res, err := t.Post(xbase.DidApiGetAddrByUid, body)
	if err != nil {
		t.Logger.Warn("post request xasset failed. url: %s, err: %v", xbase.DidApiGetAddrByUid, err)
		return nil, nil, xbase.ComErrRequsetFailed
	}
	if res.HttpCode != 200 {
		t.Logger.Warn("post req resp not 200.[http_code: %d] [url: %s] [body: %s] [trace_id: %s]",
			res.HttpCode, res.ReqUrl, res.Body, t.GetTarceId(res.Header))
		return nil, nil, xbase.ComErrRespCodeErr
	}

	var resp xbase.GetAddrByUnionIdResp
	err = json.Unmarshal([]byte(res.Body), &resp)
	if err != nil {
		t.Logger.Warn("unmarshal body failed.err:%v [http_code: %d] [url: %s] [body: %s] [trace_id: %s]",
			err, res.HttpCode, res.ReqUrl, res.Body, t.GetTarceId(res.Header))
		return nil, res, xbase.ComErrUnmarshalBodyFailed
	}
	if resp.Errno != xbase.XassetErrNoSucc {
		t.Logger.Warn("get resp failed. [url: %s] [request_id: %s] [err_no: %d] [trace_id: %s]",
			res.ReqUrl, resp.RequestId, resp.Errno, t.GetTarceId(res.Header))
		return nil, res, xbase.ComErrServRespErrnoErr
	}

	t.Logger.Trace("operate succ. [union_id: %s] [url: %s] [request_id: %s] [trace_id: %s]",
		uid, res.ReqUrl, resp.RequestId, t.GetTarceId(res.Header))
	return &resp, res, nil
}

func (t *AssetOper) aesEncodeStr(str string) (string, error) {
	return utils.AesEncode(str, t.Cfg.Credentials.SecretAccessKey)
}

func (t *AssetOper) aesDecodeStr(str string) (string, error) {
	return utils.AesDecode(str, t.Cfg.Credentials.SecretAccessKey)
}

func (t *AssetOper) genCreateOrderBody(param *xbase.HubCreateOrderParam) (string, error) {
	v := url.Values{}
	v.Set("code", fmt.Sprintf("%d", param.Code))
	v.Set("order_type", fmt.Sprintf("%d", param.OrderType))
	v.Set("executor", param.ExecutorAPI)
	v.Set("executor_data", param.ExecutorData)
	v.Set("timestamp", fmt.Sprintf("%d", param.Timestamp))
	v.Set("time_expire", fmt.Sprintf("%d", param.TimeExpire))
	v.Set("profit_sharing", fmt.Sprintf("%d", param.ProfitSharing))
	v.Set("uid", fmt.Sprintf("%d", param.Uid))
	v.Set("creator_details", param.Details)
	v.Set("app_id", fmt.Sprintf("%d", param.AppId))
	v.Set("act_id", fmt.Sprintf("%d", param.ActId))
	v.Set("asset_id", fmt.Sprintf("%d", param.AssetId))
	v.Set("buyer_addr", param.BuyerAddr)
	v.Set("seller_addr", param.SellerAddr)
	v.Set("client_type", fmt.Sprintf("%d", param.ClientType))
	v.Set("chan", fmt.Sprintf("%d", param.Chan))
	v.Set("scene", fmt.Sprintf("%d", param.Scene))
	v.Set("signed_auth", param.SignedAuth)

	return v.Encode(), nil
}

// CreateOrder creates orders.
func (t *AssetOper) CreateOrder(param *xbase.HubCreateOrderParam) (*xbase.HubCreateResp, *xbase.RequestRes, error) {
	if err := param.Valid(); err != nil {
		return nil, nil, err
	}
	body, err := t.genCreateOrderBody(param)
	if err != nil {
		t.Logger.Warn("fail to generate value for create order, err: %v, param: %+v", err, *param)
		return nil, nil, err
	}
	res, err := t.Post(xbase.HubCreateOrder, body)
	if err != nil {
		t.Logger.Warn("post request xasset failed, uri: %s, err: %v", xbase.HubCreateOrder, err)
		return nil, nil, xbase.ComErrRequsetFailed
	}
	if res.HttpCode != 200 {
		t.Logger.Warn("post request response is not 200. [http_code: %d] [url: %s] [body: %s] [trace_id: %s]",
			res.HttpCode, res.ReqUrl, res.Body, t.GetTarceId(res.Header))
		return nil, nil, xbase.ComErrRespCodeErr
	}

	var resp xbase.HubCreateResp
	err = json.Unmarshal([]byte(res.Body), &resp)
	if err != nil {
		t.Logger.Warn("unmarshal body failed. [http_code: %d] [url: %s] [body: %s] [trace_id: %s]",
			res.HttpCode, res.ReqUrl, res.Body, t.GetTarceId(res.Header))
		return nil, res, xbase.ComErrUnmarshalBodyFailed
	}
	if resp.Errno != xbase.XassetErrNoSucc {
		t.Logger.Warn("get resp failed. [url: %s] [request_id: %s] [err_no: %d] [trace_id: %s]",
			res.ReqUrl, resp.RequestId, resp.Errno, t.GetTarceId(res.Header))
		return nil, res, xbase.ComErrServRespErrnoErr
	}

	t.Logger.Trace("operate succ. [param: %+v] [url: %s] [request_id: %s] [trace_id: %s]",
		param, res.ReqUrl, resp.RequestId, t.GetTarceId(res.Header))
	return &resp, res, nil
}

// ConfirmOrder confirms orders.
func (t *AssetOper) ConfirmOrder(param *xbase.HubConfirmH5OrderParam) (*xbase.HubCreateResp, *xbase.RequestRes, error) {
	if err := param.Valid(); err != nil {
		return nil, nil, err
	}
	v := url.Values{}
	v.Set("code", fmt.Sprintf("%d", param.Code))
	v.Set("order_type", fmt.Sprintf("%d", param.OrderType))
	v.Set("app_id", fmt.Sprintf("%d", param.AppId))
	v.Set("oid", fmt.Sprintf("%d", param.Oid))
	v.Set("client_type", fmt.Sprintf("%d", param.ClientType))
	v.Set("signed_auth", param.SignAuth)
	body := v.Encode()
	res, err := t.Post(xbase.HubConfirmOrder, body)
	if err != nil {
		t.Logger.Warn("post request xasset failed, uri: %s, err: %v", xbase.HubConfirmOrder, err)
		return nil, nil, xbase.ComErrRequsetFailed
	}
	if res.HttpCode != 200 {
		t.Logger.Warn("post request response is not 200. [http_code: %d] [url: %s] [body: %s] [trace_id: %s]",
			res.HttpCode, res.ReqUrl, res.Body, t.GetTarceId(res.Header))
		return nil, nil, xbase.ComErrRespCodeErr
	}

	var resp xbase.HubCreateResp
	err = json.Unmarshal([]byte(res.Body), &resp)
	if err != nil {
		t.Logger.Warn("unmarshal body failed. [http_code: %d] [url: %s] [body: %s] [trace_id: %s]",
			res.HttpCode, res.ReqUrl, res.Body, t.GetTarceId(res.Header))
		return nil, res, xbase.ComErrUnmarshalBodyFailed
	}
	if resp.Errno != xbase.XassetErrNoSucc {
		t.Logger.Warn("get resp failed. [url: %s] [request_id: %s] [err_no: %d] [trace_id: %s]",
			res.ReqUrl, resp.RequestId, resp.Errno, t.GetTarceId(res.Header))
		return nil, res, xbase.ComErrServRespErrnoErr
	}

	t.Logger.Trace("operate succ. [param: %+v] [url: %s] [request_id: %s] [trace_id: %s]",
		param, res.ReqUrl, resp.RequestId, t.GetTarceId(res.Header))
	return &resp, res, nil
}

// QueryOrderDetail gets order info.
func (t *AssetOper) QueryOrderDetail(param *xbase.HubOrderDetailParam) (*xbase.HubOrderDetailResp, *xbase.RequestRes, error) {
	if err := param.Valid(); err != nil {
		return nil, nil, err
	}
	v := url.Values{}
	v.Set("app_id", fmt.Sprintf("%d", param.AppId))
	v.Set("oid", fmt.Sprintf("%d", param.Oid))
	body := v.Encode()
	res, err := t.Post(xbase.HubDetailOrder, body)
	if err != nil {
		t.Logger.Warn("post request xasset failed, uri: %s, err: %v", xbase.HubDetailOrder, err)
		return nil, nil, xbase.ComErrRequsetFailed
	}
	if res.HttpCode != 200 {
		t.Logger.Warn("post request response is not 200. [http_code: %d] [url: %s] [body: %s] [trace_id: %s]",
			res.HttpCode, res.ReqUrl, res.Body, t.GetTarceId(res.Header))
		return nil, nil, xbase.ComErrRespCodeErr
	}

	var resp xbase.HubOrderDetailResp
	err = json.Unmarshal([]byte(res.Body), &resp)
	if err != nil {
		t.Logger.Warn("unmarshal body failed. [http_code: %d] [url: %s] [body: %s] [trace_id: %s]",
			res.HttpCode, res.ReqUrl, res.Body, t.GetTarceId(res.Header))
		return nil, res, xbase.ComErrUnmarshalBodyFailed
	}
	if resp.Errno != xbase.XassetErrNoSucc {
		t.Logger.Warn("get resp failed. [url: %s] [request_id: %s] [err_no: %d] [trace_id: %s]",
			res.ReqUrl, resp.RequestId, resp.Errno, t.GetTarceId(res.Header))
		return nil, res, xbase.ComErrServRespErrnoErr
	}

	t.Logger.Trace("operate succ. [param: %+v] [url: %s] [request_id: %s] [trace_id: %s]",
		param, res.ReqUrl, resp.RequestId, t.GetTarceId(res.Header))
	return &resp, res, nil
}

func (t *AssetOper) genEditOrderBody(param *xbase.HubEditOrderParam) (string, error) {
	v := url.Values{}
	v.Set("app_id", fmt.Sprintf("%d", param.AppId))
	v.Set("oid", fmt.Sprintf("%d", param.Oid))
	v.Set("status", fmt.Sprintf("%d", param.Status))
	v.Set("pay_channel", fmt.Sprintf("%d", param.PayChannel))
	v.Set("third_oid", param.ThirdOid)
	v.Set("pay_info", param.PayInfo)
	v.Set("pay_time", fmt.Sprintf("%d", param.PayTime))
	v.Set("close_time", fmt.Sprintf("%d", param.CloseTime))
	v.Set("close_reason", param.CloseReason)

	return v.Encode(), nil
}

// EditOrder edits order info.
func (t *AssetOper) EditOrder(param *xbase.HubEditOrderParam) (*xbase.BaseResp, *xbase.RequestRes, error) {
	if err := param.Valid(); err != nil {
		return nil, nil, err
	}
	body, err := t.genEditOrderBody(param)
	if err != nil {
		t.Logger.Warn("fail to generate value for edit order, err: %v, param: %+v", err, *param)
		return nil, nil, err
	}
	res, err := t.Post(xbase.HubEditOrder, body)
	if err != nil {
		t.Logger.Warn("post request xasset failed, uri: %s, err: %v", xbase.HubEditOrder, err)
		return nil, nil, xbase.ComErrRequsetFailed
	}
	if res.HttpCode != 200 {
		t.Logger.Warn("post request response is not 200. [http_code: %d] [url: %s] [body: %s] [trace_id: %s]",
			res.HttpCode, res.ReqUrl, res.Body, t.GetTarceId(res.Header))
		return nil, nil, xbase.ComErrRespCodeErr
	}

	var resp xbase.BaseResp
	err = json.Unmarshal([]byte(res.Body), &resp)
	if err != nil {
		t.Logger.Warn("unmarshal body failed. [http_code: %d] [url: %s] [body: %s] [trace_id: %s]",
			res.HttpCode, res.ReqUrl, res.Body, t.GetTarceId(res.Header))
		return nil, res, xbase.ComErrUnmarshalBodyFailed
	}
	if resp.Errno != xbase.XassetErrNoSucc {
		t.Logger.Warn("get resp failed. [url: %s] [request_id: %s] [err_no: %d] [trace_id: %s]",
			res.ReqUrl, resp.RequestId, resp.Errno, t.GetTarceId(res.Header))
		return nil, res, xbase.ComErrServRespErrnoErr
	}

	t.Logger.Trace("operate succ. [param: %+v] [url: %s] [request_id: %s] [trace_id: %s]",
		param, res.ReqUrl, resp.RequestId, t.GetTarceId(res.Header))
	return &resp, res, nil
}

// QueryOrderList gets order list by address.
func (t *AssetOper) QueryOrderList(param *xbase.HubListOrderParam) (*xbase.HubListOrderResp, *xbase.RequestRes, error) {
	if err := param.Valid(); err != nil {
		return nil, nil, err
	}
	v := url.Values{}
	v.Set("app_id", fmt.Sprintf("%d", param.AppId))
	v.Set("address", param.Addr)
	v.Set("status", fmt.Sprintf("%d", param.Status))
	v.Set("cursor", param.Cursor)
	v.Set("limit", fmt.Sprintf("%d", param.Limit))

	body := v.Encode()
	res, err := t.Post(xbase.HubListOrder, body)
	if err != nil {
		t.Logger.Warn("post request xasset failed, uri: %s, err: %v", xbase.HubListOrder, err)
		return nil, nil, xbase.ComErrRequsetFailed
	}
	if res.HttpCode != 200 {
		t.Logger.Warn("post request response is not 200. [http_code: %d] [url: %s] [body: %s] [trace_id: %s]",
			res.HttpCode, res.ReqUrl, res.Body, t.GetTarceId(res.Header))
		return nil, nil, xbase.ComErrRespCodeErr
	}

	var resp xbase.HubListOrderResp
	err = json.Unmarshal([]byte(res.Body), &resp)
	if err != nil {
		t.Logger.Warn("unmarshal body failed. [http_code: %d] [url: %s] [body: %s] [trace_id: %s]",
			res.HttpCode, res.ReqUrl, res.Body, t.GetTarceId(res.Header))
		return nil, res, xbase.ComErrUnmarshalBodyFailed
	}
	if resp.Errno != xbase.XassetErrNoSucc {
		t.Logger.Warn("get resp failed. [url: %s] [request_id: %s] [err_no: %d] [trace_id: %s]",
			res.ReqUrl, resp.RequestId, resp.Errno, t.GetTarceId(res.Header))
		return nil, res, xbase.ComErrServRespErrnoErr
	}

	t.Logger.Trace("operate succ. [param: %+v] [url: %s] [request_id: %s] [trace_id: %s]",
		param, res.ReqUrl, resp.RequestId, t.GetTarceId(res.Header))
	return &resp, res, nil
}
