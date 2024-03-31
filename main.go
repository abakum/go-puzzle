package main

import (
	"fmt"
	"io"
	"os"
	"os/exec"
	// можно импортировать любые модули
)

var (
// можно объявлять любые переменные
)

func start() {
	// можно писать любой код
}

func done() {
	// можно писать любой код
}

func stdIn() io.ReadCloser {
	// можно писать любой код
	return os.Stdin
}

// Мой маленький дружок, попробуй изменить функции кроме `main`,
// чтоб можно было завершить программу дважды нажав клавишу `Enter`
//
// Когда подрастёшь, то сможешь завершить программу дважды нажав любую клавишу
func main() {
	for i := 0; i < 2; i++ {
		start()
		cmd := exec.Command("cmd", "/c", "echo Press any key to continue . . .&&pause")

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
