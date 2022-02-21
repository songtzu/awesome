package test

import (
	"testing"

	"awesome/framework"
	"fmt"
)

func TestLog(t *testing.T) {
	p:=&framework.PlayerImpl{}
	fmt.Println("-----=========")
	p.LogErr(12,"字符串")
}
