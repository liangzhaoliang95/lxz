package model

import (
	"sort"
	"strings"
)

type RedisGroupTree struct {
	Name     string                     // 当前节点名称
	Children map[string]*RedisGroupTree // 子节点
	Value    string                     // 可选：叶子节点的值（如果对应 Redis 的值）
	Type     string                     // 可选：值类型（如 string, hash 等）
}

// GetSortedChildren 获取按名称排序的子节点
func (r *RedisGroupTree) GetSortedChildren() []*RedisGroupTree {
	if r.Children == nil {
		return nil
	}

	// 获取所有子节点的名称并排序
	names := make([]string, 0, len(r.Children))
	for name := range r.Children {
		names = append(names, name)
	}

	// 自定义排序：小写字母在前，对应大写字母紧跟其后
	sort.Slice(names, func(i, j int) bool {
		a, b := names[i], names[j]

		// 如果两个字符串都是空字符串，返回false
		if len(a) == 0 && len(b) == 0 {
			return false
		}
		if len(a) == 0 {
			return true
		}
		if len(b) == 0 {
			return false
		}

		// 获取第一个字符进行比较
		charA, charB := rune(a[0]), rune(b[0])

		// 转换为小写进行比较，实现 aAbBccDD 的排序
		lowerA := charA
		lowerB := charB

		// 如果是大写字母，转换为小写进行比较
		if charA >= 'A' && charA <= 'Z' {
			lowerA = charA + 32 // 转换为小写
		}
		if charB >= 'A' && charB <= 'Z' {
			lowerB = charB + 32 // 转换为小写
		}

		// 如果小写字母不同，按小写字母顺序排序
		if lowerA != lowerB {
			return lowerA < lowerB
		}

		// 如果小写字母相同，小写字母在前，大写字母在后
		if charA != charB {
			return charA >= 'a' && charA <= 'z' // 小写字母在前
		}

		// 如果第一个字符完全相同，按整个字符串排序
		return a < b
	})

	// 按排序后的名称返回子节点
	children := make([]*RedisGroupTree, len(names))
	for i, name := range names {
		children[i] = r.Children[name]
	}

	return children
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
}

func NewRedisData(keys []string) *RedisData {
	R := &RedisData{
		Keys: keys,
	}
	R.GroupKeys()
	return R
}
