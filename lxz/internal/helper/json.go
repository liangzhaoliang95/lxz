package helper

import (
	"encoding/json"
)

func Prettify(i interface{}) string {
	// 判断i的类型 如果是string 先转map在转字符串
	if str, ok := i.(string); ok {
		var m interface{}
		if err := json.Unmarshal([]byte(str), &m); err != nil {
			return str // 如果转换失败，返回原始字符串
		}
		i = m // 将i转换为map
	}

	resp, _ := json.MarshalIndent(i, "", "  ")
	return string(resp)
}
