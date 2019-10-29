package redis

import (
	"errors"
	"fmt"
	"github.com/gomodule/redigo/redis"
	"time"
)

type RedisPool struct {
	pool *redis.Pool
}

func CreateNewRedisPool(dataSource string, database int, passwd ...string) RedisPool {
	rp := RedisPool{}
	if len(passwd) > 0 {
		rp.register(dataSource, database, passwd[0])
	} else {
		rp.register(dataSource, database)
	}

	return rp

}
func (rp *RedisPool) register(dataSource string, database int, passwd ...string) {
	if len(passwd) > 0 {
		rp.pool = &redis.Pool{
			MaxIdle:     20,
			MaxActive:   20,
			IdleTimeout: 200 * time.Second,
			Wait:        true,
			Dial: func() (redis.Conn, error) {
				con, err := redis.Dial("tcp", dataSource,
					redis.DialPassword(passwd[0]),
					redis.DialDatabase(database),
					redis.DialConnectTimeout(20*time.Second),
					redis.DialReadTimeout(20*time.Second),
					redis.DialWriteTimeout(20*time.Second))
				if err != nil {
					return nil, err
				}
				return con, nil
			},
		}
	} else {
		rp.pool = &redis.Pool{
			MaxIdle:     20,
			MaxActive:   20,
			IdleTimeout: 200 * time.Second,
			Wait:        true,
			Dial: func() (redis.Conn, error) {
				con, err := redis.Dial("tcp", dataSource,
					//redis.DialPassword(conf["Password"].(string)),
					redis.DialDatabase(database),
					redis.DialConnectTimeout(20*time.Second),
					redis.DialReadTimeout(20*time.Second),
					redis.DialWriteTimeout(20*time.Second))
				if err != nil {
					return nil, err
				}
				return con, nil
			},
		}
	}

}

func (rp *RedisPool) do(commandName string, args ...interface{}) (reply interface{}, err error) {
	if len(args) < 1 {
		return nil, errors.New("missing required arguments")
	}
	//args[0] = pool.associate(args[0])
	c := rp.pool.Get()

	//defer c.Close()
	//fmt.Println(args)
	return c.Do(commandName, args...)
}
func (rp *RedisPool) Get(key string) (reply interface{}, err error) {
	reply, err = rp.do("GET", key)
	if err != nil {
		fmt.Println("err while setting:", err)
	}
	return
}

func (rp *RedisPool) Put(key string, val interface{}, timeout int) error {

	if timeout != 0 {
		_, err := rp.do("SET", key, val, "EX", timeout)
		return err
	}
	_, err := rp.do("SET", key, val)
	return err
}

func (rp *RedisPool) Delete(key string) error {
	_, err := rp.do("DEL", key)
	return err
}
