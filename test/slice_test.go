package test

import (
	"testing"
	"fmt"
)

func TestSlice(t *testing.T) {
	s:=[]int{}
	fmt.Println(len(s))
	s = append(s,1)
	fmt.Println(len(s))
	fmt.Println(s[len(s)-1])
}