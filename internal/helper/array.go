/**
 * @author  zhaoliang.liang
 * @date  2025/8/7 14:34
 */

package helper

func Contains[T string | int | int64](list []T, key T) bool {
	for _, item := range list {
		if item == key {
			return true
		}
	}
	return false
}
