package main

import (
	"testing"
	"net/http"
)

func TestgetCommitter(t *testing.T){
	var inp = []string{"https://api.github.com/repos/apache/kafka", "https://api.github.com/repos/google/go-github", "https://github.com/"}
	var out = []string{"ijuma", "willnorris", ""}

	for i := range inp {
		json1, err := http.Get(inp[i])

		if err != nil{
			t.Fatalf("Error: ", err)
			return
		}

		owner, err := getOwner(json1.Body)

		if owner != out[i] {
			t.Fatalf("ERROR expected: %s but got: %s", out[i], owner)
		}
	}
}