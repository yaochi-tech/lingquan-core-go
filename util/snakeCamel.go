package util

import "strings"

// ToSnake 将驼峰命名转换为蛇形命名
// 如果遇到数字，则完整的数字作为一部分
func ToSnake(s string) string {
	var res []rune
	var inNumber bool
	for i, c := range s {
		cIsNumber := c >= '0' && c <= '9'
		if i > 0 {
			if c >= 'A' && c <= 'Z' {
				res = append(res, '_')
			} else if cIsNumber && !inNumber {
				res = append(res, '_')
			}
		}
		inNumber = cIsNumber
		res = append(res, c)
	}
	return strings.ToLower(string(res))
}

// ToCamel 将蛇形命名转换为驼峰命名
// 如果转换后的字符串首字母大写，则首字母不会变成小写
// 如果下划线后面是数字，则数字不会变成大写
// 例如：hello_world -> helloWorld
func ToCamel(s string) string {
	var res []rune
	nextUpper := false
	for _, c := range s {
		if c == '_' {
			nextUpper = true
			continue
		}
		if nextUpper {
			nextUpper = false
			// 应该要判断，如果是小写字母，则转换为大写字母
			if c >= 'a' && c <= 'z' {
				res = append(res, c-32)
				continue
			}
			res = append(res, c)
		} else {
			res = append(res, c)
		}
	}
	return string(res)
}
