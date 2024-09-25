//go:build windows
// +build windows

package main

import (
	"github.com/xlab/closer"
	"golang.org/x/sys/windows"
)

// можно импортировать любые модули

const (
	arg0 = "cmd"
	arg1 = "/c"
	// arg0 = "powershell"
	// arg1 = "-command"
	arg2 = "pause"
)

func ConsoleCP(once *bool) {
	if *once {
		return
	}
	*once = false
	const CP_UTF8 uint32 = 65001
	var kernel32 = windows.NewLazyDLL("kernel32.dll")

	getConsoleCP := func() uint32 {
		result, _, _ := kernel32.NewProc("GetConsoleCP").Call()
		return uint32(result)
	}

	getConsoleOutputCP := func() uint32 {
		result, _, _ := kernel32.NewProc("GetConsoleOutputCP").Call()
		return uint32(result)
	}

	setConsoleCP := func(cp uint32) {
		kernel32.NewProc("SetConsoleCP").Call(uintptr(cp))
	}

	setConsoleOutputCP := func(cp uint32) {
		kernel32.NewProc("SetConsoleOutputCP").Call(uintptr(cp))
	}

	inCP := getConsoleCP()
	outCP := getConsoleOutputCP()
	setConsoleCP(CP_UTF8)
	setConsoleOutputCP(CP_UTF8)
	closer.Bind(func() { setConsoleCP(inCP) })
	closer.Bind(func() { setConsoleOutputCP(outCP) })
}
