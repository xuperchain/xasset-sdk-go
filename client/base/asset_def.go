package base

import (
	"encoding/json"
	"fmt"

	"github.com/xuperchain/xasset-sdk-go/auth"
)

const (
	AssetApiCreate           = "/xasset/horae/v1/create"
	AssetApiAlter            = "/xasset/horae/v1/alter"
	AssetApiPublish          = "/xasset/horae/v1/publish"
	AssetApiQueryAsset       = "/xasset/horae/v1/query"
	AssetApiGrant            = "/xasset/horae/v1/grant"
	AssetApiFreeze           = "/xasset/horae/v1/freeze"
	AssetApiConsume          = "/xasset/horae/v1/consume"
	AssetApiTransfer         = "/xasset/damocles/v1/transfer"
	AssetApiQueryShard       = "/xasset/horae/v1/querysds"
	AssetApiListShardsByAddr = "/xasset/horae/v1/listsdsbyaddr"
	AssetApiListAssetByAddr  = "/xasset/horae/v1/listastbyaddr"
	AssetListShardsByAsset   = "/xasset/horae/v1/listsdsbyast"
	AssetApiGetEvidenceInfo  = "/xasset/horae/v1/getevidenceinfo"
	AssetApiListDiffByAddr   = "/xasset/horae/v1/listdiffbyaddr"
	FileApiGetStoken         = "/xasset/file/v1/getstoken"
	ListAssetHistory         = "/xasset/horae/v1/history"

	SceneListShardByAddr = "/xasset/scene/v1/listsdsbyaddr"
	SceneQueryShard      = "/xasset/scene/v1/qrysdsinfo"
	SceneListDiffByAddr  = "/xasset/scene/v1/listdiffbyaddr"
	SceneListAddr        = "/xasset/scene/v1/listaddr"
	SceneHasAstByAddr    = "/xasset/scene/v1/hasastbyaddr"

	DidApiRegister     = "/xasset/did/v1/bdboxregister"
	DidApiBind         = "/xasset/did/v1/bdboxbind"
	DidApiBindByUid    = "/xasset/did/v1/bindbyunionid"
	DidApiGetAddrByUid = "/xasset/did/v1/getaddrbyunionid"

	HubCreateOrder  = "/xasset/trade/v1/create_order"
	HubConfirmOrder = "/xasset/trade/v1/confirm_order"
	HubDetailOrder  = "/xasset/trade/v1/order_detail"
	HubEditOrder    = "/xasset/trade/v1/edit_order"
	HubListOrder    = "/xasset/trade/v1/order_list"
)

/////// Gen Token /////////
type GetStokenParam struct {
	Account *auth.Account `json:"account"`
}

func (t *GetStokenParam) Valid() error {
	if t == nil {
		return ErrNilPointer
	}
	if err := AccountValid(t.Account); err != nil {
		return err
	}
	return nil
}

type AccessInfo struct {
	Bucket       string `json:"bucket"`
	EndPoint     string `json:"endpoint"`
	ObjectPath   string `json:"object_path"`
	AK           string `json:"access_key_id"`
	SK           string `json:"secret_access_key"`
	SessionToken string `json:"session_token"`
	CreateTime   string `json:"createTime"`
	Expiration   string `json:"expiration"`
}

type GetStokenResp struct {
	BaseResp
	AccessInfo *AccessInfo `json:"accessInfo"`
}

//////// Upload File /////////////
// Account 创建资产区块链账户
// FileName 文件名称
// FilePath 文件绝对路径
// DataByte 文件二进制串
// Property 文件属性。例如图片类型文件，则为图片宽高，格式为 width_height
// 注意：文件路径和二进制串为二选一
type UploadFileParam struct {
	Account  *auth.Account `json:"account"`
	FileName string        `json:"file_name"`
	FilePath string        `json:"file_path"`
	DataByte []byte        `json:"data_byte"`
	Property string        `json:"property"`
}

