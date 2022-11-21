package base

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
	AssetId       int64   `json:"asset_id"`
	AssetCate     int      `json:"asset_cate"`
	Thumb         []string `json:"thumb"`
	Title         string   `json:"title"`
	ShortDesc     string   `json:"short_desc"`
	TxId          string   `json:"tx_id"`
	AssetUrl      []string `json:"asset_url"`
	ImgDesc       []string `json:"img_desc"`
	ActId         int64   `json:"act_id"`
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
