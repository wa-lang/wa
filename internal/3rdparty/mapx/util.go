// 版权 @2024 凹语言 作者。保留所有权利。

package mapx

func Compare(l, r interface{}) int {
	switch l.(type) {
	case int:
		l, r := l.(int), r.(int)
		switch {
		case l < r:
			return -1
		case l > r:
			return 1
		default:
			return 0
		}
	case string:
		l, r := l.(string), r.(string)
		switch {
		case l < r:
			return -1
		case l > r:
			return 1
		default:
			return 0
		}
	}
	panic("unreachable")
}
