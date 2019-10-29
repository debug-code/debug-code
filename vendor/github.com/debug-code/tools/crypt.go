package tools

import (
	"encoding/base64"
	"encoding/xml"
	"errors"
	"fmt"
	"git.xrewin.com/go/beego"
	"sort"
	"strconv"
	"strings"
	"wechat-go/models"
)

//获得加密后的密文
func GetCipherText(re models.TextInfo, encodingAESKey, appId, noce, token, timeStamp string) (string, error) {
	msg, err := xml.Marshal(re)
	if err != nil {
		return "", err
	}

	ress := string(msg)
	ress = strings.Replace(ress, "&lt;", "<", -1)
	ress = strings.Replace(ress, "&gt;", ">", -1)

	encv, err := EncryptMsgs(encodingAESKey, appId, ress, noce)
	if err != nil {
		return "", err
	}

	sig := Signature(token, timeStamp, noce, encv)

	enc := models.Encrypted{}
	enc.TimeStamp = " <![CDATA[" + timeStamp + "]]> " //int(time.Now().Unix())
	enc.Encrypt = " <![CDATA[" + encv + "]]> "
	enc.MsgSignature = " <![CDATA[" + sig + "]]> "
	enc.Nonce = " <![CDATA[" + noce + "]]> "
	ret, err := xml.Marshal(enc)
	if err != nil {
		return "", err
	}

	res := string(ret)
	res = strings.Replace(res, "&lt;", "<", -1)
	res = strings.Replace(res, "&gt;", ">", -1)

	return res, nil
}

func GetDecryptMsgs(xmlData, encodingAESKey, appId string) (models.InputInfo, error) {
	beforeDec := models.InputInfo{}
	err := xml.Unmarshal([]byte(xmlData), &beforeDec)
	if err != nil {
		beego.Error(err)
		return models.InputInfo{}, err
	}

	result, err := DecryptMsgs(encodingAESKey, appId, beforeDec.Encrypt)

	afterDec := models.InputInfo{}
	err = xml.Unmarshal([]byte(result), &afterDec)
	if err != nil {
		beego.Error(err)
		return models.InputInfo{}, err
	}
	return afterDec, nil
}

////
func DecryptMsgs(aesKey string, appId string, xmlData string) (string, error) {

	decodeBytes, err := base64.StdEncoding.DecodeString(aesKey + "=")
	if err != nil {
		beego.Error(err)
		return "", err
	}
	encrypt, err := base64.StdEncoding.DecodeString(xmlData)
	if err != nil {
		beego.Error(err)
		return "", err
	}
	content, err := AesDecrypt([]byte(encrypt), decodeBytes)

	content = content[16:]
	lenXml := SockNonhl(content[:4])

	if string(content[lenXml+4:]) != appId {
		return "", errors.New("信息有误")
	}
	re := string(content[4 : lenXml+4])

	return re, nil
}

//wechat 加密
func EncryptMsgs(aesKey string, appId string,
	xmlStr string, nonce string) (string, error) {
	decodeBytes, err := base64.StdEncoding.DecodeString(aesKey + "=")
	if err != nil {
		return "", err
	}

	//fmt.Println("decodeBytes", string(decodeBytes))
	//fmt.Println("decodeBytes", decodeBytes)
	//nonce := GetRandomStr(16, true, true, true)
	//nonce := "xT2NDivDK7qpK4wM"
	//timestamp := "" //int(time.Now().Unix())
	//明文长度，int转byte
	//lenStr := BytesToHexString(SockHtonl(len(xmlStr)),true)
	lenStr := string(SockHtonl(len(xmlStr)))

	text := nonce + lenStr + xmlStr + appId
	//fmt.Println("len", len(xmlStr))

	//加密
	cryptograph, err := AesEncrypt(text, decodeBytes)
	if err != nil {
		return "", err
	}
	//fmt.Println("cryptograph")
	//fmt.Println(cryptograph)

	//生成安全签名
	//var msg []string
	//msg = append(msg, token)
	////msg = append(msg, strconv.Itoa(timestamp))
	//msg = append(msg, nonce)
	//msg = append(msg, cryptograph)
	//sort.Strings(msg)
	//str := ArrToString(msg)
	//signature, err := SHA1(str)
	//if err != nil {
	//	return "", err
	//}

	//fmt.Println("加密时生成的签名：", signature)
	//var encrypt models.Encrypted
	//encrypt.MsgSignature = signature
	//encrypt.Nonce = nonce
	//encrypt.TimeStamp = timestamp
	//encrypt.Encrypt = cryptograph

	//data, err := xml.MarshalIndent(&encrypt, "", "\t")
	//if err != nil {
	//	//fmt.Println(err)
	//	return "", nil
	//}

	return cryptograph, nil
}

