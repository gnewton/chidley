package main

import (
	"fmt"
)

type stdoutWriter struct {
}

var doneChannel chan bool

func (w *stdoutWriter) open(s string, lineChannel chan string) error {
	doneChannel = make(chan bool)
	go w.writer(lineChannel, doneChannel)

	return nil
}

func (w *stdoutWriter) writer(lineChannel chan string, doneChannel chan bool) {
	for line := range lineChannel {
		fmt.Println(line)
	}
	doneChannel <- true
}

func (w *stdoutWriter) close() {
	_ = <-doneChannel
}
