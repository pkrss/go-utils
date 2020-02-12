package log

import (
	"encoding/json"
	"fmt"
	"io"
	logOld "log"
	"os"
	"strconv"
	"strings"
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
var logIgnoreStringList []string

// TimeFmt ...
var TimeFmt string = "2006-01-02 15:04:05.000"

// Println ...
func Println(p string) {
	if logIgnoreStringList != nil {
		for k := range logIgnoreStringList {
			if strings.Contains(p, logIgnoreStringList[k]) {
				return
			}
		}
	}
	if out != nil {
		if logRowsLimit != 0 {

			if logRowsLimit > 0 {
				logRowsCnt++
				if logRowsCnt > logRowsLimit {
					Clean()
				}
			}

			t := time.Now()
			if timeLocation != nil {
				t = t.In(timeLocation)
			}
			p2 := t.Format(TimeFmt) + p + logSplit
			out.Write([]byte(p2))
		}
	}
	if stdOut {
		logOld.Println(p)
	}
}

func Clean() {
	logRowsCnt = 0
	if out != nil {
		out.Clean()
	}
}

func PrintObj(obj interface{}) {
	if obj == nil {
		Println("obj is nil")
		return
	}
	if content, _ := json.Marshal(obj); content != nil {
		Println(string(content))
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

func SetLogIgnoreStringList(v []string) {
	logIgnoreStringList = v
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
	return NewOutLogWritter2(file, false)
}

func NewOutLogWritter2(file string, bak bool) (ret *LogWriter, e error) {
	root := pkFile.FileDir(file)

	logOld.Println("prepare createDir: " + root)

	e = pkFile.CreateDir(root)
	if e != nil {
		logOld.Fatal(e)
		return
	}

	if bak {
		// if _, err := os.Stat(file); os.IsExist(err) {

		in, err := os.Open(file)
		if err == nil {
			defer in.Close()

			out, err := os.Create(file + "_" + strconv.FormatInt(time.Now().Unix(), 10))
			if err == nil {
				defer out.Close()

				io.Copy(out, in)
			}
		}
		// }
	}

	f, err := os.OpenFile(file, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0755)
	if err != nil {
		e = err
		return
	}
	ret = &LogWriter{f: f}
	return
}
