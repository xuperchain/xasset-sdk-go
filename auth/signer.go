package auth

import (
	"fmt"
	"net/http"
	"sort"
	"strconv"
	"strings"
	"time"

	"github.com/baidubce/bce-sdk-go/util"
)

var (
	BCE_AUTH_VERSION        = "bce-auth-v1"
	SIGN_JOINER             = "\n"
	SIGN_HEADER_JOINER      = ";"
	DEFAULT_EXPIRE_SECONDS  = 1800
	DEFAULT_HEADERS_TO_SIGN = map[string]struct{}{
		strings.ToLower("Host"):           {},
		strings.ToLower("Content-Length"): {},
		strings.ToLower("Content-Type"):   {},
		strings.ToLower("Content-Md5"):    {},
	}
)

// SignOptions defines the data structure used by Signer
type SignOptions struct {
	HeadersToSign map[string]struct{}
	Timestamp     int64
	ExpireSeconds int
}

type Credentials struct {
	AppId           int64  // app id
	AccessKeyId     string // access key id to the service
	SecretAccessKey string // secret access key to the service
}

func (t *Credentials) String() string {
	return fmt.Sprintf("AppId:%d AccessKeyId:%s SecretAccessKey:%s",
		t.AppId, t.AccessKeyId, t.SecretAccessKey)
}

func (t *SignOptions) String() string {
	return fmt.Sprintf("HeadersToSign=%s Timestamp=%d ExpireSeconds=%d",
		t.HeadersToSign, t.Timestamp, t.ExpireSeconds)
}

// Sign - 签名
//
// PARAMS:
//     - req: *http.Request for this sign
//     - cred: *BceCredentials to access the serice
//     - opt: *SignOptions for this sign algorithm
func Sign(req *http.Request, cred *Credentials, opt *SignOptions) (string, error) {
	// Check param
	if req == nil || cred == nil || opt == nil {
		return "", fmt.Errorf("param error")
	}

	// Prepare parameters
	accessKeyId := cred.AccessKeyId
	secretAccessKey := cred.SecretAccessKey
	signDate := util.FormatISO8601Date(util.NowUTCSeconds())

	// Modify the sign time if it is not the default value but specified by client
	if opt.Timestamp != 0 {
		signDate = util.FormatISO8601Date(opt.Timestamp)
	}
	if opt.HeadersToSign == nil {
		opt.HeadersToSign = map[string]struct{}{"host": struct{}{}}
	}
	if opt.ExpireSeconds < 1 {
		opt.ExpireSeconds = DEFAULT_EXPIRE_SECONDS
	}

	// Prepare the canonical request components
	signKeyInfo := fmt.Sprintf("%s/%s/%s/%d", BCE_AUTH_VERSION, accessKeyId, signDate, opt.ExpireSeconds)
	signKey := util.HmacSha256Hex(secretAccessKey, signKeyInfo)

	// Generate signed head and signature
	signedHeaders, signature := getSignature(req, opt, signKey)

	// Generate auth string and add to the reqeust header
	authStr := signKeyInfo + "/" + signedHeaders + "/" + signature

	return authStr, nil
}

// CheckSign - 校验签名
func CheckSign(req *http.Request, cred *Credentials) error {
	if req == nil || cred == nil {
		return fmt.Errorf("param set error")
	}

	// 1.检验header中的Authorization格式
	author := req.Header.Get("Authorization")
	authStrs := strings.Split(author, "/")
	if len(authStrs) != 6 {
		return fmt.Errorf("author format error.auth:%s", author)
	}
	if authStrs[0] != BCE_AUTH_VERSION {
		return fmt.Errorf("author format error.auth:%s", author)
	}

	// 2. 校验签名是否过期
	expirationPeriodInSeconds, err := strconv.ParseInt(authStrs[3], 10, 32)
	if err != nil {
		return fmt.Errorf("author sign expiration set error.auth:%s", author)
	}
	if expirationPeriodInSeconds < 0 || expirationPeriodInSeconds > 3600 {
		return fmt.Errorf("author sign expiration set error.auth:%s", author)
	}
	timestamp, err := util.ParseISO8601Date(authStrs[2])
	if err != nil {
		return fmt.Errorf("author sign expiration set error.auth:%s", author)
	}
	if timestamp.Unix()+expirationPeriodInSeconds < time.Now().Unix() {
		return fmt.Errorf("author sign expiration.auth:%s", author)
	}

	// 3.校验签名头域
	HeadersToSign := map[string]struct{}{}
	for _, h := range strings.Split(authStrs[4], ";") {
		HeadersToSign[h] = struct{}{}
	}
	if _, ok := HeadersToSign["host"]; !ok {
		return fmt.Errorf("author sign headers unset host.auth:%s", author)
	}
	opt := &SignOptions{
		HeadersToSign: HeadersToSign,
	}

	// 4.校验签名摘要
	signingKey := util.HmacSha256Hex(cred.SecretAccessKey, strings.Join(authStrs[0:4], "/"))
	_, signature := getSignature(req, opt, signingKey)
	if signature != authStrs[5] {
		return fmt.Errorf("check signature failed.sign:%s-%s", signature, authStrs[5])
	}

	return nil
}

