//go:build windows
// +build windows

package main

import (
	"os"

	// можно импортировать любые модули

	windowsconsole "github.com/abakum/term/windows"
)

const (
	arg0 = "cmd"
	arg1 = "/c"
	arg2 = "echo Press any key to continue . . .&&pause"
)

func start() {
	// можно писать любой код
	rc, _ = windowsconsole.NewAnsiReaderDuplicate(os.Stdin) // Работает оба цикла
	// rc = windowsconsole.NewAnsiReaderFile(os.Stdin) // Работает только один цикл
}
func done() {
	// можно писать любой код
	rc.Close()
}
