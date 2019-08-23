package controllers

import (
	"fmt"
	"github.com/astaxie/beego"
	"regexp"
	"strings"
)

type StgsController struct {
	beego.Controller
}

func (this *StgsController) Post() {

	body := this.Ctx.Input.RequestBody

	//Stgs := models.Stgs{}
	//err := json.Unmarshal(body, &Stgs)
	//if err != nil{
	//	fmt.Println(err)
	//}
	//end := len(string(body))-4
	text := string(body)
	text = strings.Replace(text, "`", "", -1)
	text = strings.Replace(text, "(\n", "", -1)
	text = strings.Replace(text, ")", "", -1)
	//text = strings.Replace(text, "\n", "", -1)
	//text := string(body)[8:end]
	beego.Info("text", text)

	lines := strings.Split(text, "\n")

	tableName := strings.Split(lines[0], " ")[2]

	beego.Info("tableName", tableName)
	reg := regexp.MustCompile(`\(?.*(?:,)`)
	list := reg.FindAllString(text, -1)

	beego.Info("list", list)

	maps := map[string]string{}
	for _, v := range list {

		if len(v) > 2 {
			item := strings.Split(v, ",")[0]
			//reg := regexp.MustCompile(`[\t|\n|\v|\f|\r]*`)
			//ss := reg.ReplaceAllString(item, "1")
			//
			//t := strings.Split(ss, "")
			//ts := strings.Replace(t[0], " ", "", -1)
			//beego.Info("ts", ts)
			//beego.Info("s", t)
			//tt := marshal(t[0])
			//beego.Info("tt", tt)

			t := strings.Fields(item)
			tt := marshal(t[0])
			beego.Info("t", t)
			beego.Info("tt", tt)
			beego.Info("", v)
			if strings.Contains(v, "varchar") {
				maps[tt] = "string" + "\t\t`" + `orm:"column(` + t[0] + ")\";json:\"" + t[0] + "\"`"
			} else if strings.Contains(v, "longtext") {
				maps[tt] = "string" + "\t\t`" + `orm:"column(` + t[0] + ")\";json:\"" + t[0] + "\"`"
			} else if strings.Contains(v, "TINYINT") {
				maps[tt] = "bool" + "\t\t`" + `orm:"column(` + t[0] + ")\";json:\"" + t[0] + "\"`"
			} else if strings.Contains(v, "int") {
				maps[tt] = "int" + "\t\t`" + `orm:"column(` + t[0] + ")\";json:\"" + t[0] + "\"`"
			}
		}
	}
	fmt.Println(maps)
	res := `type ` + marshal(tableName) + ` struct { `
	for k, v := range maps {
		res += "\n\t" + k + "\t\t\t\t" + v
	}
	res += "\n}"
	this.Ctx.WriteString(res)
}

func marshal(str string) string {
	var res string
	//res += strings.ToUpper(str[0:1])
	flag := 0
	for k, v := range str {
		if k == 0 {
			res += strings.ToUpper(string(v))
			continue
		}
		if string(v) == "_" {
			flag = 1
			continue
		}
		if flag == 1 {
			res += strings.ToUpper(string(v))
			flag = 0
			continue
		}
		res += string(v)
	}

	return res
}
