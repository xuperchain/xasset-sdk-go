package base

import (
	"fmt"

	"github.com/xuperchain/xasset-sdk-go/auth"
	"github.com/xuperchain/xasset-sdk-go/common/config"
)

/*
const (
	// 需要修改为正确的配置
	TestAppId    = 0
	TestAK       = "xxx"
	TestSK       = "xxx"
	TestEndpoint = "http://120.48.16.137:8360"
	Openid       = "xxx"
	Appkey       = "xxx"
	UnionId      = "xxx"
)*/

const (
	// remember to del
	TestAppId    = 300100
	TestAK       = "032b9af2f1b776d69c8a55031f2ae68e"
	TestSK       = "2cb51374f71d8d274b370685d36d2280"
	TestEndpoint = "http://10.117.131.18:8360"
)
const (
	OpenId  = "0qQzshUoZjh35XnUxXFas_C_Z2"
	AppKey  = "3lsCeIYo00pErR7MBWzfZcR1nZpE42dq"
	UnionId = "uCeh1dJxQdc6LUkhphVFb95dDYsXX3k"
)

var TestAccount, _ = auth.NewXchainEcdsaAccount(auth.MnemStrgthStrong, auth.MnemLangCN)

var TestTransAccount, _ = auth.NewXchainEcdsaAccount(auth.MnemStrgthStrong, auth.MnemLangCN)

func TestGetXassetConfig() *config.XassetCliConfig {
	cfg := config.NewXassetCliConf()
	cfg.SetCredentials(int64(TestAppId), TestAK, TestSK)
	cfg.Endpoint = TestEndpoint

	return cfg
}

// mock logger
type TestLogger struct {
}

func (t *TestLogger) Error(msg string, ctx ...interface{}) {
	fmt.Println(t.logFmt("Error", msg, ctx...))
}

func (t *TestLogger) Warn(msg string, ctx ...interface{}) {
	fmt.Println(t.logFmt("Warn", msg, ctx...))
}

func (t *TestLogger) Info(msg string, ctx ...interface{}) {
	fmt.Println(t.logFmt("Info", msg, ctx...))
}

func (t *TestLogger) Trace(msg string, ctx ...interface{}) {
	fmt.Println(t.logFmt("Trace", msg, ctx...))
}

func (t *TestLogger) Debug(msg string, ctx ...interface{}) {
	fmt.Println(t.logFmt("Debug", msg, ctx...))
}

func (t *TestLogger) logFmt(lvl, msg string, ctx ...interface{}) string {
	msg = fmt.Sprintf("[lvl:%s] "+msg, lvl)
	return fmt.Sprintf(msg, ctx...)
}
