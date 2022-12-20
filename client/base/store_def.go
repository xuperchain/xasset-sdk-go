package base

import "fmt"

const (
	StoreApiCreate           = "/xasset/store/v1/create"
	StoreApiAlter            = "/xasset/store/v1/alter"
	StoreApiQuery            = "/xasset/store/v1/query"
	StoreApiList             = "/xasset/store/v1/liststorebyapp"
	StoreApiCreateAct        = "/xasset/store/v1/createact"
	StoreApiAlterAct         = "/xasset/store/v1/alteract"
	StoreApiRemoveAct        = "/xasset/store/v1/removeact"
	StoreApiQueryAct         = "/xasset/store/v1/queryact"
	StoreApiListAct          = "/xasset/store/v1/listact"
	StoreApiPubAct           = "/xasset/store/v1/pubact"
	StoreApiBindAst          = "/xasset/store/v1/bindast"
	StoreApiAlterAst         = "/xasset/store/v1/alterast"
	StoreApiCancelAst        = "/xasset/store/v1/cancelast"
	StoreApiCancelAstByActId = "/xasset/store/v1/cancelastbyact"
	StoreApiQueryAst         = "/xasset/store/v1/queryast"
	StoreApiListAst          = "/xasset/store/v1/listast"

	HubCreateOrder   = "/xasset/trade/v1/create_order"
	HubConfirmOrder  = "/xasset/trade/v1/confirm_order"
	HubDetailOrder   = "/xasset/trade/v1/order_detail"
	HubEditOrder     = "/xasset/trade/v1/edit_order"
	HubListOrder     = "/xasset/trade/v1/order_list"
	HubListOrderPage = "/xasset/trade/v1/order_page"
	CountOrder       = "/xasset/trade/v1/count_order"
)

/////// Create Store /////////
type CreateOrAlterStoreParam struct {
	StoreId   int    `json:"store_id"`
	Name      string `json:"name"`
	Logo      string `json:"logo"`
	Cover     string `json:"cover"`
	ShortDesc string `json:"short_desc"`
	Weight    int    `json:"weight,omitempty"`
	ExtInfo   string `json:"ext_info,omitempty"`
	Wechat    string `json:"wechat,omitempty"`
}

func (t *CreateOrAlterStoreParam) CreateValid() error {
	if t == nil {
		return ErrNilPointer
	}
	if t.StoreId < 1 || t.Name == "" || t.Logo == "" || t.Cover == "" {
		return ErrParamInvalid
	}
	return nil
}

////////// Query Store //////////
type BaseStoreParam struct {
	StoreId int `json:"store_id"`
}

func (t *BaseStoreParam) Valid() error {
	if t == nil {
		return ErrNilPointer
	}
	if t.StoreId < 1 {
		return ErrParamInvalid
	}
	return nil
}

type QueryStoreResp struct {
	BaseResp
	Meta *QueryStoreMeta `json:"meta"`
}

type ListStoreResp struct {
	BaseResp
	List []*QueryStoreMeta `json:"list"`
}

type QueryStoreMeta struct {
	AppId     int64  `json:"app_id"`
	StoreId   int    `json:"store_id"`
	Name      string `json:"name"`
	Logo      string `json:"logo"`
	Cover     string `json:"cover"`
	ShortDesc string `json:"short_desc"`
	Status    int    `json:"status"`
	Ctime     int64  `json:"ctime"`
	Mtime     int64  `json:"mtime"`
	Weight    int    `json:"weight"`
	ExtInfo   string `json:"ext_info"`
	Wechat    string `json:"wechat"`
	Likes     int64  `json:"int64"`
}

/////// Create Act /////////
type CreateOrAlterActParam struct {
	StoreId   int    `json:"store_id"`
	ActId     int64  `json:"act_id"`
	JumpLink  string `json:"jump_link"`
	Issuer    string `json:"issuer"`
	ActName   string `json:"act_name"`
	ShortDesc string `json:"short_desc"`
	Thumb     string `json:"thumb"`
	ImgDesc   string `json:"img_desc"`
	Start     int64  `json:"start"`
	End       int64  `json:"end"`
	Weight    int    `json:"weight,omitempty"`
	ExtInfo   string `json:"ext_info,omitempty"`
}

