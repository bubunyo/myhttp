package main

import (
	"flag"
	"fmt"
	"net/url"
	"os"
	"strconv"
	"strings"
	"sync"
)

const MaxProc = 30

var p = flag.Int("parallel", 10, "Number of jobs to run in parallel")

func main() {
	flag.Parse()
	procNum := getProcNum()
	in := make(chan string)
	out := hashWorker(procNum, in)
	var wg sync.WaitGroup
	go func() {
		for r := range out {
			wg.Done()
			fmt.Printf("%s %s\n", r.Url, r.Hash)
		}
	}()
	for i := 1; i < len(os.Args); i++ {
		if s, ok := fixUrl(os.Args[i]); ok {
			wg.Add(1)
			in <- s
		}
	}
	wg.Wait()
	close(in)
}

func fixUrl(test string) (string, bool) {
	u, err := url.Parse(test)
	if err != nil || !strings.Contains(u.String(), ".") {
		return "", false
	}
	if u.Scheme == "" {
		u.Scheme = "http"
	}
	return u.String(), true
}

func getProcNum() int {
	max := MaxProc
	v, ok := os.LookupEnv("MAX_PROC")
	if ok {
		m, err := strconv.Atoi(v)
		if err == nil {
			max = m
		}
	}
	if *p < max {
		return *p
	}
	return max
}
