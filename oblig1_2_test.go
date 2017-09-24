package main

import (
	"testing"
	"os"
)


func TestGetOwner(t *testing.T){
	var out = []string{"apache"}

	for i := range out {
		json1, err := os.Open("./test/kafka.json")

		if err != nil{
			t.Fatalf("Error: ", err)
			return
		}

		owner, err := getOwner(json1)

		if owner != out[i] {
			t.Fatalf("ERROR expected: %s but got: %s", out[i], owner)
		}
	}
}


func TestGetCommitter(t *testing.T){
	var out = []string{"ijuma"}
	var out2 = []int{315}

	for i := range out {
		json1, err := os.Open("./test/contributors.json")

		if err != nil{
			t.Fatalf("Error: %s", err)
			return
		}

		committer, commits, err := getCommitter(json1)

		if committer != out[i] && commits != out2[i] {
			t.Fatalf("ERROR expected: %s and %v but got: %s and %v", out[i], out2[i], committer, commits)
		}
	}
}

func TestGetLang(t *testing.T){
	type Data struct {
		out []string
	}
	var data Data

	data.out = []string{"Java"}

	for i := range data.out {
		json1, err := os.Open("./test/lang.json")

		if err != nil{
			t.Fatalf("Error: ", err)
			return
		}

		lang, err := getLang(json1)

		if lang[i] != data.out[i] {
			t.Fatalf("ERROR expected: %s but got: %s", data.out[i], lang[i])
		}
	}
}

func TestCheckNotFound(t *testing.T){

	data := "Not Found"

	json1, err := os.Open("./test/notfound.json")

	if err != nil{
		t.Fatalf("Error: ", err)
		return
	}

	notfound, err := checkNotFound(json1)

	if notfound != data {
		t.Fatalf("ERROR expected: %s but got: %s", data, notfound)
	}
}