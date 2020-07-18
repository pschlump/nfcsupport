package support

// This file is MIT licensed.
// Copyright (C) 2013-2019 Philip Schlump

import (
	"fmt"
	"os"
	"sync"

	"github.com/pschlump/MiscLib"
	"github.com/pschlump/godebug"
	"github.com/pschlump/radix.v2/redis"
)

type RedisConfigType struct {
	RedisConnectHost string `json:"redis_host" default:"$ENV$REDIS_HOST"`
	RedisConnectAuth string `json:"redis_auth" default:"$ENV$REDIS_AUTH"`
	RedisConnectPort string `json:"redis_port" default:"6379"`
}

// RedisClient makes a connection to the Redis datagbase and returns the client connection and a true/false flag.
// If the configuration includes an non-empty RedisConnectAuth then it will also do authenication with the AUTH
// command in the redis system.
func RedisClient(gCfg *RedisConfigType) (client *redis.Client, conFlag bool) {
	var err error
	if dbFlag["RedisClient"] {
		fmt.Printf("AT: connect to redis with: %s %s\n", godebug.LF(), gCfg.RedisConnectHost+":"+gCfg.RedisConnectPort)
	}
	client, err = redis.Dial("tcp", gCfg.RedisConnectHost+":"+gCfg.RedisConnectPort)
	if err != nil {
		fmt.Printf("Error on connect to redis:%s, fatal\n", err)
		fmt.Fprintf(os.Stderr, "%s\n\n\n-----------------------------------------------------------------------------------------------\nError on connect to redis:%s, fatal\n", MiscLib.ColorRed, err)
		fmt.Fprintf(os.Stderr, "Config Data: %s\n", godebug.SVarI(gCfg))
		fmt.Fprintf(os.Stderr, "\n-----------------------------------------------------------------------------------------------\n\n\n%s", MiscLib.ColorReset)
		return
	}
	if gCfg.RedisConnectAuth != "" {
		err = client.Cmd("AUTH", gCfg.RedisConnectAuth).Err
		if err != nil {
			fmt.Printf("Error on connect to Redis --- Invalid authentication:%s, fatal\n", err)
			fmt.Fprintf(os.Stderr, "%s\nError on connect to Redis --- Invalid authentication:%s, fatal%s\n\n", MiscLib.ColorRed, err, MiscLib.ColorReset)
			return
		}
		conFlag = true
	}
	conFlag = true
	return
}

// ----------------------------------------------------------------------------------------------------------------------
// Redis interface
// ----------------------------------------------------------------------------------------------------------------------
type RedisConnection struct {
	redisClient *redis.Client
	redisMux    *sync.Mutex
	TtlRedis    int
}

func NewRedisConnection(gCfg *RedisConfigType) (rv *RedisConnection, err error) {
	cli, ok := RedisClient(gCfg)
	if !ok {
		err = fmt.Errorf("Faild to connect to redis")
		return
	}
	rv = &RedisConnection{
		redisClient: cli,
		redisMux:    &sync.Mutex{},
		TtlRedis:    0,
	}
	return
}

//func (rCon *RedisConnection) ConnectToRedis() {
//	rCon.redisClient = redisClient
//}

// GetRedis queries redis to return a value.
func (rCon *RedisConnection) GetRedis(key string) (rv string, err error) {
	rCon.redisMux.Lock()
	defer rCon.redisMux.Unlock()
	rv, err = rCon.redisClient.Cmd("GET", key).Str()
	if err != nil {
		fmt.Fprintf(logFilePtr, "ERROR: unable to get data from redis, key=%s, err=%s\n", key, err)
	}
	return
}

// SetRedis Sets `key` to value `data` in redis.
func (rCon *RedisConnection) SetRedis(key, data string) (err error) {
	fmt.Fprintf(logFilePtr, "Redis Set key[%s] value [%s] at:%s\n", key, data, godebug.LF(-2))
	rCon.redisMux.Lock()
	defer rCon.redisMux.Unlock()
	if rCon.TtlRedis == 0 {
		err = rCon.redisClient.Cmd("SET", key, data).Err
	} else {
		err = rCon.redisClient.Cmd("SETEX", key, rCon.TtlRedis, data).Err
	}
	if err != nil {
		fmt.Fprintf(logFilePtr, "Unable to SETEX data in redis, key=%s val=%s err=%s\n", key, data, err)
	}
	return err
}

// SetRedis Sets `key` to value `data` in redis.
func (rCon *RedisConnection) SetRedisTTL(key, data string, Ttl int) (err error) {
	rCon.redisMux.Lock()
	defer rCon.redisMux.Unlock()
	err = rCon.redisClient.Cmd("SETEX", key, Ttl, data).Err
	if err != nil {
		fmt.Fprintf(logFilePtr, "Unable to SETEX data in redis, key=%s val=%s err=%s\n", key, data, err)
	}
	return err
}

func (rCon *RedisConnection) DelRedis(key string) (err error) {
	rCon.redisMux.Lock()
	defer rCon.redisMux.Unlock()
	err = rCon.redisClient.Cmd("DEL", key).Err
	if err != nil {
		fmt.Fprintf(logFilePtr, "Unable to DELETE data in redis, key=%s err=%s\n", key, err)
	}
	return err
}

/* vim: set noai ts=4 sw=4: */
