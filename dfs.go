package main

import (
	"fmt"
	"os"
)

// 记录一个边的信息，包括其 vertex_id, vertex_label, edge_label
type DFS struct {
	from, to, fromlabel, elabel, tolabel int
}

/* 一个图的 DFSCode，DFSCode 相同等价于图同构
   DFSCode 中包含了图的所有边
*/
type DFSCode []DFS

func (dc *DFSCode) push(from, to, fromlabel, elabel, tolabel int) {
	*dc = append(*dc, DFS{from, to, fromlabel, elabel, tolabel})
}

func (dc *DFSCode) pop() {
	if len(*dc) > 0 {
		*dc = (*dc)[:len(*dc)-1]
	}
}

// 从一个图中构建 DFSCode，注意需要保证图 g 连通
func (dc *DFSCode) fromGraph(g *Graph) {
	*dc = make(DFSCode, 0)
	for from, v := range g.VertexArray {
		// 找到在所有 v 的临边中，相连的节点 nei.label 比 v.label 大的边，目的是对无向图中的双连接去重
		edgeList := g.getForwardRoot(&v)
		for _, e := range edgeList {
			assert(e.from == from)
			dc.push(from, e.to, g.VertexArray[e.from].label, e.elabel, g.VertexArray[e.to].label)
		}
	}
}

// 将 DFSCode 转化为一个图
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

// 获得 DFSCode 所代表的图中的节点数
func (dc *DFSCode) nodeCount() int {
	nodeCount := 0
	for i := 0; i < len(*dc); i++ {
		nodeCount = maxint(nodeCount, maxint((*dc)[i].from, (*dc)[i].to)+1)
	}
	return nodeCount
}

// 打印 DFSCode
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

// 在 DFSCode 上获取最右路径
func (dc *DFSCode) buildRMPath() []int {
	rmpath := make([]int, 0)
	old_from := -1
	for i := len(*dc) - 1; i >= 0; i-- {
		if (*dc)[i].from < (*dc)[i].to && (len(rmpath) == 0 || old_from == (*dc)[i].to) {
			rmpath = append(rmpath, i)
			old_from = (*dc)[i].from
		}
	}
	return rmpath
}
