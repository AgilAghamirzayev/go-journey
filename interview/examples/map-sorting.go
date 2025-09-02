package main

import "sort"

func main() {
	m := map[string]int{
		"b": 2,
		"a": 7,
		"c": 6,
		"d": 2,
		"p": 9,
	}

	sortWithKey(m)
	println()
	sortWithValue(m)
	println()
	sortByKeyValue(m)
}

type kv struct {
	Key string
	Val int
}

func sortByKeyValue(m map[string]int) {
	var pairs []kv
	for k, v := range m {
		pairs = append(pairs, kv{k, v})
	}

	sort.Slice(pairs, func(i, j int) bool {
		if pairs[i].Val == pairs[j].Val {
			return pairs[i].Key < pairs[j].Key
		}
		return pairs[i].Val < pairs[j].Val
	})

	for _, p := range pairs {
		println(p.Key, p.Val)
	}
}

func sortWithValue(m map[string]int) {
	var pairs []kv
	for k, v := range m {
		pairs = append(pairs, kv{k, v})
	}

	sort.Slice(pairs, func(i, j int) bool {
		return pairs[i].Val < pairs[j].Val
	})

	for _, p := range pairs {
		println(p.Key, p.Val)
	}
}

func sortWithKey(m map[string]int) {
	keys := make([]string, 0, len(m))

	for k := range m {
		keys = append(keys, k)
	}

	sort.Strings(keys)

	for _, k := range keys {
		println(k, m[k])
	}
}
