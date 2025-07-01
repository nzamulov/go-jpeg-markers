package gojpegmarkers

import (
	"bufio"
	"io"
	"log"
	"os"
)

func ReadFile(path string) ([]byte, error) {
	file, err := os.Open(path)
	if err != nil {
		panic(err)
	}

	defer func(file *os.File) {
		if err = file.Close(); err != nil {
			log.Printf("file.Close() error: %v\n", err)
		}
	}(file)

	return io.ReadAll(bufio.NewReader(file))
}
