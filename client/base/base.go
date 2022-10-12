package base

import (
	"crypto/md5"
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"time"

	"github.com/xuperchain/xasset-sdk-go/auth"
	"github.com/xuperchain/xasset-sdk-go/common/config"
	"github.com/xuperchain/xasset-sdk-go/common/httpcli"
	"github.com/xuperchain/xasset-sdk-go/common/logs"
)

// 常用错误
var (
	ComErrParamInvalid        = errors.New("param invalid")
	ComErrAccountSignFailed   = errors.New("account sign failed")
	ComErrJsonMarFailed       = errors.New("json marhsal failed")
	ComErrRequsetFailed       = errors.New("send request failed")
	ComErrRespCodeErr         = errors.New("request resp code error")
	ComErrUnmarshalBodyFailed = errors.New("json unmarhsal body failed")
	ComErrServRespErrnoErr    = errors.New("server resp errno error")
	ComErrGenRequestFailed    = errors.New("generate http request failed")
	ComErrXassetSignFailed    = errors.New("xasser access sign failed")
	ComErrConfigErr           = errors.New("client config error")
)

// 服务端相应错误码
const (
	XassetErrNoSucc = 0
)

// 资产类型
type AssetType int

const (
	_ AssetType = iota
	// 1.艺术品
	AssetCateArt
	// 2.收藏品
	AssetCateCollect
	// 3.门票
	AssetCateTicket
	// 4.酒店
	AssetCateHotel
)

// 列表分页范围限制
const (
	// 最大范围
	MaxLimit = 50
)

// 支付渠道类型
const (
	// 微信支付
	PayByWechat = 1
	// 阿里支付
	PayByAli = 2
	// 特殊支付
	PayBySpecial = 9
)

var (
	ErrParamInvalid      = errors.New("param invalid")
	ErrAssetInvalid      = errors.New("asset invalid, must be a positive integer")
	ErrAddressInvalid    = errors.New("address invalid, empty string")
	ErrUserIdInvalid     = errors.New("user id invalid, must be a positive integer")
	ErrAmountInvalid     = errors.New("amount invalid, must be a positive integer or a zero value")
	ErrPriceInvalid      = errors.New("price invalid, must be a positive integer or a zero value")
	ErrNilPointer        = errors.New("target paramater valid, nil pointer")
	ErrDescInvalid       = errors.New("target paramater invalid, empty string")
	ErrAssetTypeInvalid  = errors.New("type invalid, must between 1 to 4")
	ErrAlterInfo         = errors.New("alter info invalid")
	ErrAlterAssetInvalid = errors.New("param for altering invalid, must contain amount or valid asset info")
	ErrShardInvalid      = errors.New("shard invalid, must be a positive integer")
	ErrImgInvalid        = errors.New("imgs invalid")
	ErrEvidenceInvalid   = errors.New("evidence type invalid, must between 0 to 1")
	ErrBytesInvalid      = errors.New("bytes invalid, nil pointer")
	ErrStatusInvalid     = errors.New("status invalid")
	ErrAssetListInvalid  = errors.New("null asset list")
	ErrUnionIdInvalid    = errors.New("union id invalid")
	ErrOpenIdInvalid     = errors.New("open id invalid")
	ErrAppKeyInvalid     = errors.New("app key invalid")
	ErrMnemInvalid       = errors.New("mnemonic invalid")
	ErrNameInvalid       = errors.New("target parameter invalid, empty string")
)

type ThumbMap struct {
	Urls   map[string]string `json:"urls"`
	Width  string            `json:"width"`
	Height string            `json:"height"`
}

func HasAssetType(t AssetType) bool {
	return t != 0
}

func AssetTypeValid(t AssetType) error {
	if t >= AssetCateArt && t <= AssetCateHotel {
		return nil
	}
	return ErrAssetTypeInvalid
}

func AssetIdValid(assetId int64) error {
	if assetId <= 0 {
		return ErrAssetInvalid
	}
	return nil
}

func AccountValid(account *auth.Account) error {
	if account == nil {
		return ErrNilPointer
	}
	return nil
}

func EvidenceValid(evidence int) error {
	if evidence != 0 && evidence != 1 {
		return ErrEvidenceInvalid
	}
	return nil
}

func AddrValid(addr string) error {
	if addr == "" {
		return ErrAddressInvalid
	}
	return nil
}

func HasId(id int64) bool {
	return id != 0
}

func IdValid(id int64) error {
	if id <= 0 {
		return ErrUserIdInvalid
	}
	return nil
}

func PriceInvalid(price int64) error {
	if price < 0 {
		return ErrPriceInvalid
	}
	return nil
}