func (t *UploadFileParam) Valid() error {
	if t == nil {
		return ErrNilPointer
	}
	if err := AccountValid(t.Account); err != nil {
		return err
	}
	if err := DescValid(t.FileName); err != nil {
		return err
	}
	pathErr := DescValid(t.FilePath)
	BytesErr := ByteValid(t.DataByte)
	if pathErr != nil && BytesErr != nil {
		return pathErr
	}
	return nil
}

type UploadFileResp struct {
	Link       string      `json:"link"`
	AccessInfo *AccessInfo `json:"accessInfo"`
}

///////// Create Asset ///////////
type CreateAssetInfo struct {
	AssetCate AssetType `json:"asset_cate"`
	Title     string    `json:"title"`
	Thumb     []string  `json:"thumb"`
	ShortDesc string    `json:"short_desc"`
	ImgDesc   []string  `json:"img_desc"`
	AssetUrl  []string  `json:"asset_url"`
	LongDesc  string    `json:"long_desc,omitempty"`
	AssetExt  string    `json:"asset_ext,omitempty"`
	GroupId   int64     `json:"group_id,omitempty"`
}

func CreateAssetInfoValid(p *CreateAssetInfo) error {
	if p == nil {
		return ErrNilPointer
	}
	if err := AssetTypeValid(p.AssetCate); err != nil {
		return err
	}
	if err := DescValid(p.Title); err != nil {
		return err
	}
	if err := DescValid(p.ShortDesc); err != nil {
		return err
	}
	if err := ImgValid(p.Thumb); err != nil {
		return err
	}
	if err := FileValid(p.AssetUrl); err != nil {
		return err
	}
	return nil
}

type CreateAssetParam struct {
	Price     int64            `json:"price,omitempty"`
	Amount    int              `json:"amount"`
	AssetInfo *CreateAssetInfo `json:"asset_info"`
	Account   *auth.Account    `json:"account"`
	UserId    int64            `json:"user_id,omitempty"`
	FileHash  string           `json:"file_hash,omitempty"`
}

func (t *CreateAssetParam) Valid() error {
	if t == nil {
		return ErrNilPointer
	}
	if err := PriceInvalid(t.Price); err != nil {
		return err
	}
	if err := AmountInvalid(t.Amount); err != nil {
		return err
	}
	if err := AccountValid(t.Account); err != nil {
		return err
	}
	if err := CreateAssetInfoValid(t.AssetInfo); err != nil {
		return err
	}
	return nil
}

type CreateAssetResp struct {
	BaseResp
	AssetId int64 `json:"asset_id"`
}

///////// Alter Asset //////////
type AlterAssetInfo struct {
	AssetCate AssetType `json:"asset_cate,omitempty"`
	Title     string    `json:"title,omitempty"`
	Thumb     []string  `json:"thumb,omitempty"`
	ShortDesc string    `json:"short_desc,omitempty"`
	ImgDesc   []string  `json:"img_desc,omitempty"`
	AssetUrl  []string  `json:"asset_url,omitempty"`
	LongDesc  string    `json:"long_desc,omitempty"`
	AssetExt  string    `json:"asset_ext,omitempty"`
	GroupId   int64     `json:"group_id,omitempty"`
}

func AlterAssetInfoValid(p *AlterAssetInfo) error {
	if p == nil {
		return ErrNilPointer
	}
	if HasAssetType(p.AssetCate) {
		if err := AssetTypeValid(p.AssetCate); err != nil {
			return err
		}
	}
	if HasId(p.GroupId) {
		if err := IdValid(p.GroupId); err != nil {
			return err
		}
	}
	if !HasAssetType(p.AssetCate) &&
		!HasDesc(p.Title) && !HasDesc(p.ShortDesc) && !HasDesc(p.LongDesc) && !HasDesc(p.AssetExt) &&
		!HasImg(p.Thumb) && !HasImg(p.ImgDesc) && !HasImg(p.AssetUrl) &&
		!HasId(p.GroupId) {
		return ErrAlterInfo
	}
	return nil
}

type AlterAssetParam struct {
	AssetId   int64           `json:"asset_id"`
	Price     int64           `json:"price,omitempty"`
	Amount    int             `json:"amount,omitempty"`
	FileHash  string          `json:"file_hash"`
	AssetInfo *AlterAssetInfo `json:"asset_info"`
	Account   *auth.Account   `json:"account"`
}

