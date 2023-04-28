package consistent_hash

import (
	"github.com/spaolacci/murmur3"
	"hash/crc32"
	"strconv"
	"testing"
)

func TestNewConsistentHash(t *testing.T) {
	ch := NewConsistentHash(5, nil)
	ch.AddNode("db1", "db2")
	if ch.replicas != 5 {
		t.Error("replicas error")
	}
	if len(ch.nodes) != 10 {
		t.Error("add node error")
	}
	if s := "test"; ch.hashFn([]byte(s)) != murmur3.Sum32([]byte(s)) {
		t.Error("hash func error")
	}
	ch.hashFn = crc32.ChecksumIEEE
	if s := "test"; ch.hashFn([]byte(s)) != crc32.ChecksumIEEE([]byte(s)) {
		t.Error("hash func(crc32.ChecksumIEEE) error")
	}
}

func TestConsistentHash_AddNode(t *testing.T) {
	ch := NewConsistentHash(5, nil)
	ch.AddNode("db1", "db2", "db3")
	if len(ch.nodes) != 15 {
		t.Error("add node error")
	}
	var preHash uint32
	for _, v := range ch.nodes {
		if v.hash <= preHash {
			t.Error("hash error")
		}
		//t.Logf("No: %d, NodeName: %s, NodeHash: %d\n", k+1, v.name, v.hash)
	}
}

func TestConsistentHash_GetNode(t *testing.T) {
	ch := NewConsistentHash(5, nil)
	ch.AddNode("db1", "db2", "db3")
	if len(ch.nodes) != 15 {
		t.Error("add node error")
	}
	nodeName := ch.GetNode("data1")
	if nodeName != "db1" && nodeName != "db2" && nodeName != "db3" {
		t.Error("get node error")
	}
}

func TestConsistentHash_RemoveNode(t *testing.T) {
	ch := NewConsistentHash(5, nil)
	ch.AddNode("db1", "db2", "db3")
	if len(ch.nodes) != 15 {
		t.Error("add node error")
	}
	ch.RemoveNode("db2")
	if len(ch.nodes) != 10 {
		t.Error("remove node error")
	}
	ch.RemoveNode("db4")
	if len(ch.nodes) != 10 {
		t.Error("remove node error")
	}
}

// 测试扩容数据分布
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

// 测试缩容数据分布
func TestConsistentHash_DataShrinkageCapacity(t *testing.T) {
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

	t.Log("--------------缩容后---------------")
	ch.RemoveNode("db2")
	nodesHash = make(map[string]int)
	for i := 0; i < dataTotal; i++ {
		dataKey := "data" + strconv.Itoa(i)
		node := ch.GetNode(dataKey)
		nodesHash[node]++
	}
	for k, v := range nodesHash {
		t.Logf("%s count: %d (%.2f%%)\n", k, v, float64(v)/float64(dataTotal)*100)
	}
}
