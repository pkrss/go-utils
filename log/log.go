package log

import (
	"fmt"
	logOld "log"
	"os"
	"time"

	pkFile "github.com/pkrss/go-utils/file"
)

type Writer interface {
	Write(p []byte) (n int, err error)
	Clean() (err error)
}

var out Writer

var logSplit = "\n"

var logRowsLimit = -1
var logRowsCnt = -1
var timeLocation *time.Location
var stdOut = false

func Println(p string) {
	if out != nil {
		if logRowsLimit != 0 {

			if logRowsLimit > 0 {
				logRowsCnt++
				if logRowsCnt > logRowsLimit {
					logRowsCnt = 0
					out.Clean()
				}
			}

			t := time.Now()
			if timeLocation != nil {
				t = t.In(timeLocation)
			}
			p2 := t.Format("2006-01-02 15:04:05 ") + p + logSplit
			out.Write([]byte(p2))
		}
	}
	if stdOut {
		logOld.Println(p)
	}
}

func Printf(format string, v ...interface{}) {
	p := fmt.Sprintf(format, v...)
	Println(p)
}

func SetOut(o Writer) Writer {
	old := out
	out = o
	return old
}

type LogWriter struct {
	f *os.File
}

func (l *LogWriter) Write(p []byte) (n int, err error) {
	if l.f != nil {
		n, err = l.f.Write(p)
	}
	return
}

func (l *LogWriter) Clean() (err error) {
	if l.f != nil {
		l.f.Truncate(0)
	}
	return
}

func (l *LogWriter) Close() {
	if l.f != nil {
		l.f.Close()
		l.f = nil
	}
	return
}

func SetLogSplitString(s string) {
	logSplit = s
}

func SetLogRowsLimit(limit int) {
	logRowsLimit = limit
}

func SetStdOut(v bool) {
	stdOut = v
}

func SetLogTimeZone(zone string) {
	loc, e := time.LoadLocation(zone)
	if e == nil {
		timeLocation = loc
	} else {
		Printf("SetLogTimeZone(%s) error: %v", zone, e)
	}
}

func NewOutLogWritter(file string) (ret *LogWriter, e error) {
	root := pkFile.FileDir(file)

	logOld.Println("prepare createDir: " + root)

	e = pkFile.CreateDir(root)
	if e != nil {
		logOld.Fatal(e)
		return
	}

	f, err := os.OpenFile(file, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0755)
	if err != nil {
		e = err
		return
	}
	ret = &LogWriter{f: f}
	return
}
