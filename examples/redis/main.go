package main

import (
	"errors"
	"fmt"
	"github.com/gomodule/redigo/redis"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func main() {
	go func() {
		for {
			if err := testConn(); err != nil {
				log.Println(err)
				time.Sleep(time.Second * 1)
			} else {
				break
			}
		}
	}()
	signalHandler()
}

func testConn() (err error) {
	redisInstance := newPool("redis-master:6379", "").Get()
	defer func() {
		if err := redisInstance.Close(); err != nil {
			log.Println(err)
		}
	}()
	if _, err := redisInstance.Do("set", "abc", 123); err != nil {
		return errors.New(fmt.Sprintf("set err:%v", err))
	}
	if res, err := redis.String(redisInstance.Do("get", "abc")); err != nil {
		return errors.New(fmt.Sprintf("get err:%v", err))
	} else {
		log.Println("res:", res)
	}
	return nil
}

func newPool(addr string, pwd string) *redis.Pool {
	return &redis.Pool{
		MaxIdle:     3,
		IdleTimeout: 240 * time.Second,
		// Dial or DialContext must be set. When both are set, DialContext takes precedence over Dial.
		Dial: func() (redis.Conn, error) {
			c, err := redis.Dial("tcp", addr)
			if err != nil {
				log.Println("redis.Dial err:", err)
				return nil, err
			}
			if pwd != "" {
				if _, err := c.Do("AUTH", pwd); err != nil {
					if err2 := c.Close(); err2 != nil {
						log.Println("close err2:", err2)
					}
					return nil, err
				}
			}
			return c, nil
		},
	}
}

func signalHandler() {
	var (
		ch = make(chan os.Signal, 1)
	)
	tick := time.NewTicker(time.Second * 5)
	signal.Notify(ch, syscall.SIGHUP, syscall.SIGQUIT, syscall.SIGTERM, syscall.SIGINT)
	for {
		select {
		case sig := <-ch:
			log.Printf("get a signal %s, stop the daemon-hook container \n", sig.String())
			switch sig {
			case syscall.SIGQUIT, syscall.SIGTERM, syscall.SIGINT:
				time.Sleep(time.Second)
				return
			case syscall.SIGHUP:
			default:
				return
			}
		case <-tick.C:
			log.Print("tick point")
		}
	}
}
