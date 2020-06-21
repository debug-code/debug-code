package lib

import (
	"debug-code/models"
	"debug-code/tools/goDB/mysql"
	"debug-code/tools/jwt"
	"errors"
	"github.com/astaxie/beego"
)

func CheckManager(manager models.Manager) (string, error) {

	//search manager from database
	db, err := mysql.Context()
	if err != nil {
		beego.Error("CheckManager mysql.Context() err:", err)
		return "get db error", err
	}
	defer mysql.ReleaseC(db)
	managerS := models.Manager{}

	err = db.Where("status = ?", 1).
		Where("account = ?", manager.Account).
		First(&managerS).Error

	if err != nil || manager.Account != managerS.Account || managerS.Account == "" {
		beego.Error("CheckManager db.first", err)

		return "未找到该用户", err
	}
	if manager.Passwd != managerS.Passwd {
		beego.Error("CheckManager pwd", err)
		return "密码错误！", errors.New("pwd error")
	}

	return managerS.Uid, nil
}

func GetToken(uid string) (string, error) {

	mySigningKey := []byte("hzwy23")
	token, err := jwt.GetNewToken(mySigningKey, uid)
	if err != nil {
		beego.Error("GetToken GetNewToken err:", err)
		return "", err
	}

	return token, nil
}
