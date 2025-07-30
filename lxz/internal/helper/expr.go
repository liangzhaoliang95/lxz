/**
 * @author  zhaoliang.liang
 * @date  2025/7/30 15:11
 */

package helper

func If[T any](cond bool, a, b T) T {
	if cond {
		return a
	}
	return b
}
