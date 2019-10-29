package tools

import (
	"bytes"
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"crypto/sha1"
	"encoding/base64"
	"fmt"
	"git.xrewin.com/go/beego"
	"sort"
	"strings"
)

// EncryptMsg 加密报文
func EncryptMsg(msg []byte, aesKey []byte, appId string, noce string) (b64Enc string, err error) {
	// 拼接完整报文
	src := SpliceFullMsg(msg, appId)
	fmt.Println(src)
	// AES CBC 加密报文

	//dst, err := AESCBCEncrypt(src, aesKey, aesKey[:aes.BlockSize])
	dst, err := AESCBCEncrypt(src, aesKey, aesKey[:16])
	if err != nil {
		return "", err
	}

	return base64.StdEncoding.EncodeToString(dst), nil
}

// DecryptMsg 解密报文
func DecryptMsg(b64Enc string, aesKey []byte, appId string) (msg []byte, err error) {
	beego.Debug("b64Enc: %s", b64Enc)
	enc, err := base64.StdEncoding.DecodeString(b64Enc)
	if err != nil {
		return nil, err
	}

	// AES CBC 解密报文
	src, err := AESCBCDecrypt(enc, aesKey, aesKey[:aes.BlockSize])
	if err != nil {
		return nil, err
	}

	beego.Debug("full message: %s", src)

	_, _, msg, appId2 := ParseFullMsg(src)
	if appId2 != appId {
		return nil, fmt.Errorf("expected appId %s, but %s", appId, appId2)
	}

	return msg, nil
}

// SpliceFullMsg 拼接完整报文，
// AES加密的buf由16个字节的随机字符串、4个字节的msg_len(网络字节序)、msg和$AppId组成，
// 其中msg_len为msg的长度，$AppId为公众帐号的AppId
func SpliceFullMsg(msg []byte, appId string) (fullMsg []byte) {
	// 16个字节的随机字符串
	//randBytes := RandBytes(16)
	randBytes := []byte(GetRandomStr(16, true, true, true))
	// 4个字节的msg_len(网络字节序)
	msgLen := len(msg)
	lenBytes := []byte{
		byte(msgLen >> 24 & 0xFF),
		byte(msgLen >> 16 & 0xFF),
		byte(msgLen >> 8 & 0xFF),
		byte(msgLen & 0xFF),
	}

	return bytes.Join([][]byte{randBytes, lenBytes, msg, []byte(appId)}, nil)
}

// ParseFullMsg 从完整报文中解析出消息内容，
// AES加密的buf由16个字节的随机字符串、4个字节的msg_len(网络字节序)、msg和$AppId组成，
// 其中msg_len为msg的长度，$AppId为公众帐号的AppId
func ParseFullMsg(fullMsg []byte) (randBytes []byte, msgLen int, msg []byte, appId string) {
	randBytes = fullMsg[:16]

	msgLen = (int(fullMsg[16]) << 24) |
		(int(fullMsg[17]) << 16) |
		(int(fullMsg[18]) << 8) |
		int(fullMsg[19])
	// beego.Debug("msgLen=[% x]=(%d %d %d %d)=%d", fullMsg[16:20], (int(fullMsg[16]) << 24),
	// 	(int(fullMsg[17]) << 16), (int(fullMsg[18]) << 8), int(fullMsg[19]), msgLen)

	msg = fullMsg[20 : 20+msgLen]

	appId = string(fullMsg[20+msgLen:])

	return
}

// RandBytes 产生 size 个长度的随机字节
func RandBytes(size int) (r []byte) {
	r = make([]byte, size)
	_, err := rand.Read(r)
	if err != nil {
		// 忽略错误，不影响其他逻辑，仅仅打印日志
		//log.Warnf("rand read error: %s", err)
	}
	return r
}

func AESCBCEncrypt(src, key, iv []byte) (enc []byte, err error) {
	beego.Debug("src: %s", string(src))
	src = PKCS7Padding(src, len(key))

	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}

	mode := cipher.NewCBCEncrypter(block, iv)

	mode.CryptBlocks(src, src)
	enc = src

	beego.Debug("enc: % x", enc)
	return enc, nil
}

// AESCBCDecrypt 采用 CBC 模式的 AES 解密
func AESCBCDecrypt(enc, key, iv []byte) (src []byte, err error) {
	beego.Debug("enc: % x", enc)
	if len(enc) < len(key) {
		return nil, fmt.Errorf("the length of encrypted message too short: %d", len(enc))
	}
	if len(enc)&(len(key)-1) != 0 { // or len(enc)%len(key) != 0
		return nil, fmt.Errorf("encrypted message is not a multiple of the key size(%d), the length is %d", len(key), len(enc))
	}

	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}

	mode := cipher.NewCBCDecrypter(block, iv)

	mode.CryptBlocks(enc, enc)
	src = PKCS7UnPadding(enc)

	beego.Debug("src: %s", src)
	return src, nil
}

// PKCS7Padding PKCS#7填充，Buf需要被填充为K的整数倍，
// 在buf的尾部填充(K-N%K)个字节，每个字节的内容是(K- N%K)
func PKCS7Padding(src []byte, k int) (padded []byte) {
	padLen := k - len(src)%k
	padding := bytes.Repeat([]byte{byte(padLen)}, padLen)
	return append(src, padding...)
}

// PKCS7UnPadding 去掉PKCS#7填充，Buf需要被填充为K的整数倍，
// 在buf的尾部填充(K-N%K)个字节，每个字节的内容是(K- N%K)
func PKCS7UnPadding(src []byte) (padded []byte) {
	padLen := int(src[len(src)-1])
	return src[:len(src)-padLen]
}

// ValidateURL 验证 URL 以判断来源是否合法
func ValidateURL(token, timestamp, nonce, signature string) bool {
	tmpArr := []string{token, timestamp, nonce}
	sort.Strings(tmpArr)

	tmpStr := strings.Join(tmpArr, "")
	actual := fmt.Sprintf("%x", sha1.Sum([]byte(tmpStr)))

	beego.Debug("%s %s", tmpArr, actual)
	return actual == signature
}

// Signature 对加密的报文计算签名
func Signature(token, timestamp, nonce, encrypt string) string {
	tmpArr := []string{token, timestamp, nonce, encrypt}
	sort.Strings(tmpArr)

	tmpStr := strings.Join(tmpArr, "")
	actual := fmt.Sprintf("%x", sha1.Sum([]byte(tmpStr)))

	beego.Debug("%s %s", tmpArr, actual)
	return actual
}

// CheckSignature 验证加密的报文的签名
func CheckSignature(token, timestamp, nonce, encrypt, sign string) bool {
	return Signature(token, timestamp, nonce, encrypt) == sign
}
