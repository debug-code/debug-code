package lib

import "github.com/astaxie/beego"

func GetTokenByContorller(con beego.Controller) (string, error) {
	tokens := con.Ctx.Input.Context.Request.Header.Get("Authorization")
	uid, err := GetTokenInfo(tokens)
	if err != nil {
		beego.Error("GetTokenByContorller GetTokenInfo err:", err)
		return "", err
	}
	token, err := GetToken(uid)
	if err != nil {
		beego.Error("GetTokenByContorller GetToken err:", err)
		return "", err
	}

	return token, nil

}
