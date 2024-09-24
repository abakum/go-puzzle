package main

import (
	"context"
	"io"
	"log"
	"os"
	"os/exec"
	"strings"
	"time"

	"github.com/containerd/console"
	"github.com/mattn/go-isatty"
	"github.com/xlab/closer"
	// можно импортировать любые модули
)

var (
	rc   io.ReadCloser
	once bool
)

// Какой ты хитренький \8^)
// Может ты сможешь изменить строчку
// `cmd := exec.Command("cmd", "/c", "echo Press any key to continue . . .&&pause")`
// на
// `cmd := exec.Command("cmd")`
// и выйти из программы введя дважды команду `exit`?
func main() {
	defer closer.Close()
	// setRaw(&once)
	log.SetFlags(log.Lmicroseconds | log.Lshortfile)
	for i := 0; i < 2; i++ {
		// start()
		ctx, cancel := context.WithCancel(context.Background())
		// cmd := exec.Command(arg0, arg1, arg2) // Работает
		cmd := exec.Command(arg0) // Работает
		log.Println(cmd)
		// cmd.Stdin = os.Stdin
		// cmd.Stdout = os.Stdout
		// cmd.Run()
		// continue
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
			cancel()
			log.Println("Stdout done")
		}()

		io.Copy(in, NewReader(ctx, os.Stdin))
		log.Println("Stdin done")
		cmd.Wait()
		log.Println("Wait done")

	}

}

type reader struct {
	ctx context.Context
	r   io.Reader
}

func NewReader(ctx context.Context, r io.Reader) io.Reader {
	if r, ok := r.(*reader); ok && ctx == r.ctx {
		return r
	}
	return &reader{ctx: ctx, r: r}
}

func (r *reader) Read(p []byte) (n int, err error) {
	select {
	case <-r.ctx.Done():
		return 0, r.ctx.Err()
	case <-time.After(time.Millisecond * 3):
		return r.r.Read(p)
	}
}

func setRaw(already *bool) {
	if *already {
		return
	}
	*already = true

	var (
		err      error
		current  console.Console
		settings string
	)

	current, err = console.ConsoleFromFile(os.Stdin)
	if err == nil {
		err = current.SetRaw()
		if err == nil {
			closer.Bind(func() { current.Reset() })
			log.Println("Set raw by go")
			return
		}
	}

	if isatty.IsCygwinTerminal(os.Stdin.Fd()) {
		settings, err = sttySettings()
		if err == nil {
			err = sttyMakeRaw()
			if err == nil {
				closer.Bind(func() { sttyReset(settings) })
				log.Println("Set raw by stty")
				return
			}
		}
	}
	log.Println(err)

}

func sttyMakeRaw() error {
	cmd := exec.Command("stty", "raw", "-echo")
	cmd.Stdin = os.Stdin
	return cmd.Run()
}

func sttySettings() (string, error) {
	cmd := exec.Command("stty", "-g")
	cmd.Stdin = os.Stdin
	out, err := cmd.Output()
	if err != nil {
		return "", err
	}
	return strings.TrimSpace(string(out)), nil
}

func sttyReset(settings string) {
	cmd := exec.Command("stty", settings)
	cmd.Stdin = os.Stdin
	_ = cmd.Run()
}
