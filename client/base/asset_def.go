package base

import (
	"github.com/xuperchain/xasset-sdk-go/auth"
)

const (
	AssetApiCreate           = "/xasset/horae/v1/create"
	AssetApiAlter            = "/xasset/horae/v1/alter"
	AssetApiPublish          = "/xasset/horae/v1/publish"
	AssetApiQueryAsset       = "/xasset/horae/v1/query"
	AssetApiGrant            = "/xasset/horae/v1/grant"
	AssetApiTransfer         = "/xasset/damocles/v1/transfer"
	AssetApiQueryShard       = "/xasset/horae/v1/querysds"
	AssetApiListShardsByAddr = "/xasset/horae/v1/listsdsbyaddr"
	AssetApiGetEvidenceInfo  = "/xasset/horae/v1/getevidenceinfo"
	FileApiGetStoken         = "/xasset/file/v1/getstoken"
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
	AccessInfo interface{} `json:"accessInfo"`
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
	Amount    int              `json:"amount"`
	AssetInfo *CreateAssetInfo `json:"asset_info"`
	Account   *auth.Account    `json:"account"`
	UserId    int64            `json:"user_id,omitempty"`
}

func (t *CreateAssetParam) Valid() error {
	if t == nil {
		return ErrNilPointer
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
	Amount    int             `json:"amount,omitempty"`
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
	errAmount := AmountInvalid(t.Amount)
	if errInfo != nil && errAmount != nil {
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
}

///////// List Shard By Address //////////
type ListCursorResp struct {
	BaseResp
	List    interface{} `json:"list"`
	HasMore int         `json:"has_more"`
	Cursor  string      `json:"cursor"`
}

type ListShardsByAddrParam struct {
	Addr  string `json:"addr"`
	Page  int    `json:"page"`
	Limit int    `json:"limit"`
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

type ListPageResp struct {
	BaseResp
	List     interface{} `json:"list"`
	TotalCnt int         `json:"total_cnt"`
}

///////////// Get Evidence Info /////////////
type GetEvidenceInfoParam struct {
	AssetId int64 `json:"asset_id"`
	ShardId int64 `json:"shard_id"`
}

func (t *GetEvidenceInfoParam) Valid() error {
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

type GetEvidenceInfoResp struct {
	BaseResp
	CreateAddr string            `json:"create_addr"`
	TxId       string            `json:"tx_id"`
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
