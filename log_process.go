package main

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"strings"
	"time"
)

type Reader interface {
	Read(rc chan []byte)
}

type Writer interface {
	Write(wc chan string)
}

type LogProcess struct {
	rc    chan []byte
	wc    chan string
	read  Reader
	write Writer
}

type ReadFromFile struct {
	path string // 读取文件路径
}

func (r *ReadFromFile) Read(rc chan []byte) {
	// 读取模块

	f, err := os.Open(r.path)
	if err != nil {
		panic(fmt.Sprintf("open file error: %s", err.Error()))
	}

	f.Seek(0, 2)
	rd := bufio.NewReader(f)

	for {
		line, err := rd.ReadBytes('\n')
		if err == io.EOF {
			time.Sleep(500 * time.Millisecond)
			continue
		} else if err != nil {
			panic(fmt.Sprintf("ReadBytes error: %s", err.Error()))
		}
		rc <- line
	}
}

func (l *LogProcess) Process() {
	// 解析模块

	for v := range l.rc {
		l.wc <- strings.ToUpper(string(v))
	}

}

type WriteToInfluxDB struct {
	influxDBDsn string
}

func (w *WriteToInfluxDB) Write(wc chan string) {
	// 写入模块
	for v := range wc {
		fmt.Print(v)
	}
}

func main() {

	r := &ReadFromFile{
		path: "./access.log",
	}

	w := &WriteToInfluxDB{
		influxDBDsn: "user&pass",
	}

	lp := &LogProcess{
		rc:    make(chan []byte),
		wc:    make(chan string),
		read:  r,
		write: w,
	}
	go lp.read.Read(lp.rc)
	go lp.Process()
	go lp.write.Write(lp.wc)

	time.Sleep(30 * time.Second)
}
