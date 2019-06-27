package main

import (
	"testing"
)

var modules = []*Module{
	&Module{Id: 1, Name: "Order", Path: "order.rb"},
	&Module{Id: 2, Name: "OrderItem", Path: "order_item.rb"},
	&Module{Id: 3, Name: "User", Path: "user.rb"},
	&Module{Id: 4, Name: "Item", Path: "item.rb"},
	&Module{Id: 5, Name: "CircularA", Path: "circular_a.rb"},
	&Module{Id: 6, Name: "CircularB", Path: "circular_b.rb"},
	&Module{Id: 7, Name: "DeepCircularA", Path: "deep_circular_a.rb"},
	&Module{Id: 8, Name: "DeepCircularB", Path: "deep_circular_b.rb"},
	&Module{Id: 9, Name: "DeepCircularC", Path: "deep_circular_c.rb"},
}

// - 3: User
//   - 1: Order
//     - 2: OrderItem
// - 4: Item
//     - 2: OrderItem
// - 5: CircularA
//     - 6: CircularB
// - 6: CircularA
//     - 5: CircularB
var dep = map[int][]int{
	2: []int{1, 4},
	1: []int{3},
	5: []int{6},
	6: []int{5},
	7: []int{8},
	8: []int{9},
	9: []int{7},
}

var table = []struct {
	diff     []int
	depth    int
	expected []int
}{
	{
		diff:     []int{2},
		depth:    0,
		expected: []int{2},
	},
	{
		diff:     []int{4},
		depth:    0,
		expected: []int{2, 4},
	},
	// specify depth
	{
		diff:     []int{3},
		depth:    1,
		expected: []int{1, 3},
	},
	{
		diff:     []int{3},
		depth:    0,
		expected: []int{1, 2, 3},
	},
	{
		diff:     []int{1},
		depth:    0,
		expected: []int{1, 2},
	},
	// multiple diff
	{
		diff:     []int{1, 4},
		depth:    0,
		expected: []int{1, 2, 4},
	},
	// simple circular dependency
	{
		diff:     []int{5, 6},
		depth:    0,
		expected: []int{5, 6},
	},
	{
		diff:     []int{5},
		depth:    0,
		expected: []int{5, 6},
	},
	// complext circular dependency
	{
		diff:     []int{7},
		depth:    0,
		expected: []int{7, 8, 9},
	},
	{
		diff:     []int{8, 9},
		depth:    0,
		expected: []int{7, 8, 9},
	},
}

func intSliceEqual(a, b []int) bool {
	if len(a) != len(b) {
		return false
	}

	for i, _ := range a {
		if a[i] != b[i] {
			return false
		}
	}

	return true
}

func TestSequelFinder(t *testing.T) {
	for _, v := range table {
		actual := sequalFinder(modules, reverseDep(dep), v.diff, v.depth)
		expected := v.expected
		if !intSliceEqual(expected, actual) {
			t.Errorf("expected %v, actual %v for %+v\n", expected, actual, v)
		}
	}
}