func getCanonicalURIPath(path string) string {
	if len(path) == 0 {
		return "/"
	}
	canonical_path := path
	if strings.HasPrefix(path, "/") {
		canonical_path = path[1:]
	}
	canonical_path = util.UriEncode(canonical_path, false)
	return "/" + canonical_path
}

func getCanonicalQueryString(params map[string]string) string {
	if len(params) == 0 {
		return ""
	}

	result := make([]string, 0, len(params))
	for k, v := range params {
		if strings.ToLower(k) == strings.ToLower("Authorization") {
			continue
		}
		item := ""
		if len(v) == 0 {
			item = fmt.Sprintf("%s=", util.UriEncode(k, true))
		} else {
			item = fmt.Sprintf("%s=%s", util.UriEncode(k, true), util.UriEncode(v, true))
		}
		result = append(result, item)
	}
	sort.Strings(result)
	return strings.Join(result, "&")
}

func getCanonicalHeaders(headers map[string]string, headersToSign map[string]struct{}) (string, []string) {
	canonicalHeaders := make([]string, 0, len(headers))
	signHeaders := make([]string, 0, len(headersToSign))
	for k, v := range headers {
		headKey := strings.ToLower(k)
		if headKey == strings.ToLower("Authorization") {
			continue
		}
		_, headExists := headersToSign[headKey]
		if headExists || (headKey == "host") {
			headVal := strings.TrimSpace(v)
			encoded := util.UriEncode(headKey, true) + ":" + util.UriEncode(headVal, true)
			canonicalHeaders = append(canonicalHeaders, encoded)
			signHeaders = append(signHeaders, headKey)
		}
	}
	sort.Strings(canonicalHeaders)
	sort.Strings(signHeaders)
	return strings.Join(canonicalHeaders, SIGN_JOINER), signHeaders
}

func getSignature(req *http.Request, opt *SignOptions, signingKey string) (string, string) {
	canonicalUri := getCanonicalURIPath(req.URL.Path)
	queryParams := make(map[string]string, 0)
	for _, query := range strings.Split(req.URL.RawQuery, "&") {
		if len(query) == 0 {
			break
		}
		param := strings.Split(query, "=")
		if len(param) == 1 {
			queryParams[param[0]] = ""
		} else {
			queryParams[param[0]] = param[1]
		}
	}
	canonicalQueryString := getCanonicalQueryString(queryParams)

	headerParams := make(map[string]string, 0)
	for k, v := range req.Header {
		headerParams[k] = strings.Join(v, ";")
	}
	canonicalHeaders, signedHeadersArr := getCanonicalHeaders(headerParams, opt.HeadersToSign)

	// Generate signed headers string
	signedHeaders := ""
	if len(signedHeadersArr) > 0 {
		sort.Strings(signedHeadersArr)
		signedHeaders = strings.Join(signedHeadersArr, SIGN_HEADER_JOINER)
	}

	// Generate signature
	canonicalParts := []string{req.Method, canonicalUri, canonicalQueryString, canonicalHeaders}
	canonicalReq := strings.Join(canonicalParts, SIGN_JOINER)
	return signedHeaders, util.HmacSha256Hex(signingKey, canonicalReq)
}
