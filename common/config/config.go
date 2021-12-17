package config

import (
	"fmt"

	"github.com/xuperchain/xasset-sdk-go/auth"
)

const (
	EndpointDefault       = "http://120.48.16.137:8360"
	UserAgentDefault      = "xasset-sdk-go"
	ConnectTimeoutMsDef   = 1000
	ReadWriteTimeoutMsDef = 3000
)

type XassetCliConfig struct {
	Endpoint           string
	UserAgent          string
	Credentials        *auth.Credentials
	SignOption         *auth.SignOptions
	ConnectTimeoutMs   int
	ReadWriteTimeoutMs int
}

func NewXassetCliConf() *XassetCliConfig {
	return &XassetCliConfig{
		Endpoint:  EndpointDefault,
		UserAgent: UserAgentDefault,
		SignOption: &auth.SignOptions{
			HeadersToSign: auth.DEFAULT_HEADERS_TO_SIGN,
			Timestamp:     0,
			ExpireSeconds: auth.DEFAULT_EXPIRE_SECONDS,
		},
		ConnectTimeoutMs:   ConnectTimeoutMsDef,
		ReadWriteTimeoutMs: ReadWriteTimeoutMsDef,
	}
}

func (t *XassetCliConfig) SetCredentials(appId int64, ak, sk string) {
	t.Credentials = &auth.Credentials{
		AppId:           appId,
		AccessKeyId:     ak,
		SecretAccessKey: sk,
	}
}

func (t *XassetCliConfig) String() string {
	return fmt.Sprintf("[Endpoint:%s] [UserAgent:%s] [Credentials:%v] [SignOption:%v] "+
		"[ConnectTimeoutMs:%dms] [ReadWriteTimeoutMs:%dms]", t.Endpoint, t.UserAgent,
		t.Credentials, t.SignOption, t.ConnectTimeoutMs, t.ReadWriteTimeoutMs)
}

func (t *XassetCliConfig) IsVaild() bool {
	if t.Endpoint == "" || t.Credentials == nil || t.SignOption == nil {
		return false
	}

	if t.UserAgent == "" {
		t.UserAgent = UserAgentDefault
	}
	if t.ConnectTimeoutMs == 0 {
		t.ConnectTimeoutMs = ConnectTimeoutMsDef
	}
	if t.ReadWriteTimeoutMs == 0 {
		t.ReadWriteTimeoutMs = ReadWriteTimeoutMsDef
	}

	return true
}
