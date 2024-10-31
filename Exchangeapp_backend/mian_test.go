package main

import (
	"fmt"
	"regexp"
	"strings"
	"testing"
)

func removeEscapeChars(s string) string {
	return regexp.MustCompile(`\\(.)`).ReplaceAllStringFunc(s, func(m string) string {
		return string([]byte(m)[1:])
	})
}

func Test1(t *testing.T) {
	escapedStr := "This\\nis\\na\\tstring\\nwith\\tescape\\ncharacters\\t\\n"
	unescapedStr := removeEscapeChars(escapedStr)
	fmt.Println(unescapedStr)
	fmt.Println("\\") //输出一个反斜杠
	s := "\\\"你好\""   //   \你好
	fmt.Println(s)
	fmt.Println(strings.ReplaceAll(s, "\\", "")) //把反斜杠换成=

}
