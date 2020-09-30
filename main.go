package main

import (
	"bufio"
	"fmt"
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
			file := new(os.File)
			file, err = os.OpenFile(f, flg, 0666)
			defer file.Close()
			if err != nil {
				panic(err)
			}
			files = append(files, file)
		}
	}

	ts := time.Now()
	stdin := bufio.NewScanner(os.Stdin)
	for stdin.Scan() {
		str := ""
		if opts.ClockTime {
			str = fmt.Sprintf("%s %s", time.Now().Format("[2006/01/02 15:04:05.000]"), stdin.Text())
		} else {
			d := time.Since(ts)
			h := int(d.Hours())
			m := int(d.Minutes()) % 60
			s := int(d.Seconds()) % 60
			ms := d.Milliseconds() % 1000
			str = fmt.Sprintf("[%02d:%02d:%02d.%03d] %s", h, m, s, ms, stdin.Text())
		}
		fmt.Println(str)
		for _, f := range files {
			if f != nil {
				fmt.Fprintln(f, str)
			}
		}
	}
}
