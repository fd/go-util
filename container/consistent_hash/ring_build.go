package consistent_hash

import (
	"crypto/sha1"
	"fmt"
	"sort"
)

func New(l []Node, buckets int) Ring {
	e := wrap_nodes(l, buckets)
	e = make_entrie_hashes(e)
	e = sort_entries(e)
	e = remove_duplicate_entries(e)
	e = make_entry_rings(e, l)
	return Ring{l, e}
}

func wrap_nodes(l []Node, buckets int) []entry_t {
	o := make([]entry_t, len(l)*buckets)

	for i, n := range l {
		node_id := n.HashID()

		for j := 0; j < buckets; j++ {
			o[i*buckets+j] = entry_t{
				Node:      n,
				node_id:   node_id,
				partition: j,
			}
		}
	}

	return o
}

func make_entrie_hashes(entries []entry_t) []entry_t {
	var (
		b   [20]byte
		l   []byte
		sha = sha1.New()
	)

	for i, e := range entries {
		l = b[:0]

		sha.Reset()
		fmt.Fprintln(sha, e.node_id)
		l = sha.Sum(l)

		e.node_hash = uint64(l[0])<<56 |
			uint64(l[1])<<48 |
			uint64(l[2])<<40 |
			uint64(l[3])<<32 |
			uint64(l[4])<<24 |
			uint64(l[5])<<16 |
			uint64(l[6])<<8 |
			uint64(l[7])

		l = b[:0]
		fmt.Fprintln(sha, e.partition)
		l = sha.Sum(l)

		e.entry_hash = uint64(l[0])<<56 |
			uint64(l[1])<<48 |
			uint64(l[2])<<40 |
			uint64(l[3])<<32 |
			uint64(l[4])<<24 |
			uint64(l[5])<<16 |
			uint64(l[6])<<8 |
			uint64(l[7])

		entries[i] = e
	}

	return entries
}

func sort_entries(l []entry_t) []entry_t {
	sort.Sort(entry_sorter(l))
	return l
}

func remove_duplicate_entries(l []entry_t) []entry_t {
	last := l[len(l)-1].node_id
	o := make([]entry_t, 0, len(l))

	for _, e := range l {
		if last == e.node_id {
			continue
		}

		o = append(o, e)
		last = e.node_id
	}

	return o
}

func make_entry_rings(entries []entry_t, nodes []Node) []entry_t {
	var (
		l_ring_len = len(nodes)
		g_ring_len = len(entries)
		known_a    = make(map[uint64]int, l_ring_len)
		known_b    = make(map[uint64]int, l_ring_len)
		o          = make([]Node, len(entries)*l_ring_len)
		i          = 0
	)

	// build first l_ring
	for _, e := range entries {
		if _, p := known_a[e.node_hash]; !p {
			o[i] = e.Node
			known_a[e.node_hash] = i
			i++
		}
	}
	entries[0].ring = o[:l_ring_len]
	l := o[:l_ring_len]

	// build other rings
	for i := g_ring_len - 1; i > 0; i-- {
		e := entries[i]
		r := o[l_ring_len*i : l_ring_len*(i+1)]

		idx := known_a[e.node_hash]
		if idx == 0 {
			copy(r, l)
		} else {
			r[0] = l[idx]
			copy(r[1:idx+1], l[:idx])
			if idx+1 < l_ring_len {
				copy(r[idx+1:], l[idx+1:])
			}

			for h, i := range known_a {
				if i == idx {
					known_b[h] = 0
				} else if i < idx {
					known_b[h] = i + 1
				} else {
					known_b[h] = i
				}
			}

			known_a, known_b = known_b, known_a
		}

		e.ring = r
		entries[i] = e
		r = r
	}

	return entries
}

type entry_sorter []entry_t

func (s entry_sorter) Len() int           { return len(s) }
func (s entry_sorter) Less(i, j int) bool { return s[i].entry_hash < s[j].entry_hash }
func (s entry_sorter) Swap(i, j int)      { s[i], s[j] = s[j], s[i] }
