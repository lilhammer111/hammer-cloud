package mysql

import (
	"fmt"
	"testing"
)

type st struct {
	i int
	s string
}

func TestConn(t *testing.T) {
	a := 1
	b := 1
	c := "1"
	d := "1"
	st1 := st{i: 1, s: "1"}
	st2 := st{i: 1, s: "1"}
	fmt.Println(a == b)
	fmt.Println(c == d)
	fmt.Println(st1 == st2)
}
