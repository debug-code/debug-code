package timedTask

import (
	"errors"
	"wechat-go/tools/timedTask/redis"
)

//type CallBack func(key string)

var (
	Client redis.RedisPool
)

func RegisterTimedTask(callBack redis.CallBack, driverName, dataSource string, database int, passwd ...string) error {
	switch driverName {
	case "REDIS":
		if len(passwd) > 0 {
			Client = redis.CreateNewRedisPool(dataSource, database, passwd[0])
			err := redis.Redister(callBack, dataSource, database, passwd[0])
			if err != nil {
				return err
			}
		} else {
			Client = redis.CreateNewRedisPool(dataSource, database)
			err := redis.Redister(callBack, dataSource, database)
			if err != nil {
				return err
			}
		}

	case "MEMORY":
		return errors.New("not supported yet")
	default:
		return errors.New("invalid driverName")
	}
	return nil
}
