/**
 * @author  zhaoliang.liang
 * @date  2025/8/6 14:12
 */

package redis_drivers

import (
	"context"
	"fmt"
	"log/slog"
	"lxz/internal/config"
	"strconv"
	"sync"

	"github.com/go-redis/redis/v8"
)

var connMap sync.Map

type RedisClient struct {
	rdb *redis.Client
}

func _initRedis(cfg *config.RedisConnConfig) *RedisClient {
	options := &redis.Options{
		Addr:     fmt.Sprintf("%s:%d", cfg.Host, cfg.Port),
		Password: cfg.Password,
	}
	if cfg.UserName != "" {
		options.Username = cfg.UserName
	}
	rdb := redis.NewClient(options)
	rdbClient := &RedisClient{
		rdb: rdb,
	}
	return rdbClient

}

func GetConnect(cfg *config.RedisConnConfig) (*RedisClient, error) {
	if db, exists := connMap.Load(cfg.Name); exists {
		return db.(*RedisClient), nil
	} else {
		return nil, fmt.Errorf("redis connection not found for key: %s", cfg.Name)
	}
}

func GetConnectOrInit(cfg *config.RedisConnConfig) (*RedisClient, error) {
	defer func() {
		if r := recover(); r != nil {
			slog.Error("panic in GetConnectOrInit", "error", r)
		}
	}()

	key := cfg.Name
	if db, exists := connMap.Load(key); exists {
		if conn, ok := db.(*RedisClient); ok {
			return conn, nil
		} else {
			return nil, fmt.Errorf("invalid type stored in connMap for key %s", key)
		}
	}

	iDriver := _initRedis(cfg)
	connMap.Store(key, iDriver)
	return GetConnect(cfg)
}

func TestConnection(cfg *config.RedisConnConfig) error {
	if cfg == nil {
		return fmt.Errorf("redis connection configuration is nil")
	}
	iDriver := _initRedis(cfg)
	pong, err := iDriver.rdb.Ping(context.Background()).Result()
	if err != nil {
		return fmt.Errorf("failed to get sql.DB from gorm.DB: %w", err)
	}
	slog.Info("Redis connection test successful", "pong", pong)
	if pong != "PONG" {
		return fmt.Errorf("unexpected ping response: %s", pong)
	}

	return nil
}

func (_this *RedisClient) ListDB() (int, error) {

	dbs, err := _this.rdb.ConfigGet(context.Background(), "databases").Result()
	if err != nil {
		return 0, fmt.Errorf("failed to list Redis databases: %w", err)
	}
	slog.Info("Redis databases listed successfully", "dbs", dbs)
	if len(dbs) == 2 {
		fmt.Printf("当前 Redis 实例支持的数据库数为: %s\n", dbs[1])
	} else {
		fmt.Println("无法解析 DB 数量")
	}
	num := dbs[1].(string)
	numInt, _ := strconv.ParseInt(num, 10, 64)
	return int(numInt), nil
}

// GetRecords
func (_this *RedisClient) GetRecords(db int, key string) ([]string, error) {
	val, err := _this.rdb.Do(context.Background(), "LRANGE", key, 0, -1).Result()
	if err != nil {
		return nil, fmt.Errorf("failed to get records from Redis: %w", err)
	}
	slog.Info("Records retrieved successfully", "key", key, "records", val)
	return val.([]string), nil

}
