package consistent_hash

import (
	"fmt"
	"testing"
	"testing/quick"
)

func TestBuild(t *testing.T) {
	nodes := build_nodes(16)
	ring := New(nodes, 2)

	t.Log(ring.Lookup([]byte("hello"), ring.MakeBuffer(3)))
	t.Log(ring.Lookup([]byte("hello"), ring.MakeBuffer(-1)))

	buf := ring.MakeBuffer(-1)
	f := func(k []byte) bool {
		nodes := ring.Lookup(k, buf)

		if len(nodes) != 16 {
			return false
		}

		for _, n := range nodes {
			if n == nil {
				return false
			}
		}

		return true
	}

	if e := quick.Check(f, nil); e != nil {
		t.Fatal(e)
	}

	buf = ring.MakeBuffer(3)
	f = func(k []byte) bool {
		nodes := ring.Lookup(k, buf)

		if len(nodes) != 3 {
			return false
		}

		for _, n := range nodes {
			if n == nil {
				return false
			}
		}

		return true
	}

	if e := quick.Check(f, nil); e != nil {
		t.Fatal(e)
	}
}

func BenchmarkLookup_128_25(b *testing.B) {
	nodes := build_nodes(128)
	ring := New(nodes, 25)
	k := []byte("hello")
	b.ResetTimer()

	buf := ring.MakeBuffer(-1)

	for i := 0; i < b.N; i++ {
		ring.Lookup(k, buf)
	}
}

func BenchmarkLookup_128_50(b *testing.B) {
	nodes := build_nodes(128)
	ring := New(nodes, 50)
	k := []byte("hello")
	b.ResetTimer()

	buf := ring.MakeBuffer(-1)

	for i := 0; i < b.N; i++ {
		ring.Lookup(k, buf)
	}
}

func BenchmarkLookup_128_100(b *testing.B) {
	nodes := build_nodes(128)
	ring := New(nodes, 100)
	k := []byte("hello")
	b.ResetTimer()

	buf := ring.MakeBuffer(-1)

	for i := 0; i < b.N; i++ {
		ring.Lookup(k, buf)
	}
}

func BenchmarkLookup_128_200(b *testing.B) {
	nodes := build_nodes(128)
	ring := New(nodes, 200)
	k := []byte("hello")
	b.ResetTimer()

	buf := ring.MakeBuffer(-1)

	for i := 0; i < b.N; i++ {
		ring.Lookup(k, buf)
	}
}

func BenchmarkLookup_128_400(b *testing.B) {
	nodes := build_nodes(128)
	ring := New(nodes, 400)
	k := []byte("hello")
	b.ResetTimer()

	buf := ring.MakeBuffer(-1)

	for i := 0; i < b.N; i++ {
		ring.Lookup(k, buf)
	}
}

func BenchmarkLookup_128_1024(b *testing.B) {
	nodes := build_nodes(128)
	ring := New(nodes, 1024)
	k := []byte("hello")
	b.ResetTimer()

	buf := ring.MakeBuffer(-1)

	for i := 0; i < b.N; i++ {
		ring.Lookup(k, buf)
	}
}

func BenchmarkLookup_256_25(b *testing.B) {
	nodes := build_nodes(256)
	ring := New(nodes, 25)
	k := []byte("hello")
	b.ResetTimer()

	buf := ring.MakeBuffer(-1)

	for i := 0; i < b.N; i++ {
		ring.Lookup(k, buf)
	}
}

func BenchmarkBuild_128_25(b *testing.B) {
	for i := 0; i < b.N; i++ {
		nodes := build_nodes(128)
		New(nodes, 25)
	}
}

func BenchmarkBuild_128_50(b *testing.B) {
	for i := 0; i < b.N; i++ {
		nodes := build_nodes(128)
		New(nodes, 50)
	}
}

func BenchmarkBuild_128_100(b *testing.B) {
	for i := 0; i < b.N; i++ {
		nodes := build_nodes(128)
		New(nodes, 100)
	}
}

func BenchmarkBuild_128_200(b *testing.B) {
	for i := 0; i < b.N; i++ {
		nodes := build_nodes(128)
		New(nodes, 200)
	}
}

func BenchmarkBuild_128_400(b *testing.B) {
	for i := 0; i < b.N; i++ {
		nodes := build_nodes(128)
		New(nodes, 400)
	}
}

func BenchmarkBuild_128_1024(b *testing.B) {
	for i := 0; i < b.N; i++ {
		nodes := build_nodes(128)
		New(nodes, 1024)
	}
}

func BenchmarkBuild_256_25(b *testing.B) {
	for i := 0; i < b.N; i++ {
		nodes := build_nodes(256)
		New(nodes, 25)
	}
}

func build_nodes(l int) []Node {
	o := make([]Node, l)
	for i := 0; i < l; i++ {
		o[i] = &mock_node{i}
	}

	return o
}

type mock_node struct {
	i int
}

func (m *mock_node) HashID() string {
	return fmt.Sprintf("%d", m.i)
}

func (m *mock_node) String() string {
	return fmt.Sprintf("(node:%03d)", m.i)
}
