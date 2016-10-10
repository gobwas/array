package array

import (
	"math/rand"
	"reflect"
	"sort"
	"testing"
)

func TestAscendRange(t *testing.T) {
	for i, test := range []struct {
		items []Item
		x, y  Item
		exp   []Item
	}{
		{
			[]Item{Int(0), Int(1), Int(2), Int(3), Int(4), Int(5)},
			Int(3), Int(5),
			[]Item{Int(3), Int(4), Int(5)},
		},
		{
			[]Item{Int(0), Int(1), Int(2), Int(3), Int(4), Int(5)},
			Int(3), Int(3),
			[]Item{Int(3)},
		},
	} {
		a := Array{test.items}
		var data []Item
		a.AscendRange(test.x, test.y, func(x Item) bool {
			data = append(data, x)
			return true
		})
		if !reflect.DeepEqual(data, test.exp) {
			t.Errorf("[%d] AscendRange not equal: got %v; want %v\n", i, data, test.exp)
		}
	}
}

// getUniq returns int that not present in dup.
func getUniq(dup map[int]bool) int {
	for {
		v := rand.Int()
		if _, ok := dup[v]; !ok {
			dup[v] = true
			return v
		}
	}
}

// getArrayItems returns n unique Int items that is not present in dup.
func getArrayItems(n int, dup map[int]bool) (ret []Item) {
	ret = make([]Item, n)
	for i := 0; i < n; i++ {
		ret[i] = Int(getUniq(dup))
	}
	return
}

// initArray initializes array with n unique items.
func initArray(n int) (ret Array, items []Item, dup map[int]bool) {
	data := make([]int, n)
	dup = make(map[int]bool, n)
	for i := 0; i < n; i++ {
		data[i] = getUniq(dup)
	}

	sort.Sort(sort.IntSlice(data))
	items = make([]Item, n)
	for i := 0; i < n; i++ {
		items[i] = Int(data[i])
	}

	return Array{items}, items, dup
}

func insertGen(n, k int) func() (Array, []Item) {
	return func() (ret Array, items []Item) {
		ret, _, dup := initArray(n)
		items = getArrayItems(k, dup)
		return
	}
}

func deleteGen(n, k int) func() (Array, []Item) {
	return func() (ret Array, remove []Item) {
		ret, items, _ := initArray(n)
		remove = make([]Item, k)
		for i, v := range rand.Perm(n) {
			remove[i] = items[v]
			if i+1 == k {
				break
			}
		}
		return
	}
}

func benchmarkArrayInsert(b *testing.B, n int) {
	ops := 10
	gen := insertGen(n, ops)
	var prev Item
	var arr Array
	var insert []Item

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		next := i % ops
		if next == 0 {
			b.StopTimer()
			arr, insert = gen()
			b.StartTimer()
		}
		arr, prev = arr.Upsert(insert[next])
		if prev != nil {
			b.Fatalf("could not insert new item %v: already exists", insert[next])
		}
	}
}

func benchmarkArrayDelete(b *testing.B, n int) {
	ops := 10
	gen := deleteGen(n, ops)
	var prev Item
	var arr Array
	var remove []Item

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		next := i % ops
		if next == 0 {
			b.StopTimer()
			arr, remove = gen()
			b.StartTimer()
		}
		arr, prev = arr.Delete(remove[next])
		if prev == nil {
			b.Fatalf("could not remove item %v: not exists", remove[next])
		}
	}
}

func benchmarkArrayHas(b *testing.B, n int) {
	arr, items, _ := initArray(n)
	check := rand.Perm(n)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		if !arr.Has(items[check[i%n]]) {
			b.Fatalf("@_O arr has no such item")
		}
	}
}

func benchmarkArrayHasMiss(b *testing.B, n int) {
	arr, _, dup := initArray(n)
	check := getArrayItems(b.N, dup)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		if arr.Has(check[i]) {
			b.Fatalf("@_@ arr has item")
		}
	}
}

func BenchmarkArrayInsert_0(b *testing.B)      { benchmarkArrayInsert(b, 0) }
func BenchmarkArrayInsert_1000(b *testing.B)   { benchmarkArrayInsert(b, 1000) }
func BenchmarkArrayInsert_100000(b *testing.B) { benchmarkArrayInsert(b, 100000) }

func BenchmarkArrayDelete_1000(b *testing.B)   { benchmarkArrayDelete(b, 1000) }
func BenchmarkArrayDelete_100000(b *testing.B) { benchmarkArrayDelete(b, 100000) }

func BenchmarkArrayHas_1000(b *testing.B)    { benchmarkArrayHas(b, 1000) }
func BenchmarkArrayHas_100000(b *testing.B)  { benchmarkArrayHas(b, 100000) }
func BenchmarkArrayHas_1000000(b *testing.B) { benchmarkArrayHas(b, 1000000) }

func BenchmarkArrayHasMiss_1000(b *testing.B)    { benchmarkArrayHasMiss(b, 1000) }
func BenchmarkArrayHasMiss_100000(b *testing.B)  { benchmarkArrayHasMiss(b, 100000) }
func BenchmarkArrayHasMiss_1000000(b *testing.B) { benchmarkArrayHasMiss(b, 1000000) }
