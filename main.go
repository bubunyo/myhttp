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
const DefaultProc = 10

func main() {
	var proc int
	flag.IntVar(&proc, "parallel", DefaultProc, "Number of jobs to run in parallel")
	flag.Parse()
	procNum := getProcNum(proc)
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
	s := u.String()
	if err != nil || !strings.Contains(s, ".") || strings.HasSuffix(s, ".") {
		return "", false
	}
	if u.Scheme == "" {
		u.Scheme = "http"
	}
	return u.String(), true
}

func getProcNum(proc int) int {
	max := MaxProc
	v, ok := os.LookupEnv("MAX_PROC")
	if ok {
		m, err := strconv.Atoi(v)
		if err == nil {
			max = m
		}
	}
	if proc == 0 {
		return DefaultProc
	}
	if proc > 0 && proc < max {
		return proc
	}
	return max
}