func (t *CreateOrAlterActParam) CreateValid() error {
	if t == nil {
		return ErrNilPointer
	}
	if t.StoreId < 1 || t.ActId < 1 || t.Issuer == "" || t.ActName == "" || t.Thumb == "" || t.Start < 1 || t.End < 1 || t.Start > t.End {
		return ErrParamInvalid
	}
	return nil
}

type BaseActParam struct {
	ActId  int64 `json:"act_id"`
	OpType int   `json:"op_type"`
	IsBox  int   `json:"is_box"`
}

func (t *BaseActParam) Valid() error {
	if t == nil {
		return ErrNilPointer
	}
	if t.ActId < 1 {
		return ErrParamInvalid
	}
	return nil
}

type QueryActResp struct {
	BaseResp
	Meta *QueryActMeta `json:"meta"`
}

type QueryActMeta struct {
	AppId         int64    `json:"app_id"`
	StoreId       int      `json:"store_id"`
	StoreName     string   `json:"store_name"`
	ActId         int64    `json:"act_id"`
	JumpLink      string   `json:"jump_link"`
	Issuer        string   `json:"issuer"`
	ActName       string   `json:"act_name"`
	Thumb         []string `json:"thumb"`
	ImgDesc       []string `json:"img_desc"`
	ShortDesc     string   `json:"short_desc"`
	Status        int      `json:"status"`
	PublishStatus int      `json:"publish_status"`
	Start         int64    `json:"start"`
	End           int64    `json:"end"`
	Weight        int      `json:"weight"`
	ExtInfo       string   `json:"ext_info"`
	Ctime         int64    `json:"ctime"`
	Mtime         int64    `json:"mtime"`
}

type ListActParam struct {
	StoreId int    `json:"store_id"`
	Cursor  string `json:"cursor"`
	Limit   int    `json:"limie"`
}

func (t *ListActParam) Valid() error {
	if t == nil {
		return ErrNilPointer
	}
	if t.StoreId < 1 {
		return ErrParamInvalid
	}
	return nil
}

type ListActResp struct {
	BaseResp
	List    []*QueryActMeta `json:"list"`
	Cursor  string          `json:"cursor"`
	HasMore int             `json:"has_more"`
}

type BindOrAlterAstParam struct {
	AssetId   int64  `json:"asset_id"`
	Amount    int    `json:"amount"`
	ApplyForm int    `json:"apply_form"`
	AssetCate int    `json:"asset_cate"`
	Price     int64  `json:"price"`
	OriPrice  int64  `json:"ori_price"`
	ActId     int64  `json:"act_id"`
	Start     int64  `json:"start"`
	End       int64  `json:"end"`
	GrantMode int    `json:"grant_mode"`
	JumpLink  string `json:"jump_link"`
	ExtInfo   string `json:"ext_info"`
	IsBox     int    `json:"is_box"`
}

func (t *BindOrAlterAstParam) CreateValid() error {
	if t == nil {
		return ErrNilPointer
	}
	if t.ActId < 1 && t.AssetId < 1 && t.End < t.Start {
		return ErrParamInvalid
	}

	return nil
}

func (t *BindOrAlterAstParam) AlterValid() error {
	if t == nil {
		return ErrNilPointer
	}

	if t.AssetId < 1 || t.ActId < 1 {
		return ErrParamInvalid
	}
	return nil
}

type BaseAstParam struct {
	AssetId int64 `json:"asset_id"`
	ActId   int64 `json:"act_id"`
}

func (t *BaseAstParam) Valid() error {
	if t == nil {
		return ErrNilPointer
	}

	if t.AssetId < 1 || t.ActId < 1 {
		return ErrParamInvalid
	}
	return nil
}

type QueryActAstResp struct {
	BaseResp
	Meta *QueryActAstMeta `json:"meta"`
}

type ListActAstResp struct {
	BaseResp
	List []*QueryActAstMeta `json:"list"`
}

