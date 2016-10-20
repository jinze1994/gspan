package main

import (
	"fmt"
	"os"
)

type DFS struct {
	from, to, fromlabel, elabel, tolabel int
}

func (d *DFS) Equal(d2 *DFS) bool {
	if d.from == d2.from && d.to == d2.to && d.fromlabel == d2.fromlabel && d.elabel == d2.elabel && d.tolabel == d2.tolabel {
		return true
	} else {
		return false
	}
}

type DFSCode []DFS

func (dc *DFSCode) push(from, to, fromlabel, elabel, tolabel int) {
	*dc = append(*dc, DFS{from, to, fromlabel, elabel, tolabel})
}

func (dc *DFSCode) pop() {
	if len(*dc) > 0 {
		*dc = (*dc)[:len(*dc)-1]
	}
}

func (dc *DFSCode) fromGraph(g *Graph) {
	*dc = make(DFSCode, 0)
	for from, v := range g.VertexArray {
		edgeList := getForwardRoot(g, &v)
		for _, e := range edgeList {
			assert(e.from == from)
			dc.push(from, e.to, g.VertexArray[e.from].label, e.elabel, g.VertexArray[e.to].label)
		}
	}
}

func (dc *DFSCode) toGraph(g *Graph) {
	g.VertexArray = make(VertexArray, dc.nodeCount())

	for _, dfs := range *dc {
		if dfs.fromlabel != -1 {
			g.VertexArray[dfs.from].label = dfs.fromlabel
		}
		if dfs.tolabel != -1 {
			g.VertexArray[dfs.to].label = dfs.tolabel
		}

		g.VertexArray[dfs.from].push(dfs.from, dfs.to, dfs.elabel)
		g.VertexArray[dfs.to].push(dfs.to, dfs.from, dfs.elabel)
	}

	g.buildEdge()
}

func (dc DFSCode) nodeCount() int {
	nodeCount := 0
	for i := 0; i < len(dc); i++ {
		nodeCount = maxint(nodeCount, maxint(dc[i].from, dc[i].to)+1)
	}
	return nodeCount
}

func (dc DFSCode) write(fout *os.File) {
	if len(dc) == 0 {
		return
	}

	fmt.Fprintf(fout, "(%d) %d (0f%d)", dc[0].fromlabel, dc[0].elabel, dc[0].tolabel)
	for i := 1; i < len(dc); i++ {
		if dc[i].from < dc[i].to {
			fmt.Fprintf(fout, " %d (%df%d)", dc[i].elabel, dc[i].from, dc[i].tolabel)
		} else {
			fmt.Fprintf(fout, " %d (b%d)", dc[i].elabel, dc[i].to)
		}
	}
}
