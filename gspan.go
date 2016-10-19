package main

type PDFS struct {
	id   int
	edge *Edge
	prev *PDFS
}

type Projected []PDFS

func (p Projected) push(id int, edge *Edge, prev *PDFS) {
	p = append(p, PDFS{id, edge, prev})
}

type Projected_map1 map[int]Projected
type Projected_map2 map[int]Projected_map1
type Projected_map3 map[int]Projected_map2

func (p Projected_map1) get(key int) Projected {
	m, ok := p[key]
	if !ok {
		m = make(Projected, 0)
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

func run() {
	gs, err := BuildGraphFromFile("graph.data")
	assert(err == nil, err)

	// 对每个图的每个点上的前向边，组成 [from.label][edge.label][to.label] 做划分，每个类下保存图和边的位置
	root := make(Projected_map3)
	for gid, g := range gs {
		for from, v := range g.VertexArray {
			edgeList := getForwardRoot(&g, &v)
			for _, e := range edgeList {
				root.get(g.VertexArray[from].label).get(e.elabel).get(g.VertexArray[e.to].label).push(gid, e, nil)
			}
		}
	}

	dfscode := make(DFSCode, 0)
	for fromlabel, m2 := range root {
		for edgelabel, m1 := range m2 {
			for tolabel, p := range m1 {
				dfscode.push(0, 1, fromlabel, edgelabel, tolabel)
				project(p)
				dfscode.pop()
			}
		}
	}
}

func project(projected Projected) {
}
