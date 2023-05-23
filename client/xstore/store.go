package xstore

import (
	"crypto/md5"
	"encoding/json"
	"fmt"
	"io"
	"net/url"

	xbase "github.com/xuperchain/xasset-sdk-go/client/base"
	"github.com/xuperchain/xasset-sdk-go/common/config"
	"github.com/xuperchain/xasset-sdk-go/common/logs"
	"github.com/xuperchain/xasset-sdk-go/utils"
)

type StoreOper struct {
	xbase.XassetBaseClient
}

func NewXstoreOper(cfg *config.XassetCliConfig, logger logs.LogDriver) (*StoreOper, error) {
	obj := &StoreOper{}
	err := obj.InitClient(cfg, logger)
	if err != nil {
		return nil, err
	}

	return obj, nil
}

func (t *StoreOper) genCreateOrAlterStoreBody(param *xbase.CreateOrAlterStoreParam) (string, error) {
	v := url.Values{}
	v.Set("store_id", fmt.Sprintf("%d", param.StoreId))
	v.Set("name", param.Name)
	v.Set("logo", param.Logo)
	v.Set("cover", param.Cover)
	v.Set("short_desc", param.ShortDesc)
	v.Set("weight", fmt.Sprintf("%d", param.Weight))
	v.Set("ext_info", param.ExtInfo)
	v.Set("wechat", param.Wechat)
	body := v.Encode()
	return body, nil
}

