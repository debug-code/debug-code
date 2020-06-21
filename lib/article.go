package lib

import (
	"debug-code/models"
	ep "debug-code/tools/encryption"
	"debug-code/tools/goDB/mysql"
	"github.com/astaxie/beego"
	"strconv"
	"time"
)

func SaveArticle(article models.Article) (models.Article, error) {

	article.CreateTime = int(time.Now().Unix())
	article.UpdateTime = article.CreateTime
	article.Uid = ep.StringToSha256(article.Title +
		strconv.Itoa(article.CreateTime))
	article.Status = 1
	article.Views = 0

	db, err := mysql.Context()
	if err != nil {
		beego.Error("fac err:", err)
		return article, err
	}
	defer mysql.ReleaseC(db)

	db.Create(&article)
	if db.Error != nil {
		beego.Error(db.Error)
		return article, err
	}

	return article, nil
}

func GetArticle(uid string, offset, limit int) (interface{}, error) {

	list := []models.Article{}
	maps := map[string]interface{}{}
	db, err := mysql.Context()
	if err != nil {
		beego.Error("GetArticle() mysql.Context() err:", err)
		return maps, err
	}
	defer mysql.ReleaseC(db)

	if uid != "" {
		db.Where("status != 0 ").Where("uid = ?", uid).Find(&list)
		return list, nil
	}

	counts := 0
	temp := db.Where("status != 0")
	temp.Find(&list).Count(&counts)

	temp = temp.Offset(offset)

	if limit != 0 {
		temp = temp.Limit(limit)
	}

	temp.Order("update_time desc").Find(&list)
	if db.Error != nil {
		beego.Error("GetArticle() db.Find(&list)err:", db.Error)
		return list, err
	}

	maps["list"] = list
	maps["counts"] = counts

	return maps, nil
}

func UpdateArticle(article models.Article) error {
	db, err := mysql.Context()
	if err != nil {
		beego.Error("UpdateArticle Context err:", err)
		return err
	}
	defer mysql.ReleaseC(db)

	articleOld := models.Article{}

	db.Find(&articleOld, "uid = ?", article.Uid)
	if db.Error != nil {
		beego.Error("UpdateArticle Find", db.Error)
		return err
	}

	//set update data
	article.UpdateTime = int(time.Now().Unix())
	article.Key = articleOld.Key
	article.Id = articleOld.Id
	article.CreateTime = articleOld.CreateTime
	article.Status = articleOld.Status
	article.Url = articleOld.Url
	article.Uid = articleOld.Uid
	article.Views = articleOld.Views

	db.Save(&article)
	if db.Error != nil {
		beego.Error("UpdateArticle Update", db.Error)
		return err
	}

	return nil
}

func DeleteArticle(uid string) error {
	db, err := mysql.Context()
	if err != nil {
		beego.Error("DeleteArticle Context err:", err)
		return err
	}
	defer mysql.ReleaseC(db)

	article := models.Article{}
	db.Find(&article, "uid = ?", uid)
	if db.Error != nil {
		beego.Error("DeleteArticle Find", db.Error)
		return err
	}

	db.Model(&article).Update("status", 0)
	if db.Error != nil {
		beego.Error("DeleteArticle Update", db.Error)
		return err
	}
	return nil
}