// AlterAssetParam be valid where has the amount or thr asset info to be altered.
func (t *AlterAssetParam) Valid() error {
	if t == nil {
		return ErrNilPointer
	}
	if err := AssetIdValid(t.AssetId); err != nil {
		return err
	}
	errInfo := AlterAssetInfoValid(t.AssetInfo)
	errPrice := PriceInvalid(t.Price)
	errAmount := AmountInvalid(t.Amount)
	if errInfo != nil && errAmount != nil && errPrice != nil {
		return ErrAlterAssetInvalid
	}
	if err := AccountValid(t.Account); err != nil {
		return err
	}
	return nil
}

////////// Publish Asset ////////////
type PublishAssetParam struct {
	AssetId    int64         `json:"asset_id"`
	Account    *auth.Account `json:"account"`
	IsEvidence int           `json:"is_evidence,omitempty"`
}

func (t *PublishAssetParam) Valid() error {
	if t == nil {
		return ErrNilPointer
	}
	if err := AssetIdValid(t.AssetId); err != nil {
		return err
	}
	if err := AccountValid(t.Account); err != nil {
		return err
	}
	if err := EvidenceValid(t.IsEvidence); err != nil {
		return err
	}
	return nil
}

////////// Query Asset //////////
type QueryAssetParam struct {
	AssetId int64 `json:"asset_id"`
}

func (t *QueryAssetParam) Valid() error {
	if t == nil {
		return ErrNilPointer
	}
	if err := AssetIdValid(t.AssetId); err != nil {
		return err
	}
	return nil
}

type QueryAssetResp struct {
	BaseResp
	Meta *QueryAssetMeta `json:"meta"`
}

type QueryAssetMeta struct {
	AssetId    int64      `json:"asset_id"`
	GroupId    int64      `json:"group_id"`
	AssetCate  int        `json:"asset_cate"`
	Title      string     `json:"title"`
	Thumb      []ThumbMap `json:"thumb"`
	ShortDesc  string     `json:"short_desc"`
	LongDesc   string     `json:"long_desc"`
	ImgDesc    []string   `json:"img_desc"`
	AssetUrl   []string   `json:"asset_url"`
	AssetExt   string     `json:"asset_ext"`
	Price      int64      `json:"price"`
	Amount     int        `json:"amount"`
	Status     int        `json:"status"`
	CreateAddr string     `json:"create_addr"`
	Ctime      int64      `json:"ctime"`
	Mtime      int64      `json:"mtime"`
	TxId       string     `json:"tx_id"`
}

////////// Grant Asset /////////////
type GrantAssetParam struct {
	AssetId  int64         `json:"asset_id"`
	ShardId  int64         `json:"shard_id"`
	Price    int64         `json:"price,omitempty"`
	Account  *auth.Account `json:"account"`
	Addr     string        `json:"addr"`
	ToAddr   string        `json:"to_addr"`
	ToUserId int64         `json:"to_userid,omitempty"`
}

func (p *GrantAssetParam) Valid() error {
	if p == nil {
		return ErrNilPointer
	}
	if err := AssetIdValid(p.AssetId); err != nil {
		return err
	}
	if err := PriceInvalid(p.Price); err != nil {
		return err
	}
	if err := AccountValid(p.Account); err != nil {
		return err
	}
	if err := AddrValid(p.Addr); err != nil {
		return err
	}
	if err := AddrValid(p.ToAddr); err != nil {
		return err
	}
	return nil
}

type GrantAssetResp struct {
	BaseResp
	AssetId int64 `json:"asset_id"`
	ShardId int64 `json:"shard_id"`
}

////////// Query Shard ////////////
type QueryShardParam struct {
	AssetId int64 `json:"asset_id"`
	ShardId int64 `json:"shard_id"`
}

func (t *QueryShardParam) Valid() error {
	if t == nil {
		return ErrNilPointer
	}
	if err := AssetIdValid(t.AssetId); err != nil {
		return err
	}
	if err := ShardIdValid(t.ShardId); err != nil {
		return err
	}
	return nil
}

