package main

import (
	"fmt"
	"io"
	"os"
	"os/exec"

	"github.com/abakum/term"
	// можно импортировать любые модули
)

var (
	// можно объявлять любые переменные
	ioe *term.IOE
)

func start() {
	// можно писать любой код
	ioe = term.NewIOE()
}

func done() {
	// можно писать любой код
	ioe.Close()
}

func stdIn() io.ReadCloser {
	// можно писать любой код
	return ioe.ReadCloser()
}

// Какой ты хитренький \8^)
// Может ты сможешь изменить строчку
// `cmd := exec.Command("cmd", "/c", "echo Press any key to continue . . .&&pause")`
// на
// `cmd := exec.Command("cmd")`
// и выйти из программы введя дважды команду `exit`?
func main() {

	for i := 0; i < 2; i++ {
		start()
		cmd := exec.Command("cmd", "/c", "echo Press any key to continue . . .&&pause")
		// cmd := exec.Command("cmd")

		out, err := cmd.StdoutPipe()
		if err != nil {
			panic(err)
		}
		in, err := cmd.StdinPipe()
		if err != nil {
			panic(err)
		}

		err = cmd.Start()
		if err != nil {
			panic(err)
		}

		go func() {
			io.Copy(os.Stdout, out)
			done()
			fmt.Println("Stdout done")
		}()

		go func() {
			cmd.Wait()
			fmt.Println("Wait done")
		}()

		io.Copy(in, stdIn())
		fmt.Println("Stdin done")
	}
}
