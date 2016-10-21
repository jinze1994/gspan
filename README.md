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

- `Graph`: 图的邻接表数据结构，包含一列 `Vertex`，每个 `Vertex` 又回包含一列 `Edge`。全局变量 `TRANS` 记录了所有图。
	- `getForwardRmpath`
	- `getForwardPure`
	- `getBackward`
	- `getForwardRoot`

- `DFSCode`: DFSCode 的本质是图中所有边的信息在某种全序关系定义下的排列，每条边上记录顶点 ID、顶点和边的 Label。图同构等价于 DFSCode 相同。全局变量 `DFS_CODE` 记录了当前正在处理（当前栈下）的频繁子图。
	- `fromGraph(g *Graph)` 将一个图转换为 DFSCode
	- `toGraph(g *Graph)` 将 DFSCode 转换为图
	- `buildRMPath() []int` 在 DFSCode 上获取最右路径

- `Projected`: 

- `History`: 在递归过程中，


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
