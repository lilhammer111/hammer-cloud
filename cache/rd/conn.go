package rd

import (
	"github.com/gomodule/redigo/redis"
	"time"
)

var (
	pool      *redis.Pool
	redisHost = "127.0.0.1:6379"
	//redisPWD = ""
)

func newRedisPool() *redis.Pool {
	return &redis.Pool{
		Dial: func() (redis.Conn, error) {
			dial, err := redis.Dial("tcp", redisHost)
			if err != nil {
				return nil, err
			}

			//_, err = dial.Do("AUTH", "")
			//if err != nil {
			//	log.Println(err)
			//	_ = dial.Close()
			//	return nil, err
			//}
			return dial, nil
		},
		MaxIdle:     50,
		MaxActive:   30,
		IdleTimeout: 300 * time.Second,
		TestOnBorrow: func(c redis.Conn, t time.Time) error {
			if time.Since(t) < time.Minute {
				return nil
			}
			_, err := c.Do("PING")
			return err
		},
	}
}

func init() {
	pool = newRedisPool()
}

func Pool() *redis.Pool {
	return pool
}
