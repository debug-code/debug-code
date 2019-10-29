package tools

import (
	"git.xrewin.com/go/beego"
	"github.com/zheng-ji/gophone"
	"regexp"
	"strings"
)

func CheckPhone(phone string) bool {
	reg := `^1([38][0-9]|14[579]|5[^4]|16[6]|7[1-35-8]|9[189])\d{8}$`
	rgx := regexp.MustCompile(reg)
	return rgx.MatchString(phone)

}

func CheckPhoneIsp(phone string) (res string) {

	ispMap := map[string]string{
		"130": "联通",
		"131": "联通",
		"132": "联通",
		"145": "联通",
		"146": "联通",
		"155": "联通",
		"156": "联通",
		"166": "联通",
		"175": "联通",
		"176": "联通",
		"185": "联通",
		"186": "联通",
		"171": "联通",
		"133": "电信",
		"149": "电信",
		"173": "电信",
		"174": "电信",
		"153": "电信",
		"177": "电信",
		"180": "电信",
		"181": "电信",
		"189": "电信",
		"199": "电信",
		"134": "移动",
		"135": "移动",
		"136": "移动",
		"137": "移动",
		"138": "移动",
		"139": "移动",
		"147": "移动",
		"148": "移动",
		"150": "移动",
		"151": "移动",
		"152": "移动",
		"157": "移动",
		"158": "移动",
		"159": "移动",
		"172": "移动",
		"178": "移动",
		"182": "移动",
		"183": "移动",
		"184": "移动",
		"187": "移动",
		"188": "移动",
		"198": "移动",
		"170": "其他",
	}

	res = "其他"
	if v, s := ispMap[phone[:3]]; s {
		res = v
	}
	return
}
func CheckPhoneArea(phones string) string {
	phone := strings.TrimSpace(phones)
	pr, err := gophone.Find(phone)
	if err != nil {
		beego.Error(err)
		return ""
	}

	return pr.Province + "/" + pr.City
}

func CheckEmail(email string) bool {

	reg := `^([a-zA-Z0-9_-])+@([a-zA-Z0-9_-])+(.[a-zA-Z0-9_-])+`
	rgx := regexp.MustCompile(reg)
	return rgx.MatchString(email)

}
func CheckEmailIsp(emails string) string {
	email := strings.TrimSpace(emails)
	res := "其他"
	if strings.HasSuffix(email, "163.com") ||
		strings.HasSuffix(email, "126.com") ||
		strings.HasSuffix(email, "188.com") ||
		strings.HasSuffix(email, "yeah.net") {
		res = "网易"
	}
	if strings.HasSuffix(email, "qq.com") {
		res = "腾讯"
	}
	if strings.HasSuffix(email, "sina.com") {
		res = "新浪"
	}
	return res
}
