# consistent_hash

[![Go Report Card](https://goreportcard.com/badge/github.com/ZYallers/consistent_hash)](https://goreportcard.com/report/github.com/ZYallers/consistent_hash)
[![MIT license](https://img.shields.io/badge/license-MIT-brightgreen.svg)](https://opensource.org/licenses/MIT)
[![Build Status](https://travis-ci.org/ZYallers/consistent_hash.svg?branch=master)](https://travis-ci.org/ZYallers/consistent_hash)
[![Foundation](https://img.shields.io/badge/Golang-Foundation-green.svg)](http://golangfoundation.org)
[![GoDoc](https://pkg.go.dev/badge/github.com/ZYallers/consistent_hash?status.svg)](https://pkg.go.dev/github.com/ZYallers/consistent_hash?tab=doc)
[![Sourcegraph](https://sourcegraph.com/github.com/ZYallers/consistent_hash/-/badge.svg)](https://sourcegraph.com/github.com/ZYallers/consistent_hash?badge)
[![Release](https://img.shields.io/github/release/ZYallers/consistent_hash.svg?style=flat-square)](https://github.com/ZYallers/consistent_hash/releases)
[![TODOs](https://badgen.net/https/api.tickgit.com/badgen/github.com/ZYallers/consistent_hash)](https://www.tickgit.com/browse?repo=github.com/ZYallers/consistent_hash)
[![goproxy.cn](https://goproxy.cn/stats/github.com/ZYallers/consistent_hash/badges/download-count.svg)](https://goproxy.cn)

一致性哈希算法实现数据分布, 它支持动态添加、删除节点, 并引入虚拟节点保证数据在节点间的均匀分布。

## 特点
1. 动态添加、删除节点
2. 引入虚拟节点保证数据均匀分布
3. 自定义hash函数

## 扩容数据分布
```go
func TestConsistentHash_DataDistribution(t *testing.T) {
	ch := NewConsistentHash(300, nil)
	ch.AddNode("db1", "db2", "db3", "db4")

	dataTotal := 50000
	nodesHash := make(map[string]int)
	for i := 0; i < dataTotal; i++ {
		dataKey := "data" + strconv.Itoa(i)
		node := ch.GetNode(dataKey)
		nodesHash[node]++
	}
	for k, v := range nodesHash {
		t.Logf("%s count: %d (%.2f%%)\n", k, v, float64(v)/float64(dataTotal)*100)
	}

	t.Log("--------------扩容后---------------")
	ch.AddNode("db5", "db6")

	countPrint := func(total, incr int) int {
		t.Logf("--------------数据增长%d---------------\n", incr)
		dataTotal := incr + total
		for i := total; i < dataTotal; i++ {
			dataKey := "data" + strconv.Itoa(i)
			node := ch.GetNode(dataKey)
			nodesHash[node]++
		}
		for k, v := range nodesHash {
			t.Logf("%s count: %d (%.2f%%)\n", k, v, float64(v)/float64(dataTotal)*100)
		}
		return dataTotal
	}

	dataTotal = countPrint(dataTotal, 100000)
	dataTotal = countPrint(dataTotal, 300000)
	dataTotal = countPrint(dataTotal, 500000)
}
```
测试结果：
```bash
db1 count: 12220 (24.44%)
db3 count: 13211 (26.42%)
db4 count: 12153 (24.31%)
db2 count: 12416 (24.83%)
--------------扩容后---------------
--------------数据增长100000---------------
db6 count: 15861 (10.57%)
db5 count: 16160 (10.77%)
db1 count: 29222 (19.48%)
db3 count: 31524 (21.02%)
db4 count: 29421 (19.61%)
db2 count: 27812 (18.54%)
--------------数据增长300000---------------
db6 count: 63630 (14.14%)
db5 count: 64501 (14.33%)
db1 count: 80014 (17.78%)
db3 count: 85481 (19.00%)
db4 count: 81780 (18.17%)
db2 count: 74594 (16.58%)
--------------数据增长500000---------------
db5 count: 145391 (15.30%)
db1 count: 164249 (17.29%)
db3 count: 175432 (18.47%)
db4 count: 169568 (17.85%)
db2 count: 152161 (16.02%)
db6 count: 143199 (15.07%)
```
# License
Released under the [MIT License](https://github.com/ZYallers/consistent_hash/blob/master/LICENSE)