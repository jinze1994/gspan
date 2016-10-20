package main

import ()

func TestReadGraph() {
	assert(len(TRANS) == 10000)

	for _, g := range TRANS {
		// g.write(os.Stdout)
		g.check()
	}
}

func main() {
	var err error
	TRANS, err = BuildGraphFromFile("graph.data")
	assert(err == nil, err)

	TestReadGraph()
}
