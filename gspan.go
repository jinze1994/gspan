package main

import (
	"fmt"
	"os"
)

type PDFS struct {
	id   int
	edge *Edge
	prev *PDFS
}

type Projected []PDFS

func (p *Projected) push(id int, edge *Edge, prev *PDFS) {
	*p = append(*p, PDFS{id, edge, prev})
}

type Projected_map1 map[int]*Projected
type Projected_map2 map[int]Projected_map1
type Projected_map3 map[int]Projected_map2

func (p Projected_map1) get(key int) *Projected {
	m, ok := p[key]
	if !ok {
		pro := make(Projected, 0)
		m = &pro
		p[key] = m
	}
	return m
}

func (p Projected_map2) get(key int) Projected_map1 {
	m, ok := p[key]
	if !ok {
		m = make(Projected_map1)
		p[key] = m
	}
	return m
}

func (p Projected_map3) get(key int) Projected_map2 {
	m, ok := p[key]
	if !ok {
		m = make(Projected_map2)
		p[key] = m
	}
	return m
}

func Run() {
	// 对每个图的每个点上的前向边，组成 [from.label][edge.label][to.label] 做划分，每个类下保存图和边的位置
	root := make(Projected_map3)
	for gid, g := range TRANS {
		for _, v := range g.VertexArray {
			edgeList := g.getForwardRoot(&v)
			for _, e := range edgeList {
				root.get(v.label).get(e.elabel).get(g.VertexArray[e.to].label).push(gid, e, nil)
			}
		}
	}

	for fromlabel, m2 := range root {
		for edgelabel, m1 := range m2 {
			for tolabel, p := range m1 {
				DFS_CODE.push(0, 1, fromlabel, edgelabel, tolabel)
				sub_mining(*p)
				DFS_CODE.pop()
			}
		}
	}
}

func report(projected Projected, sup int) {
	if maxpat_max > maxpat_min && DFS_CODE.nodeCount() > maxpat_max {
		return
	}
	if maxpat_min > 0 && DFS_CODE.nodeCount() < maxpat_min {
		return
	}
	var g Graph
	DFS_CODE.toGraph(&g)
	os.Stdout.WriteString(fmt.Sprintf("t # %d * %d\n", ID, sup))
	ID++
	g.write(os.Stdout)
	os.Stdout.WriteString("\n")
}

func sub_mining(projected Projected) {
	sup := support(projected)
	if sup < minsup || isMin() == false {
		return
	}
	report(projected, sup)
	if maxpat_max > maxpat_min && DFS_CODE.nodeCount() > maxpat_max {
		return
	}

	rmpath := DFS_CODE.buildRMPath()
	minlabel := DFS_CODE[0].fromlabel
	maxtoc := DFS_CODE[rmpath[0]].to

	newFwdRoot := make(Projected_map3)
	newBckRoot := make(Projected_map2)
	for n := range projected {
		id := projected[n].id
		cur := &projected[n]
		his := newHistory(&TRANS[id], cur)

		for i := len(rmpath) - 1; i >= 1; i-- {
			if e := TRANS[id].getBackward(his.earray[rmpath[i]], his.earray[rmpath[0]], his); e != nil {
				newBckRoot.get(DFS_CODE[rmpath[i]].from).get(e.elabel).push(id, e, cur)
			}
		}

		edges := TRANS[id].getForwardPure(his.earray[rmpath[0]], minlabel, his)
		for _, e := range edges {
			newFwdRoot.get(maxtoc).get(e.elabel).get(TRANS[id].VertexArray[e.to].label).push(id, e, cur)
		}

		for _, rm := range rmpath {
			edges = TRANS[id].getForwardRmpath(his.earray[rm], minlabel, his)
			for _, e := range edges {
				newFwdRoot.get(DFS_CODE[rm].from).get(e.elabel).get(TRANS[id].VertexArray[e.to].label).push(id, e, cur)
			}
		}
	}

	for to, m1 := range newBckRoot {
		for elabel, p := range m1 {
			DFS_CODE.push(maxtoc, to, -1, elabel, -1)
			sub_mining(*p)
			DFS_CODE.pop()
		}
	}

	for from, m2 := range newFwdRoot {
		for elabel, m1 := range m2 {
			for tolabel, p := range m1 {
				DFS_CODE.push(from, maxtoc+1, -1, elabel, tolabel)
				sub_mining(*p)
				DFS_CODE.pop()
			}
		}
	}

	return
}

func support(projected Projected) int {
	n, oid := 0, -1
	for i := range projected {
		if projected[i].id != oid {
			n++
		}
		oid = projected[i].id
	}
	return n
}

var ID int
var TRANS []Graph
var DFS_CODE DFSCode

var minsup int = 8000
var maxpat_min int = 2
var maxpat_max int = 100000
