# XASSET GO SDK

## 概述

本项目提供Xasset Go语言版的开发者工具包（SDK），开发者可以基于该SDK使用Go语言接入到Xasset平台。

## 使用说明

- 1.从平台申请获得到API准入AK/SK。注意AK/SK是准入凭证，不要泄露，不要下发或配置在客户端使用。
- 2.在开发中导入（import）对应包，开发业务程序。建议使用go mod管理包依赖。
- 3.接入联调环境联调测试，测试通过后更换到线上环境，完成接入。

### 运行环境

GO SDK可以在go1.3及以上环境下运行。

### 导入SDK包

```

// 按需导入对应包使用，使用go mod管理（go mod tidy）
import (
    "github.com/xuperchain/xasset-sdk-go/client/xasset"
)

```

### 配置说明

```

type XassetCliConfig struct {
    // 请求接入Host
    Endpoint           string
    // 请求UA
    UserAgent          string
    // 准入授权信息，注意保密
    Credentials        *auth.Credentials
    // 准入签名Header
    SignOption         *auth.SignOptions
    // Http请求链接超时
    ConnectTimeoutMs   int
    // Http请求读写超时
    ReadWriteTimeoutMs int
}

// 使用示例
cfg := config.NewXassetCliConf()
cfg.Endpoint = "http://127.0.0.1:8360"
cfg.SetCredentials(appId, ak, sk)

```

### 使用示例

```

// 导入包
import (
    "github.com/xuperchain/xasset-sdk-go/client/xasset"        
    "github.com/xuperchain/xasset-sdk-go/common/config"
)

// 配置SDK
cfg := config.NewXassetCliConf()
cfg.Endpoint = "http://127.0.0.1:8360"
cfg.SetCredentials(appId, ak, sk)

// 调用SDK方法，可以参考单元测试
handle, _ := xasset.NewAssetOperCli(cfg, &Logger{})
handle.CreateAsset()

```

### sk加解密
```
//导入包
import (
    github.com/xuperchain/xasset-sdk-go/utils
)

// 使用sk对union_id加密
signedUnionId, _ := utils.AesEncode(unionId, sk)

// 部分应用场景调用需要先加密后传输，请参考client/xasset下的单元测试
```