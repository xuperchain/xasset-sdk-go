package store

import (
	"encoding/json"
	"fmt"
	"net/url"

	"github.com/xuperchain/xasset-sdk-go/common/config"
	"github.com/xuperchain/xasset-sdk-go/common/logs"
	xbase "github.com/xuperchain/xasset-sdk-go/client/base"
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

func (t *StoreOper) genPubActBody(param *xbase.BaseActParam) (string, error) {
	v := url.Values{}
	v.Set("act_id", fmt.Sprintf("%d", param.ActId))
	v.Set("op_type", fmt.Sprintf("%d", param.OpType))
	body := v.Encode()
	return body, nil
}

func (t *StoreOper) PubAct(param *xbase.BaseActParam) (*xbase.BaseResp, *xbase.RequestRes, error) {
	if err := param.Valid(); err != nil {
		return nil, nil, xbase.ErrParamInvalid
	}

	body, err := t.genPubActBody(param)
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
