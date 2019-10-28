package controllers

import (
	"fmt"
	"geekerblog/lib"
	"github.com/astaxie/beego"
	"strconv"
)

type OpenArticlesController struct {
	beego.Controller
}

func (this *OpenArticlesController) Get() {
	//get data

	uid := this.Ctx.Input.Param(":id")
	offset, _ := strconv.Atoi(this.GetString("offset"))
	limit, _ := strconv.Atoi(this.GetString("limit"))

	fmt.Println("uid", uid)

	list, err := lib.GetArticle(uid, offset, limit)
	if err != nil {
		beego.Error("ArticleController Get lib.GetArticleType() err:", err)
		lib.ReturnJson(this.Controller, "444002", nil)
		return
	}

	lib.ReturnJson(this.Controller, "1", list)
	return
}