type QueryShardResp struct {
	BaseResp
	Meta *QueryShardMeta `json:"meta"`
}

type QueryShardMeta struct {
	AssetId   int64           `json:"asset_id"`
	ShardId   int64           `json:"shard_id"`
	Price     int64           `json:"price"`
	OwnerAddr string          `json:"owner_addr"`
	Status    int             `json:"status"`
	TxId      string          `json:"tx_id"`
	AssetInfo *ShardAssetInfo `json:"asset_info"`
	Ctime     int64           `json:"ctime"`
}

type ShardAssetInfo struct {
	Title      string     `json:"title"`
	AssetCate  int        `json:"asset_cate"`
	Thumb      []ThumbMap `json:"thumb"`
	ShortDesc  string     `json:"short_desc"`
	CreateAddr string     `json:"create_addr"`
	GroupId    int64      `json:"group_id"`
}

///////// List Shard By Address //////////
type ListShardsByAddrParam struct {
	Addr  string `json:"addr"`
	Page  int    `json:"page"`
	Limit int    `json:"limit"`
	// 可选
	AssetId int64 `json:"asset_id"`
}

func (t *ListShardsByAddrParam) Valid() error {
	if t == nil {
		return ErrNilPointer
	}
	if err := AddrValid(t.Addr); err != nil {
		return err
	}
	if err := IdValid(int64(t.Page)); err != nil {
		return err
	}
	return nil
}

type ListShardsByAddrResp struct {
	BaseResp
	List     []*QueryShardMeta `json:"list"`
	TotalCnt int               `json:"total_cnt"`
}

///////// List Assets By Address //////////
type ListAssetsByAddrParam struct {
	Addr   string `json:"addr"`
	Status int    `json:"status"`
	Page   int    `json:"page"`
	Limit  int    `json:"limit"`
}

func (t *ListAssetsByAddrParam) Valid() error {
	if t == nil {
		return ErrNilPointer
	}
	if err := AddrValid(t.Addr); err != nil {
		return err
	}
	if err := StatusValid(t.Status); err != nil {
		return err
	}
	return nil
}

type ListAssetsByAddrResp struct {
	BaseResp
	List     []*QueryAssetMeta `json:"list"`
	TotalCnt int               `json:"total_cnt"`
}

//////////// listdiffbyaddr /////////////////
type ListDiffByAddrParam struct {
	Addr string `json:"addr"`
	// 可选参数
	Limit  int    `json:"limit"`
	Cursor string `json:"cursor"`
	OpTyps string `json:"op_types"`
}

func (t *ListDiffByAddrParam) Valid() error {
	if t == nil || t.Addr == "" || t.Limit > 50 {
		return ErrParamInvalid
	}

	if t.OpTyps == "" {
		return nil
	}

	var arr []int
	err := json.Unmarshal([]byte(t.OpTyps), &arr)
	if err != nil {
		return ErrParamInvalid
	}

	return nil
}

type ListDiffByAddrNode struct {
	AssetId int64      `json:"asset_id"`
	ShardId int64      `json:"shard_id"`
	Operate int        `json:"operate"`
	Title   string     `json:"title"`
	Thumb   []ThumbMap `json:"thumb"`
	Ctime   int64      `json:"ctime"`
}

type ListDiffByAddrResp struct {
	BaseResp
	List    []*ListDiffByAddrNode `json:"list"`
	Cursor  string                `json:"cursor"`
	HasMore int                   `json:"has_more"`
}

///////// List Shards By Asset //////////
type ListShardsByAssetParam struct {
	AssetId int64  `json:"asset_id"`
	Cursor  string `json:"cursor"`
	Limit   int    `json:"limit"`
}

func (t *ListShardsByAssetParam) Valid() error {
	if t == nil {
		return ErrNilPointer
	}
	if err := AssetIdValid(t.AssetId); err != nil {
		return err
	}
	if err := AmountInvalid(t.Limit); err != nil {
		return err
	}
	return nil
}

