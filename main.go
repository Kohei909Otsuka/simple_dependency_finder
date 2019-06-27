package main

import (
	"sort"
)

type Module struct {
	Id   int
	Name string
	Path string
}

func uniq(org []int) []int {
	output := make([]int, 0)
	marker := make(map[int]bool)
	for _, v := range org {
		if !marker[v] {
			marker[v] = true
			output = append(output, v)
		}
	}
	return output
}

// Original Dep is about who knows who
// Reverse is Who is known by who
func reverseDep(dep map[int][]int) map[int][]int {
	output := make(map[int][]int)
	for k, v := range dep {
		for _, vv := range v {
			output[vv] = append(output[vv], k)
		}
	}
	return output
}

func sequalFinder(ms []*Module, rdep map[int][]int, diff []int) []int {
	output := make([]int, 0)
	output = append(output, diff...)
	for _, v := range diff {
		output = append(output, rdep[v]...)
	}

	// GoにはSetがないのでuniqしてsortして返す
	uniq := uniq(output)
	sort.Ints(uniq)
	return uniq
}

func main() {
}
