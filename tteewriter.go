package main

import (
	"fmt"
	"io"
	"os"
	"strings"
	"time"
)

type TteeWriter struct {
	dest      io.Writer
	files     []*os.File
	baseTime  time.Time
	clockTime bool
}

func NewTteeWriter(dest io.Writer, files []*os.File, baseTime time.Time, clockTime bool) *TteeWriter {
	return &TteeWriter{dest, files, baseTime, clockTime}
}

func (ttw *TteeWriter) Write(data []byte) (int, error) {
	s := string(data)
	n := 0
	if strings.Contains(s, "\n") {
		s := strings.TrimSuffix(s, "\n")
		for _, l := range strings.Split(s, "\n") {
			lf := fmt.Sprintf("[%s] %s\n", ttw.getTime(), string(l))
			m, _err := ttw.dest.Write([]byte(lf))
			if _err != nil {
				return n, _err
			}
			n += m

			for _, f := range ttw.files {
				if f != nil {
					fmt.Fprint(f, lf)
				}
			}
		}
	} else {
		ttw.dest.Write(data)

		for _, f := range ttw.files {
			if f != nil {
				fmt.Fprint(f, string(data))
			}
		}
	}
	return n, nil
}

func (ttw *TteeWriter) getTime() string {
	str := ""
	if ttw.clockTime {
		str = time.Now().Format("2006/01/02 15:04:05.000")
	} else {
		d := time.Since(ttw.baseTime)
		h := int(d.Hours())
		m := int(d.Minutes()) % 60
		s := int(d.Seconds()) % 60
		ms := d.Milliseconds() % 1000
		str = fmt.Sprintf("%02d:%02d:%02d.%03d", h, m, s, ms)
	}
	return str
}
