package main

import (
	"bytes"
	"io/ioutil"
	"net/http"
	"strings"
	"testing"
)

type RoundTripFunc func(req *http.Request) *http.Response

// RoundTrip .
func (f RoundTripFunc) RoundTrip(req *http.Request) (*http.Response, error) {
	return f(req), nil
}

//NewTestClient returns *http.Client with Transport replaced to avoid making real calls
func NewTestClient(fn RoundTripFunc) *http.Client {
	return &http.Client{
		Transport: RoundTripFunc(fn),
	}
}

func TestFetchUrl(t *testing.T) {
	url := "http://example.com"

	td := map[string]string{
		url: "<html><body>Hello World!</body></html>",
	}
	client := NewTestClient(func(req *http.Request) *http.Response {
		// Test request parameters
		return &http.Response{
			StatusCode: 200,
			// Send response to be tested
			Body: ioutil.NopCloser(bytes.NewBufferString(td[req.URL.String()])),
			// Must be set to non-nil value or it panics
			Header: make(http.Header),
		}
	})

	n := Network{Client: client}
	for in, out := range td {
		res, _ := n.fetchUrl(in)
		if res != hashResp(strings.NewReader(out)) {
			t.Errorf("Expected: '%v',  Got: '%v' ", out, res)
		}
	}
}

func TestHashResp(t *testing.T) {
	td := map[string][]string{
		"http://example.com": {"<html><body>Hello World!</body></html>", "aa89faed549b8b6424790201ad3a4a3f"},
	}

	for in, out := range td {
		h := hashResp(strings.NewReader(td[in][0]))
		if h != out[1] {
			t.Errorf("Expected: '%v',  Got: '%v' ", out, h)
		}
	}
}
