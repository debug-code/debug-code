package controllers

import (
	"encoding/json"
	"geekerblog/lib"
	"geekerblog/models"
	"github.com/astaxie/beego"
	"strconv"
)

type ArticleController struct {
	beego.Controller
}

func (this *ArticleController) Post() {
	//接收数据
	body := this.Ctx.Input.RequestBody
	article := models.Article{}
	err := json.Unmarshal(body, &article)
	if err != nil {
		beego.Error("ArticleController Post json.Unmarshal err:", err)
		lib.ReturnJson(this.Controller, "444003", nil)
		return
	}

	err = lib.SaveArticle(article)
	if err != nil {
		beego.Error("lib.SaveArticle(articleType) err:", err)
		lib.ReturnJson(this.Controller, "444003", nil)
		return
	}

	lib.ReturnJson(this.Controller, "1", nil)
	return
}

func (this *ArticleController) Get() {
	//get data
	uid := this.Ctx.Input.Param(":id")
	offset, _ := strconv.Atoi(this.GetString("offset"))
	limit, _ := strconv.Atoi(this.GetString("limit"))

	list, err := lib.GetArticle(uid, offset, limit)
	if err != nil {
		beego.Error("ArticleController Get lib.GetArticleType() err:", err)
		lib.ReturnJson(this.Controller, "444002", nil)
		return
	}

	lib.ReturnJson(this.Controller, "1", list)
	return
}

func (this *ArticleController) Put() {
	body := this.Ctx.Input.RequestBody
	article := models.Article{}
	err := json.Unmarshal(body, &article)
	if err != nil {
		beego.Error("ArticleController Post json.Unmarshal err:", err)
		lib.ReturnJson(this.Controller, "444003", nil)
		return
	}

	err = lib.UpdateArticle(article)
	if err != nil {
		beego.Error("lib.SaveArticle(articleType) err:", err)
		lib.ReturnJson(this.Controller, "444003", nil)
		return
	}

	lib.ReturnJson(this.Controller, "1", nil)
	return
}

func (this *ArticleController) Delete() {
	uid := this.Ctx.Input.Param(":id")

	err := lib.DeleteArticle(uid)
	if err != nil {
		beego.Error("lib.SaveArticle(articleType) err:", err)
		lib.ReturnJson(this.Controller, "444003", nil)
		return
	}

	lib.ReturnJson(this.Controller, "1", nil)
	return
}