func AmountInvalid(amount int) error {
	if amount < 0 {
		return ErrAmountInvalid
	}
	return nil
}

func ByteValid(b []byte) error {
	if b == nil {
		return ErrBytesInvalid
	}
	return nil
}

func HasDesc(desc string) bool {
	return desc != ""
}

func DescValid(desc string) error {
	if desc == "" {
		return ErrDescInvalid
	}
	return nil
}

func HasImg(imgs []string) bool {
	return len(imgs) != 0
}

func ImgValid(imgs []string) error {
	if imgs == nil {
		return ErrNilPointer
	}
	if len(imgs) == 0 {
		return ErrImgInvalid
	}
	return nil
}

func FileValid(files []string) error {
	if files == nil {
		return ErrNilPointer
	}
	if len(files) == 0 {
		return ErrImgInvalid
	}
	return nil
}

func ShardIdValid(id int64) error {
	if id <= 0 {
		return ErrShardInvalid
	}
	return nil
}

func StatusValid(status int) error {
	if status < 0 {
		return ErrStatusInvalid
	}
	return nil
}

func UnionIdValid(uid string) error {
	if uid == "" {
		return ErrUnionIdInvalid
	}
	return nil
}

func OpenIdValid(oid string) error {
	if oid == "" {
		return ErrOpenIdInvalid
	}
	return nil
}

func AppKeyValid(ak string) error {
	if ak == "" {
		return ErrAppKeyInvalid
	}
	return nil
}

func MnemonicValid(mnem string) error {
	if mnem == "" {
		return ErrMnemInvalid
	}
	return nil
}

// ///// General Client ///////
type RequestRes struct {
	HttpCode int         `json:"http_code"`
	ReqUrl   string      `json:"req_url"`
	Header   http.Header `json:"header"`
	Body     string      `json:"body"`
}

type BaseResp struct {
	RequestId string `json:"request_id"`
	Errno     int    `json:"errno"`
	Errmsg    string `json:"errmsg"`
}

// XassetBaseClient
type XassetBaseClient struct {
	Cfg    *config.XassetCliConfig
	Logger *logs.Logger
}

func (t *XassetBaseClient) InitClient(cfg *config.XassetCliConfig, logger logs.LogDriver) error {
	if cfg == nil || !cfg.IsVaild() {
		return ComErrParamInvalid
	}

	t.Cfg = cfg
	t.Logger = logs.NewLogger(logger)

	return nil
}

func (t *XassetBaseClient) GetConfig() *config.XassetCliConfig {
	return t.Cfg
}

func (t *XassetBaseClient) Post(uri, data string) (*RequestRes, error) {
	reqUrl := fmt.Sprintf("%s%s", t.GetConfig().Endpoint, uri)
	u, err := url.Parse(reqUrl)
	if err != nil {
		t.Logger.Warn("url error.[url:%s] [err:%v]", reqUrl, err)
		return nil, ComErrConfigErr
	}
	header := map[string]string{
		"Content-Type": "application/x-www-form-urlencoded;charset=utf-8",
		"Host":         u.Hostname(),
		"Timestamp":    fmt.Sprintf("%d", time.Now().Unix()),
		"Content-Md5":  fmt.Sprintf("%x", md5.Sum([]byte(data))),
	}

	req, err := httpcli.GenRequest("POST", reqUrl, header, data)
	if err != nil {
		t.Logger.Warn("generate request failed.[err:%v]", err)
		return nil, ComErrGenRequestFailed
	}
	sign, err := auth.Sign(req, t.GetConfig().Credentials, t.GetConfig().SignOption)
	if err != nil {
		return nil, ComErrXassetSignFailed
	}
	req.Header.Set("Authorization", sign)

	opts := make(map[string]string)
	if httpcli.IsHttps(reqUrl) {
		opts[httpcli.OptTlsSipVerify] = "1"
	}
	resp, err := httpcli.SendRequest(req, t.GetConfig().ConnectTimeoutMs,
		t.GetConfig().ReadWriteTimeoutMs, opts)
	if err != nil {
		t.Logger.Warn("send http request failed.[url:%s] [err:%v]", reqUrl, err)
		return nil, ComErrRequsetFailed
	}

	result := &RequestRes{
		HttpCode: resp.StatusCode,
		ReqUrl:   reqUrl,
		Header:   resp.Header,
		Body:     string(resp.Body),
	}
	return result, nil
}

func (t *XassetBaseClient) GetTarceId(header http.Header) string {
	var traceId string
	if header != nil {
		traceId = header.Get("xasset-trace-id")
	}

	if traceId == "" {
		traceId = "0"
	}
	return traceId
}
