package crypto

import (
	"crypto/cipher"
	"crypto/des"
	"crypto/hmac"
	"crypto/md5"
	"crypto/sha1"
	"crypto/sha256"
	"crypto/sha512"
	"encoding/base64"
	"encoding/hex"
	"fmt"
	"io"
)

// EncodeSHA1HMACBase64 : encrypt to SHA1HMAC input key, data String. Output to String in Base64 format
func EncodeSHA1HMACBase64(key string, data ...string) string {
	return EncodeBASE64(ComputeSHA1HMAC(key, data...))
}

// EncodeSHA1HMAC : encrypt to SHA1HMAC input key, data String. Output to String in Base16/Hex format
func EncodeSHA1HMAC(key string, data ...string) string {
	return fmt.Sprintf("%x", ComputeSHA1HMAC(key, data...))
}

//ComputeSHA1HMAC : encrypt to SHA1HMAC input key, data String. Output to String
func ComputeSHA1HMAC(key string, data ...string) string {
	h := hmac.New(sha1.New, []byte(key))
	for _, v := range data {
		io.WriteString(h, v)
	}
	return string(h.Sum(nil))
}

func EncodeSHA256HMACBase64(key string, data ...string) string {
	return EncodeBASE64(ComputeSHA256HMAC(key, data...))
}

func EncodeSHA256HMAC(key string, data ...string) string {
	return fmt.Sprintf("%x", ComputeSHA256HMAC(key, data...))
}

func ComputeSHA256HMAC(key string, data ...string) string {
	h := hmac.New(sha256.New, []byte(key))
	for _, v := range data {
		io.WriteString(h, v)
	}
	return string(h.Sum(nil))
}

func EncodeSHA512HMACBase64(key string, data ...string) string {
	return EncodeBASE64(ComputeSHA512HMAC(key, data...))
}

func EncodeSHA512HMAC(key string, data ...string) string {
	return fmt.Sprintf("%x", ComputeSHA512HMAC(key, data...))
}

func ComputeSHA512HMAC(key string, data ...string) string {
	h := hmac.New(sha512.New, []byte(key))
	for _, v := range data {
		io.WriteString(h, v)
	}
	return string(h.Sum(nil))
}

//EncodeMD5 : encrypt to MD5 input string, output to string
func EncodeMD5(text string) string {
	h := md5.New()
	h.Write([]byte(text))
	return hex.EncodeToString(h.Sum(nil))
}

func EncodeMD5Base64(text string) string {
	h := md5.New()
	h.Write([]byte(text))
	// return EncodeBASE64(hex.EncodeToString(h.Sum(nil)))
	return base64.StdEncoding.EncodeToString((h.Sum(nil)))
}

//EncodeBASE64 : Encrypt to Base64. Input string, output string
func EncodeBASE64(text string) string {
	return base64.StdEncoding.EncodeToString([]byte(text))
}

//DecodeBASE64 : Decrypt Base64. Input string, output string
func DecodeBASE64(text string) (string, error) {
	byt, err := base64.StdEncoding.DecodeString(text)
	return string(byt), err
}

//EncodeBASE64URL : Encrypt to Base64URL. Input string, output text
func EncodeBASE64URL(text string) string {
	return base64.URLEncoding.EncodeToString([]byte(text))
}

//EncodeDES : Encrypt to DES. input string, output chiper
func EncodeDES(text string) (cipher.Block, error) {
	desKey, _ := hex.DecodeString(text)
	cipher, err := des.NewTripleDESCipher(desKey)
	return cipher, err
}

//EncodeSHA256: Encrypt to SHA256. input string, output text
func EncodeSHA256(text string) string {
	h := sha256.Sum256([]byte(text))
	return fmt.Sprintf("%x", h)
}

//EncodeSHA512 Encrypt to SHA512. input string, output text
func EncodeSHA512(text string) string {
	h := sha512.Sum512([]byte(text))
	return fmt.Sprintf("%x", h)
}