type ListShardsByAssetResp struct {
	BaseResp
	List    []*QueryShardMeta `json:"list"`
	Cursor  string            `json:"cursor"`
	HasMore int               `json:"has_more"`
}

///////////// Get Evidence Info /////////////
type GetEvidenceInfoParam struct {
	AssetId int64 `json:"asset_id"`
}

func (t *GetEvidenceInfoParam) Valid() error {
	if t == nil {
		return ErrNilPointer
	}
	if err := AssetIdValid(t.AssetId); err != nil {
		return err
	}
	return nil
}

type GetEvidenceInfoResp struct {
	BaseResp
	CreateAddr string            `json:"create_addr"`
	TxId       string            `json:"tx_id"`
	FileHash   string            `json:"file_hash"`
	GhCertId   string            `json:"gh_cert_id"`
	AssetInfo  *HoraeAssetObject `json:"asset_info"`
	Ctime      int64             `json:"ctime"`
}

type HoraeAssetObject struct {
	AssetId   int64      `json:"asset_id"`
	AssetCate int        `json:"asset_cate"`
	Title     string     `json:"title"`
	Thumb     []ThumbMap `json:"thumb"`
	ShortDesc string     `json:"short_desc"`
}

////////// Transfer Asset //////////
type TransferAssetParam struct {
	AssetId  int64         `json:"asset_id"`
	ShardId  int64         `json:"shard_id"`
	Price    int64         `json:"price,omitempty"`
	Account  *auth.Account `json:"account"`
	Addr     string        `json:"addr"`
	ToAddr   string        `json:"to_addr"`
	ToUserId int64         `json:"to_userid,omitempty"`
}

func (p *TransferAssetParam) Valid() error {
	if p == nil {
		return ErrNilPointer
	}
	if err := AssetIdValid(p.AssetId); err != nil {
		return err
	}
	if err := ShardIdValid(p.ShardId); err != nil {
		return err
	}
	if err := PriceInvalid(p.Price); err != nil {
		return err
	}
	if err := AccountValid(p.Account); err != nil {
		return err
	}
	if err := AddrValid(p.Addr); err != nil {
		return err
	}
	if err := AddrValid(p.ToAddr); err != nil {
		return err
	}
	return nil
}

////////// Freeze Asset ////////////
type FreezeAssetParam struct {
	AssetId int64         `json:"asset_id"`
	Account *auth.Account `json:"account"`
}

func (t *FreezeAssetParam) Valid() error {
	if t == nil {
		return ErrNilPointer
	}
	if err := AssetIdValid(t.AssetId); err != nil {
		return err
	}
	if err := AccountValid(t.Account); err != nil {
		return err
	}
	return nil
}

////////// Consume Shard ////////////
type ConsumeShardParam struct {
	AssetId  int64         `json:"asset_id"`
	ShardId  int64         `json:"shard_id"`
	Nonce    int64         `json:"nonce"`
	UAddr    string        `json:"user_addr"`
	USign    string        `json:"user_sign"`
	UPKey    string        `json:"user_pkey"`
	CAccount *auth.Account `json:"create_account"`
}

func (t *ConsumeShardParam) Valid() error {
	if t == nil {
		return ErrNilPointer
	}
	if err := AssetIdValid(t.AssetId); err != nil {
		return err
	}
	if err := ShardIdValid(t.ShardId); err != nil {
		return err
	}
	if err := IdValid(t.Nonce); err != nil {
		return err
	}
	if err := AccountValid(t.CAccount); err != nil {
		return err
	}
	if err := AddrValid(t.UAddr); err != nil {
		return err
	}
	return nil
}

////////// Get History ////////////
type ListAssetHisParam struct {
	AssetId int64 `json:"asset_id"`
	ShardId int64 `json:"shard_id"`
	Page    int   `json:"page"`
	Limit   int   `json:"limit"`
}

func (t *ListAssetHisParam) Valid() error {
	if t == nil {
		return ErrNilPointer
	}
	if err := AssetIdValid(t.AssetId); err != nil {
		return err
	}
	if err := IdValid(int64(t.Page)); err != nil {
		return err
	}
	return nil
}

