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

func includes(con []int, e int) bool {
	for _, v := range con {
		if v == e {
			return true
		}
	}
	return false
}

// ある整数n int[]をうけとって、それに依存しているm []intを返す関数
// counter must be pointer to 0
func recursive(nums, checked []int, rdep map[int][]int, depth int, counter *int) []int {

	output := make([]int, 0)

	// 再帰処理の終了条件
	if len(nums) == 0 {
		return output
	}

	for _, v := range nums {
		// output = append(output, rdep[v]...)
		for _, vv := range rdep[v] {
			if includes(checked, vv) {
				continue
			}
			checked = append(checked, vv)
			output = append(output, vv)
		}
	}

	*counter += 1
	if *counter < depth || depth == 0 {
		return append(output, recursive(output, checked, rdep, depth, counter)...)
	} else {
		return output
	}
}

// Original Dep is about who knows who
// Reverse is Who is known by who
func reverseDep(dep map[int][]int) map[int][]int {
	output := make(map[int][]int)
	for k, v := range dep {
		for _, vv := range v {
			output[vv] = append(output[vv], k)
		}

		if output[k] == nil {
			output[k] = []int{}
		}
	}
	return output
}

func sequalFinder(ms []*Module, rdep map[int][]int, diff []int, depth int) []int {
	counter := 0
	checked := make([]int, 0)
	effected := recursive(diff, checked, rdep, depth, &counter)
	effected = append(effected, diff...) // 変更があったものは検索するまでもなく追加

	// GoにはSetがないのでuniqしてsortして返す
	uniq := uniq(effected)
	sort.Ints(uniq)
	return uniq
}

func main() {
}
