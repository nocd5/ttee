package main

import (
	"bufio"
	"io"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/jessevdk/go-flags"
)

type Options struct {
	Append          bool `short:"a" long:"append" description:"Append to the given FILEs, do not overwrite"`
	IgnoreInterrupt bool `short:"i" long:"ignore-interrupts" description:"Ignore interrupt signals"`
	ClockTime       bool `short:"c" long:"clock-time" description:"Display clock time"`
}

var opts Options

func main() {
	parser := flags.NewParser(&opts, flags.Default)
	parser.Name = "ttee"
	parser.Usage = "[OPTION]... [FILE]..."
	args, err := parser.Parse()
	if err != nil {
		os.Exit(1)
	}

	if opts.IgnoreInterrupt {
		signal.Ignore(syscall.SIGINT)
	}

	var files []*os.File

	if len(args) > 0 {
		flg := os.O_WRONLY | os.O_CREATE
		if opts.Append {
			flg |= os.O_APPEND
		} else {
			flg |= os.O_TRUNC
		}
		for _, f := range args {
			file, err := os.OpenFile(f, flg, 0666)
			defer file.Close()
			if err != nil {
				panic(err)
			}
			files = append(files, file)
		}
	}

	ttw := NewTteeWriter(os.Stdout, files, time.Now(), opts.ClockTime)
	tr := io.TeeReader(os.Stdin, ttw)
	s := bufio.NewScanner(tr)
	for s.Scan() {
	}
	if err := s.Err(); err != nil {
		panic(err)
	}
}
