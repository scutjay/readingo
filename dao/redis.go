package dao

import (
	"context"
	"errors"
	"github.com/alecthomas/log4go"
	"github.com/go-redis/redis/v8"
	"readingo/conf"
	"readingo/constant"
	"readingo/model"
	"readingo/util"
	"strconv"
	"sync"
	"time"
)

var redisDBMap = make(map[string]map[string]*redis.Client, 0)
var lastUpdateTime time.Time
var hostCache []model.HostInfo
var refreshDBTreeMutex sync.Mutex

func init() {
	for _, _db := range conf.Redis {
		redisDBMap[_db.NAME] = newRedisDB(_db)
	}
	go autoRefreshDBTree()
}

func autoRefreshDBTree() {
	duration := constant.DBTreeAutoRefreshNormalDuration
	failCount := 0
	for {
		time.Sleep(duration)
		if err := RefreshDBTree(context.Background()); err != nil {
			log4go.Error("Auto refresh db tree failed, error: %v", err)
			failCount++
		}
		if failCount > constant.DBTreeRefreshFailTimes {
			duration = constant.DBTreeAutoRefreshLongDuration
		}
	}
}

func GetConn(dbKey, partitionKey string) (*redis.Client, error) {
	if db, ok := redisDBMap[dbKey]; ok {
		if par, ok := db[partitionKey]; ok {
			return par, nil
		}
	}
	return nil, errors.New("not found")
}

func RefreshDBTree(ctx context.Context) error {
	refreshDBTreeMutex.Lock()
	defer refreshDBTreeMutex.Unlock()

	if !timeToRefresh() {
		return errors.New("refresh too often, please try later")
	}

	hosts := make([]model.HostInfo, 0)
	for host, db := range redisDBMap {
		hostInfo := model.HostInfo{
			Host: host,
			Type: "redis",
		}
		partitionInfos := make([]model.PartitionInfo, 0)
		for i := 0; i < 16; i++ {
			log4go.Info("Get DB size for host[%s@%d] failed", host, i)
			conn := db[strconv.Itoa(i)]
			reply, err := conn.Do(ctx, "DBSIZE").Result()
			if err != nil {
				return err
			} else {
				partitionInfo := model.PartitionInfo{
					Index:  i,
					DBSize: reply.(int64),
				}
				partitionInfos = append(partitionInfos, partitionInfo)
			}
		}
		hostInfo.Partitions = partitionInfos
		hosts = append(hosts, hostInfo)
	}
	lastUpdateTime = time.Now()
	hostCache = hosts
	log4go.Info("Refresh DB tree at %s", time.Now().String())
	return nil
}

func newRedisDB(db conf.RedisConf) map[string]*redis.Client {
	redisDBs := make(map[string]*redis.Client)
	for i := 0; i < 16; i++ {
		redisDBs[strconv.Itoa(i)] = newRedisClient(db, i)
	}
	return redisDBs
}

func newRedisClient(db conf.RedisConf, idx int) *redis.Client {
	return redis.NewClient(&redis.Options{
		Addr:        db.HOST,
		DB:          idx,
		PoolSize:    db.MaxActive,
		PoolTimeout: util.ParseStringToDuration(db.ConnectTimeout),
		IdleTimeout: util.ParseStringToDuration(db.IdleTimeout),
		Password:    db.PASSWORD,
	})
}

func timeToRefresh() bool {
	return lastUpdateTime.Add(constant.DBTreeRefreshMinInterval).Before(time.Now())
}

func GetLastUpdateTime() time.Time {
	return lastUpdateTime
}

func GetHostsCache(ctx context.Context) []model.HostInfo {
	if len(hostCache) == 0 || timeToRefresh() {
		_ = RefreshDBTree(ctx)
	}
	return hostCache
}
