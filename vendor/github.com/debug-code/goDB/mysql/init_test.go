package mysql

import (
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"strconv"
	"sync"
	"testing"
	"time"
)

// 基本模型的定义
type Model struct {
	ID        uint `gorm:"primary_key"`
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt *time.Time
}

// 添加字段 `ID`, `CreatedAt`, `UpdatedAt`, `DeletedAt`
type Managers struct {
	Model
	Name string
}

type Manager struct {
	Id         int
	Account    string
	Passwd     string
	Phone      string
	Status     int
	UserName   string
	LastTime   int
	Character  string
	Email      string
	CreateTime int `gorm:"column:create_time"`
}

func TestRegist(t *testing.T) {
	n, err := Regist("mysql",
		"alex:123qwe=-0@tcp(101.200.148.236:13306)/geekerblog?charset=utf8mb4",
		20, 50, time.Second*1800)
	if err != nil {
		fmt.Println("fac err:", err)
		//return db, err
	}
	fmt.Println("连接个数", n)
	time.Sleep(time.Second * 2)

	start := time.Now().UnixNano()
	var wg sync.WaitGroup
	for i := 0; i < 10000; i++ {
		wg.Add(1)
		//time.Sleep(time.Millisecond * 2)
		//godb.g
		//fmt.Println("index", i)
		go func() {
			defer wg.Done()
			db, err := Context()
			if err != nil {
				fmt.Println("fac err:", err)
				//return db, err
			}

			manager1 := Manager{Account: "sdf" + strconv.Itoa(i), Passwd: "sdsdff" + strconv.Itoa(i)}
			db.Create(&manager1)

			err = ReleaseC(db)
			if err != nil {
				fmt.Println("fac err:", err)
				//return db, err
			}

		}()

		//db.First(&manager1)
		//fmt.Println("sdf", manager1)
	}
	wg.Wait()
	end := time.Now().UnixNano()

	fmt.Println(end - start)
}
