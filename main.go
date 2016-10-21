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

	var tmp int
	switch len(os.Args) {
	case 4:
		if tmp, err = strconv.Atoi(os.Args[3]); err != nil {
			maxpat_max = tmp
		}
		fallthrough
	case 3:
		if tmp, err = strconv.Atoi(os.Args[2]); err != nil {
			maxpat_min = tmp
		}
		fallthrough
	case 2:
		if tmp, err = strconv.Atoi(os.Args[1]); err != nil {
			minsup = tmp
		}
	}

	t := time.Now()
	Run()
	fmt.Println(time.Now().Sub(t))
}
