package tools

import (
	"encoding/json"
	"errors"
	"fmt"
	"git.xrewin.com/go/beego"
	"git.xrewin.com/go/beego/httplib"
	"math/rand"
	"net/http"
	"strconv"
	"time"
	"wechat-go/models"
)

//获取全局ak
func GetAccess_token(appid string, appsecret string) (string, error) {
	url_head := "https://api.weixin.qq.com/cgi-bin/token?"
	grant_type := "client_credential"
	AppSecret := appsecret
	AppID := appid

	url := url_head +
		"grant_type=" + grant_type +
		"&appid=" + AppID +
		"&secret=" + AppSecret

	req := httplib.Get(url)
	str, err := req.String()
	fmt.Println("Access_token", str)
	if err != nil {
		fmt.Println(err)
	}
	maps := make(map[string]interface{})
	err = json.Unmarshal([]byte(str), &maps)
	if err != nil {

		fmt.Println("err", err)

	}

	access_token := maps["access_token"].(string)

	return access_token, nil
}

func GetTicket(access_token string) (string, error) {

	url_head := "https://api.weixin.qq.com/cgi-bin/ticket/getticket?"
	getType := "jsapi"

	url := url_head +
		"&access_token=" + access_token +
		"&type=" + getType

	req := httplib.Get(url)
	str, err := req.String()
	fmt.Println("Ticket", str)
	if err != nil {
		fmt.Println(err)
	}
	maps := models.Ticket{}
	err = json.Unmarshal([]byte(str), &maps)
	if err != nil {
		fmt.Println("err", err)
	}
	if maps.Errcode != 0 {
		return "no", errors.New(maps.Errmsg)
	}

	ticket := maps.Ticket

	return ticket, nil
}

func GetJsConfig(appid string, appsecret string, url string) (models.JsConfig, error) {
	access_token, err := GetAccess_token(appid, appsecret)
	if err != nil {
		fmt.Println(err)
		return models.JsConfig{}, err
	}

	ticket, err := GetTicket(access_token)
	if err != nil {
		fmt.Println(err)
		return models.JsConfig{}, err
	}

	noncestr := GetRandomStr(16, true, true, true)
	timestamp := time.Now().Unix()
	timestampInt := strconv.FormatInt(time.Now().Unix(), 10)

	str := "jsapi_ticket=" + ticket + "&noncestr=" + noncestr + "&timestamp=" + timestampInt + "&url=" + url

	signature, err := SHA1(str)
	if err != nil {
		return models.JsConfig{}, err
	}

	jsData := models.JsConfig{}
	jsData.Appid = appid
	jsData.NonceStr = noncestr
	jsData.Signature = signature
	jsData.Timestamp = timestamp

	fmt.Println(ticket)
	fmt.Println(jsData)
	return jsData, nil
}

//get redirect_url
func GetRUrl(appid string, redirect_uri string) (string, error) {
	RUrl := ""

	urlHead := "https://open.weixin.qq.com/connect/oauth2/authorize?"
	state := "123"

	RUrl = urlHead +
		"appid=" + appid +
		"&redirect_uri=" + redirect_uri +
		"&response_type=code&scope=snsapi_userinfo&state=" + state +
		"#wechat_redirect"

	return RUrl, nil
}

func GetAccess_tokenWeb(appid string, appsecret string, code string, grant_type string) (models.Acess_tokenWeb, error) {
	url_head := "https://api.weixin.qq.com/sns/oauth2/access_token?"
	//"appid=APPID&secret=SECRET&code=CODE&grant_type=authorization_code"

	url := url_head +
		"appid=" + appid +
		"&secret=" + appsecret +
		"&code=" + code +
		"&grant_type=" + grant_type

	req := httplib.Get(url)
	str, err := req.String()
	if err != nil {
		fmt.Println(err)
	}
	maps := models.Acess_tokenWeb{}
	err = json.Unmarshal([]byte(str), &maps)
	if err != nil {

		fmt.Println("err", err)

	}
	return maps, nil
}

