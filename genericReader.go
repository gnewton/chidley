package main

import (
	"bufio"
	"compress/bzip2"
	"compress/gzip"
	"io"
	"log"
	"net/http"
	"os"
	"strings"
)

type GenericReader struct {
	io.Reader
	reader io.Reader
}

func (ge *GenericReader) Read(p []byte) (n int, err error) {
	return ge.reader.Read(p)
}

func (ge *GenericReader) Open(name string) error {
	if name == "" {
		return nil
	}
	if strings.HasPrefix(name, "http://") || strings.HasPrefix(name, "https://") {
		return ge.openUrl(name)
	} else {
		return ge.openLocalFile(name)
	}
}

func (ge *GenericReader) Close() error {
	return nil
}

func (ge *GenericReader) openLocalFile(filename string) error {
	file, err := openIfExistsIsFileIsReadable(filename)
	if err != nil {
		return err
	}

	if strings.HasSuffix(filename, "bz2") {
		ge.reader = bufio.NewReader(bzip2.NewReader(bufio.NewReader(file)))
	} else {
		if strings.HasSuffix(filename, "gz") {
			reader, err := gzip.NewReader(bufio.NewReader(file))
			if err != nil {
				return err
			}
			ge.reader = bufio.NewReader(reader)
		} else {
			ge.reader = bufio.NewReader(file)
		}
	}
	return nil
}

func (gr *GenericReader) openUrl(url string) error {
	res, err := http.Get(url)
	if err != nil {
		return err
	}

	if res.StatusCode != 200 {
		log.Fatal("ERROR: bad http status code != 200: ", res.StatusCode, "   ", url)
		return nil

	}
	gr.reader = res.Body
	return nil
}

func openIfExistsIsFileIsReadable(fileName string) (*os.File, error) {
	file, err := os.Open(fileName) // For read access.
	log.Print(fileName)
	if err != nil {
		log.Print("Problem opening file: [" + fileName + "]")
		log.Print(err)
		return nil, err
	}
	if _, err := os.Stat(fileName); os.IsNotExist(err) {
		log.Print("Problem stat'ing file: [" + fileName + "]")
		return nil, err
	}

	fi, err := file.Stat()
	if err != nil {
		log.Print("Problem stat'ing file: [" + fileName + "]")
		return nil, err
	}

	fm := fi.Mode()
	if !fm.IsRegular() {
		error := new(InternalError)
		error.ErrorString = "Is directory, needs to be file: " + fileName
		log.Print(error.ErrorString)
		return nil, error
	}

	log.Print(fm.Perm().String())
	if fm.Perm().String()[7] != 'r' {
		error := new(InternalError)
		error.ErrorString = "Exists but unable to read: " + fileName
		log.Print(error.ErrorString)
		return nil, error
	}
	return file, nil
}

type InternalError struct {
	ErrorString string
}

func (ie *InternalError) Error() string {
	return "Error: " + ie.ErrorString
}
