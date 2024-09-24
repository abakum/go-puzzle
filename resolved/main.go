package main

import (
	"fmt"
	"io"
	"os"
	"os/exec"
	// можно импортировать любые модули
)

var (
	rc io.ReadCloser
)

// Какой ты хитренький \8^)
// Может ты сможешь изменить строчку
// `cmd := exec.Command("cmd", "/c", "echo Press any key to continue . . .&&pause")`
// на
// `cmd := exec.Command("cmd")`
// и выйти из программы введя дважды команду `exit`?
func main() {

	for i := 0; i < 2; i++ {
		start()
		cmd := exec.Command(arg0, arg1, arg2) // Работает
		// cmd := exec.Command("cmd") // Не работает
		// cmd := exec.Command("powershell") // Работает

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

		io.Copy(in, rc)
		fmt.Println("Stdin done")
	}
}