func (t *StoreOper) CreateStore(param *xbase.CreateOrAlterStoreParam) (*xbase.BaseResp, *xbase.RequestRes, error) {
	if err := param.CreateValid(); err != nil {
		return nil, nil, err
	}

	body, err := t.genCreateOrAlterStoreBody(param)
	if err != nil {
		t.Logger.Warn("fail to generate value for create store, err: %v, param: %+v", err, *param)
		return nil, nil, err
	}
	res, err := t.Post(xbase.StoreApiCreate, body)
	if err != nil {
		t.Logger.Warn("post request xasset failed.err:%v", err)
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

	t.Logger.Trace("operate succ. [store_id: %v] [url: %s] [request_id: %s] [trace_id: %s]",
		param.StoreId, res.ReqUrl, resp.RequestId, t.GetTarceId(res.Header))
	return &resp, res, nil
}

func (t *StoreOper) AlterStore(param *xbase.CreateOrAlterStoreParam) (*xbase.BaseResp, *xbase.RequestRes, error) {
	if param.StoreId < 1 {
		return nil, nil, xbase.ErrParamInvalid
	}

	body, err := t.genCreateOrAlterStoreBody(param)
	if err != nil {
		t.Logger.Warn("fail to generate value for alter store, err: %v, param: %+v", err, *param)
		return nil, nil, err
	}
	res, err := t.Post(xbase.StoreApiAlter, body)
	if err != nil {
		t.Logger.Warn("post request xasset failed.err:%v", err)
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

	t.Logger.Trace("operate succ. [store_id: %v] [url: %s] [request_id: %s] [trace_id: %s]",
		param.StoreId, res.ReqUrl, resp.RequestId, t.GetTarceId(res.Header))
	return &resp, res, nil
}

func (t *StoreOper) genQueryStoreBody(param *xbase.BaseStoreParam) (string, error) {
	v := url.Values{}
	v.Set("store_id", fmt.Sprintf("%d", param.StoreId))
	body := v.Encode()
	return body, nil
}

func (t *StoreOper) QueryStore(param *xbase.BaseStoreParam) (*xbase.QueryStoreResp, *xbase.RequestRes, error) {
	if param.StoreId < 1 {
		return nil, nil, xbase.ErrParamInvalid
	}

	body, err := t.genQueryStoreBody(param)
	if err != nil {
		t.Logger.Warn("fail to generate value for query store, err: %v, param: %+v", err, *param)
		return nil, nil, err
	}
	res, err := t.Post(xbase.StoreApiQuery, body)
	if err != nil {
		t.Logger.Warn("post request xasset failed.err:%v", err)
		return nil, nil, xbase.ComErrRequsetFailed
	}
	if res.HttpCode != 200 {
		t.Logger.Warn("post request response is not 200. [http_code: %d] [url: %s] [body: %s] [trace_id: %s]",
			res.HttpCode, res.ReqUrl, res.Body, t.GetTarceId(res.Header))
		return nil, nil, xbase.ComErrRespCodeErr
	}

	var resp xbase.QueryStoreResp
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

	t.Logger.Trace("operate succ. [store_id: %v] [url: %s] [request_id: %s] [trace_id: %s]",
		param.StoreId, res.ReqUrl, resp.RequestId, t.GetTarceId(res.Header))
	return &resp, res, nil
}

func (t *StoreOper) ListStore() (*xbase.ListStoreResp, *xbase.RequestRes, error) {
	res, err := t.Post(xbase.StoreApiList, "")
	if err != nil {
		t.Logger.Warn("post request xasset failed.err:%v", err)
		return nil, nil, xbase.ComErrRequsetFailed
	}
	if res.HttpCode != 200 {
		t.Logger.Warn("post request response is not 200. [http_code: %d] [url: %s] [body: %s] [trace_id: %s]",
			res.HttpCode, res.ReqUrl, res.Body, t.GetTarceId(res.Header))
		return nil, nil, xbase.ComErrRespCodeErr
	}

	var resp xbase.ListStoreResp
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

	t.Logger.Trace("operate succ. [url: %s] [request_id: %s] [trace_id: %s]",
		res.ReqUrl, resp.RequestId, t.GetTarceId(res.Header))
	return &resp, res, nil
}

func (t *StoreOper) genCreateOrAlterActBody(param *xbase.CreateOrAlterActParam) (string, error) {
	v := url.Values{}
	v.Set("store_id", fmt.Sprintf("%d", param.StoreId))
	v.Set("act_id", fmt.Sprintf("%d", param.ActId))
	v.Set("jump_link", param.JumpLink)
	v.Set("issuer", param.Issuer)
	v.Set("act_name", param.ActName)
	v.Set("short_desc", param.ShortDesc)
	v.Set("thumb", param.Thumb)
	v.Set("img_desc", param.ImgDesc)
	v.Set("bookable", fmt.Sprintf("%d", param.Bookable))
	v.Set("start", fmt.Sprintf("%d", param.Start))
	v.Set("end", fmt.Sprintf("%d", param.End))
	v.Set("weight", fmt.Sprintf("%d", param.Weight))
	v.Set("ext_info", param.ExtInfo)
	body := v.Encode()
	return body, nil
}

func (t *StoreOper) CreateAct(param *xbase.CreateOrAlterActParam) (*xbase.BaseResp, *xbase.RequestRes, error) {
	if err := param.CreateValid(); err != nil {
		return nil, nil, err
	}

	body, err := t.genCreateOrAlterActBody(param)
	if err != nil {
		t.Logger.Warn("fail to generate value for create act, err: %v, param: %+v", err, *param)
		return nil, nil, err
	}
	res, err := t.Post(xbase.StoreApiCreateAct, body)
	if err != nil {
		t.Logger.Warn("post request xasset failed.err:%v", err)
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

	t.Logger.Trace("operate succ. [store_id: %v] [url: %s] [request_id: %s] [trace_id: %s]",
		param.StoreId, res.ReqUrl, resp.RequestId, t.GetTarceId(res.Header))
	return &resp, res, nil
}

func (t *StoreOper) AlterAct(param *xbase.CreateOrAlterActParam) (*xbase.BaseResp, *xbase.RequestRes, error) {
	if param.ActId < 1 {
		return nil, nil, xbase.ErrParamInvalid
	}

	body, err := t.genCreateOrAlterActBody(param)
	if err != nil {
		t.Logger.Warn("fail to generate value for alter act, err: %v, param: %+v", err, *param)
		return nil, nil, err
	}
	res, err := t.Post(xbase.StoreApiAlterAct, body)
	if err != nil {
		t.Logger.Warn("post request xasset failed.err:%v", err)
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

	t.Logger.Trace("operate succ. [store_id: %v] [url: %s] [request_id: %s] [trace_id: %s]",
		param.StoreId, res.ReqUrl, resp.RequestId, t.GetTarceId(res.Header))
	return &resp, res, nil
}

func (t *StoreOper) genBaseActBody(param *xbase.BaseActParam) (string, error) {
	v := url.Values{}
	v.Set("act_id", fmt.Sprintf("%d", param.ActId))
	v.Set("op_type", fmt.Sprintf("%d", param.OpType))
	body := v.Encode()
	return body, nil
}

func (t *StoreOper) RemoveAct(param *xbase.BaseActParam) (*xbase.BaseResp, *xbase.RequestRes, error) {
	if err := param.Valid(); err != nil {
		return nil, nil, xbase.ErrParamInvalid
	}

	body, err := t.genBaseActBody(param)
	if err != nil {
		t.Logger.Warn("fail to generate value for remove act, err: %v, param: %+v", err, *param)
		return nil, nil, err
	}
	res, err := t.Post(xbase.StoreApiRemoveAct, body)
	if err != nil {
		t.Logger.Warn("post request xasset failed.err:%v", err)
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

	t.Logger.Trace("operate succ. [act_id: %v] [url: %s] [request_id: %s] [trace_id: %s]",
		param.ActId, res.ReqUrl, resp.RequestId, t.GetTarceId(res.Header))
	return &resp, res, nil
}

func (t *StoreOper) QueryAct(param *xbase.BaseActParam) (*xbase.QueryActResp, *xbase.RequestRes, error) {
	if err := param.Valid(); err != nil {
		return nil, nil, xbase.ErrParamInvalid
	}

	body, err := t.genBaseActBody(param)
	if err != nil {
		t.Logger.Warn("fail to generate value for query act, err: %v, param: %+v", err, *param)
		return nil, nil, err
	}
	res, err := t.Post(xbase.StoreApiQueryAct, body)
	if err != nil {
		t.Logger.Warn("post request xasset failed.err:%v", err)
		return nil, nil, xbase.ComErrRequsetFailed
	}
	if res.HttpCode != 200 {
		t.Logger.Warn("post request response is not 200. [http_code: %d] [url: %s] [body: %s] [trace_id: %s]",
			res.HttpCode, res.ReqUrl, res.Body, t.GetTarceId(res.Header))
		return nil, nil, xbase.ComErrRespCodeErr
	}

	var resp xbase.QueryActResp
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

	t.Logger.Trace("operate succ. [act_id: %v] [url: %s] [request_id: %s] [trace_id: %s]",
		param.ActId, res.ReqUrl, resp.RequestId, t.GetTarceId(res.Header))
	return &resp, res, nil
}

func (t *StoreOper) genListActBody(param *xbase.ListActParam) (string, error) {
	v := url.Values{}
	v.Set("store_id", fmt.Sprintf("%d", param.StoreId))
	v.Set("limit", fmt.Sprintf("%d", param.Limit))
	v.Set("cursor", param.Cursor)
	body := v.Encode()
	return body, nil
}

func (t *StoreOper) ListAct(param *xbase.ListActParam) (*xbase.ListActResp, *xbase.RequestRes, error) {
	if err := param.Valid(); err != nil {
		return nil, nil, xbase.ErrParamInvalid
	}

	body, err := t.genListActBody(param)
	if err != nil {
		t.Logger.Warn("fail to generate value for list act, err: %v, param: %+v", err, *param)
		return nil, nil, err
	}
	res, err := t.Post(xbase.StoreApiListAct, body)
	if err != nil {
		t.Logger.Warn("post request xasset failed.err:%v", err)
		return nil, nil, xbase.ComErrRequsetFailed
	}
	if res.HttpCode != 200 {
		t.Logger.Warn("post request response is not 200. [http_code: %d] [url: %s] [body: %s] [trace_id: %s]",
			res.HttpCode, res.ReqUrl, res.Body, t.GetTarceId(res.Header))
		return nil, nil, xbase.ComErrRespCodeErr
	}

	var resp xbase.ListActResp
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

	t.Logger.Trace("operate succ. [store_id: %v] [url: %s] [request_id: %s] [trace_id: %s]",
		param.StoreId, res.ReqUrl, resp.RequestId, t.GetTarceId(res.Header))
	return &resp, res, nil
}

func (t *StoreOper) PubAct(param *xbase.BaseActParam) (*xbase.BaseResp, *xbase.RequestRes, error) {
	if err := param.Valid(); err != nil {
		return nil, nil, xbase.ErrParamInvalid
	}

	body, err := t.genBaseActBody(param)
	if err != nil {
		t.Logger.Warn("fail to generate value for pub act, err: %v, param: %+v", err, *param)
		return nil, nil, err
	}
	res, err := t.Post(xbase.StoreApiPubAct, body)
	if err != nil {
		t.Logger.Warn("post request xasset failed.err:%v", err)
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

	t.Logger.Trace("operate succ. [act_id: %v] [url: %s] [request_id: %s] [trace_id: %s]",
		param.ActId, res.ReqUrl, resp.RequestId, t.GetTarceId(res.Header))
	return &resp, res, nil
}

func (t *StoreOper) genBindOrAlterAstBody(param *xbase.BindOrAlterAstParam) (string, error) {
	v := url.Values{}
	v.Set("act_id", fmt.Sprintf("%d", param.ActId))
	v.Set("asset_id", fmt.Sprintf("%d", param.AssetId))
	v.Set("split_id", fmt.Sprintf("%d", param.SplitId))
	v.Set("asset_cate", fmt.Sprintf("%d", param.AssetCate))
	v.Set("start", fmt.Sprintf("%d", param.Start))
	v.Set("end", fmt.Sprintf("%d", param.End))
	v.Set("grant_mode", fmt.Sprintf("%d", param.GrantMode))
	v.Set("jump_link", param.JumpLink)
	v.Set("ext_info", param.ExtInfo)
	v.Set("price", fmt.Sprintf("%d", param.Price))
	v.Set("ori_price", fmt.Sprintf("%d", param.OriPrice))
	v.Set("amount", fmt.Sprintf("%d", param.Amount))
	v.Set("apply_form", fmt.Sprintf("%d", param.ApplyForm))
	v.Set("is_box", fmt.Sprintf("%d", param.IsBox))
	body := v.Encode()
	return body, nil
}

func (t *StoreOper) BindAst(param *xbase.BindOrAlterAstParam) (*xbase.BaseResp, *xbase.RequestRes, error) {
	if err := param.CreateValid(); err != nil {
		return nil, nil, err
	}

	body, err := t.genBindOrAlterAstBody(param)
	if err != nil {
		t.Logger.Warn("fail to generate value for bind ast, err: %v, param: %+v", err, *param)
		return nil, nil, err
	}
	res, err := t.Post(xbase.StoreApiBindAst, body)
	if err != nil {
		t.Logger.Warn("post request xasset failed.err:%v", err)
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

	t.Logger.Trace("operate succ. [asset_id: %v] [url: %s] [request_id: %s] [trace_id: %s]",
		param.AssetId, res.ReqUrl, resp.RequestId, t.GetTarceId(res.Header))
	return &resp, res, nil
}

func (t *StoreOper) AlterAst(param *xbase.BindOrAlterAstParam) (*xbase.BaseResp, *xbase.RequestRes, error) {
	if err := param.AlterValid(); err != nil {
		return nil, nil, err
	}

	body, err := t.genBindOrAlterAstBody(param)
	if err != nil {
		t.Logger.Warn("fail to generate value for alter ast, err: %v, param: %+v", err, *param)
		return nil, nil, err
	}
	res, err := t.Post(xbase.StoreApiAlterAst, body)
	if err != nil {
		t.Logger.Warn("post request xasset failed.err:%v", err)
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

	t.Logger.Trace("operate succ. [asset_id: %v] [url: %s] [request_id: %s] [trace_id: %s]",
		param.AssetId, res.ReqUrl, resp.RequestId, t.GetTarceId(res.Header))
	return &resp, res, nil
}

func (t *StoreOper) genCancelAstBody(param *xbase.BaseAstParam) (string, error) {
	v := url.Values{}
	v.Set("act_id", fmt.Sprintf("%d", param.ActId))
	v.Set("asset_id", fmt.Sprintf("%d", param.AssetId))
	body := v.Encode()
	return body, nil
}

func (t *StoreOper) CancelAst(param *xbase.BaseAstParam) (*xbase.BaseResp, *xbase.RequestRes, error) {
	if err := param.Valid(); err != nil {
		return nil, nil, err
	}

	body, err := t.genCancelAstBody(param)
	if err != nil {
		t.Logger.Warn("fail to generate value for cancel ast, err: %v, param: %+v", err, *param)
		return nil, nil, err
	}
	res, err := t.Post(xbase.StoreApiCancelAst, body)
	if err != nil {
		t.Logger.Warn("post request xasset failed.err:%v", err)
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

	t.Logger.Trace("operate succ. [asset_id: %v] [url: %s] [request_id: %s] [trace_id: %s]",
		param.AssetId, res.ReqUrl, resp.RequestId, t.GetTarceId(res.Header))
	return &resp, res, nil
}

func (t *StoreOper) genCancelAstByActIdBody(param *xbase.BaseActParam) (string, error) {
	v := url.Values{}
	v.Set("act_id", fmt.Sprintf("%d", param.ActId))
	v.Set("is_box", fmt.Sprintf("%d", param.IsBox))
	body := v.Encode()
	return body, nil
}

func (t *StoreOper) CancelAstByActId(param *xbase.BaseActParam) (*xbase.BaseResp, *xbase.RequestRes, error) {
	if err := param.Valid(); err != nil {
		return nil, nil, err
	}

	body, err := t.genCancelAstByActIdBody(param)
	if err != nil {
		t.Logger.Warn("fail to generate value for cancel ast by act_id, err: %v, param: %+v", err, *param)
		return nil, nil, err
	}
	res, err := t.Post(xbase.StoreApiCancelAstByActId, body)
	if err != nil {
		t.Logger.Warn("post request xasset failed.err:%v", err)
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

	t.Logger.Trace("operate succ. [act_id: %v] [url: %s] [request_id: %s] [trace_id: %s]",
		param.ActId, res.ReqUrl, resp.RequestId, t.GetTarceId(res.Header))
	return &resp, res, nil
}

func (t *StoreOper) genQueryActAstBody(param *xbase.BaseAstParam) (string, error) {
	v := url.Values{}
	v.Set("asset_id", fmt.Sprintf("%d", param.AssetId))
	v.Set("act_id", fmt.Sprintf("%d", param.ActId))
	body := v.Encode()
	return body, nil
}

func (t *StoreOper) QueryActAst(param *xbase.BaseAstParam) (*xbase.QueryActAstResp, *xbase.RequestRes, error) {
	if param.ActId < 1 || param.AssetId < 1 {
		return nil, nil, xbase.ErrParamInvalid
	}

	body, err := t.genQueryActAstBody(param)
	if err != nil {
		t.Logger.Warn("fail to generate value for query act ast, err: %v, param: %+v", err, *param)
		return nil, nil, err
	}
	res, err := t.Post(xbase.StoreApiQueryAst, body)
	if err != nil {
		t.Logger.Warn("post request xasset failed.err:%v", err)
		return nil, nil, xbase.ComErrRequsetFailed
	}
	if res.HttpCode != 200 {
		t.Logger.Warn("post request response is not 200. [http_code: %d] [url: %s] [body: %s] [trace_id: %s]",
			res.HttpCode, res.ReqUrl, res.Body, t.GetTarceId(res.Header))
		return nil, nil, xbase.ComErrRespCodeErr
	}

	var resp xbase.QueryActAstResp
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

	t.Logger.Trace("operate succ. [asset_id: %v] [url: %s] [request_id: %s] [trace_id: %s]",
		param.AssetId, res.ReqUrl, resp.RequestId, t.GetTarceId(res.Header))
	return &resp, res, nil
}

func (t *StoreOper) genListActAstBody(param *xbase.BaseActParam) (string, error) {
	v := url.Values{}
	v.Set("act_id", fmt.Sprintf("%d", param.ActId))
	body := v.Encode()
	return body, nil
}

func (t *StoreOper) ListActAst(param *xbase.BaseActParam) (*xbase.ListActAstResp, *xbase.RequestRes, error) {
	if param.ActId < 1 {
		return nil, nil, xbase.ErrParamInvalid
	}

	body, err := t.genListActAstBody(param)
	if err != nil {
		t.Logger.Warn("fail to generate value for list act ast, err: %v, param: %+v", err, *param)
		return nil, nil, err
	}
	res, err := t.Post(xbase.StoreApiListAst, body)
	if err != nil {
		t.Logger.Warn("post request xasset failed.err:%v", err)
		return nil, nil, xbase.ComErrRequsetFailed
	}
	if res.HttpCode != 200 {
		t.Logger.Warn("post request response is not 200. [http_code: %d] [url: %s] [body: %s] [trace_id: %s]",
			res.HttpCode, res.ReqUrl, res.Body, t.GetTarceId(res.Header))
		return nil, nil, xbase.ComErrRespCodeErr
	}

	var resp xbase.ListActAstResp
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

	t.Logger.Trace("operate succ. [act_id: %v] [url: %s] [request_id: %s] [trace_id: %s]",
		param.ActId, res.ReqUrl, resp.RequestId, t.GetTarceId(res.Header))
	return &resp, res, nil
}

// CreateOrder creates orders.
func (t *StoreOper) CreateOrder(param *xbase.HubCreateOrderParam, uid int64, auth string) (*xbase.HubCreateResp, *xbase.RequestRes, error) {
	var err error
	if err = param.Valid(); err != nil {
		return nil, nil, xbase.ErrParamInvalid
	}
	// 使用百度收银台H5支付组件请务必携带鉴权串
	if param.Code == xbase.CodeBaiduH5 && auth == "" {
		return nil, nil, xbase.ErrParamInvalid
	}
	// 使用百度收银台请务必携带uk
	if _, ok := xbase.BaiduCashierCode[param.Code]; ok {
		if uid <= 0 {
			return nil, nil, xbase.ErrParamInvalid
		}
	}
	var secretAuth, uk string
	if auth != "" {
		secretAuth, err = t.GenSecretData(auth)
		if err != nil {
			return nil, nil, err
		}
	}
	if uid > 0 {
		uk, err = t.GenSecretData(fmt.Sprintf("%d", uid))
		if err != nil {
			return nil, nil, err
		}
	}

	v := url.Values{}
	v.Set("code", fmt.Sprintf("%d", param.Code))
	v.Set("order_type", fmt.Sprintf("%d", param.OrderType))
	v.Set("executor", param.ExecutorAPI)
	v.Set("executor_data", param.ExecutorData)
	v.Set("timestamp", fmt.Sprintf("%d", param.Timestamp))
	v.Set("time_expire", fmt.Sprintf("%d", param.TimeExpire))
	v.Set("profit_sharing", fmt.Sprintf("%d", param.ProfitSharing))
	v.Set("uk", uk)
	v.Set("creator_details", param.Details)
	v.Set("act_id", fmt.Sprintf("%d", param.ActId))
	v.Set("asset_id", fmt.Sprintf("%d", param.AssetId))
	v.Set("buyer_addr", param.BuyerAddr)
	v.Set("seller_addr", param.SellerAddr)
	v.Set("client_type", fmt.Sprintf("%d", param.ClientType))
	v.Set("chan", fmt.Sprintf("%d", param.Chan))
	v.Set("scene", fmt.Sprintf("%d", param.Scene))
	v.Set("signed_auth", secretAuth)
	v.Set("buy_count", fmt.Sprintf("%d", param.BuyCount))
	body := v.Encode()

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
func (t *StoreOper) ConfirmOrder(param *xbase.HubConfirmH5OrderParam, auth string) (*xbase.HubCreateResp, *xbase.RequestRes, error) {
	var err error
	if err = param.Valid(); err != nil {
		return nil, nil, xbase.ErrParamInvalid
	}
	// 使用百度收银台H5支付时，请提供鉴权串
	var secretAuth string
	if auth != "" {
		secretAuth, err = t.GenSecretData(auth)
		if err != nil {
			return nil, nil, err
		}
	}
	v := url.Values{}
	v.Set("code", fmt.Sprintf("%d", param.Code))
	v.Set("order_type", fmt.Sprintf("%d", param.OrderType))
	v.Set("oid", fmt.Sprintf("%d", param.Oid))
	v.Set("client_type", fmt.Sprintf("%d", param.ClientType))
	v.Set("creator_details", param.Details)
	v.Set("signed_auth", secretAuth)
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
func (t *StoreOper) QueryOrderDetail(param *xbase.HubOrderDetailParam) (*xbase.HubOrderDetailResp, *xbase.RequestRes, error) {
	if err := param.Valid(); err != nil {
		return nil, nil, xbase.ErrParamInvalid
	}
	v := url.Values{}
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

// EditOrder edits order info.
func (t *StoreOper) EditOrder(param *xbase.HubEditOrderParam) (*xbase.BaseResp, *xbase.RequestRes, error) {
	if err := param.Valid(); err != nil {
		return nil, nil, xbase.ErrParamInvalid
	}
	v := url.Values{}
	v.Set("oid", fmt.Sprintf("%d", param.Oid))
	v.Set("status", fmt.Sprintf("%d", param.Status))
	v.Set("pay_channel", fmt.Sprintf("%d", param.PayChannel))
	v.Set("third_oid", param.ThirdOid)
	v.Set("pay_info", param.PayInfo)
	v.Set("pay_time", fmt.Sprintf("%d", param.PayTime))
	v.Set("close_time", fmt.Sprintf("%d", param.CloseTime))
	v.Set("close_reason", param.CloseReason)
	body := v.Encode()

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
func (t *StoreOper) QueryOrderList(param *xbase.HubListOrderParam) (*xbase.HubListOrderResp, *xbase.RequestRes, error) {
	if err := param.Valid(); err != nil {
		return nil, nil, xbase.ErrParamInvalid
	}
	v := url.Values{}
	v.Set("address", param.Addr)
	v.Set("status", fmt.Sprintf("%d", param.Status))
	v.Set("cursor", param.Cursor)
	v.Set("limit", fmt.Sprintf("%d", param.Limit))
	v.Set("time_begin", fmt.Sprintf("%d", param.TimeBegin))
	v.Set("time_end", fmt.Sprintf("%d", param.TimeEnd))
	v.Set("monotonicity", fmt.Sprintf("%d", param.Mono))

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

// QueryOrderPage gets order pages by address.
func (t *StoreOper) QueryOrderPage(param *xbase.HubOrderPageParam) (*xbase.HubOrderPageResp, *xbase.RequestRes, error) {
	if err := param.Valid(); err != nil {
		return nil, nil, xbase.ErrParamInvalid
	}
	v := url.Values{}
	v.Set("address", param.Addr)
	v.Set("status", fmt.Sprintf("%d", param.Status))
	v.Set("page", fmt.Sprintf("%d", param.Page))
	v.Set("size", fmt.Sprintf("%d", param.Size))
	v.Set("time_begin", fmt.Sprintf("%d", param.TimeBegin))
	v.Set("time_end", fmt.Sprintf("%d", param.TimeEnd))

	body := v.Encode()
	res, err := t.Post(xbase.HubListOrderPage, body)
	if err != nil {
		t.Logger.Warn("post request xasset failed, uri: %s, err: %v", xbase.HubListOrderPage, err)
		return nil, nil, xbase.ComErrRequsetFailed
	}
	if res.HttpCode != 200 {
		t.Logger.Warn("post request response is not 200. [http_code: %d] [url: %s] [body: %s] [trace_id: %s]",
			res.HttpCode, res.ReqUrl, res.Body, t.GetTarceId(res.Header))
		return nil, nil, xbase.ComErrRespCodeErr
	}

	var resp xbase.HubOrderPageResp
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

// CountOrder count valid orders.
func (t *StoreOper) CountOrder(param *xbase.CountOrderParam) (*xbase.CountOrderResp, *xbase.RequestRes, error) {
	if err := param.Valid(); err != nil {
		return nil, nil, xbase.ErrParamInvalid
	}
	v := url.Values{}
	v.Set("asset_id", fmt.Sprintf("%d", param.AssetId))
	v.Set("status", fmt.Sprintf("%d", param.Status))

	body := v.Encode()
	res, err := t.Post(xbase.CountOrder, body)
	if err != nil {
		t.Logger.Warn("post request xasset failed, uri: %s, err: %v", xbase.CountOrder, err)
		return nil, nil, xbase.ComErrRequsetFailed
	}
	if res.HttpCode != 200 {
		t.Logger.Warn("post request response is not 200. [http_code: %d] [url: %s] [body: %s] [trace_id: %s]",
			res.HttpCode, res.ReqUrl, res.Body, t.GetTarceId(res.Header))
		return nil, nil, xbase.ComErrRespCodeErr
	}

	var resp xbase.CountOrderResp
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

// SumOrderPrice sum valid orders price.
func (t *StoreOper) SumOrderPrice(param *xbase.SumOrderPriceParam) (*xbase.SumOrderPriceResp, *xbase.RequestRes, error) {
	if err := param.Valid(); err != nil {
		return nil, nil, xbase.ErrParamInvalid
	}
	v := url.Values{}
	v.Set("status", fmt.Sprintf("%d", param.Status))

	body := v.Encode()
	res, err := t.Post(xbase.SumOrderPrice, body)
	if err != nil {
		t.Logger.Warn("post request xasset failed, uri: %s, err: %v", xbase.SumOrderPrice, err)
		return nil, nil, xbase.ComErrRequsetFailed
	}
	if res.HttpCode != 200 {
		t.Logger.Warn("post request response is not 200. [http_code: %d] [url: %s] [body: %s] [trace_id: %s]",
			res.HttpCode, res.ReqUrl, res.Body, t.GetTarceId(res.Header))
		return nil, nil, xbase.ComErrRespCodeErr
	}

	var resp xbase.SumOrderPriceResp
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

func (t *StoreOper) GenSecretData(data string) (string, error) {
	input := fmt.Sprintf("%d_%s_%s", t.Cfg.Credentials.AppId, t.Cfg.Credentials.AccessKeyId, t.Cfg.Credentials.SecretAccessKey)
	h := md5.New()
	io.WriteString(h, input)
	digest := h.Sum(nil)
	key := fmt.Sprintf("%X", digest)
	return utils.AesEncode(data, key)
}
