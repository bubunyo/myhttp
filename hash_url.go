package main

import (
	"crypto/md5"
	"encoding/hex"
	"fmt"
	"io"
	"log"
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
type Network struct {
	Client *http.Client
}

//Create a pool of n workers with a maximum pool size of WorkerCount if size exceeds WorkerCount
//if not, pool size == size
func NewWorkerPool(n Network, proc int) (chan<- string, <-chan Result) {
	job := Job{
		Req: make(chan string),
		Res: make(chan Result),
	}
	for wid := 1; wid <= proc; wid++ {
		// dispatch number of workers equal to proc
		go Worker(n, wid, job)
	}
	return job.Req, job.Res
}

// Workers executes a task in a pool by passing arguments to the task executor
func Worker(n Network, id int, job Job) {
	for url := range job.Req {
		res, err := n.fetchUrl(url)
		if err != nil {
			err = fmt.Errorf("[[Error]] Worker Id: %v, Url: %v, Msg: %v", id, url, err.Error())
			job.Res <- Result{id, url, "", err}
		} else {
			job.Res <- Result{id, url, res, nil}
		}
	}
}

func (n Network) fetchUrl(url string) (string, error) {
	resp, err := n.Client.Get(url)
	if err != nil {
		return "", fmt.Errorf("error fetching url: %v", err.Error())
	}
	defer resp.Body.Close()
	return hashResp(resp.Body), nil
}

func hashResp(r io.Reader) string {
	hash := md5.New()
	if _, err := io.Copy(hash, r); err != nil {
		log.Fatalf("error creating hash: %v", err.Error())
	}
	return hex.EncodeToString(hash.Sum(nil))
}
