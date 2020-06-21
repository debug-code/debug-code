package controllers

import (
	"debug-code/lib"
	"debug-code/models"
	"encoding/json"
	"github.com/astaxie/beego"
)

type ArticleTypeController struct {
	beego.Controller
}

func (this *ArticleTypeController) Post() {
	//接收数据
	body := this.Ctx.Input.RequestBody
	articleType := models.ArticleType{}
	err := json.Unmarshal(body, &articleType)
	if err != nil {
		beego.Error("ArticleTypeController Post json.Unmarshal err:", err)
		lib.ReturnJson(this.Controller, "444003", nil)
		return
	}

	err = lib.SaveArticleType(articleType)
	if err != nil {
		beego.Error("lib.SaveArticleType(articleType) err:", err)
		lib.ReturnJson(this.Controller, "444003", nil)
		return
	}

	lib.ReturnJson(this.Controller, "1", nil)
	return
}

func (this *ArticleTypeController) Get() {

	list, err := lib.GetArticleType()
	if err != nil {
		beego.Error("ArticleTypeController Get lib.GetArticleType() err:", err)
		lib.ReturnJson(this.Controller, "444002", nil)
		return
	}

	lib.ReturnJson(this.Controller, "1", list)
	return
}

func (this *ArticleTypeController) Put() {
	//接收数据
	body := this.Ctx.Input.RequestBody
	articleType := models.ArticleType{}
	err := json.Unmarshal(body, &articleType)
	if err != nil {
		beego.Error("ArticleTypeController put json.Unmarshal err:", err)
		lib.ReturnJson(this.Controller, "444003", nil)
		return
	}
	err = lib.UpdateArticleType(articleType)
	if err != nil {
		beego.Error("ArticleTypeController put UpdateArticleType err:", err)
		lib.ReturnJson(this.Controller, "444003", nil)
		return
	}

	lib.ReturnJson(this.Controller, "1", nil)
	return
}

func (this *ArticleTypeController) Delete() {

	uid := this.Ctx.Input.Param(":id")
	err := lib.DeleteArticleType(uid)
	if err != nil {
		beego.Error("ArticleTypeController Delete DeleteArticleType err:", err)
		lib.ReturnJson(this.Controller, "444003", nil)
		return
	}

	lib.ReturnJson(this.Controller, "1", nil)
	return

}
