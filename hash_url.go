package main

import (
	"crypto/md5"
	"encoding/hex"
	"fmt"
	"io"
	"net/http"
)

type Result struct {
	WorkerId int
	Url      string
	Hash     string
	Err      error
}

type Job struct {
	Req chan string
	Res chan Result
}

func hashWorker(proc int, in chan string) chan Result {
	out := make(chan Result)
	job, result := NewWorkerPool(proc)
	go func() {
		go func() {
			for r := range result {
				out <- r
			}
		}()
		for s := range in {
			job <- s
		}
		close(out)
	}()
	return out
}

//Create a pool of n workers with a maximum pool size of WorkerCount if size exceeds WorkerCount
//if not, pool size == size
func NewWorkerPool(proc int) (chan<- string, <-chan Result) {
	job := Job{
		Req: make(chan string),
		Res: make(chan Result),
	}
	for wid := 1; wid <= proc; wid++ {
		go Worker(wid, job)
	}
	return job.Req, job.Res
}

// Workers executes a task in a pool by passing arguments to the task executor
func Worker(id int, job Job) {
	for url := range job.Req {
		res, err := fetchHashUrl(url)
		if err != nil {
			err = fmt.Errorf("[[Error]] Worker Id: %v, Url: %v, Msg: %v", id, url, err.Error())
			job.Res <- Result{id, url, "", err}
		} else {
			job.Res <- Result{id, url, res, nil}
		}
	}
}

func fetchHashUrl(url string) (string, error) {
	resp, err := http.Get(url)
	if err != nil {
		return "", fmt.Errorf("error fetching url: %v", err.Error())
	}
	defer resp.Body.Close()
	hash := md5.New()
	if _, err := io.Copy(hash, resp.Body); err != nil {
		return "", fmt.Errorf("error creating hash: %v", err.Error())
	}
	return hex.EncodeToString(hash.Sum(nil)), nil
}