type HistoryMeta struct {
	AssetId int64  `json:"asset_id"`
	Type    int    `json:"type"`
	ShardId int64  `json:"shard_id"`
	Price   int64  `json:"price"`
	TxId    string `json:"tx_id"`
	From    string `json:"from"`
	To      string `json:"to"`
	Ctime   int64  `json:"ctime"`
}

type ListAssetHistoryResp struct {
	BaseResp
	List     []*HistoryMeta `json:"list"`
	TotalCnt int            `json:"total_cnt"`
	Cursor   string         `json:"cursor"`
	HasMore  int            `json:"has_more"`
}

////////// Scene ListShardByAddr ////////////
type SceneListShardByAddrParam struct {
	Addr   string `json:"addr"`
	Token  string `json:"token"`
	Cursor string `json:"cursor"`
	Limit  int    `json:"limit"`
}

func (t *SceneListShardByAddrParam) Valid() error {
	if t == nil {
		return ErrNilPointer
	}
	if err := AddrValid(t.Addr); err != nil {
		return err
	}
	if err := DescValid(t.Token); err != nil {
		return err
	}
	return nil
}

type SceneListShardByAddrResp struct {
	BaseResp
	List    []*SceneListMeta `json:"list"`
	Cursor  string           `json:"cursor"`
	HasMore int              `json:"has_more"`
}

type SceneListMeta struct {
	AssetId int64      `json:"asset_id"`
	ShardId int64      `json:"shard_id"`
	Title   string     `json:"title"`
	Thumb   []ThumbMap `json:"thumb"`
}

////////// Scene QueryShard ////////////
type SceneQueryShardParam struct {
	Addr    string `json:"addr"`
	Token   string `json:"token"`
	AssetId int64  `json:"asset_id"`
	ShardId int64  `json:"shard_id"`
}

func (t *SceneQueryShardParam) Valid() error {
	if t == nil {
		return ErrNilPointer
	}
	if err := AddrValid(t.Addr); err != nil {
		return err
	}
	if err := DescValid(t.Token); err != nil {
		return err
	}
	if err := IdValid(t.AssetId); err != nil {
		return err
	}
	if err := IdValid(t.ShardId); err != nil {
		return err
	}
	return nil
}

type SceneQueryShardResp struct {
	BaseResp
	Meta *SceneQueryMeta `json:"meta"`
}

type SceneQueryMeta struct {
	AssetId    int64      `json:"asset_id"`
	ShardId    int64      `json:"shard_id"`
	OwnerAddr  string     `json:"owner_addr"`
	Status     int        `json:"status"`
	TxId       string     `json:"tx_id"`
	Ctime      int64      `json:"ctime"`
	JumpLink   string     `json:"jump_link"`
	Price      int64      `json:"price"`
	Title      string     `json:"title"`
	Thumb      []ThumbMap `json:"thumb"`
	AssetUrl   []string   `json:"asset_url"`
	ImgDesc    []string   `json:"img_desc"`
	ShortDesc  string     `json:"short_desc"`
	CreateAddr string     `json:"create_addr"`
}

//////////// Scene listdiffbyaddr /////////////////
type SceneListDiffByAddrParam struct {
	Addr  string `json:"addr"`
	Token string `json:"token"`
	// 可选参数
	Limit  int    `json:"limit"`
	Cursor string `json:"cursor"`
	OpTyps string `json:"op_types"`
}

func (t *SceneListDiffByAddrParam) Valid() error {
	if t == nil || t.Addr == "" || t.Token == "" || t.Limit > 50 {
		return ErrParamInvalid
	}

	if t.OpTyps == "" {
		return nil
	}

	var arr []int
	err := json.Unmarshal([]byte(t.OpTyps), &arr)
	if err != nil {
		return ErrParamInvalid
	}

	return nil
}

//////////// Scene hasassetbyaddr ///////////////
type SceneHasAssetByAddrParam struct {
	Addr     string `json:"addr"`
	Token    string `json:"token"`
	AssetIds string `json:"asset_ids"`
}

