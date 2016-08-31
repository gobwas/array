package array

import (
	"math/rand"
	"sort"
	"testing"
)

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

func benchmarkArrayInsert(b *testing.B, n int) {
	arr, _, dup := initArray(n)
	insert := getArrayItems(b.N, dup)
	var prev Item

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		arr, prev = arr.Upsert(insert[i])
		if prev != nil {
			b.Fatalf("could not insert new item %v: already exists", insert[i])
		}
	}
}

func benchmarkArrayDelete(b *testing.B, n int) {
	arr, items, _ := initArray(b.N + n)
	remove := rand.Perm(b.N + n)
	var prev Item

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		arr, prev = arr.Delete(items[remove[i]])
		if prev == nil {
			b.Fatalf("could not remove item %v: not exists", items[remove[i]])
		}
	}
}

func benchmarkArrayHas(b *testing.B, n int) {
	arr, items, _ := initArray(b.N + n)
	check := rand.Perm(b.N + n)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		if !arr.Has(items[check[i]]) {
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

func BenchmarkArrayDelete_0(b *testing.B)      { benchmarkArrayDelete(b, 0) }
func BenchmarkArrayDelete_1000(b *testing.B)   { benchmarkArrayDelete(b, 1000) }
func BenchmarkArrayDelete_100000(b *testing.B) { benchmarkArrayDelete(b, 100000) }

func BenchmarkArrayHas_0(b *testing.B)      { benchmarkArrayHas(b, 0) }
func BenchmarkArrayHas_1000(b *testing.B)   { benchmarkArrayHas(b, 1000) }
func BenchmarkArrayHas_100000(b *testing.B) { benchmarkArrayHas(b, 100000) }

func BenchmarkArrayHasMiss_1000(b *testing.B)   { benchmarkArrayHasMiss(b, 1000) }
func BenchmarkArrayHasMiss_100000(b *testing.B) { benchmarkArrayHasMiss(b, 100000) }
