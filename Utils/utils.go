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
