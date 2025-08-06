/**
 * @author  zhaoliang.liang
 * @date  2025/8/6 14:12
 */

package redis_drivers

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"lxz/internal/config"
	"strconv"
	"sync"

	"github.com/go-redis/redis/v8"
)

var connMap sync.Map

func connMapKey(name string, db int) string {
	return fmt.Sprintf("%s@%d", name, db)
}

type RedisClient struct {
	dbNum  int // 数据库编号
	rdb    *redis.Client
	config *config.RedisConnConfig // Redis连接配置
}

func _initRedis(cfg *config.RedisConnConfig, dbNum int) *RedisClient {
	options := &redis.Options{
		Addr:     fmt.Sprintf("%s:%d", cfg.Host, cfg.Port),
		Password: cfg.Password,
		DB:       dbNum,
	}
	if cfg.UserName != "" {
		options.Username = cfg.UserName
	}
	rdb := redis.NewClient(options)
	rdbClient := &RedisClient{
		rdb:    rdb,
		config: cfg,
		dbNum:  dbNum,
	}
	return rdbClient

}

func GetConnect(cfg *config.RedisConnConfig, dbNum int) (*RedisClient, error) {
	if db, exists := connMap.Load(connMapKey(cfg.Name, dbNum)); exists {
		return db.(*RedisClient), nil
	} else {
		return nil, fmt.Errorf("redis connection not found for key: %s", cfg.Name)
	}
}

func GetConnectOrInit(cfg *config.RedisConnConfig, dbNum int) (*RedisClient, error) {
	defer func() {
		if r := recover(); r != nil {
			slog.Error("panic in GetConnectOrInit", "error", r)
		}
	}()

	key := connMapKey(cfg.Name, dbNum)
	if db, exists := connMap.Load(key); exists {
		if conn, ok := db.(*RedisClient); ok {
			return conn, nil
		} else {
			return nil, fmt.Errorf("invalid type stored in connMap for key %s", key)
		}
	}

	iDriver := _initRedis(cfg, dbNum)
	connMap.Store(key, iDriver)
	return GetConnect(cfg, dbNum)
}

func TestConnection(cfg *config.RedisConnConfig) error {
	if cfg == nil {
		return fmt.Errorf("redis connection configuration is nil")
	}
	iDriver := _initRedis(cfg, 0)
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
		slog.Info("Redis ConfigGet response", "dbs", dbs)
	} else {
		return 0, fmt.Errorf("unexpected response from Redis ConfigGet: %v", dbs)
	}
	num := dbs[1].(string)
	numInt, _ := strconv.ParseInt(num, 10, 64)
	return int(numInt), nil
}

// GetRecords 获取指定数据库的记录
func (_this *RedisClient) GetRecords(key string) ([]string, error) {
	var cursor uint64 = 0
	var allKeys = make([]string, 0)
	var search = "*"
	if key != "" {
		search = "*" + key + "*"
	}
	for {
		keys, nextCursor, err := _this.rdb.Scan(context.Background(), cursor, search, 100).Result()
		if err != nil {
			return nil, fmt.Errorf("failed to scan Redis keys: %w", err)
		}
		allKeys = append(allKeys, keys...)

		if nextCursor == 0 {
			break
		}
		cursor = nextCursor
	}

	slog.Info("Keys retrieved successfully", "db", _this.dbNum, "search", search, "keys", allKeys)
	return allKeys, nil

}

// GetDBKeyNum 获取指定数据库的键数量
func (_this *RedisClient) GetDBKeyNum() (int64, error) {
	val, err := _this.rdb.DBSize(context.Background()).Result()
	if err != nil {
		return 0, fmt.Errorf("failed to get key count from Redis: %w", err)
	}
	slog.Info("Key count retrieved successfully", "db", _this.dbNum, "key_count", val)
	return val, nil
}

// GetKeyType 获取指定键的类型
func (_this *RedisClient) GetKeyType(key string) (string, error) {
	val, err := _this.rdb.Type(context.Background(), key).Result()
	if err != nil {
		return "", fmt.Errorf("failed to get key type from Redis: %w", err)
	}
	slog.Info("Key type retrieved successfully", "key", key, "type", val)
	return val, nil
}

// GetKeyValue 获取指定键的值
func (_this *RedisClient) GetKeyValue(key string) (string, error) {
	val, err := _this.rdb.Get(context.Background(), key).Result()
	if err != nil {
		if errors.Is(err, redis.Nil) {
			return "", fmt.Errorf("key %s does not exist", key)
		}
		return "", fmt.Errorf("failed to get key value from Redis: %w", err)
	}
	return val, nil
}

// GetKeyTTL 获取指定键的生存时间
func (_this *RedisClient) GetKeyTTL(key string) (int64, error) {
	ttl, err := _this.rdb.TTL(context.Background(), key).Result()
	if err != nil {
		return 0, fmt.Errorf("failed to get key TTL from Redis: %w", err)
	}
	if ttl < 0 {
		return -1, nil // -1表示键没有设置过期时间
	}
	slog.Info("Key TTL retrieved successfully", "key", key, "ttl", ttl.Seconds())
	return int64(ttl.Seconds()), nil
}
