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

func str2int(str string) int {
	n, err := strconv.Atoi(str)
	assert(err == nil)
	return n
}
