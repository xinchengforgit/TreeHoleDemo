package main

import (
	"fmt"
	"testing"
)

func Test_pkg(t *testing.T) {
	a := ErrorResponse("123")
	fmt.Println(a.Success, a.Hint, a.Data)
}
