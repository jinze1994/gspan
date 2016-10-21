gspan
===

[gspan](http://www.cs.ucsb.edu/~xyan/software/gSpan.htm) 是一个频繁子图挖掘算法，代码实现参照了 [gboost](http://www.nowozin.net/sebastian/gboost/)，输出结果与其进行了多轮对比，可保证本程序的正确性。

本程序完全用 golang 实现。运行速度上，比用 C++ 实现的 [gboost](http://www.nowozin.net/sebastian/gboost/) 慢 3 倍左右。

运行环境
---
* go version >= go1.7
* arch linux x86_64 GNU/Linux (optional)

运行说明
---
```bash
$ ./gspan [min_support] [min_node] [max_node]
```
