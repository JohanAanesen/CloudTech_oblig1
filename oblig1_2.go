/*
Cloud Technologies Assignment 1
Name: Johan Aanesen
Studnr: 473182

Q-Why is the file named oblig1_2?
A-Because i first finished off the assignment using the package go-github (github.com/google/go-github),
  but heroku went apeshit when trying to deploy it (after 4 hours of trying to deploy it I gave up).
  And that is why I redid the assignment without using go-github (hooray it deploys!)

 */

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

//payload structure, as in assignment spec
type Payload struct {
	Project string 		`json:"project"`
	Owner 	string 		`json:"owner"`
	Committer string	`json:"committer"`
	Commits int			`json:"commits"`
	Language []string 	`json:"language"`
}

func getCommitter(r io.Reader)(string, int, error){

	//Data structure, a string and an int
	type Data struct{
		Login string		`json:"login"`
		Contributions int	`json:"contributions"`
	}

	//data object
	var data []Data

	//json decoder
	err := json.NewDecoder(r).Decode(&data)

	//err handler
	if err != nil{
		fmt.Printf("Something went wrong with the JSON decoder\n", err)
		return "error", 0, err
	}

	//I only need the info from first instance
	return data[0].Login, data[0].Contributions, nil
}

func getOwner(r io.Reader)(string, error){

	//this is shady, but it works i guess 0:)
	//json reads an object inside an object, so we need to decode it into object inside an object
	//yes that's right I'm tired and it's 3AM.
	type own1 struct{
		Login string	`json:"login"`
	}

	type Data struct{
		Owner own1
	}

	//data object
	var data Data

	//json decoder
	err := json.NewDecoder(r).Decode(&data)

	//err handler
	if err != nil{
		fmt.Printf("Something went wrong with the JSON decoder\n", err)
		return "error", err
	}

	//returns data login and no error
	return data.Owner.Login, nil
}

func getLang(r io.Reader)([]string, error){

	//I'm using map because it lets me assign a object without knowing the variables, different from getCommitter()
	type Data map[string]interface{}

	//data object
	var data Data

	//json decoder
	err := json.NewDecoder(r).Decode(&data)

	//err handler
	if err != nil{
		fmt.Printf("Something went wrong with the JSON decoder\n", err)
		return nil, err
	}

	//lang array to hold all the languages
	var lang []string
	//loops through and adds the languages
	for r:= range data {
		lang = append(lang, r)
	}

	//return array with languages
	return lang, nil
}

func HandleOblig(w http.ResponseWriter, r *http.Request){
	//content-type because firefox and prettyprint
	http.Header.Add(w.Header(), "content-type", "application/json")

	//URL parts, 1 is projectinfo, 2 is v1, 3 is github.com and then the 2 variables
	URL := strings.Split(r.URL.Path, "/")


	//GET requests, URL[4] and URL [5] is APACHE and KAFKA
	json1, err := http.Get(GITHUB_URL + URL[4] + "/" + URL[5])
	json2, err := http.Get(GITHUB_URL + URL[4] + "/" + URL[5] + COMMIT_URL)
	json3, err := http.Get(GITHUB_URL + URL[4] + "/" + URL[5] + LANG_URL)

	//populating variables
	owner, err := getOwner(json1.Body)
	committer, commits, err := getCommitter(json2.Body)
	language, err := getLang(json3.Body)

	//error handler, not sure if it is needed here(?) - I'm new to this language
	if err != nil{
		fmt.Printf("Something went wrong\n", err)
		return
	}

	//payload object
	var payload Payload

	//populates the payload object
	payload.Project = URL[3] + "/" + URL[4] + "/" + URL[5]
	payload.Owner = owner
	payload.Committer = committer
	payload.Commits = commits
	payload.Language = language

	//encodes payload into "beautiful" json
	json.NewEncoder(w).Encode(payload)
}

func main() {
	port := os.Getenv("PORT")
	http.HandleFunc("/projectinfo/v1/", HandleOblig)
	http.ListenAndServe(":"+port, nil)
//	http.ListenAndServe(":8080", nil)
}