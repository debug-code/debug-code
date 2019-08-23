package models

//用户
type Manager struct {
	Id         int    `json:"id"`
	Uid        string `json:"uid"`
	Account    string `json:"account"`
	Passwd     string `json:"passwd"`
	Phone      string `json:"phone"`
	Status     int    `json:"status"`
	UserName   string `json:"user_name"`
	LastTime   int    `json:"last_time"`
	ThisTime   int    `json:"this_time"`
	Character  string `json:"character"`
	Email      string `json:"email"`
	CreateTime int    `json:"create_time"`
}

//文章
type Article struct {
	Title           string `json:"title"`
	Manager         int    `json:"manager"`
	Url             string `json:"url"`
	Key             string `json:"key"`
	Uid             string `json:"uid"`
	Views           int    `json:"views"`
	Id              int    `json:"id"`
	Status          int    `json:"status"`
	Summary         string `json:"summary"`
	Tag             string `json:"tag"`
	Type            string `json:"type"`
	CreateTime      int    `json:"create_time"`
	UpdateTime      int    `json:"update_time"`
	ContentMarkdown string `json:"content_markdown"`
	ContentHtml     string `json:"content_html"`
}

//个人分类
type ArticleType struct {
	Id         int    `json:"id"`
	Uid        string `json:"uid"`
	Type       string `json:"type"`
	Remark     string `json:"remark"`
	CreateTime int    `json:"create_time"`
	UpdateTime int    `json:"update_time"`
	Status     int    `json:"status"`
}
