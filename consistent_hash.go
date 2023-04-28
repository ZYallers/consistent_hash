package consistent_hash

import (
	"github.com/spaolacci/murmur3"
	"sort"
	"strconv"
)

type DbNode struct {
	name string
	hash uint32
}

type HashFunc func([]byte) uint32

// 一致性哈希算法实现数据分布, 它支持动态添加、删除节点, 并引入虚拟节点保证数据在节点间的均匀分布
type ConsistentHash struct {
	// 虚拟节点数量
	replicas int

	// hash函数, 默认为"murmur3.Sum32"
	// 已匹配支持的hash函数如下:
	// 1. "hash/crc32" crc32.ChecksumIEEE
	// 2. "https://github.com/spaolacci/murmur3" murmur3.Sum32
	// 3. "https://github.com/cespare/xxhash" xxhash.Sum64String
	hashFn HashFunc

	// 节点列表
	nodes []*DbNode
}

func NewConsistentHash(replicas int, hashFn HashFunc) *ConsistentHash {
	if hashFn == nil {
		hashFn = murmur3.Sum32
	}
	return &ConsistentHash{
		hashFn:   hashFn,
		replicas: replicas,
		nodes:    make([]*DbNode, 0),
	}
}

func (ch *ConsistentHash) AddNode(names ...string) {
	if len(names) == 0 {
		return
	}
	for _, name := range names {
		var existBool bool
		for _, node := range ch.nodes {
			if node.name == name {
				existBool = true
				break
			}
		}
		if existBool {
			continue
		}
		for i := 0; i < ch.replicas; i++ {
			hash := ch.hashFn([]byte(name + strconv.Itoa(i)))
			node := &DbNode{name: name, hash: hash}
			ch.nodes = append(ch.nodes, node)
		}
	}
	sort.Slice(ch.nodes, func(i, j int) bool {
		return ch.nodes[i].hash < ch.nodes[j].hash
	})
}

func (ch *ConsistentHash) RemoveNode(names ...string) {
	if len(names) == 0 {
		return
	}
	newNodes := make([]*DbNode, 0)
	for _, name := range names {
		for _, node := range ch.nodes {
			if node.name != name {
				newNodes = append(newNodes, node)
			}
		}
	}
	ch.nodes = newNodes
}

func (ch *ConsistentHash) GetNode(key string) string {
	hash := ch.hashFn([]byte(key))
	index := sort.Search(len(ch.nodes), func(i int) bool {
		return ch.nodes[i].hash >= hash
	})
	if index == len(ch.nodes) {
		index = 0
	}
	return ch.nodes[index].name
}
