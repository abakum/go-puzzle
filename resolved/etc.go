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
	arg2 = "echo Press any key to continue . . .; read -rn1"
)

func start() {
	// можно писать любой код
	rc = os.Stdin
}
func done() {
	// можно писать любой код
}