func GetRefresh_tokenWeb(appid string, grant_type string, refresh_token string) (models.Acess_tokenWeb, error) {
	url_head := "https://api.weixin.qq.com/sns/oauth2/refresh_token?"
	//appid=APPID&grant_type=refresh_token&refresh_token=REFRESH_TOKEN

	url := url_head +
		"appid=" + appid +
		"&grant_type=" + grant_type +
		"&refresh_token=" + refresh_token

	req := httplib.Get(url)
	str, err := req.String()
	fmt.Println(" refresh Access_token", str)
	if err != nil {
		fmt.Println(err)
	}
	maps := models.Acess_tokenWeb{}
	err = json.Unmarshal([]byte(str), &maps)
	if err != nil {

		fmt.Println("err", err)

	}
	return maps, nil
}

func GetUserInfoWeb(openid string, access_token string, lang string) (models.UserInfo, error) {
	url_head := "https://api.weixin.qq.com/sns/userinfo?"
	// https://api.weixin.qq.com/sns/userinfo?access_token=ACCESS_TOKEN&openid=OPENID&lang=zh_CN

	url := url_head +
		"access_token=" + access_token +
		"&openid=" + openid +
		"&lang=" + lang
	req := httplib.Get(url)
	str, err := req.String()
	if err != nil {
		fmt.Println(err)
	}
	maps := models.UserInfo{}
	err = json.Unmarshal([]byte(str), &maps)
	if err != nil {

		fmt.Println("err", err)

	}
	return maps, nil
}

func GetUserInfoWebMore(openid string, access_token string, lang string) (models.UserInfo, error) {
	url_head := "https://api.weixin.qq.com/cgi-bin/user/info?"
	//https://api.weixin.qq.com/cgi-bin/user/info?access_token=ACCESS_TOKEN&openid=OPENID&lang=zh_CN
	//https://api.weixin.qq.com/cgi-bin/user/info?access_token=ACCESS_TOKEN&openid=OPENID&lang=zh_CN
	url := url_head +
		"access_token=" + access_token +
		"&openid=" + openid +
		"&lang=" + lang
	req := httplib.Get(url)
	str, err := req.String()
	if err != nil {
		fmt.Println(err)
	}
	maps := models.UserInfo{}
	err = json.Unmarshal([]byte(str), &maps)
	if err != nil {

		fmt.Println("err", err)

	}
	return maps, nil
}

//获取用户列表
func GetMemberAll(access_token string) (models.Acess_tokenWeb, error) {
	url_head := "https://api.weixin.qq.com/cgi-bin/user/get?"
	//https://api.weixin.qq.com/cgi-bin/user/get?access_token=ACCESS_TOKEN&next_openid=NEXT_OPENID
	url := url_head +
		"access_token=" + access_token
	//"&openid=" + openid +
	//"&lang=" + lang

	req := httplib.Get(url)
	_, err := req.String()
	if err != nil {
		fmt.Println(err)
	}
	maps := models.Acess_tokenWeb{}
	//err = json.Unmarshal([]byte(str), &maps)
	//if err != nil {
	//
	//	fmt.Println("err",err)
	//
	//}
	return maps, nil
}

func GetMembersOpenid(next_openid string, access_token string) (string, error) {
	url_head := "https://api.weixin.qq.com/cgi-bin/user/get?"
	url := url_head +
		"access_token=" + access_token

	if next_openid != "" {
		url += "next_openid=" + next_openid
	}
	req := httplib.Get(url)
	str, err := req.String()

	if err != nil {
		//fmt.Println(err)
		beego.Error(err)

		return "", err
	}
	//fmt.Println("GetMemberOpenid :", str)

	return str, nil
}

func GetMembersInfo(jsonData string, access_token string) (string, error) {
	url_head := "https://api.weixin.qq.com/cgi-bin/user/info/batchget?"
	url := url_head +
		"access_token=" + access_token

	req := httplib.Post(url)
	req.Body(jsonData)

	str, err := req.String()

	if err != nil {
		//fmt.Println(err)
		return "", err
	}
	//fmt.Println("GetCardUserInfo :", str)
	return str, nil
}

