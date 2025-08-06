package model

import (
	"fmt"
	"log/slog"
	"lxz/internal/helper"
	"strings"
)

type RedisGroupTree struct {
	Name     string                     // 当前节点名称
	Children map[string]*RedisGroupTree // 子节点
	Value    string                     // 可选：叶子节点的值（如果对应 Redis 的值）
	Type     string                     // 可选：值类型（如 string, hash 等）
}
type RedisData struct {
	Keys []string        `json:"keys,omitempty"` // Redis keys
	Tree *RedisGroupTree `json:"tree,omitempty"` // Redis group tree
}

// splitKey splits a Redis key into parts based on the delimiter.
func splitKey(key string) []string {
	// Assuming the delimiter is ':', you can change it if needed
	parts := make([]string, 0)

	for _, part := range strings.Split(key, ":") {
		if part != "" {
			parts = append(parts, part)
		}
	}
	return parts
}

// 将keys分组获得分组的数据
func (_this *RedisData) GroupKeys() {
	root := &RedisGroupTree{Name: "root"}
	for _, key := range _this.Keys {
		parts := strings.Split(key, ":")
		curr := root
		for _, part := range parts {
			if curr.Children == nil {
				curr.Children = make(map[string]*RedisGroupTree)
			}
			if _, exists := curr.Children[part]; !exists {
				curr.Children[part] = &RedisGroupTree{Name: part}
			}
			curr = curr.Children[part]
		}
		// 最后的 leaf 存储 redis key 的值和类型
		//curr.Value = value
		//curr.Type = typ
	}
	_this.Tree = root
	slog.Info(fmt.Sprintf("RedisData GroupKeys:\n %s", helper.Prettify(_this.Tree)))
}

func NewRedisData(keys []string) *RedisData {
	R := &RedisData{
		Keys: keys,
	}
	R.GroupKeys()
	return R
}
