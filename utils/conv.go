package utils

import "strconv"

func ParseInt(str string) int {
	i, _ := strconv.Atoi(str)
	return i
}
