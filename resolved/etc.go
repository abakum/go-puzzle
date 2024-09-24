//go:build !windows
// +build !windows

package main

import (
	"os"
	// можно импортировать любые модули
)

const (
	arg0 = "bash"
	arg1 = "-c"
	arg2 = `read -n1 -r -p "Press any key to continue..."`
)

func start() {
	// можно писать любой код
	rc = os.Stdin
}
func done() {
	// можно писать любой код
}
