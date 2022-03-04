package dao

import (
	"readingo/conf"
	"readingo/util"
	"errors"
	"github.com/go-redis/redis/v8"
)

var redisClusterClients = make(map[string]*redis.ClusterClient, 0)

func init() {
	for _, _db := range conf.RedisClusters {
		redisClusterClients[_db.NAME] = newRedisClusterClient(_db)
	}
}

func GetAllRedisClusterClients() map[string]*redis.ClusterClient {
	return redisClusterClients
}

func GetClusterConn(host string) (*redis.ClusterClient, error) {
	if client, ok := redisClusterClients[host]; ok {
		return client, nil
	}
	return nil, errors.New("not found")
}

func newRedisClusterClient(db conf.RedisClusterConf) *redis.ClusterClient {
	conn := redis.ClusterOptions{
		Addrs:       db.HOST,
		PoolSize:    db.MaxActive,
		IdleTimeout: util.ParseStringToDuration(db.IdleTimeout),
	}

	if db.PASSWORD != "" {
		conn.Password = db.PASSWORD
	}

	return redis.NewClusterClient(&conn)
}
