package lib

import (
	"debug-code/models"
	"debug-code/tools/jwt"
	"errors"
	"github.com/astaxie/beego"
	"github.com/debug-code/goDB/mysql"
)

func GetManager(this beego.Controller) (models.Manager, error) {
	tokens := this.Ctx.Input.Context.Request.Header.Get("Authorization")
	manager := models.Manager{}
	uid, err := GetTokenInfo(tokens)
	if err != nil {
		beego.Error("GetManager GetTokenInfo err:", err)
		return manager, err
	}

	db, err := mysql.Context()
	if err != nil {
		beego.Error("GetManager mysql.Context() err:", err)
		return manager, err
	}
	defer mysql.ReleaseC(db)

	err = db.Where("uid = ?", uid).First(&manager).Error
	if err != nil {
		beego.Error("GetManager First err:", err)
		return manager, err
	}

	return manager, nil
}

func GetTokenInfo(token string) (string, error) {
	mySigningKey := []byte("hzwy23")
	t, err := jwt.GetTokenInfo(mySigningKey, token)
	if err != nil {
		return "", errors.New("wrong")
	}
	return t, nil
}