func GetMemberByOpenid(openid string, access_token string) (string, error) {
	url_head := "https://api.weixin.qq.com/cgi-bin/user/info?"
	url := url_head +
		"access_token=" + access_token
	url = url +
		"&openid=" + openid

	req := httplib.Get(url)

	str, err := req.String()

	if err != nil {
		//fmt.Println(err)
		return "", err
	}
	//fmt.Println("GetCardUserInfo :", str)
	return str, nil
}

func GetUserCard(access_token string, openid string, card_id string) (models.OpenidCard, error) {
	url_head := "https://api.weixin.qq.com/card/user/getcardlist?"
	//https://api.weixin.qq.com/card/user/getcardlist?access_token=TOKEN
	url := url_head +
		"access_token=" + access_token
	//"&openid=" + openid +
	//"&lang=" + lang

	req := httplib.Post(url)
	body := `{"openid":"` + openid + `","card_id":"` + card_id + `"}`
	req.Body(body)

	str, err := req.String()
	//fmt.Println("GetUserCard :", str)
	if err != nil {
		fmt.Println(err)
	}
	maps := models.OpenidCard{}
	err = json.Unmarshal([]byte(str), &maps)
	if err != nil {

		fmt.Println("err", err)

	}
	return maps, nil
}

func GetCardUserInfo(access_token string, code string, card_id string) (models.CardUserInfo, error) {
	url_head := "https://api.weixin.qq.com/card/membercard/userinfo/get?"
	//https://api.weixin.qq.com/card/membercard/userinfo/get?access_token=TOKEN
	url := url_head +
		"access_token=" + access_token
	//"&openid=" + openid +
	//"&lang=" + lang

	req := httplib.Post(url)
	body := `{"code":"` + code + `","card_id":"` + card_id + `"}`
	req.Body(body)

	str, err := req.String()
	fmt.Println("GetCardUserInfo :", str)

	if err != nil {
		fmt.Println(err)
	}
	maps := models.CardUserInfo{}
	err = json.Unmarshal([]byte(str), &maps)
	if err != nil {

		fmt.Println("err", err)

	}
	return maps, nil
}

func Send(jsonData string, access_token string) (string, error) {
	//url_head := "https://api.weixin.qq.com/cgi-bin/user/info/batchget?"
	url_head := "https://api.weixin.qq.com/cgi-bin/message/mass/send?"
	url := url_head +
		"access_token=" + access_token

	req := httplib.Post(url)
	req.Body(jsonData)

	str, err := req.String()

	if err != nil {
		//fmt.Println(err)
		return "", err
	}
	//fmt.Println("GetCardUserInfo :", str)
	return str, nil

}

//卡卷
//创建会员卡
func CreateCard(jsonData string, access_token string) (string, error) {
	url_head := "https://api.weixin.qq.com/card/create?"
	url := url_head +
		"access_token=" + access_token

	req := httplib.Post(url)
	req.Body(jsonData)

	str, err := req.String()

	if err != nil {
		//fmt.Println(err)
		return "", err
	}
	//fmt.Println("GetCardUserInfo :", str)
	return str, nil
}

//删除会员卡
func DeleteVipCard(jsonData string, access_token string) (string, error) {
	//https://api.weixin.qq.com/card/delete?access_token=TOKEN
	url_head := "https://api.weixin.qq.com/card/delete?"
	url := url_head +
		"access_token=" + access_token

	req := httplib.Post(url)
	req.Body(jsonData)

	str, err := req.String()

	if err != nil {

		return "", err
	}

	return str, nil

}

//修改会员卡信息
func UpdateVipCard(jsonData string, access_token string) (string, error) {
	//https://api.weixin.qq.com/card/update?access_token=TOKEN
	url_head := "https://api.weixin.qq.com/card/update?"
	url := url_head +
		"access_token=" + access_token

	req := httplib.Post(url)
	req.Body(jsonData)

	str, err := req.String()

	if err != nil {
		return "", err
	}
	return str, nil
}

