package main

type History struct {
	earray []*Edge
	edge   []int
	vertex []int
}

func (h *History) hasEdge(id int) bool {
	return h.edge[id] != 0
}

func (h *History) hasVertec(id int) bool {
	return h.vertex[id] != 0
}

func newHistory(g *Graph, e *PDFS) (h *History) {
	h = new(History)
	h.earray = make([]*Edge, 0)
	h.edge = make([]int, g.edge_size)
	h.vertex = make([]int, len(g.VertexArray))

	if e != nil {
		h.earray = append(h.earray, e.edge)
		h.edge[e.edge.id], h.vertex[e.edge.from], h.vertex[e.edge.to] = 1, 1, 1

		for p := e.prev; p != nil; p = p.prev {
			h.earray = append(h.earray, p.edge)
			h.edge[p.edge.id], h.vertex[p.edge.from], h.vertex[p.edge.to] = 1, 1, 1
		}

		n := len(h.earray)
		for i := 0; i < n/2; i++ {
			h.earray[i], h.earray[n-i-1] = h.earray[n-i-1], h.earray[i]
		}
	}
	return h
}
