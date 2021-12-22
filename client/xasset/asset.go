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

// genGetStokenBody Grant uses the general parameter for getting stoken as follows,
//    {
// 		   Addr     string `json:"addr"`
// 		   Sign     string `json:"sign"`
// 		   PKey     string `json:"pkey"`
// 		   Nonce    int64  `json:"nonce"`
// 	  }
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
// {
// 		AssetId   int64  `json:"asset_id"`
// 		Amount    int    `json:"amount"`
// 		AssetInfo string `json:"asset_info"`
// 		Addr      string `json:"addr"`
// 		Sign      string `json:"sign"`
// 		PKey      string `json:"pkey"`
// 		Nonce     int64  `json:"nonce"`
// 		UserId    int64  `json:"user_id,omitempty"`
// }
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
	v.Set("amount", fmt.Sprintf("%d", param.Amount))
	v.Set("asset_info", string(assetInfo))
	v.Set("addr", param.Account.Address)
	v.Set("sign", sign)
	v.Set("pkey", param.Account.PublicKey)
	v.Set("nonce", fmt.Sprintf("%d", nonce))
	if err := xbase.IdValid(param.UserId); err == nil {
		v.Set("user_id", fmt.Sprintf("%d", param.UserId))
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
// {
// 		AssetId   int64  `json:"asset_id"`
// 		Addr      string `json:"addr"`
// 		Sign      string `json:"sign"`
// 		PKey      string `json:"pkey"`
// 		Nonce     int64  `json:"nonce"`
// 		Amount    int    `json:"amount"`
// 		AssetInfo string `json:"asset_info"`
// }
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
	body := v.Encode()
	return body, nil
}

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
// {
// 		AssetId    int64  `json:"asset_id"`
// 		Addr       string `json:"addr"`
// 		Sign       string `json:"sign"`
// 		PKey       string `json:"pkey"`
// 		Nonce      int64  `json:"nonce"`
//	    IsEvidence int    `json:"is_evidence,omitempty"`
// }
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

// GenAlterAssetBody uses the parameter as follows,
// {
//		AssetId int64 `json:"asset_id"`
// }
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

// GenGrantAssetBody Grant uses the general parameter for granting as follows,
//    {
// 		   AssetId  int64  `json:"asset_id"`
// 		   ShardId  int64  `json:"shard_id"`
// 		   Addr     string `json:"addr"`
// 		   Sign     string `json:"sign"`
// 		   PKey     string `json:"pkey"`
// 		   Nonce    int64  `json:"nonce"`
// 		   ToAddr   string `json:"to_addr"`
// 		   ToUserId int64  `json:"to_userid,omitempty"`
// 	  }
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

// GenQueryShardsBody uses the general parameter for getting shard as follows,
//    {
// 		   AssetId  int64  `json:"asset_id"`
// 		   ShardId  int64  `json:"shard_id"`
// 	  }
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
//    {
//		   Addr  string `json:"addr"`
//		   Page  int    `json:"page"`
//		   Limit int    `json:"limit"`
// 	  }
func (t *AssetOper) genListShardsByAddrBody(param *xbase.ListShardsByAddrParam) (string, error) {
	v := url.Values{}
	v.Set("addr", param.Addr)
	v.Set("page", fmt.Sprintf("%d", param.Page))
	v.Set("limit", fmt.Sprintf("%d", param.Limit))
	body := v.Encode()
	return body, nil
}

func (t *AssetOper) ListShardsByAddr(param *xbase.ListShardsByAddrParam) (*xbase.ListPageResp, *xbase.RequestRes, error) {
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

	var resp xbase.ListPageResp
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

// GenEvidenceBody uses the general parameter for getting shard as follows,
//    {
// 		   AssetId  int64  `json:"asset_id"`
// 		   ShardId  int64  `json:"shard_id"`
// 	  }
func (t *AssetOper) genEvidenceBody(param *xbase.GetEvidenceInfoParam) (string, error) {
	v := url.Values{}
	v.Set("asset_id", fmt.Sprintf("%d", param.AssetId))
	v.Set("shard_id", fmt.Sprintf("%d", param.ShardId))
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

// GenTransferAssetBody uses the general parameter for transferring as follows,
//    {
// 		   AssetId  int64  `json:"asset_id"`
// 		   ShardId  int64  `json:"shard_id"`
// 		   Addr     string `json:"addr"`
// 		   Sign     string `json:"sign"`
// 		   PKey     string `json:"pkey"`
// 		   Nonce    int64  `json:"nonce"`
// 		   ToAddr   string `json:"to_addr"`
// 		   ToUserId int64  `json:"to_userid,omitempty"`
// 	  }
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
	v.Set("addr", param.Addr)
	v.Set("sign", sign)
	v.Set("pkey", param.Account.PublicKey)
	v.Set("nonce", fmt.Sprintf("%d", nonce))
	v.Set("to_addr", param.ToAddr)
	if err := xbase.IdValid(param.ToUserId); err == nil {
		v.Set("to_userid", param.Addr)
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
		t.Logger.Warn("post request xasset failed, uri: %s, err: %v", xbase.AssetApiGrant, err)
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
