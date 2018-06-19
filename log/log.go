package log

import (
	"fmt"
	"io"
	logOld "log"
	"os"
	"time"

	pkFile "github.com/pkrss/go-utils/file"
)

var out io.Writer

func Println(p string) {
	if out != nil {
		p2 := time.Now().Format("2006-01-02 15:04:05 ") + p + "\n"
		out.Write([]byte(p2))
	}
	logOld.Println(p)
}

func Printf(format string, v ...interface{}) {
	p := fmt.Sprintf(format, v...)
	Println(p)
}

func SetOut(o io.Writer) io.Writer {
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

func (l *LogWriter) Close() {
	if l.f != nil {
		l.f.Close()
		l.f = nil
	}
	return
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
