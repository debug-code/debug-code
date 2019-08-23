package lib

import (
	"geekerblog/models"
	ep "geekerblog/tools/encryption"
	"github.com/astaxie/beego"
	"github.com/debug-code/goDB/mysql"
	"strconv"
	"time"
)

func SaveArticleType(articleType models.ArticleType) error {

	articleType.CreateTime = int(time.Now().Unix())
	articleType.UpdateTime = articleType.CreateTime
	articleType.Uid = ep.StringToSha256(articleType.Type +
		strconv.Itoa(articleType.CreateTime))
	articleType.Status = 1

	db, err := mysql.Context()
	if err != nil {
		beego.Error("fac err:", err)
		return err
	}
	defer mysql.ReleaseC(db)

	db.Create(&articleType)
	if db.Error != nil {
		beego.Error(db.Error)
		return err
	}

	return nil
}

func GetArticleType() ([]models.ArticleType, error) {

	list := []models.ArticleType{}
	db, err := mysql.Context()
	if err != nil {
		beego.Error("GetArticleType() mysql.Context() err:", err)
		return list, err
	}
	defer mysql.ReleaseC(db)

	db.Where("status != 0").Find(&list)
	if db.Error != nil {
		beego.Error("GetArticleType() db.Find(&list)err:", db.Error)
		return list, err
	}

	return list, nil
}

func UpdateArticleType(articleType models.ArticleType) error {
	db, err := mysql.Context()
	if err != nil {
		beego.Error("UpdateArticleType Context err:", err)
		return err
	}
	defer mysql.ReleaseC(db)

	articleTypeOld := models.ArticleType{}

	db.Find(&articleTypeOld, "uid = ?", articleType.Uid)
	if db.Error != nil {
		beego.Error("UpdateArticleType Find", db.Error)
		return err
	}

	articleTypeOld.Type = articleType.Type
	articleTypeOld.Remark = articleType.Remark
	articleTypeOld.UpdateTime = int(time.Now().Unix())

	db.Save(&articleTypeOld)
	if db.Error != nil {
		beego.Error("UpdateArticleType Update", db.Error)
		return err
	}

	return nil
}

func DeleteArticleType(uid string) error {
	db, err := mysql.Context()
	if err != nil {
		beego.Error("DeleteArticleType Context err:", err)
		return err
	}
	defer mysql.ReleaseC(db)

	article := models.ArticleType{}
	db.Find(&article, "uid = ?", uid)
	if db.Error != nil {
		beego.Error("DeleteArticleType Find", db.Error)
		return err
	}

	db.Model(&article).Update("status", 0)
	if db.Error != nil {
		beego.Error("DeleteArticleType Update", db.Error)
		return err
	}
	return nil
}
