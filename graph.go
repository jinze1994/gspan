package main

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"sort"
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

// 赋予每条边一个唯一的标号 id
func (g *Graph) buildEdge() {
	tmp := make(map[string]int)
	id := 0
	for from, v := range g.VertexArray {
		for i, e := range v.edge {
			var buf string
			if from <= e.to {
				buf = fmt.Sprintf("%d %d %d", from, e.to, e.elabel)
			} else {
				buf = fmt.Sprintf("%d %d %d", e.to, from, e.elabel)
			}
			if l, ok := tmp[buf]; ok {
				g.VertexArray[from].edge[i].id = l
				e.id = l
			} else {
				g.VertexArray[from].edge[i].id = id
				e.id = id
				tmp[buf] = id
				id++
			}
		}
	}
	g.edge_size = id
}

// 从输入中读取一个图
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

// 打印图
func (g *Graph) write(fout *os.File) {
	mstr := make(map[string]bool)
	for from, v := range g.VertexArray {
		fout.WriteString(fmt.Sprintf("v %d %d\n", from, v.label))
		for _, e := range v.edge {
			if from <= e.to {
				mstr[fmt.Sprintf("e %d %d %d\n", from, e.to, e.elabel)] = true
			} else {
				mstr[fmt.Sprintf("e %d %d %d\n", e.to, from, e.elabel)] = true
			}
		}
	}
	estr := make([]string, 0, len(mstr))
	for str, _ := range mstr {
		estr = append(estr, str)
	}
	sort.Strings(estr)
	for _, str := range estr {
		fout.WriteString(str)
	}
}

// 检查图的构造合理性
func (g *Graph) check() {
	b := make([][]int, len(g.VertexArray))
	for i := range b {
		b[i] = make([]int, len(g.VertexArray))
	}
	for from, v := range g.VertexArray {
		assert(v.label >= 2)
		eid := make(map[int]bool)
		for _, e := range v.edge {
			assert(e.from == from)
			assert(eid[e.id] == false)
			eid[e.id] = true
			assert(b[e.from][e.to] == 0)
			b[e.from][e.to] = e.elabel
		}
	}
	for i := range b {
		for j := i; j < len(b[i]); j++ {
			assert(b[i][j] == b[j][i])
		}
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

func (g *Graph) getForwardRmpath(edge *Edge, minlabel int, his *History) (edgeList []*Edge) {
	edgeList = make([]*Edge, 0)
	assert(edge.from >= 0 && edge.from < len(g.VertexArray))
	assert(edge.to >= 0 && edge.to < len(g.VertexArray))
	tolabel := g.VertexArray[edge.to].label

	for i, e := range g.VertexArray[edge.from].edge {
		tolabel2 := g.VertexArray[e.to].label
		if edge.to == e.to || minlabel > tolabel2 || his.hasVertec(e.to) {
			continue
		}
		if edge.elabel < e.elabel || (edge.elabel == e.elabel && tolabel <= tolabel2) {
			edgeList = append(edgeList, &g.VertexArray[edge.from].edge[i])
		}
	}

	return edgeList
}

func (g *Graph) getForwardPure(edge *Edge, minlabel int, his *History) (edgeList []*Edge) {
	edgeList = make([]*Edge, 0)
	assert(edge.from >= 0 && edge.from < len(g.VertexArray))
	assert(edge.to >= 0 && edge.to < len(g.VertexArray))

	for i, e := range g.VertexArray[edge.to].edge {
		if minlabel > g.VertexArray[e.to].label || his.hasVertec(e.to) {
			continue
		}
		edgeList = append(edgeList, &g.VertexArray[edge.to].edge[i])
	}

	return edgeList
}

func (g *Graph) getForwardRoot(v *Vertex) (edgeList []*Edge) {
	edgeList = make([]*Edge, 0)
	for i, e := range v.edge {
		assert(e.to >= 0 && e.to < len(g.VertexArray))
		if v.label <= g.VertexArray[e.to].label {
			edgeList = append(edgeList, &v.edge[i])
		}
	}
	return edgeList
}

func (g *Graph) getBackward(e1, e2 *Edge, his *History) *Edge {
	if e1 == e2 {
		return nil
	}
	assert(e1.from >= 0 && e1.from < len(g.VertexArray))
	assert(e1.to >= 0 && e1.to < len(g.VertexArray))
	assert(e2.from >= 0 && e2.from < len(g.VertexArray))
	assert(e2.to >= 0 && e2.to < len(g.VertexArray))

	for i, nei := range g.VertexArray[e2.to].edge {
		if his.hasEdge(nei.id) || nei.to != e1.from {
			continue
		}
		if e1.elabel < nei.elabel ||
			(e1.elabel == nei.elabel && g.VertexArray[e1.to].label <= g.VertexArray[e2.to].label) {
			return &g.VertexArray[e2.to].edge[i]
		}
	}

	return nil
}
