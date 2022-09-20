package auth

import (
	"crypto/ecdsa"
	"encoding/base64"
	"encoding/hex"
	"fmt"
	"bytes"
	"crypto/aes"
	"crypto/cipher"

	"github.com/xuperchain/crypto/core/account"
	"github.com/xuperchain/crypto/core/hash"
	"github.com/xuperchain/crypto/core/sign"
)

// xasset签名完整方法
// @jsPrivtKey: json格式的private key
// @oriMsg: 签名的原始数据
func XassetSignECDSA(jsPrivtKey string, oriMsg []byte) (string, error) {
	// 1.对消息统一做SHA256
	msg := HashBySha256(oriMsg)

	// 2.使用ECC私钥来签名
	k, err := GetEcdsaPriKeyByJsStr(jsPrivtKey)
	if err != nil {
		return "", err
	}
	signature, err := SignECDSA(k, msg)
	if err != nil {
		return "", err
	}

	// 3.对签名转化为16进制字符串显示
	return EncodeSign(signature), nil
}

// xasset校验签名完整方法
// @jsPubKey: json格式public key
// @oriMsg: 签名的原始数据
func XassetVerifyECDSA(jsPubKey, signature string, oriMsg []byte) (bool, error) {
	// 1.对消息统一做SHA256
	msg := HashBySha256(oriMsg)

	// 2.使用ECC共钥验签
	k, err := GetEcdsaPubKeyByJsStr(jsPubKey)
	if err != nil {
		return false, err
	}
	sigBytes, err := DecodeSign(signature)
	if err != nil {
		return false, err
	}

	return VerifyECDSA(k, sigBytes, msg)
}

// 使用SHA256做单次哈希运算
func HashBySha256(data []byte) []byte {
	return hash.HashUsingSha256(data)
}

// 将byte转换为16进制字符串显示
func EncodeSign(src []byte) string {
	return hex.EncodeToString(src)
}

// 将16进制字符串显示转换为byte
func DecodeSign(signature string) ([]byte, error) {
	return hex.DecodeString(signature)
}

// 使用ECC私钥来签名
func SignECDSA(k *ecdsa.PrivateKey, msg []byte) ([]byte, error) {
	if k == nil {
		return nil, fmt.Errorf("sign private key unset")
	}

	signature, err := sign.SignECDSA(k, msg)
	if err != nil {
		return nil, fmt.Errorf("ecdsa sign failed.err:%v", err)
	}

	return signature, nil
}

// 使用ECC公钥来验证签名，验证统一签名的新签名函数
func VerifyECDSA(k *ecdsa.PublicKey, signature, msg []byte) (bool, error) {
	if k == nil {
		return false, fmt.Errorf("sign public key unset")
	}

	result, err := sign.VerifyECDSA(k, signature, msg)
	if err != nil {
		return false, fmt.Errorf("verify ecdsa failed.err:%v", err)
	}

	return result, nil
}

// 从json格式私钥内容字符串产生ECC私钥
func GetEcdsaPriKeyByJsStr(keyStr string) (*ecdsa.PrivateKey, error) {
	return account.GetEcdsaPrivateKeyFromJson([]byte(keyStr))
}

// 从json格式公钥内容字符串产生ECC公钥
func GetEcdsaPubKeyByJsStr(keyStr string) (*ecdsa.PublicKey, error) {
	return account.GetEcdsaPublicKeyFromJson([]byte(keyStr))
}

// 使用单个公钥来生成钱包地址
func GetAddrByPubKey(key *ecdsa.PublicKey) (string, error) {
	if key == nil {
		return "", fmt.Errorf("public key unset")
	}

	return account.GetAddressFromPublicKey(key)
}

// 验证钱包地址和公钥是否匹配
func VerifyAddrByPubKey(address string, pub *ecdsa.PublicKey) (bool, uint8) {
	if pub == nil {
		return false, 0
	}

	return account.VerifyAddressUsingPublicKey(address, pub)
}

func PKCS7Padding(ciphertext []byte, blockSize int) []byte {
	padding := blockSize - len(ciphertext)%blockSize
	padtext := bytes.Repeat([]byte{byte(padding)}, padding)
	return append(ciphertext, padtext...)
}

func PKCS7UnPadding(origData []byte) []byte {
	length := len(origData)
	unpadding := int(origData[length-1])
	return origData[:(length - unpadding)]
}

// 加密
func AesEncrypt(origData, key []byte) ([]byte, error) {
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}
	blockSize := block.BlockSize()

	// 填充补齐
	origData = PKCS7Padding(origData, blockSize)

	// 加密
	blockMode := cipher.NewCBCEncrypter(block, key[:blockSize])
	crypted := make([]byte, len(origData))
	blockMode.CryptBlocks(crypted, origData)

	return crypted, nil
}

// 解密
func AesDecrypt(crypted, key []byte) ([]byte, error) {
	// 获取block size
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}
	blockSize := block.BlockSize()
	if len(crypted)%blockSize != 0 {
		return nil, fmt.Errorf("crypted not full blocks")
	}

	// 解密
	blockMode := cipher.NewCBCDecrypter(block, key[:blockSize])
	origData := make([]byte, len(crypted))
	blockMode.CryptBlocks(origData, crypted)

	// 校验解密后数据合法性，防止panic
	unpadding := int(origData[len(origData)-1])
	if len(origData) < blockSize || len(origData)%blockSize != 0 || unpadding > blockSize {
		return nil, fmt.Errorf("origin data error.block_size:%d length:%d unpadding:%d",
			blockSize, len(origData), unpadding)
	}

	// 去掉填充字符
	origData = PKCS7UnPadding(origData)
	return origData, nil
}

func Base64Encode(content []byte) string {
	return base64.StdEncoding.EncodeToString(content)
}

func Base64Decode(content string) ([]byte, error) {
	return base64.StdEncoding.DecodeString(content)
}

func Base64UrlEncode(content []byte) string {
	return base64.RawURLEncoding.EncodeToString(content)
}

func Base64UrlDecode(content string) ([]byte, error) {
	return base64.RawURLEncoding.DecodeString(content)
}