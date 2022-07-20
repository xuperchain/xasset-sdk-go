package httpcli

import (
	"crypto/tls"
	"errors"
	"io/ioutil"
	"net"
	"net/http"
	"net/url"
	"strings"
	"time"
)

const (
	OptDisableFollowLocation = "DisableFollowLocation"
	OptDisableCompression    = "DisableCompression"
	OptTlsSipVerify          = "TlsSkipVerify"
)

type HttpResponse struct {
	StatusCode int
	Header     http.Header
	Body       []byte
}

var DisableRedirectError = errors.New("Don't redirect!")

func noRedirect(req *http.Request, via []*http.Request) error {
	return DisableRedirectError
}

func SendRequest(req *http.Request, ConnTimeoutMs, RWTimeoutMs int,
	opt map[string]string) (HttpResponse, error) {

	disableCompression := false
	checkRedirect := noRedirect
	if v, ok := opt[OptDisableFollowLocation]; !ok || v != "1" {
		checkRedirect = nil
	}
	if v, ok := opt[OptDisableCompression]; ok && v == "1" {
		disableCompression = true
	}

	transport := &http.Transport{
		Dial: func(netw, addr string) (net.Conn, error) {
			conn, err := net.DialTimeout(netw, addr, time.Duration(ConnTimeoutMs)*time.Millisecond)
			if err != nil {
				return nil, err
			}
			conn.SetDeadline(time.Now().Add(time.Duration(RWTimeoutMs) * time.Millisecond))
			return conn, nil
		},
		MaxIdleConnsPerHost: -1,
		DisableCompression:  disableCompression,
		DisableKeepAlives:   true,
	}

	// tls is skip verify
	if v, ok := opt[OptTlsSipVerify]; ok && v == "1" {
		transport.TLSClientConfig = &tls.Config{
			CipherSuites: []uint16{
				tls.TLS_ECDHE_RSA_WITH_AES_128_GCM_SHA256,
			},
			PreferServerCipherSuites: true,
			InsecureSkipVerify:       true,
			MinVersion:               tls.VersionTLS12,
			MaxVersion:               tls.VersionTLS12,
		}
	}

	// do request
	var res HttpResponse
	client := &http.Client{
		Transport:     transport,
		CheckRedirect: checkRedirect,
	}
	response, err := client.Do(req)
	if response != nil {
		res.StatusCode = response.StatusCode
		res.Header = response.Header
	}
	if err != nil {
		if urlError, ok := err.(*url.Error); ok && urlError.Err == DisableRedirectError {
			// Discard body and response.Body.Close() in go src when status is 302, ft
			return res, nil
		}
		return res, err
	}

	defer response.Body.Close()
	res.Body, err = ioutil.ReadAll(response.Body)
	if err != nil {
		return res, err
	}
	return res, nil
}

func GenRequest(method, url string, header map[string]string, data string) (*http.Request, error) {
	var req *http.Request
	var err error
	if data != "" {
		req, err = http.NewRequest(method, url, strings.NewReader(data))
	} else {
		req, err = http.NewRequest(method, url, nil)
	}
	if err != nil {
		return nil, err
	}

	for k, v := range header {
		if strings.EqualFold(k, "host") {
			req.Host = v
		}
		req.Header.Set(k, v)
	}

	return req, nil
}

func IsHttps(url string) bool {
	if strings.HasPrefix(strings.ToLower(url), "https") {
		return true
	}

	return false
}