//修改会员卡信息
func UpdateStock(jsonData string, access_token string) (string, error) {
	//https://api.weixin.qq.com/card/update?access_token=TOKEN
	url_head := "https://api.weixin.qq.com/card/modifystock?"
	url := url_head +
		"access_token=" + access_token

	req := httplib.Post(url)
	req.Body(jsonData)

	str, err := req.String()

	if err != nil {
		return "", err
	}
	return str, nil
}

//查看卡卷详情
func GetVipCard(jsonData string, access_token string) {
	//https://api.weixin.qq.com/card/get?access_token=TOKEN
	url_head := "https://api.weixin.qq.com/card/get?"
	url := url_head +
		"access_token=" + access_token

	req := httplib.Post(url)
	req.Body(jsonData)

	str, err := req.String()

	if err != nil {
		fmt.Println(err)
	}
	fmt.Println("GetCardUserInfo :", str)
}

func SetCard(jsonData string, access_token string) (string, error) {
	//https://api.weixin.qq.com/card/membercard/activateuserform/set?access_token=TOKEN
	url_head := "https://api.weixin.qq.com/card/membercard/activateuserform/set?"
	url := url_head +
		"access_token=" + access_token

	req := httplib.Post(url)
	req.Body(jsonData)

	str, err := req.String()

	if err != nil {

		return "", err
	}

	return str, nil
}

func SetIndustry(jsonData string, access_token string) (string, error) {
	//https://api.weixin.qq.com/card/membercard/activateuserform/set?access_token=TOKEN
	url_head := "https://api.weixin.qq.com/cgi-bin/template/api_set_industry?"
	url := url_head +
		"access_token=" + access_token

	req := httplib.Post(url)
	req.Body(jsonData)

	str, err := req.String()

	if err != nil {

		return "", err
	}

	return str, nil
}

//添加优惠卷
//func CreateCoupon(jsonData string , access_token string)(string, error){
//
//	url_head := "https://api.weixin.qq.com/card/create?"
//	url := url_head +
//		"access_token=" + access_token
//
//	req := httplib.Post(url)
//	req.Body(jsonData)
//
//	str, err := req.String()
//
//	if err != nil {
//		//fmt.Println(err)
//		return "", err
//	}
//	//fmt.Println("GetCardUserInfo :", str)
//	return str, nil
//}

//素材

//图文消息中的图片，不占用素材库中的数量
func UploadImg(name string, filePath, aK string) (string, error) {
	url_head := "https://api.weixin.qq.com/cgi-bin/media/uploadimg?"

	url := url_head +
		"access_token=" + aK

	req, err := PostFile(url, nil, name, filePath)
	if err != nil {
		beego.Error(err)
	}
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		beego.Error(err)
	}
	defer resp.Body.Close()
	ret := make(map[string]interface{})
	if err := json.NewDecoder(resp.Body).Decode(&ret); err != nil {
		beego.Error(err)
	}
	str, _ := json.Marshal(ret)

	return string(str), err
}
func UploadSource(name string, params map[string]string, filePath, aK string, ty string) (string, error) {
	url_head := "https://api.weixin.qq.com/cgi-bin/material/add_material?"

	url := url_head +
		"access_token=" + aK
	url += "&type=" + ty

	req, err := PostFile(url, params, name, filePath)
	if err != nil {
		beego.Error(err)
	}
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		beego.Error(err)
	}
	defer resp.Body.Close()
	ret := make(map[string]interface{})
	if err := json.NewDecoder(resp.Body).Decode(&ret); err != nil {
		beego.Error(err)
	}
	str, _ := json.Marshal(ret)

	return string(str), err
}

func DeleteSource(jsonData string, access_token string) (string, error) {
	url_head := "https://api.weixin.qq.com/cgi-bin/material/del_material?"
	url := url_head +
		"access_token=" + access_token

	req := httplib.Post(url)
	req.Body(jsonData)

	str, err := req.String()

	if err != nil {
		fmt.Println(err)
		return "", err
	}
	//fmt.Println("getswlist :", str)
	return str, nil
}

