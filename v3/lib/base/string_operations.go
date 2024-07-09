package base

import (
	"fmt"
	"strings"
)

func PrettyPrintValue(v interface{}) (str string) {
	st_1 := fmt.Sprintf("%v", v)
	indentLevel := 0

	buildString := func(m int) string {
		s := ""

		for i := 0; i < (m * 2); i += 1 {
			s += " "
		}
		return s
	}
	for _, char := range st_1 {

		switch char {
		case '{':
			fallthrough
		case '[':
			indentLevel += 1
			str += string(char) + "\n" + buildString(indentLevel)
		case '}':
			fallthrough
		case ']':
			indentLevel -= 1
			str += "\n" + buildString(indentLevel) + string(char)
		case ' ':
			str += "\n" + buildString(indentLevel)
		default:
			str += string(char)
		}
	}

	return
}


func TrimUrlParams(s string, sub ...interface{}) string {
	st_0 := fmt.Sprintf(s, sub...)
	st_1 := strings.Replace(st_0, "\n", "", 0)
	st_2 := strings.Replace(st_1, " ", "", 0)

	return st_2
}