type QueryActAstMeta struct {
	AppId         int64    `json:"app_id"`
	Addr          string   `json:"addr"`
	AssetId       int64    `json:"asset_id"`
	AssetCate     int      `json:"asset_cate"`
	Thumb         []string `json:"thumb"`
	Title         string   `json:"title"`
	ShortDesc     string   `json:"short_desc"`
	TxId          string   `json:"tx_id"`
	AssetUrl      []string `json:"asset_url"`
	ImgDesc       []string `json:"img_desc"`
	ActId         int64    `json:"act_id"`
	Start         int64    `json:"start"`
	End           int64    `json:"end"`
	Price         int64    `json:"price"`
	OriPrice      int64    `json:"ori_price"`
	Amount        int64    `json:"amount"`
	ApplyForm     int      `json:"apply_form"`
	GrantMode     int      `json:"grant_mode"`
	Status        int      `json:"status"`
	PublishStatus int      `json:"publish_status"`
	JumpLink      string   `json:"jump_link"`
	ExtInfo       string   `json:"ext_info"`
	Ctime         int64    `json:"ctime"`
	Mtime         int64    `json:"mtime"`
}

////////////// Trade Orders API /////////////
const (
	// ---------- Code Enum -----------
	CodeBaiduSmartApp = 1001 // 百度收银台安卓端
	CodeBaiduIOS      = 1002 // 百度收银台IOS端
	CodeBaiduH5       = 1003 // 百度收银台H5
)

var (
	BaiduCashierCode = map[int]interface{}{
		CodeBaiduSmartApp: struct{}{},
		CodeBaiduIOS:      struct{}{},
		CodeBaiduH5:       struct{}{},
	}
)

//////////// Create Orders ///////////////
type HubCreateOrderParam struct {
	// 标记支付平台, 1001: 百度收银台-小程序 1002: 百度收银台-IOS 1003: 百度收银台-H5  2001: 第三方微信支付
	Code int `json:"code"`
	// 订单类型, 1:PGC商城订单
	OrderType int `json:"order_type"`
	// 使用百度收银台，且有回调通知的业务方服务请填写回调地址，支付成功后，执行器回调接口，应用方需保证接口幂等
	ExecutorAPI string `json:"executor"`
	// object, 使用百度收银台时，支付成功后，执行器回调时携带参数
	ExecutorData string `json:"executor_data"`
	// 下单时间，建议填写now time
	Timestamp int64 `json:"timestamp"`
	// 订单失效时间，针对使用百度收银台的订单有效
	// 当TimeExpire为0时，表示永久不过期
	// TimeExpire为秒偏移量，以create_time时间为准偏移
	TimeExpire int64 `json:"time_expire"`
	// 是否需要分账，0: 不分账，1: 分账
	ProfitSharing int `json:"profit_sharing"`
	// object，其余支付参数的序列化值，标记剩余的非通用参数
	Details    string `json:"creator_details"`
	ActId      int64  `json:"act_id"`
	AssetId    int64  `json:"asset_id"`
	BuyerAddr  string `json:"buyer_addr"`
	SellerAddr string `json:"seller_addr"`
	ClientType int    `json:"client_type"`
	Chan       int64  `json:"chan"`
	Scene      int64  `json:"scene"`
}

func (p *HubCreateOrderParam) Valid() error {
	if p.SellerAddr == "" {
		return fmt.Errorf("seller_addr empty")
	}

	if p.AssetId <= 0 {
		return fmt.Errorf("asset_id invalid, must be a positive integer")
	}
	return nil
}

type HubCreateResp struct {
	BaseResp
	Data HubCreateData `json:"data"`
}

type HubCreateData struct {
	Code      int    `json:"code"`
	OrderType int    `json:"order_type"`
	Details   string `json:"details"`
	CTime     int64  `json:"ctime"`
}

type HubConfirmH5OrderParam struct {
	Code       int   `json:"code"`
	OrderType  int   `json:"order_type"`
	Oid        int64 `json:"oid"`
	ClientType int   `json:"client_type"`
}

func (p *HubConfirmH5OrderParam) Valid() error {
	if p.Oid < 0 {
		return fmt.Errorf("oid empty")
	}
	if p.ClientType < 0 {
		return fmt.Errorf("client_type empty")
	}
	return nil
}

type HubOrderDetailParam struct {
	Oid int64 `json:"oid"`
}