//menu

func AddMenu(jsonData string, access_token string) (string, error) {
	url_head := "https://api.weixin.qq.com/cgi-bin/menu/create?"
	url := url_head +
		"access_token=" + access_token

	req := httplib.Post(url)
	req.Body(jsonData)

	str, err := req.String()

	if err != nil {
		fmt.Println(err)
		return "", err
	}
	//fmt.Println("getswlist :", str)
	return str, nil
}
func DeleteMenu(access_token string) (string, error) {
	url_head := "https://api.weixin.qq.com/cgi-bin/menu/delete?"
	url := url_head +
		"access_token=" + access_token

	req := httplib.Get(url)

	str, err := req.String()
	if err != nil {
		fmt.Println(err)
		return "", err
	}
	//fmt.Println("getswlist :", str)
	return str, nil
}

// 素材列表
func GetSWList(jsonData string, access_token string) (string, error) {
	//https://api.weixin.qq.com/cgi-bin/material/batchget_material?access_token=ACCESS_TOKEN
	url_head := "https://api.weixin.qq.com/cgi-bin/material/batchget_material?"
	url := url_head +
		"access_token=" + access_token

	req := httplib.Post(url)
	req.Body(jsonData)

	str, err := req.String()

	if err != nil {
		fmt.Println(err)
		return "", err
	}
	//fmt.Println("getswlist :", str)
	return str, nil
}
func GetSource(jsonData string, access_token string) (string, error) {
	//https://api.weixin.qq.com/cgi-bin/material/batchget_material?access_token=ACCESS_TOKEN
	url_head := "https://api.weixin.qq.com/cgi-bin/material/get_material?"
	url := url_head +
		"access_token=" + access_token

	req := httplib.Post(url)
	req.Body(jsonData)

	str, err := req.String()

	if err != nil {
		fmt.Println(err)
		return "", err
	}
	//fmt.Println("getswlist :", str)
	return str, nil
}

//数据分析

//图文分析数据接口
func GetFormData(formType string, jsonData string, access_token string) (string, error) {
	var url string

	switch formType {

	case "Get_Article_Summary":
		url_head := "https://api.weixin.qq.com/datacube/getarticlesummary?"
		url = url_head + "access_token=" + access_token
		break
	case "Get_Article_Total":
		url_head := "https://api.weixin.qq.com/datacube/getarticletotal?"
		url = url_head + "access_token=" + access_token
		break
	case "Get_Serread":
		url_head := "https://api.weixin.qq.com/datacube/getuserread?"
		url = url_head + "access_token=" + access_token
		break
	case "Get_Serread_Hour":
		url_head := "https://api.weixin.qq.com/datacube/getuserreadhour?"
		url = url_head + "access_token=" + access_token
		break
	case "Get_User_Share":
		url_head := "https://api.weixin.qq.com/datacube/getusershare?"
		url = url_head + "access_token=" + access_token
		break
	case "Get_User_Share_Hour":
		url_head := "https://api.weixin.qq.com/datacube/getusersharehour?"
		url = url_head + "access_token=" + access_token
		break
	default:
		err := errors.New("暂不支持改类型数据或类型错误")
		return "", err
	}

	req := httplib.Post(url)
	req.Body(jsonData)

	str, err := req.String()

	if err != nil {
		fmt.Println(err)
		return "", err
	}
	return str, nil
}

//get access_token
func GetAK() string {
	//get access_token
	AppSecret := beego.AppConfig.String("AppSecret")
	AppID := beego.AppConfig.String("AppID")

	access_tokenGloble, err := GetAccess_token(AppID, AppSecret) //全局
	if err != nil {
		beego.Error(err)
		return ""
	}
	beego.Debug("AK: ", access_tokenGloble)
	return access_tokenGloble

}

