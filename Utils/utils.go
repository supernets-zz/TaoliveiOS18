package Utils

import (
	"fmt"
	"math/rand"
	"regexp"
)

var R *rand.Rand

func ExtractNumber(s string) (int, error) {
	re, err := regexp.Compile(`(\d+)`)
	if err != nil {
		return 0, err
	}
	matches := re.FindStringSubmatch(s)
	if len(matches) == 0 {
		return 0, fmt.Errorf("no number found in the string")
	}
	var number int
	_, err = fmt.Sscan(matches[0], &number)
	if err != nil {
		return 0, err
	}
	return number, nil
}

func IsNumeric(str string) bool {
	// 正则表达式，匹配数字
	numericRegex := regexp.MustCompile(`^\d+$`)
	return numericRegex.MatchString(str)
}

func IndexOf(array []string, item string) int {
	for i := 0; i < len(array); i++ {
		if array[i] == item {
			return i
		}
	}
	return -1
}