func (p *HubOrderDetailParam) Valid() error {
	if p.Oid < 0 {
		return fmt.Errorf("value empty")
	}
	return nil
}

type HubOrderDetailResp struct {
	BaseResp
	Data HubOrderDetail `json:"data"`
}

type HubOrderDetail struct {
	Code        int      `json:"code"`
	OrderType   int      `json:"order_type"`
	Oid         int64    `json:"oid"`
	ActId       int64    `json:"act_id"`
	AssetId     int64    `json:"asset_id"`
	ShardId     int64    `json:"shard_id"`
	AssetCate   int      `json:"asset_cate"`
	BuyerAddr   string   `json:"buyer_addr"`
	Status      int      `json:"status"`
	Title       string   `json:"title"`
	Thumb       []string `json:"thumb"`
	OriginPrice int      `json:"origin_price"`
	PayPrice    int      `json:"pay_price"`
	TimeExpire  int64    `json:"time_expire"`
	PayTime     int64    `json:"pay_time"`
	CloseTime   int64    `json:"close_time"`
	Ctime       int64    `json:"ctime"`
}

type HubEditOrderParam struct {
	Oid         int64  `json:"oid"`
	Status      int    `json:"status"`
	PayChannel  int    `json:"pay_channel"`
	ThirdOid    string `json:"third_oid"`
	PayInfo     string `json:"pay_info"`
	PayTime     int64  `json:"pay_time"`
	CloseTime   int64  `json:"close_time"`
	CloseReason string `json:"close_reason"`
}

func (p *HubEditOrderParam) Valid() error {
	if p.Oid < 0 {
		return fmt.Errorf("oid invalid")
	}
	return nil
}

type HubListOrderParam struct {
	Addr      string `json:"address"`
	Status    int    `json:"status"`
	Cursor    string `json:"cursor"`
	Limit     int    `json:"limit"`
	TimeBegin int64  `json:"time_begin"`
	TimeEnd   int64  `json:"time_end"`
}

func (p *HubListOrderParam) Valid() error {
	if p.Status < 0 {
		return fmt.Errorf("status invalid")
	}
	if p.Limit < 0 {
		return fmt.Errorf("cursor limit invalid")
	}

	return nil
}

type HubListOrderResp struct {
	BaseResp
	Data HubListOrderData `json:"data"`
}

type HubListOrderData struct {
	List    []HubOrderDetail `json:"list"`
	Cursor  string           `json:"cursor"`
	HasMore int              `json:"has_more"`
}

type H5OrderItem struct {
	TpOrderId    int64  `json:"oid"`
	OrderInfoUrl string `json:"order_url"`
	H5PayInfo    string `json:"pay_info"`
	TotalAmount  string `json:"total_amount"`
	CTime        int64  `json:"ctime"`
}

type HubOrderPageParam struct {
	Addr      string `json:"address"`
	Status    int    `json:"status"`
	Page      int    `json:"page"`
	Size      int    `json:"size"`
	TimeBegin int64  `json:"time_begin"`
	TimeEnd   int64  `json:"time_end"`
}

func (p *HubOrderPageParam) Valid() error {
	if p.Status < 0 {
		return fmt.Errorf("status invalid")
	}
	if p.Page < 0 {
		return fmt.Errorf("cursor page invalid")
	}
	if p.Size < 0 {
		return fmt.Errorf("cursor size invalid")
	}

	return nil
}

type HubOrderPageData struct {
	List  []HubOrderDetail `json:"list"`
	Total int64            `json:"total"`
}

type HubOrderPageResp struct {
	BaseResp
	Data HubOrderPageData `json:"data"`
}

type CountOrderParam struct {
	AssetId int64 `json:"asset_id"`
	Status  int   `json:"status"`
}

func (p *CountOrderParam) Valid() error {
	if p.AssetId <= 0 {
		return fmt.Errorf("asset_id invalid, must be a positive integer")
	}
	if p.Status < 0 {
		return fmt.Errorf("status invalid")
	}
	return nil
}

type CountOrderData struct {
	Total int64 `json:"total"`
}

type CountOrderResp struct {
	BaseResp
	Data CountOrderData `json:"data"`
}
