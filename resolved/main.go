package main

import (
	"context"
	"fmt"
	"io"
	"log"
	"os"
	"os/exec"
	"runtime"
	"strings"
	"time"

	"github.com/containerd/console"
	"github.com/mattn/go-isatty"
	"github.com/xlab/closer"
)

func main() {
	var (
		raw    bool
		once   bool
		reset  = func(*bool) {}
		cmd    *exec.Cmd
		parent = context.Background()
		arg0   = "bash"
		arg1   = "-c"
		arg2   = "echo Press any key to continue . . .;read -rn1"
		delay  = time.Millisecond * 4
	)

	defer func() {
		reset(&raw)
		closer.Close()
	}()
	if isatty.IsCygwinTerminal(os.Stdin.Fd()) {
		ConsoleCP(&once)
	} else if runtime.GOOS == "windows" {
		delay *= 2
		arg0 = "cmd"
		arg1 = "/c"

		// delay *= 4
		// arg0 = "powershell"
		// arg1 = "-command"

		arg2 = "pause"
	}
	log.SetFlags(log.Lmicroseconds | log.Lshortfile)
	log.SetPrefix("\r")

	for i := 0; i < 8; i++ {
		ctx, cancel := context.WithCancel(parent)

		if i%4 > 1 {
			reset(&raw)
			cmd = exec.Command(arg0)
		} else {
			reset = setRaw(&raw, reset)
			cmd = exec.Command(arg0, arg1, arg2)
		}
		log.Println(cmd)
		if i < 4 {
			log.Println("---without pipes", i)
			cmd.Stdin = os.Stdin
			cmd.Stdout = os.Stdout
			fmt.Print("\r")
			if i%4 > 1 {
				fmt.Println("Type exit<Enter>\r")
			}
			cmd.Run()
		} else {
			log.Println("--with pipes", i)
			ConsoleCP(&once)

			in, err := cmd.StdinPipe()
			if err != nil {
				panic(err)
			}

			out, err := cmd.StdoutPipe()
			if err != nil {
				panic(err)
			}

			fmt.Print("\r")
			if i%4 > 1 {
				fmt.Println("Type exit<Enter>\r")
			}
			err = cmd.Start()
			if err != nil {
				panic(err)
			}

			go func() {
				defer log.Println("Stdin done", i)
				io.Copy(in, NewReader(ctx, os.Stdin, delay))
				// io.Copy(in, NewBufReader(ctx, os.Stdin))
				// io.Copy(in, bufio.NewReader(os.Stdin))
				// io.Copy(in, os.Stdin)
			}()
			io.Copy(os.Stdout, out)
			log.Println("Stdout done", i)
			cancel()
			cmd.Process.Release()
		}
	}
	time.Sleep(delay)
}

type reader struct {
	ctx context.Context
	r   io.Reader
	d   time.Duration
}

func NewReader(ctx context.Context, r io.Reader, d time.Duration) io.Reader {
	if r, ok := r.(*reader); ok && ctx == r.ctx && d == r.d {
		return r
	}
	return &reader{
		ctx: ctx,
		r:   r,
		d:   d,
	}
}

func (r *reader) Read(p []byte) (n int, err error) {
	if r.d > 0 {
		select {
		case <-r.ctx.Done():
			return 0, r.ctx.Err()
		case <-time.After(r.d):
			return r.r.Read(p)
		}
	}
	select {
	case <-r.ctx.Done():
		return 0, r.ctx.Err()
	default:
		return r.r.Read(p)
	}
}

func setRaw(raw *bool, old func(*bool)) (reset func(*bool)) {
	reset = old
	if *raw {
		return
	}
	var (
		err      error
		current  console.Console
		settings string
	)

	current, err = console.ConsoleFromFile(os.Stdin)
	if err == nil {
		err = current.SetRaw()
		if err == nil {
			*raw = true
			reset = func(raw *bool) {
				if *raw {
					err := current.Reset()
					log.Println("Restores the console to its original state by go", err)
				}
				*raw = err != nil
			}
			log.Println("Sets the console in raw mode by go")
			return
		}
	}

	if isatty.IsCygwinTerminal(os.Stdin.Fd()) {
		settings, err = sttySettings()
		if err == nil {
			err = sttyMakeRaw()
			if err == nil {
				*raw = true
				reset = func(raw *bool) {
					if *raw {
						sttyReset(settings)
						log.Println("Restores the console to its original state by stty")
					}
					*raw = false
				}
				log.Println("Sets the console in raw mode by stty")
				return
			}
		}
	}
	log.Println(err)
	return
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

// type bufReader struct {
// 	ctx context.Context
// 	r   *bufio.Reader
// }

// func NewBufReader(ctx context.Context, r io.Reader) io.Reader {
// 	if r, ok := r.(*bufReader); ok && ctx == r.ctx {
// 		return r
// 	}
// 	return &bufReader{
// 		ctx: ctx,
// 		r:   bufio.NewReader(r),
// 	}
// }

// func (r *bufReader) Read(p []byte) (n int, err error) {
// 	for {
// 		select {
// 		case <-r.ctx.Done():
// 			return 0, r.ctx.Err()
// 		case <-time.After(delay):
// 			if r.r.Buffered() == 0 {
// 				continue
// 			}
// 			return r.r.Read(p)
// 		}
// 	}
// }