//h获取模板列表
func GetModelList(access_token string) (string, error) {
	//https://api.weixin.qq.com/cgi-bin/material/batchget_material?access_token=ACCESS_TOKEN
	url_head := "https://api.weixin.qq.com/cgi-bin/template/get_all_private_template?"
	url := url_head +
		"access_token=" + access_token

	req := httplib.Get(url)
	str, err := req.String()

	if err != nil {
		//fmt.Println(err)
		beego.Error(err)

		return "", err
	}
	//fmt.Println("getswlist :", str)
	return str, nil
}

//h获取article url
func GetArticleUrl(jsonData string, access_token string) (string, error) {
	//https://api.weixin.qq.com/cgi-bin/material/batchget_material?access_token=ACCESS_TOKEN
	url_head := "https://api.weixin.qq.com/cgi-bin/material/get_material?"
	url := url_head +
		"access_token=" + access_token

	req := httplib.Post(url)
	req.Body(jsonData)
	str, err := req.String()

	if err != nil {
		//fmt.Println(err)
		beego.Error(err)

		return "", err
	}
	//fmt.Println("getswlist :", str)
	return str, nil
}

// 发送模板消息

func SendModel(jsonData string, access_token string) (string, error) {
	url_head := "https://api.weixin.qq.com/cgi-bin/message/template/send?"
	url := url_head +
		"access_token=" + access_token

	req := httplib.Post(url)
	req.Body(jsonData)

	str, err := req.String()

	if err != nil {
		fmt.Println(err)
		return "", err
	}
	//fmt.Println("getswlist :", str)
	return str, nil
}

//news
func AddNews(jsonData string, access_token string) (string, error) {
	url_head := "https://api.weixin.qq.com/cgi-bin/material/add_news?"
	url := url_head +
		"access_token=" + access_token

	//jsonData = strings.Replace(jsonData, `\u003c`, "<", -1)
	//jsonData = strings.Replace(jsonData, `\u003e`, ">", -1)
	//strig, _ := strconv.Unquote(jsonData)
	//beego.Error("replace jsonData2", strig)
	//beego.Error("replace jsonData1", fmt.Sprintf(jsonData))
	req := httplib.Post(url)
	req.Body(jsonData)

	str, err := req.String()

	if err != nil {
		fmt.Println(err)
		return "", err
	}
	//fmt.Println("getswlist :", str)
	return str, nil
}
func UpdateNews(jsonData string, access_token string) (string, error) {
	url_head := "https://api.weixin.qq.com/cgi-bin/material/update_news?"
	url := url_head +
		"access_token=" + access_token

	//jsonData = strings.Replace(jsonData, `\u003c`, "<", -1)
	//jsonData = strings.Replace(jsonData, `\u003e`, ">", -1)
	req := httplib.Post(url)
	req.Body(jsonData)

	str, err := req.String()

	if err != nil {
		fmt.Println(err)
		return "", err
	}
	//fmt.Println("getswlist :", str)
	return str, nil
}

//delete

func GetCardIframe(jsonData string, access_token string) (string, error) {
	url_head := "https://api.weixin.qq.com/card/mpnews/gethtml?"
	url := url_head +
		"access_token=" + access_token

	req := httplib.Post(url)
	req.Body(jsonData)

	str, err := req.String()

	if err != nil {
		fmt.Println(err)
		return "", err
	}
	//fmt.Println("getswlist :", str)
	return str, nil
}

//客服消息
// 添加账号
func SetKFAccount(jsonData string, access_token string) (string, error) {
	url_head := "https://api.weixin.qq.com/customservice/kfaccount/add?"
	url := url_head +
		"access_token=" + access_token

	req := httplib.Post(url)
	req.Body(jsonData)

	str, err := req.String()

	if err != nil {
		fmt.Println(err)
		return "", err
	}
	//fmt.Println("getswlist :", str)
	return str, nil
}

