package main

import (
	"bufio"
	"fmt"
	"os"
)

type fileWriter struct {
}

var file *os.File
var bwriter *bufio.Writer

func (w *fileWriter) open(path string, lineChannel chan string) error {
	var err error
	file, err = os.Create(path)
	if err != nil {
		return err
	}

	doneChannel = make(chan bool)
	bwriter = bufio.NewWriter(file)
	go w.writer(lineChannel, doneChannel)
	return nil
}

func (w *fileWriter) close() {
	_ = <-doneChannel
}

func (w *fileWriter) writer(lineChannel chan string, doneChannel chan bool) {
	for line := range lineChannel {
		fmt.Fprintln(bwriter, line)
	}
	bwriter.Flush()
	file.Close()
	doneChannel <- true

}
