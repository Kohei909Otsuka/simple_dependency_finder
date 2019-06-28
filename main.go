package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"os"
	"sort"
	"strconv"
	"strings"
)

type Module struct {
	Id   int    `json:"id"`
	Name string `json:"name"`
	Path string `json:"path"`
}

func findModule(
	modules []Module,
	fun func(Module) bool,
) Module {
	for _, m := range modules {
		if fun(m) {
			return m
		}
	}
	panic("Failed not find module")
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

func sequalFinder(rdep map[int][]int, diff []int, depth int) []int {
	counter := 0
	checked := make([]int, 0)
	effected := recursive(diff, checked, rdep, depth, &counter)
	effected = append(effected, diff...) // 変更があったものは検索するまでもなく追加

	// GoにはSetがないのでuniqしてsortして返す
	uniq := uniq(effected)
	sort.Ints(uniq)
	return uniq
}

func validateParams(mPath, rPath, diffs string, depth int) error {
	if mPath == "" {
		return fmt.Errorf("mpath option must be given")
	}
	if rPath == "" {
		return fmt.Errorf("rpath option must be given")
	}
	if diffs == "" {
		return fmt.Errorf("diffs option must be given")
	}

	return nil
}

func parseModuleFile(path string) ([]Module, error) {
	bytes, err := ioutil.ReadFile(path)
	if err != nil {
		return []Module{}, err
	}

	var modules []Module
	if err := json.Unmarshal(bytes, &modules); err != nil {
		return []Module{}, err
	}

	return modules, nil
}

func parseRelationFile(path string) (map[int][]int, error) {
	bytes, err := ioutil.ReadFile(path)
	if err != nil {
		return nil, err
	}

	var relation map[string][]int
	if err := json.Unmarshal(bytes, &relation); err != nil {
		return nil, err
	}

	intRelation := make(map[int][]int)
	for k, v := range relation {
		i, err := strconv.Atoi(k)
		if err != nil {
			return nil, err
		}
		intRelation[i] = v
	}

	return intRelation, nil
}

func parseDiffs(srtDiffs string, modules []Module) ([]int, error) {
	diffs := strings.Split(srtDiffs, ",")
	intDiffs := make([]int, len(diffs))
	for i, d := range diffs {
		var targetModuleId int
		for _, m := range modules {
			if m.Path == d {
				targetModuleId = m.Id
			}
		}

		if targetModuleId == 0 {
			return intDiffs,
				fmt.Errorf("changed file %v is not defined in modules %+v\n", d, modules)
		}

		intDiffs[i] = targetModuleId
	}

	return intDiffs, nil
}

func main() {
	var modulesPath, relationPath, diffs string
	var depth int
	var debug bool

	flag.StringVar(&modulesPath, "mpath", "", "path to modules.json")
	flag.StringVar(&relationPath, "rpath", "", "path to relations.json")
	flag.StringVar(&diffs, "diffs", "", "file diffs comma separated")
	flag.IntVar(&depth, "depth", 0, "depth to search default is unlimited")
	flag.BoolVar(&debug, "debug", false, "debug flag to show info")
	flag.Parse()

	if debug {
		fmt.Printf("arg modulesPath: %v\n", modulesPath)
		fmt.Printf("arg relationPath: %v\n", relationPath)
		fmt.Printf("arg diffs: %v\n", diffs)
		fmt.Printf("arg depth: %v\n", depth)
		fmt.Printf("arg debug: %v\n", debug)
	}

	err := validateParams(modulesPath, relationPath, diffs, depth)
	if err != nil {
		fmt.Printf("Argument validatoin err: %v\n", err)
		os.Exit(1)
	}

	modules, err := parseModuleFile(modulesPath)
	if err != nil {
		fmt.Printf("Failed to parse module json: %v\n", err)
		os.Exit(1)
	}

	relation, err := parseRelationFile(relationPath)
	if err != nil {
		fmt.Printf("Failed to parse relation json: %v\n", err)
		os.Exit(1)
	}

	intDiffs, err := parseDiffs(diffs, modules)
	if err != nil {
		fmt.Printf("Failed to parse diff by modules: %v\n", err)
		os.Exit(1)
	}

	// ここがメインの検索処理
	effectedIds := sequalFinder(reverseDep(relation), intDiffs, depth)
	effectedModules := make([]Module, len(effectedIds))
	for i, v := range effectedIds {
		effectedModules[i] = findModule(modules, func(m Module) bool {
			return m.Id == v
		})
	}

	if debug {
		fmt.Printf("modules: %+v\n", modules)
		fmt.Printf("relations: %v\n", relation)
		fmt.Printf("intDiffs: %v\n", intDiffs)
		fmt.Printf("effectedModules: %+v\n", effectedModules)
	}

	// 出力
	bytes, err := json.Marshal(effectedModules)
	if err != nil {
		fmt.Printf("Failed to write output: %v\n", err)
		os.Exit(1)
	}
	os.Stdout.Write(bytes)
}
