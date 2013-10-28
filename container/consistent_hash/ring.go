package consistent_hash

import (
	"hash/crc32"
	"sort"
)

type Ring struct {
	nodes   []Node
	entries []entry_t
}

type Node interface {
	// A unique id for this node like an IP:PORT
	HashID() string
}

type entry_t struct {
	node_idx   uint8
	entry_hash uint32
	ring       []uint8
}

func (r *Ring) MakeBuffer(n int) []Node {
	l := len(r.nodes)

	if n < 1 {
		n = l
	} else if n > l {
		n = l
	}

	return make([]Node, 0, n)
}

func (r Ring) Lookup(key []byte, b []Node) []Node {
	hash := crc32.ChecksumIEEE(key)

	idx := sort.Search(len(r.entries), func(i int) bool {
		return r.entries[i].entry_hash >= hash
	})

	// idx == len(r.entries) when the hash is before the first entry
	// in this case the last entry must be used
	if idx == len(r.entries) {
		idx--
	}

	ring := r.entries[idx].ring
	ring_len := len(ring)
	n := cap(b)

	if n > ring_len {
		n = ring_len
	} else if n < ring_len {
		ring = ring[:n]
	}

	b = b[:n]

	for i, idx := range ring {
		b[i] = r.nodes[idx]
	}

	return b
}
