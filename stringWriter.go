package main

type stringWriter struct {
	s string
}

func (w *stringWriter) open(s string, lineChannel chan string) error {
	doneChannel = make(chan bool)
	go w.writer(lineChannel, doneChannel)
	return nil
}

func (w *stringWriter) writer(lineChannel chan string, doneChannel chan bool) {
	for line := range lineChannel {
		w.s += line + "\n"
	}
	doneChannel <- true
}

func (w *stringWriter) close() {
	_ = <-doneChannel
}
