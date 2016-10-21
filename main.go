package main

import (
	"fmt"
	"os"
	"strconv"
	"time"
)

func CheckReadGraph() {
	assert(len(TRANS) == 10000)
	for _, g := range TRANS {
		g.check()
	}
}

func PrintReadGraph() {
	for _, g := range TRANS {
		g.write(os.Stdout)
	}
}

func PrintDFSCode() {
	var dfs DFSCode
	for _, g := range TRANS {
		dfs.fromGraph(&g)
		dfs.write(os.Stdout)
		os.Stdout.WriteString("\n")
	}
}

func main() {
	var err error
	TRANS, err = BuildGraphFromFile("graph.data")
	assert(err == nil, err)

	if len(os.Args) == 2 {
		if m, err := strconv.Atoi(os.Args[1]); err == nil {
			minsup = m
		}
	}

	t := time.Now()
	Run()
	fmt.Println(time.Now().Sub(t))
}
