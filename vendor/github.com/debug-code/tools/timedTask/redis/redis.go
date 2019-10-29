package redis

import (
	"fmt"
	"github.com/gomodule/redigo/redis"
	"strconv"
)

type CallBack func(key string)

func Redister(callBack CallBack, dataSource string, database int, passwd ...string) error {
	//创建链接
	//option := redis.DialOption{}
	//if len(passwd) > 0 {
	//	option = redis.DialPassword(passwd[0])
	//}
	conn, err := redis.Dial("tcp", dataSource)
	if err != nil {
		//fmt.Println("redis dial failed.")
		return err
	}
	//defer conn.Close()
	db := strconv.Itoa(database)
	//设置订阅链接
	client := redis.PubSubConn{conn}
	channel := "__keyevent@" + db + "__:expired"
	//选择订阅频道
	err = client.PSubscribe(channel)
	if err != nil {
		//fmt.Println(err)
		return err
	}

	//监听订阅频道
	go func() {
		for {

			switch res := client.Receive().(type) {

			case redis.Message:
				//fmt.Println("过期", "    key", string(res.Data), time.Now().Second())
				//
				if res.Channel == channel {
					go callBack(string(res.Data))
				}
				//else {
				//	beego.Debug("其他", "类型", res.Channel, "    key", string(res.Data))
				//}
			case redis.Subscription:
				fmt.Println(res.Channel, res.Kind, res.Count)
			case error:
				fmt.Println("Receivd failed:", res)

			}
		}
	}()
	return nil
}