func (t *SceneHasAssetByAddrParam) Valid() error {
	if t == nil {
		return ErrNilPointer
	}
	if err := AddrValid(t.Addr); err != nil {
		return err
	}
	if err := DescValid(t.Token); err != nil {
		return err
	}
	if t.AssetIds == "" {
		return ErrAssetListInvalid
	}
	return nil
}

type SceneHasAssetByAddrResp struct {
	BaseResp
	Result map[string]int `json:"result"`
}

//////////// Scene listaddr /////////////////////
type AddrGroupToken struct {
	Addr    string `json:"addr"`
	GroupId int64  `json:"group_id"`
	Token   string `json:"token"`
}
type SceneListAddrResp struct {
	BaseResp
	List []*AddrGroupToken `json:"list"`
}

//////////// Bdbox register ////////////////////
type BdBoxRegisterParam struct {
	OpenId string `json:"open_id"`
	AppKey string `json:"app_key"`
}

func (t *BdBoxRegisterParam) Valid() error {
	if t == nil {
		return ErrNilPointer
	}
	if err := OpenIdValid(t.OpenId); err != nil {
		return err
	}
	if err := AppKeyValid(t.AppKey); err != nil {
		return err
	}
	return nil
}

type BdBoxRegisterResp struct {
	BaseResp
	Address  string `json:"addr"`
	Mnemonic string `json:"mnemonic"`
	IsNew    int    `json:"is_new"`
}

///////////// Bdbox bind ////////////////////////
type BdBoxBindParam struct {
	OpenId   string `json:"open_id"`
	AppKey   string `json:"app_key"`
	Mnemonic string `json:"mnemonic"`
}

func (t *BdBoxBindParam) Valid() error {
	if t == nil {
		return ErrNilPointer
	}
	if err := OpenIdValid(t.OpenId); err != nil {
		return err
	}
	if err := AppKeyValid(t.AppKey); err != nil {
		return err
	}
	if err := MnemonicValid(t.Mnemonic); err != nil {
		return err
	}
	return nil
}

///////////// Bind by union id /////////////////
type BindByUnionIdParam struct {
	UnionId  string `json:"union_id"`
	Mnemonic string `json:"mnemonic"`
}

func (t *BindByUnionIdParam) Valid() error {
	if t == nil {
		return ErrNilPointer
	}
	if err := UnionIdValid(t.UnionId); err != nil {
		return err
	}
	if err := MnemonicValid(t.Mnemonic); err != nil {
		return err
	}
	return nil
}

//////////// Get addr by union id /////////////
type GetAddrByUnionIdResp struct {
	BaseResp
	Address string `json:"address"`
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
	//	用户标识，使用百度收银台时，该标识必须为baiduUID
	Uid int64 `json:"uid"`
	// object，其余支付参数的序列化值，标记剩余的非通用参数
	Details    string `json:"creator_details"`
	AppId      int64  `json:"app_id"`
	ActId      int64  `json:"act_id"`
	AssetId    int64  `json:"asset_id"`
	BuyerAddr  string `json:"buyer_addr"`
	SellerAddr string `json:"seller_addr"`
	ClientType int    `json:"client_type"`
	Chan       int64  `json:"chan"`
	Scene      int64  `json:"scene"`
	SignedAuth string `json:"signed_auth"`
}

func (p *HubCreateOrderParam) Valid() error {
	if p.SellerAddr == "" {
		return fmt.Errorf("seller_addr empty")
	}

	// 使用百度收银台请务必携带uid
	if _, ok := BaiduCashierCode[p.Code]; ok {
		if p.Uid <= 0 {
			return fmt.Errorf("baidu cashiers need passport id")
		}
		// 使用百度收银台H5支付必须提供加密串
		if p.Code == CodeBaiduH5 && p.SignedAuth == "" {
			return fmt.Errorf("invalid auth")
		}
	}

	if p.AssetId <= 0 {
		return fmt.Errorf("asset_id invalid, must be a positive integer")
	}
	if p.AppId <= 0 {
		return fmt.Errorf("app_id invalid, must be a positive integer")
	}
	return nil
}

