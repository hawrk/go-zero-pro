// Package tools
/*
 Author: hawrkchen
 Date: 2023/6/27 9:40
 Desc:  一致性哈希算法
*/
package tools

import (
	"hash/crc32"
	"sort"
	"strconv"
)

type Hash func(data []byte) uint32

type HashRing struct {
	hash     Hash
	replicas int
	ring     []int          // hash环 排序
	nodes    map[int]string // key ->ringid    value-> 节点值
}

func NewHashRing(rep int, fn Hash) *HashRing {
	r := &HashRing{
		hash:     fn,
		replicas: rep,
		ring:     nil,
		nodes:    make(map[int]string),
	}
	if r.hash == nil {
		r.hash = crc32.ChecksumIEEE
	}
	return r
}

// Empty 判断 hash环上是否有值
func (r *HashRing) Empty() bool {
	if len(r.ring) == 0 {
		return true
	}
	return false
}

// AddNodes hash环上增加节点
func (r *HashRing) AddNodes(nodes ...string) {
	for _, node := range nodes {
		for i := 0; i < r.replicas; i++ {
			hash := int(r.hash([]byte(strconv.Itoa(i) + node)))
			// 加入哈希环
			r.ring = append(r.ring, hash)
			//设置哈希环 对应节点的映射
			r.nodes[hash] = node
		}
	}
	sort.Ints(r.ring)
}

// Reset 重置哈希环
func (r *HashRing) Reset(nodes ...string) {
	r.ring = nil
	r.nodes = map[int]string{}
	// 重置
	r.AddNodes(nodes...)
}

func (r *HashRing) GetNode(key string) string {
	if r.Empty() {
		return ""
	}
	hash := int(r.hash([]byte(key)))

	idx := sort.Search(len(r.ring), func(i int) bool {
		return r.ring[i] >= hash
	})
	if idx == len(r.ring) {
		idx = 0
	}

	return r.nodes[r.ring[idx]]
}
