package consistent_hash

import (
	"hash/crc32"
	"sort"
)

func New(l []Node, buckets uint16) Ring {
	e := wrap_nodes(l, buckets)
	e = sort_entries(e)
	e = make_entry_rings(e, len(l))
	return Ring{l, e}
}

func wrap_nodes(l []Node, buckets uint16) []entry_t {
	var (
		o = make([]entry_t, len(l)*int(buckets))
		b [1024]byte
	)

	for i, n := range l {
		node_id := n.HashID()
		node_id_bytes := b[:len(node_id)+3]
		copy(node_id_bytes[2:], node_id+"•")

		for j := 0; j < int(buckets); j++ {
			node_id_bytes[0] = byte(uint16(j)>>8 | 0xFF)
			node_id_bytes[1] = byte(uint16(j) | 0xFF)

			o[i*int(buckets)+j] = entry_t{
				node_idx:   uint8(i),
				entry_hash: crc32.ChecksumIEEE(node_id_bytes),
			}
		}
	}

	return o
}

func sort_entries(l []entry_t) []entry_t {
	sort.Sort(entry_sorter(l))
	return l
}

func make_entry_rings(entries []entry_t, l_ring_len int) []entry_t {
	var (
		g_ring_len = len(entries)
		o          = make([]uint8, len(entries)*l_ring_len)
		i          = 0
	)

	// build first l_ring
	l := o[:l_ring_len]
FIRST_RING:
	for _, e := range entries {
		for _, node_idx := range l {
			if node_idx == e.node_idx {
				continue FIRST_RING
			}
		}

		l[i] = e.node_idx
		i++
	}
	entries[0].ring = l

	// build other rings
	for i := g_ring_len - 1; i > 0; i-- {
		e := entries[i]
		r := o[l_ring_len*i : l_ring_len*(i+1)]

		if l[0] == e.node_idx {
			copy(r, l)
		} else {
			// find idx
			idx := 0
			node_idx := uint8(0)
			for idx, node_idx = range l {
				if node_idx == e.node_idx {
					break
				}
			}

			r[0] = l[idx]
			copy(r[1:idx+1], l[:idx])
			if idx+1 < l_ring_len {
				copy(r[idx+1:], l[idx+1:])
			}
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