func SendMsgByKF(jsonData string, access_token string) (string, error) {
	url_head := "https://api.weixin.qq.com/cgi-bin/message/custom/send?"
	url := url_head +
		"access_token=" + access_token

	req := httplib.Post(url)
	req.Body(jsonData)

	str, err := req.String()

	if err != nil {
		fmt.Println(err)
		return "", err
	}
	//fmt.Println("getswlist :", str)
	return str, nil
}

func Preview(jsonData string, access_token string) (string, error) {
	url_head := "https://api.weixin.qq.com/cgi-bin/message/mass/preview?"
	url := url_head +
		"access_token=" + access_token

	req := httplib.Post(url)
	req.Body(jsonData)

	str, err := req.String()

	if err != nil {
		fmt.Println(err)
		return "", err
	}
	//fmt.Println("getswlist :", str)
	return str, nil
}

func SendKfMsg(toUserOpenId, ak, replyStr string, isAll bool) {
	//beego.Debug("replyStr", replyStr)
	reply := []models.ReplyContent{}

	err := json.Unmarshal([]byte(replyStr), &reply)
	if err != nil {
		beego.Error(err)
	}

	for k, v := range reply {
		temp := UnicodeEmojiDecode(v.Content)
		beego.Debug(temp)
		reply[k].Content = temp
	}

	if !isAll {
		rand.Seed(time.Now().UnixNano())
		if len(reply) < 2 {
			jsonData, err := models.GetReply(toUserOpenId, reply[0])
			if err != nil {
				beego.Error(err)
			}
			beego.Debug("!isall", jsonData)
			res, err := SendMsgByKF(jsonData, ak)
			if err != nil {
				beego.Error(err)

			}
			beego.Debug(res)
			return
		}
		random := rand.Intn(len(reply) - 1)
		jsonData, err := models.GetReply(toUserOpenId, reply[random])
		if err != nil {
			beego.Error(err)
		}
		beego.Debug("!isall", jsonData)
		res, err := SendMsgByKF(jsonData, ak)
		if err != nil {
			beego.Error(err)

		}
		beego.Debug(res)
		return
	}

	//beego.Debug("reply", reply)
	for _, v := range reply {
		jsonData, err := models.GetReply(toUserOpenId, v)
		if err != nil {
			beego.Error(err)
		}
		beego.Debug("isall", jsonData)
		res, err := SendMsgByKF(jsonData, ak)
		if err != nil {
			beego.Error(err)

		}
		beego.Debug("res", res)
	}
	return
}

// 图文评论
func GetComment(jsonData string, access_token string) (string, error) {
	url_head := "https://api.weixin.qq.com/cgi-bin/comment/list?"
	url := url_head +
		"access_token=" + access_token

	req := httplib.Post(url)
	req.Body(jsonData)

	str, err := req.String()

	if err != nil {
		fmt.Println(err)
		return "", err
	}
	//fmt.Println("getswlist :", str)
	return str, nil
}

func DelComment(jsonData string, access_token string) (string, error) {
	url_head := "https://api.weixin.qq.com/cgi-bin/comment/delete?"
	url := url_head +
		"access_token=" + access_token

	req := httplib.Post(url)
	req.Body(jsonData)

	str, err := req.String()

	if err != nil {
		fmt.Println(err)
		return "", err
	}
	//fmt.Println("getswlist :", str)
	return str, nil
}

func AddCommentReply(jsonData string, access_token string) (string, error) {
	url_head := "https://api.weixin.qq.com/cgi-bin/comment/reply/add?"
	url := url_head +
		"access_token=" + access_token

	req := httplib.Post(url)
	req.Body(jsonData)

	str, err := req.String()

	if err != nil {
		fmt.Println(err)
		return "", err
	}
	//fmt.Println("getswlist :", str)
	return str, nil
}
func DelCommentReply(jsonData string, access_token string) (string, error) {
	url_head := "https://api.weixin.qq.com/cgi-bin/comment/reply/delete?"
	url := url_head +
		"access_token=" + access_token

	req := httplib.Post(url)
	req.Body(jsonData)

	str, err := req.String()

	if err != nil {
		fmt.Println(err)
		return "", err
	}
	//fmt.Println("getswlist :", str)
	return str, nil
}

