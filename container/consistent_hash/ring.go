package consistent_hash

import (
	"crypto/sha1"
	"io"
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
	Node
	node_id    string
	partition  int
	node_hash  uint64
	entry_hash uint64
	ring       []Node
}

func (r Ring) Lookup(key string, n int) []Node {
	hash := make_hash(key)

	idx := sort.Search(len(r.entries), func(i int) bool {
		return r.entries[i].entry_hash >= hash
	})

	// idx == len(r.entries) when the hash is before the first entry
	// in this case the last entry must be used
	if idx == len(r.entries) {
		idx--
	}

	ring := r.entries[idx].ring
	if 0 < n && n < len(ring) {
		ring = ring[:n]
	}

	return ring
}

func make_hash(s string) uint64 {
	sha := sha1.New()
	io.WriteString(sha, s)
	l := sha.Sum(nil)

	return uint64(l[0])<<56 |
		uint64(l[1])<<48 |
		uint64(l[2])<<40 |
		uint64(l[3])<<32 |
		uint64(l[4])<<24 |
		uint64(l[5])<<16 |
		uint64(l[6])<<8 |
		uint64(l[7])
}
