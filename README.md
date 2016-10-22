gspan
===

[gspan](http://www.cs.ucsb.edu/~xyan/software/gSpan.htm) 是一个频繁子图挖掘算法，代码实现参照了 [gboost](http://www.nowozin.net/sebastian/gboost/)，输出结果与其进行了多轮对比，可保证本程序的正确性。

本程序完全用 [golang](golang.org) 实现。运行速度上，比用 C++ 实现的 [gboost](http://www.nowozin.net/sebastian/gboost/) 慢 3 倍左右。

运行说明
---
```
./gspan [min_support] [min_node] [max_node]
```

程序从标准输入读取图文件，向标准输出打印频繁子图。有以下三个可选参数

* min_support: 输出的频繁子图应满足的最小支持度，默认为 2
* min_node: 输出的频繁子图所具有的最小节点数，默认为 2
* max_node: 输出的频繁子图所具有的最大节点数，默认为极大值

重要数据结构说明
---

- `Graph`: 图的邻接表数据结构，包含一列 `Vertex`，每个 `Vertex` 包含一列 `Edge`。全局变量 `TRANS` 记录了所有图。
	- `getBackward`:	获取图中最右下节点到最右路径上的 Backward 边
	- `getForwardPure`:	获取图中最右下节点引出的所有 Forward 边
	- `getForwardRmpath`:	获取图中最右路径引出的所有 Forward 边
	- `getForwardRoot`:	获取图中某节点引出的有效边（a.label <= a.nei.label)，用于构造 `DFSCode`

- `DFSCode`: `DFSCode` 的本质是图中所有边的信息在某种全序关系定义下的排列，每条边上记录顶点 ID、顶点和边的 Label。图同构等价于 `DFSCode` 相同。全局变量 `DFS_CODE` 记录了当前正在处理（当前栈下）的频繁子图。
	- `fromGraph(g *Graph)`: 将一个图转换为 `DFSCode`
	- `toGraph(g *Graph)`: 将 `DFSCode` 转换为图
	- `buildRMPath() []int`: 在 `DFSCode` 上获取最右路径，保存了最右路径上节点在 `DFSCode` 上的索引

- `PDFS`: `PDFS` 的数据结构是一个链表，其本质代表了深度优先搜索中， `DFSCode` 在搜索栈中在某一个出现位置上的投影(projection)。由于每个 `child DFSCode` 都是在 `parent DFSCode` 上增加一条边的结果，如果将每个图或每个图的 `DFSCode` 保存在搜索栈中就会浪费大量空间，因此栈中只保存增加的边即 `PDFS.edge`，根据 `PDFS.prev` 的链表指针向前寻找可构造出该 `DFSCode` 每一条边的添加顺序。

- `Projected`: `Projected` 最主要的作用是在栈中保存所有的 `PDFS`，它是一个 `PDFS` 的数组。在递归调用的搜索栈中，每次传入一个 `Projected`，代表当前的 `DFSCode` 在所有原图中的“投影”（出现位置及每条边的被添加顺序），然后在所有原图的每个出现位置上尝试拓展相应的边，构造出下一层的很多 `Projected`，然后对这些 `Projected` 依次递归调用。

- `History`: 在递归函数中，根据当前 `DFSCode` 构造出的 `rmpath` 保存在最右路径上的节点索引，我们需要根据这些索引找到原图中最右路径上的边的指针。利用 `PDFS.prev` 的链表指针向搜索栈上方一个个寻找，可构造出整个图上的边被添加的顺序，该顺序与 `DFSCode` 中边的排列顺序相同。这样 `History.earray` 中就恢复出了按照 `DFSCode` 的排列形式在对应出现位置上的所有边的指针，即可再利用 `rmpath` 的索引直接定位到对应出现位置上的边。`History.vertex` 和 `History.edge` 是一个 hash 表，如果该点和边已经出现在 `DFSCode` 中，则相应位置被置为 1。

- `Projected_mapx`: 维护了一条边的类型（两个顶点的 label 和边的 label）到 `Projected` 的映射。使得相同类型的边被一次同时扩展。

算法框架
---

编译运行环境
---
* go version >= go1.7
* Intel(R) Xeon(R) CPU E5-2609 v3 @ 1.90GHz * 6
* memory >= 4GB
* arch linux x86_64 GNU/Linux (optional)

运行结果
---
| 最小支持度 | 频繁子图个数 | 时间开销 |
| --- | --- | --- |
| 5000 | 26 | 6.82s |                                                                                                                                                                                              
| 3000 | 52 | 8.89s |
| 1000 | 455 | 25.78s |
| 800 | 724 | 33.00s |
| 600 | 1235 | 45.64s |
| 400 | 2710 | 1m26s |
| 200 | 10621 | 3m35s |