//加密
func EnCodeXml(aesKey string, appId string,
	xmlStr string, nonce string) (string, error) {
	decodeBytes, err := base64.StdEncoding.DecodeString(aesKey + "=")
	if err != nil {
		fmt.Println(err)
	}
	//     fmt.Println("a")
	//明文长度，int转byte
	//lenStr := BytesToHexString(SockHtonl(len(xmlStr)),true)
	lenStr := string(SockHtonl(len(xmlStr)))

	text := nonce + lenStr + xmlStr + appId
	//fmt.Println("len", len(xmlStr))

	//加密
	cryptograph, err := AesEncrypt(text, decodeBytes)
	if err != nil {
		return "", err
	}
	//fmt.Println("cryptograph")
	//fmt.Println(cryptograph)

	//生成安全签名
	//var msg []string
	//msg = append(msg, token)
	//msg = append(msg, strconv.Itoa(timestamp))
	//msg = append(msg, nonce)
	//msg = append(msg, cryptograph)
	//sort.Strings(msg)
	//str := ArrToString(msg)
	//signature, err := SHA1(str)
	//if err != nil {
	//	return "", err
	//}

	//fmt.Println("加密时生成的签名：", signature)
	var encrypt models.Encrypted
	//encrypt.MsgSignature = signature
	encrypt.Nonce = nonce
	//encrypt.Timestamp = timestamp
	encrypt.Encrypt = cryptograph
	fmt.Println("加密时生成的签名：", encrypt)

	data, err := xml.MarshalIndent(&encrypt, "", "\t")
	if err != nil {
		//fmt.Println(err)
		return "", nil
	}
	//text = string(data)
	//fmt.Println("test")
	//fmt.Println(string(data))
	//fmt.Println()

	return string(data), nil
}

//解密
func DeCodeXml(postMsg string, token string, aesKey string,
	appId string, nonce string, timestamp int) (string, error) {

	//var decrypt models.Encrypt

	//解码xml
	//xml.Unmarshal([]byte(postMsg), &decrypt)
	//fmt.Println("decrypt:		", postMsg)
	//fmt.Println("decrypt.ToUserName:		", decrypt.ToUserName)
	//fmt.Println("decrypt.Encrypt:		", decrypt.Encrypt)

	//AES解密
	decodeBytes, err := base64.StdEncoding.DecodeString(aesKey + "=")
	if err != nil {
		fmt.Println(err)
	}
	encrypt, err := base64.StdEncoding.DecodeString(postMsg)
	if err != nil {
		fmt.Println(err)
	}
	content, err := AesDecrypt([]byte(encrypt), decodeBytes)

	//fmt.Println("解密后的明文")
	//fmt.Println("miwen:	", postMsg)
	fmt.Println("明文:		", string(content))

	//去除16随机字符
	content = content[16:]
	//fmt.Println("去除16个随机字符")
	//fmt.Println(string(content))
	//
	//fmt.Println("bytes content:	", content)
	//
	//fmt.Println(string(content[:4]))

	//byteS, err := HexStringToBytes(string(content[:4]), true)
	//if err != nil{
	//	return "", err
	//}
	lenXml := SockNonhl(content[:4])
	//fmt.Println(lenXml)
	//fmt.Println("appid:	", string(content[lenXml + 4:]))
	//fmt.Println("appId:	", appId)
	if string(content[lenXml+4:]) != appId {
		fmt.Println("cuole")
	}
	re := string(content[4 : lenXml+4])

	//验签
	var msg []string
	msg = append(msg, token)
	msg = append(msg, strconv.Itoa(timestamp))
	msg = append(msg, nonce)
	msg = append(msg, postMsg)
	sort.Strings(msg)
	str := ArrToString(msg)
	signature, err := SHA1(str)
	fmt.Println("解密时生成的签名：", signature)
	//fmt.Println(signature)
	//fmt.Println(decrypt.Signature)
	if err != nil {
		return "", err
	}
	//if signature != decrypt.Signature{
	//	return "验证失败", nil
	//}

	return re, nil

}
