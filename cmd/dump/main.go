package main

import (
	"bytes"
	"errors"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"strconv"
)

type service url.URL

type serviceReader struct {
	s      *service
	page   int
	size   int
	buffer bytes.Buffer
}

func (s *service) getPage(page, size int, buffer *bytes.Buffer) error {
	return nil
}

func (s *service) NewReader() *serviceReader {
	return &serviceReader{
		s: s,
	}
}

func (sr *serviceReader) Read(b []byte) (int, error) {
	if sr.buffer.Len() > 0 {

	}
	err := sr.s.getPage(sr.page, sr.size, &sr.buffer)

	return 0, err
}

var (
	address  = flag.String("address", "192.168.178.72:8080", "address and port of the service")
	pagesize = flag.Int("pagesize", 1000, "number of measurements to collect per request, maximum size is 10000")
	dir      = flag.String("dir", "./", "directory to save the measurements to, defaults to current directory")
	prefix   = flag.String("prefix", "Measure", "filename prefix")
)

func targetOk() bool {
	fi, err := os.Stat(*dir)
	return err == nil && fi.IsDir()
}

var baseUrl string

func getBaseUrl() string {
	if len(baseUrl) == 0 {
		baseUrl = fmt.Sprintf("http://%s/measurements/?size=%d", *address, *pagesize)
	}
	return baseUrl
}

func addrOk() bool {
	_, err := url.Parse(getBaseUrl())
	return err == nil
}

var (
	EOF          = errors.New("EOF")
	ErrGetFailed = errors.New("http get failed, status code was other than 200")
)

func getPage(page int) ([]byte, error) {
	u := getBaseUrl() + "&page=" + strconv.Itoa(page)
	resp, err := http.Get(u)
	if err != nil {
		return nil, err
	}
	if resp.StatusCode != http.StatusOK {
		return nil, ErrGetFailed
	}
	defer resp.Body.Close()
	return ioutil.ReadAll(resp.Body)
}

func writePage(page int, content []byte) error {
	fn := filepath.Join(*dir, fmt.Sprintf("%s-%04d.json", *prefix, page))
	return ioutil.WriteFile(fn, content, os.ModePerm)
}

func export() error {
	for page := 0; ; page++ {
		c, err := getPage(page)
		if err != nil {
			if err == EOF {
				return nil
			}
			return err
		}
		log.Printf("getPage: page: %d, %d bytes", page, len(c))
		// Write content to file
		if err := writePage(page, c); err != nil {
			return err
		}
	}
}

func main() {
	flag.Parse()
	if *pagesize > 1000 {
		*pagesize = 10000
	}
	// Test if target dir exists
	if !targetOk() {
		log.Fatal("dir is invalid")
	}
	// Test the address
	if !addrOk() {
		log.Fatal("bad address")
	}
	err := export()
	if err != nil {
		log.Fatalf("export failed: %s", err.Error())
	}
	log.Println("done")
}
