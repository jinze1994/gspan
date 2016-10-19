package main

import (
	"fmt"
	"strconv"
)

func assert(ok bool, args ...interface{}) {
	if !ok {
		panic(fmt.Sprint(args...))
	}
}

func maxint(a, b int) int {
	if a > b {
		return a
	} else {
		return b
	}
}

func minint(a, b int) int {
	if a < b {
		return a
	} else {
		return b
	}
}

func str2int(str string) int {
	n, err := strconv.Atoi(str)
	assert(err == nil)
	return n
}
