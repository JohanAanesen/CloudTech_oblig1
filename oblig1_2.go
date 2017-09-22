package main

import (
	"os"
	"net/http"
	"strings"
	"fmt"
	"encoding/json"
	"io"
)

const GITHUB_URL = "https://api.github.com/repos/"
const COMMIT_URL = "/contributors"
const LANG_URL = "/languages"

type Payload struct {
	Project string 		`json:"project"`
	Owner 	string 		`json:"owner"`
	Committer string	`json:"committer"`
	Commits int			`json:"commits"`
	Language []string 	`json:"language"`
}

func getCommitter(r io.Reader)(string, int, error){

	type Data struct{
		Login string
		Contributions int
	}

	var data []Data

	err := json.NewDecoder(r).Decode(&data)

	if err != nil{
		fmt.Printf("rip committer \n", err)
		os.Exit(1)
	}

	//I only need the info from first instance
	return data[0].Login, data[0].Contributions, nil
}

func getOwner(r io.Reader)(string, error){

	//this is shady, but it works i guess 0:)
	//json reads an object inside an object, so we need to decode it into object inside an object
	type own1 struct{
		Login string
	}

	type Data struct{
		Owner own1
	}

	var data Data

	err := json.NewDecoder(r).Decode(&data)

	if err != nil{
		fmt.Printf("rip owner\n", err)
		os.Exit(1)
	}

	return data.Owner.Login, nil
}

func getLang(r io.Reader)([]string, error){

	type Data map[string]interface{}

	var data Data

	err := json.NewDecoder(r).Decode(&data)

	if err != nil{
		fmt.Printf("rip lang\n", err)
		os.Exit(1)
	}

	var lang []string

	for r:= range data {
		lang = append(lang, r)
	}

	//return array with languages
	return lang, nil
}

func HandleOblig(w http.ResponseWriter, r *http.Request){
	http.Header.Add(w.Header(), "content-type", "application/json")
	URL := strings.Split(r.URL.Path, "/")

	json1, err := http.Get(GITHUB_URL + URL[4] + "/" + URL[5])
	json1io := json1.Body
	json2, err := http.Get(GITHUB_URL + URL[4] + "/" + URL[5] + COMMIT_URL)
	json2io := json2.Body
	json3, err := http.Get(GITHUB_URL + URL[4] + "/" + URL[5] + LANG_URL)
	json3io := json3.Body


	var payload Payload

	owner, err := getOwner(json1io)
	committer, commits, err := getCommitter(json2io)
	language, err := getLang(json3io)


	if err != nil{
		fmt.Printf("rip payload\n", err)
		return
	}

	payload.Project = "github.com/" + URL[4] + "/" + URL[5]
	payload.Owner = owner
	payload.Committer = committer
	payload.Commits = commits
	payload.Language = language

	json.NewEncoder(w).Encode(payload)
}

func main() {
	port := os.Getenv("PORT")
	http.HandleFunc("/projectinfo/v1/", HandleOblig)
	http.ListenAndServe(":"+port, nil)
//	http.ListenAndServe(":8080", nil)
}