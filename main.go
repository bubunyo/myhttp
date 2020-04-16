package main

import (
	"flag"
	"fmt"
	"net/http"
	"net/url"
	"os"
	"strconv"
	"strings"
	"sync"
	"time"
)

const (
	MaxProc       = 30
	DefaultProc   = 10
	ClientTimeout = 5 * time.Second
)

func main() {
	var proc int
	flag.IntVar(&proc, "parallel", DefaultProc, "Number of jobs to run in parallel")
	flag.Parse()
	proc = getProcNum(proc)

	n := Network{
		Client: &http.Client{
			Timeout: ClientTimeout,
		},
	}

	job, result := NewWorkerPool(n, proc)

	var wg sync.WaitGroup
	go func() {
		for r := range result {
			wg.Done()
			fmt.Printf("%s %s\n", r.Url, r.Hash)
		}
	}()
	for i := 1; i < len(os.Args); i++ {
		if s, ok := fixUrl(os.Args[i]); ok {
			wg.Add(1)
			job <- s
		}
	}
	wg.Wait()
	close(job)
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
