package mysql

import (
	"fmt"
	"github.com/debug-code/goDB"
	"github.com/jinzhu/gorm"
	"time"
)

var pools *goDB.GenericPool

// 功能:初始化连接池
// 参数:
// dialect：数据库类型（mysql...） ,
// connect：链接字符串 (account:pwd@tcp(xxx.xxx.xxx.xxx:3306)/xxx?charset=utf8mb4)
// min：最小连接数
// max：最大连接数
// maxLifeTime：生命周期
func Regist(dialect, connect string, min, max int, maxLifeTime time.Duration) (int, error) {

	factory := func() (*gorm.DB, error) {
		db, err := gorm.Open(dialect, connect)
		if err != nil {
			return db, err
		}

		//设置表名补位负数  xxx -> xxxs
		db.SingularTable(true)
		return db, nil
	}
	close := func(db *gorm.DB) error {

		err := db.Close()
		if err != nil {
			return err
		}
		return nil
	}

	var err error
	pools, err = goDB.NewGenericPool(min, max,
		maxLifeTime, factory, close)
	if err != nil {
		return 0, err
	}

	return pools.GetLife(), nil

}

func Context() (*gorm.DB, error) {
	db, err := pools.Acquire()
	if err != nil {
		return nil, err
	}

	//db := dbi.(*gorm.DB)

	//db.Find()
	//fmt.Println("numOpen", pools.GetLife())

	return db, err
}

func CloseC(db *gorm.DB) error {
	err := pools.Close(db)
	if err != nil {
		return err
	}
	return nil
}

func ReleaseC(db *gorm.DB) {
	err := pools.Release(db)
	if err != nil {
		fmt.Println("ReleaseC error:", err)
	}
}
