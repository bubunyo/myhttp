package main

import (
	"os"
	"strconv"
	"testing"
)

func TestGetProcNum(t *testing.T) {
	td := map[int]int{
		10:  10,
		5:   5,
		0:   DefaultProc,
		500: MaxProc,
	}

	for in, out := range td {
		p := getProcNum(in)
		if p != out {
			t.Errorf("Expected: %v, Got: %v", out, p)
		}
	}
}

func TestGetProcNum_Env(t *testing.T) {
	max := 25
	os.Setenv("MAX_PROC", strconv.Itoa(max))
	td := map[int]int{
		500: max,
	}
	for in, out := range td {
		res := getProcNum(in)
		if res != out {
			t.Errorf("Expected: %v, Got: %v", out, res)
		}
	}
}

func TestFixUrl(t *testing.T) {
	td := map[string]string{
		"adjust.com":         "http://adjust.com",
		"https://adjust.com": "https://adjust.com",
		"git://adjust.com":   "git://adjust.com",
	}
	for in, out := range td {
		res, ok := fixUrl(in)
		if res != out || !ok {
			t.Errorf("Expected: %v, %v, Got: %v, %v", out, true, res, ok)
		}
	}
}

func TestFixUrl_NotOK(t *testing.T) {
	td := map[string]string{
		"test":           "",
		"http://google":  "",
		"http://google.": "",
	}
	for in, out := range td {
		res, ok := fixUrl(in)
		if res != out || ok {
			t.Errorf("Expected: '%v', %v, Got: '%v', %v", out, false, res, ok)
		}
	}
}
