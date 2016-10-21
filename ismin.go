package main

import (
	"math"
)

type ProjectedLabel struct {
	fromlabel, elabel, tolabel int
	Projected
}

func (p *ProjectedLabel) rebuild(fromlabel, elabel, tolabel int) {
	p.fromlabel, p.elabel, p.tolabel = fromlabel, elabel, tolabel
	p.Projected = nil
}

func (p1 *ProjectedLabel) Less(p2 *ProjectedLabel) int {
	if p1.fromlabel != p2.fromlabel {
		return p1.fromlabel - p2.fromlabel
	} else if p1.elabel != p2.elabel {
		return p1.elabel - p2.elabel
	} else if p1.tolabel != p2.tolabel {
		return p1.tolabel - p2.tolabel
	} else {
		return 0
	}
}

type ProjectedLabelArray []ProjectedLabel

func (a ProjectedLabelArray) Len() int           { return len(a) }
func (a ProjectedLabelArray) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }
func (a ProjectedLabelArray) Less(i, j int) bool { return a[i].Less(&a[j]) < 0 }

func isMin() bool {
	if len(DFS_CODE) == 1 {
		return true
	}
	var gmin Graph
	DFS_CODE.toGraph(&gmin)

	var pl, pl2 ProjectedLabel
	pl.rebuild(int(math.MaxInt32), int(math.MaxInt32), int(math.MaxInt32))

	for _, v := range gmin.VertexArray {
		edgeList := gmin.getForwardRoot(&v)
		for _, e := range edgeList {
			pl2.rebuild(v.label, e.elabel, gmin.VertexArray[e.to].label)
			if cmp := pl2.Less(&pl); cmp == 0 {
				pl.Projected.push(0, e, nil)
			} else if cmp < 0 {
				pl = pl2
				pl.Projected.push(0, e, nil)
			}
		}
	}

	var dcmin DFSCode
	dcmin.push(0, 1, pl.fromlabel, pl.elabel, pl.tolabel)

	return pl.Projected.isMin(&dcmin, &gmin)
}

func (projected *Projected) isMin(dcmin *DFSCode, gmin *Graph) bool {
	rmpath := dcmin.buildRMPath()
	minlabel := (*dcmin)[0].fromlabel
	maxtoc := (*dcmin)[rmpath[0]].to

	flag := false
	newto := 0

	var pl, pl2 ProjectedLabel
	pl.rebuild(0, int(math.MaxInt32), 0)

	for i := len(rmpath) - 1; flag == false && i >= 1; i-- {
		for n := range *projected {
			cur := &(*projected)[n]
			his := newHistory(gmin, cur)
			if e := gmin.getBackward(his.earray[rmpath[i]], his.earray[rmpath[0]], his); e != nil {
				pl2.rebuild(0, e.elabel, 0)
				if cmp := pl2.Less(&pl); cmp == 0 {
					pl.Projected.push(0, e, cur)
				} else if cmp < 0 {
					pl = pl2
					pl.Projected.push(0, e, cur)
				}
				newto = (*dcmin)[rmpath[i]].from
				flag = true
			}
		}
	}

	if flag {
		dcmin.push(maxtoc, newto, -1, pl.elabel, -1)
		if DFS_CODE[len(*dcmin)-1] != (*dcmin)[len(*dcmin)-1] {
			return false
		}
		return pl.Projected.isMin(dcmin, gmin)
	}

	flag = false
	newfrom := 0
	pl.rebuild(0, int(math.MaxInt32), int(math.MaxInt32))

	for n := range *projected {
		cur := &(*projected)[n]
		his := newHistory(gmin, cur)
		if edges := gmin.getForwardPure(his.earray[rmpath[0]], minlabel, his); len(edges) != 0 {
			flag = true
			newfrom = maxtoc
			for _, e := range edges {
				pl2.rebuild(0, e.elabel, gmin.VertexArray[e.to].label)
				if cmp := pl2.Less(&pl); cmp == 0 {
					pl.Projected.push(0, e, cur)
				} else if cmp < 0 {
					pl = pl2
					pl.Projected.push(0, e, cur)
				}
			}
		}
	}

	for i := 0; flag == false && i < len(rmpath); i++ {
		for n := range *projected {
			cur := &(*projected)[n]
			his := newHistory(gmin, cur)
			if edges := gmin.getForwardRmpath(his.earray[rmpath[i]], minlabel, his); len(edges) != 0 {
				flag = true
				newfrom = (*dcmin)[rmpath[i]].from
				for _, e := range edges {
					pl2.rebuild(0, e.elabel, gmin.VertexArray[e.to].label)
					if cmp := pl2.Less(&pl); cmp == 0 {
						pl.Projected.push(0, e, cur)
					} else if cmp < 0 {
						pl = pl2
						pl.Projected.push(0, e, cur)
					}

				}
			}
		}
	}

	if flag {
		dcmin.push(newfrom, maxtoc+1, -1, pl.elabel, pl.tolabel)
		if DFS_CODE[len(*dcmin)-1] != (*dcmin)[len(*dcmin)-1] {
			return false
		}
		return pl.Projected.isMin(dcmin, gmin)
	}

	return true
}
