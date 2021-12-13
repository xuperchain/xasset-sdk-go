package auth

import (
	"io"
	"net/http"
	"strings"
	"testing"
	"time"
)

var cred = &Credentials{
	AccessKeyId:     "xxx",
	SecretAccessKey: "xxx",
}

func TestSign(t *testing.T) {
	var reqInfos = []struct {
		method string
		url    string
		body   io.Reader
		auth   string
	}{
		{
			method: "POST",
			url:    "http://www.baidu.com",
			body:   nil,
		},
		{
			method: "POST",
			url:    "http://www.baidu.com/",
			body:   nil,
		},
		{
			method: "POST",
			url:    "http://www.baidu.com/?toke=",
			body:   nil,
		},
		{
			method: "GET",
			url:    "http://www.baidu.com/?toke=123",
			body:   nil,
		},
		{
			method: "POST",
			url:    "http://www.baidu.com/?toke=123&name=林&age=",
			body:   strings.NewReader("{\"addr\":\"0xsss\"}"),
		},
	}

	opt := &SignOptions{
		HeadersToSign: nil,
		Timestamp:     time.Now().Unix(),
		ExpireSeconds: 1800,
	}
	for _, arg := range reqInfos {
		req, _ := http.NewRequest(arg.method, arg.url, arg.body)
		req.Header.Add("Content-Type", "application/json")
		req.Header.Add("X-Bce-Request-Id", "15304c2e-381b-45ab-bc1e-488b183f4293")
		req.Header.Add("Host", "www.baidu.com")
		sign, _ := Sign(req, cred, opt)
		t.Logf(sign)
	}
}

func TestCheckSign(t *testing.T) {
	var reqInfos = []struct {
		method string
		url    string
		body   io.Reader
	}{
		{
			method: "POST",
			url:    "http://www.baidu.com",
			body:   nil,
		},
		{
			method: "POST",
			url:    "http://www.baidu.com/?toke=",
			body:   nil,
		},
		{
			method: "GET",
			url:    "http://www.baidu.com/?toke=123",
			body:   nil,
		},
		{
			method: "POST",
			url:    "http://www.baidu.com/?toke=123&name=林&age=",
			body:   strings.NewReader("{\"addr\":\"0xsss\"}"),
		},
	}
	opt := &SignOptions{
		HeadersToSign: map[string]struct{}{"host": struct{}{}},
		Timestamp:     time.Now().Unix(),
		ExpireSeconds: 1800,
	}
	for idx, arg := range reqInfos {
		req, _ := http.NewRequest(arg.method, arg.url, arg.body)
		req.Header.Add("Content-Type", "application/json")
		req.Header.Add("X-Bce-Request-Id", "15304c2e-381b-45ab-bc1e-488b183f4293")
		req.Header.Add("Host", "www.baidu.com")

		sign, _ := Sign(req, cred, opt)
		req.Header.Add("Authorization", sign)

		err := CheckSign(req, cred)
		if err != nil {
			t.Logf("[index:%d] check sign failed.err:%v", idx, err)
			t.FailNow()
		}
	}
}