type HubCreateResp struct {
	BaseResp
	Data struct {
		Code      int    `json:"code"`
		OrderType int    `json:"order_type"`
		Details   string `json:"details"`
		CTime     int64  `json:"ctime"`
	} `json:"data"`
}

type HubConfirmH5OrderParam struct {
	Code       int    `json:"code"`
	OrderType  int    `json:"order_type"`
	AppId      int64  `json:"app_id"`
	Oid        int64  `json:"oid"`
	ClientType int    `json:"client_type"`
	SignAuth   string `json:"signed_auth"`
}

func (p *HubConfirmH5OrderParam) Valid() error {
	if p.Oid < 0 {
		return fmt.Errorf("oid empty")
	}
	if p.ClientType < 0 {
		return fmt.Errorf("client_type empty")
	}
	if p.AppId < 0 {
		return fmt.Errorf("app_id empty")
	}
	return nil
}

type HubOrderDetailParam struct {
	AppId int64 `json:"app_id"`
	Oid   int64 `json:"oid"`
}

func (p *HubOrderDetailParam) Valid() error {
	if p.Oid < 0 || p.AppId < 0 {
		return fmt.Errorf("value empty")
	}
	return nil
}

type HubOrderDetailResp struct {
	BaseResp
	Data struct {
		HubOrderDetail
	} `json:"data"`
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
	AppId       int64  `json:"app_id"`
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
	if p.AppId < 0 {
		return fmt.Errorf("app_id invalid")
	}
	if p.Oid < 0 {
		return fmt.Errorf("oid invalid")
	}
	return nil
}

type HubListOrderParam struct {
	AppId  int64  `json:"app_id"`
	Addr   string `json:"address"`
	Status int    `json:"status"`
	Cursor string `json:"cursor"`
	Limit  int    `json:"limit"`
}

func (p *HubListOrderParam) Valid() error {
	if p.AppId < 0 {
		return fmt.Errorf("app_id invalid")
	}
	if p.Addr == "" {
		return fmt.Errorf("addr invalid")
	}
	if p.Status < 0 {
		return fmt.Errorf("addr invalid")
	}
	if p.Limit < 0 {
		return fmt.Errorf("cursor limit invalid")
	}

	return nil
}

type HubListOrderResp struct {
	BaseResp
	Data struct {
		List    []HubOrderDetail `json:"list"`
		Cursor  string           `json:"cursor"`
		HasMore int              `json:"has_more"`
	} `json:"data"`
}

type HubKnockOrderParam struct {
	Code      int   `json:"code"`
	OrderType int   `json:"order_type"`
	AppId     int64 `json:"app_id"`
	Oid       int64 `json:"oid"`
}

func (p *HubKnockOrderParam) Valid() error {
	if p.AppId <= 0 {
		return fmt.Errorf("app_id invalid")
	}
	if p.Oid <= 0 {
		return fmt.Errorf("oid invalid")
	}
	return nil
}

type HubKnockOrderResp struct {
	BaseResp
	Data struct {
		Code    int                 `json:"code"`
		Status  int                 `json:"status"`
		PayInfo *OwnedMallOrderItem `json:"pay_info"`
		AssetId int64               `json:"asset_id"`
		ShardId int64               `json:"shard_id"`
	} `json:"data"`
}

type OwnedMallOrderItem struct {
	DealId          string `json:"dealId"`
	AppKey          string `json:"appKey"`
	TotalAmount     string `json:"totalAmount"`
	TpOrderId       string `json:"tpOrderId"`
	NotifyUrl       string `json:"notifyUrl"`
	RsaSign         string `json:"rsaSign"`
	SignFieldsRange string `json:"signFieldsRange"`
	CTime           int64  `json:"ctime"`
}

type H5OrderItem struct {
	TpOrderId    int64  `json:"oid"`
	OrderInfoUrl string `json:"order_url"`
	H5PayInfo    string `json:"pay_info"`
	TotalAmount  string `json:"total_amount"`
	CTime        int64  `json:"ctime"`
}