func OpenComment(jsonData string, access_token string) (string, error) {
	url_head := "https://api.weixin.qq.com/cgi-bin/comment/open?"
	url := url_head +
		"access_token=" + access_token

	req := httplib.Post(url)
	req.Body(jsonData)

	str, err := req.String()

	if err != nil {
		fmt.Println(err)
		return "", err
	}
	//fmt.Println("getswlist :", str)
	return str, nil
}
func CloseComment(jsonData string, access_token string) (string, error) {
	url_head := "https://api.weixin.qq.com/cgi-bin/comment/close?"
	url := url_head +
		"access_token=" + access_token

	req := httplib.Post(url)
	req.Body(jsonData)

	str, err := req.String()

	if err != nil {
		fmt.Println(err)
		return "", err
	}
	//fmt.Println("getswlist :", str)
	return str, nil
}
func CommentMarkelect(jsonData string, access_token string) (string, error) {
	url_head := "https://api.weixin.qq.com/cgi-bin/comment/markelect?"
	url := url_head +
		"access_token=" + access_token

	req := httplib.Post(url)
	req.Body(jsonData)

	str, err := req.String()

	if err != nil {
		fmt.Println(err)
		return "", err
	}
	//fmt.Println("getswlist :", str)
	return str, nil
}
func CommentUnmarkelectt(jsonData string, access_token string) (string, error) {
	url_head := "https://api.weixin.qq.com/cgi-bin/comment/unmarkelect?"
	url := url_head +
		"access_token=" + access_token

	req := httplib.Post(url)
	req.Body(jsonData)

	str, err := req.String()

	if err != nil {
		fmt.Println(err)
		return "", err
	}
	//fmt.Println("getswlist :", str)
	return str, nil
}

//user growth data
func GetUserSummary(jsonData string, access_token string) (string, error) {
	url_head := "https://api.weixin.qq.com/datacube/getusersummary?"
	url := url_head +
		"access_token=" + access_token

	req := httplib.Post(url)
	req.Body(jsonData)

	str, err := req.String()

	if err != nil {
		fmt.Println(err)
		return "", err
	}
	return str, nil
}
func GetUserCumulate(jsonData string, access_token string) (string, error) {
	url_head := "https://api.weixin.qq.com/datacube/getusercumulate?"
	url := url_head +
		"access_token=" + access_token

	req := httplib.Post(url)
	req.Body(jsonData)

	str, err := req.String()

	if err != nil {
		fmt.Println(err)
		return "", err
	}
	return str, nil
}

// getarticletotal

func GetArticleTotal(jsonData string, access_token string) (string, error) {
	url_head := "https://api.weixin.qq.com/datacube/getarticletotal?"
	url := url_head +
		"access_token=" + access_token

	req := httplib.Post(url)
	req.Body(jsonData)

	str, err := req.String()

	if err != nil {
		return "", err
	}
	return str, nil
}

func GetUserRead(jsonData string, access_token string) (string, error) {
	url_head := "https://api.weixin.qq.com/datacube/getuserread?"
	url := url_head +
		"access_token=" + access_token

	req := httplib.Post(url)
	req.Body(jsonData)

	str, err := req.String()

	if err != nil {
		return "", err
	}
	return str, nil
}
func GetArticleSummary(jsonData string, access_token string) (string, error) {
	url_head := "https://api.weixin.qq.com/datacube/getarticlesummary?"
	url := url_head +
		"access_token=" + access_token

	req := httplib.Post(url)
	req.Body(jsonData)

	str, err := req.String()

	if err != nil {
		return "", err
	}
	return str, nil
}

func GetUpstreamMsgHour(jsonData string, access_token string) (string, error) {
	url_head := "https://api.weixin.qq.com/datacube/getupstreammsghour?"
	url := url_head +
		"access_token=" + access_token

	req := httplib.Post(url)
	req.Body(jsonData)

	str, err := req.String()

	if err != nil {
		return "", err
	}
	return str, nil
}
