package main

import (
	"bufio"
	"fmt"
	"io"
	"os"
)

type Edge struct {
	from, to, elabel, id int
}

type Vertex struct {
	label int
	edge  []Edge
}

func (v *Vertex) push(from, to, elabel int) {
	v.edge = append(v.edge, Edge{from, to, elabel, 0})
}

type VertexArray []Vertex

type Graph struct {
	VertexArray
	edge_size int
}

func (g *Graph) buildEdge() {
	tmp := make(map[string]int)
	id := 0
	for from, v := range g.VertexArray {
		for _, e := range v.edge {
			var buf string
			if from <= e.to {
				buf = fmt.Sprintf("%d %d %d", from, e.to, e.elabel)
			} else {
				buf = fmt.Sprintf("%d %d %d", e.to, from, e.elabel)
			}
			if v, ok := tmp[buf]; ok {
				e.id = v
			} else {
				e.id = id
				tmp[buf] = id
				id++
			}
		}
	}
	g.edge_size = id
}

func (g *Graph) read(br *bufio.Reader) (eof bool) {
	eof = false
	g.VertexArray = make(VertexArray, 0)
	g.edge_size = 0
	for {
		line, err := br.ReadString('\n')
		assert(err == nil || err == io.EOF, err)
		var from, to, label, id int
		if err == io.EOF {
			eof = true
			break
		} else if line[0] == 't' {
			break
		} else if line[0] == 'v' {
			n, err := fmt.Sscanf(line, "v%d%d", &id, &label)
			assert(err == nil && n == 2, "读取节点格式错误")
			assert(id == len(g.VertexArray), "读取节点顺序错误")
			g.VertexArray = append(g.VertexArray, Vertex{label, make([]Edge, 0)})
		} else if line[0] == 'e' {
			n, err := fmt.Sscanf(line, "e%d%d%d", &from, &to, &label)
			assert(err == nil && n == 3, "读取边格式错误")
			g.VertexArray[from].push(from, to, label)
			g.VertexArray[to].push(to, from, label)
		} else {
			assert(false, "未知输入标签")
		}
	}

	g.buildEdge()
	return eof
}

func (g *Graph) write(fout *os.File) {
	estr := make([]string, 0)
	for from, v := range g.VertexArray {
		fout.WriteString(fmt.Sprintf("v %d %d\n", from, v.label))
		for _, e := range v.edge {
			if from < e.to {
				estr = append(estr, fmt.Sprintf("e %d %d %d\n", from, e.to, e.elabel))
			}
		}
	}
	for _, str := range estr {
		fout.WriteString(str)
	}
}

func BuildGraphFromFile(fileName string) (gs []Graph, err error) {
	fin, err := os.Open(fileName)
	if err != nil {
		return nil, err
	}
	defer fin.Close()

	buf := bufio.NewReader(fin)
	_, err = buf.ReadString('\n')
	assert(err == nil)

	gs = make([]Graph, 0)
	for {
		var g Graph
		eof := g.read(buf)
		gs = append(gs, g)
		if eof == true {
			break
		}
	}

	return gs, nil
}

func getForwardRoot(g *Graph, v *Vertex) (edgeList []*Edge) {
	edgeList = make([]*Edge, 0)
	for i, e := range v.edge {
		assert(e.to >= 0 && e.to < len(g.VertexArray))
		if v.label <= g.VertexArray[e.to].label {
			edgeList = append(edgeList, &v.edge[i])
		}
	}
	return edgeList
}
