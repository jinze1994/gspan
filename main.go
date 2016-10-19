package main

import (
	"fmt"
	"os"
)

func main() {
	gs, err := BuildGraphFromFile("graph.data")
	assert(err == nil, err)
	fmt.Println(len(gs))
	for t, g := range gs {
		os.Stdout.WriteString(fmt.Sprintf("t # %d\n", t))
		g.write(os.Stdout)
	}
	os.Stdout.WriteString("t # -1")
}
