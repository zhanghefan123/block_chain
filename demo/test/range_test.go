package test

import "testing"

func TestRange(t *testing.T) {
	var array = []int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10}
	for i, value := range array {
		t.Log(i, value)
	}
}
