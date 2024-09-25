package main

import (
	"context"
	"fmt"
	"io"
	"log"
	"os"
	"os/exec"
	"strings"
	"time"

	"github.com/containerd/console"
	"github.com/mattn/go-isatty"
	"github.com/xlab/closer"
)

var (
	raw  bool
	once bool
)

const delay = time.Millisecond * 77

var reset = func(*bool) {}

func main() {
	defer func() {
		reset(&raw)
		closer.Close()
	}()
	log.SetFlags(log.Lmicroseconds | log.Lshortfile)
	log.SetPrefix("\r")
	parent := context.Background()
	var cmd *exec.Cmd
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
			log.Println("without pipes")
			cmd.Stdin = os.Stdin
			cmd.Stdout = os.Stdout
		} else {
			ConsoleCP(&once)

			log.Println("with pipes")

			out, err := cmd.StdoutPipe()
			if err != nil {
				panic(err)
			}

			in, err := cmd.StdinPipe()
			if err != nil {
				panic(err)
			}

			go func() {
				io.Copy(os.Stdout, NewReader(ctx, out, 0))
				cancel()
				log.Println("Stdout done")
			}()

			go func() {
				io.Copy(in, NewReader(ctx, os.Stdin, delay))
				cancel()
				log.Println("Stdin done")

			}()
		}
		fmt.Print("\r")
		if i%4 > 1 {
			fmt.Println("Type exit<Enter>\r")
		}
		cmd.Run()

	}
}

type reader struct {
	ctx context.Context
	r   io.Reader
	d   time.Duration
}

func NewReader(ctx context.Context, r io.Reader, d time.Duration) io.Reader {
	if r, ok := r.(*reader); ok && ctx == r.ctx {
		return r
	}
	return &reader{ctx: ctx, r: r, d: d}
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
