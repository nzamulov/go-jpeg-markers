package main

import (
	"bufio"
	"context"
	"flag"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"strings"
	"time"
)

const HTTPTimeout = 3 * time.Minute

func isURL(path string) bool {
	return strings.HasPrefix(path, "http://") || strings.HasPrefix(path, "https://")
}

func output(markers []Marker) {
	for _, m := range markers {
		log.Printf("offset: %6x - %s\n", m.Offset, m.Comment)
	}
}

func readFile(path string) ([]byte, error) {
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

func downloadFile(link string) ([]byte, error) {
	httpClient := http.Client{Timeout: HTTPTimeout}

	req, err := http.NewRequestWithContext(context.Background(), http.MethodGet, link, nil)
	if err != nil {
		return []byte{}, err
	}

	req.Header.Set("User-Agent", "")

	resp, err := httpClient.Do(req)
	if err != nil {
		return []byte{}, err
	}

	defer func(Body io.ReadCloser) {
		if err = Body.Close(); err != nil {
			log.Printf("body.Close() error: %v\n", err)
		}
	}(resp.Body)

	if resp.StatusCode != http.StatusOK {
		return []byte{}, fmt.Errorf("error downloading %s: status code is %d", link, resp.StatusCode)
	}

	return io.ReadAll(resp.Body)
}

func run(path string) {
	var (
		img []byte
		err error
	)

	if isURL(path) {
		img, err = downloadFile(path)
	} else {
		img, err = readFile(path)
	}
	if err != nil {
		panic(err)
	}

	output(Scan(img))
}

func main() {
	log.SetPrefix("[go-jpeg-markers] ")
	log.SetFlags(log.Ldate | log.Ltime | log.Lmicroseconds | log.Lmsgprefix)

	path := flag.String("p", "", "path to JPEG file to parse (on FS or just link from Internet)")
	flag.Parse()

	if *path == "" {
		panic("no path was provided")
	}

	run(*path)
}
