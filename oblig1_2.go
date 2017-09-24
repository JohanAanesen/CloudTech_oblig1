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
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"
)

//GitHubURL https://api.github.com/repos/
const GitHubURL = "https://api.github.com/repos/"

//CommitURL /contributors
const CommitURL = "/contributors"

//LangURL /languages
const LangURL = "/languages"

//Payload structure, as in assignment spec
type Payload struct {
	Project   string   `json:"project"`
	Owner     string   `json:"owner"`
	Committer string   `json:"committer"`
	Commits   int      `json:"commits"`
	Language  []string `json:"language"`
}

func getCommitter(r io.Reader) (string, int, error) {

	//Data structure, a string and an int
	type Data struct {
		Login         string `json:"login"`
		Contributions int    `json:"contributions"`
	}

	//data object
	var data []Data

	//json decoder
	err := json.NewDecoder(r).Decode(&data)

	//err handler
	if err != nil {
		fmt.Printf("Something went wrong with the JSON decoder: %s\n", err)
	}

	//I only need the info from first instance
	return data[0].Login, data[0].Contributions, err
}

func getOwner(r io.Reader) (string, error) {

	//this is shady, but it works i guess 0:)
	//json reads an object inside an object, so we need to decode it into object inside an object
	//yes that's right I'm tired and it's 3AM.
	type own1 struct {
		Login string `json:"login"`
	}

	type Data struct {
		Owner own1
	}

	//data object
	var data Data

	//json decoder
	err := json.NewDecoder(r).Decode(&data)

	//err handler
	if err != nil {
		fmt.Printf("Something went wrong with the JSON decoder: %s\n", err)

	}

	//returns data login and error (if there is any error that is)
	return data.Owner.Login, err
}

func getLang(r io.Reader) ([]string, error) {

	//I'm using map because it lets me assign a object without knowing the variables, different from getCommitter()
	type Data map[string]interface{}

	//data object
	var data Data

	//json decoder
	err := json.NewDecoder(r).Decode(&data)

	//err handler
	if err != nil {
		fmt.Printf("Something went wrong with the JSON decoder: %s\n", err)
	}

	//lang array to hold all the languages
	var lang []string
	//loops through and adds the languages
	for r := range data {
		lang = append(lang, r)
	}

	//return array with languages
	return lang, err
}

func checkNotFound(r io.Reader) (string, error) {
	//map
	type Data map[string]string
	//data object
	var data Data
	//json decoder
	err := json.NewDecoder(r).Decode(&data)
	//err handler
	if err != nil {
		fmt.Printf("Something went wrong with the JSON decoder: %s\n", err)
	}

	//sent just the message part of the map
	check := data["message"]
	//return string and error
	return check, err

}

//HandleOblig primary function
func HandleOblig(w http.ResponseWriter, r *http.Request) {
	//content-type set to JSON
	http.Header.Add(w.Header(), "content-type", "application/json")

	//URL parts, 1 is projectinfo, 2 is v1, 3 is github.com and then the 2 variables
	URL := strings.Split(r.URL.Path, "/")

	//error failsafe
	if URL[3] != "github.com" {
		http.Error(w, "Need github.com in the url after v1", http.StatusBadRequest)
		return
	}

	//more failsafes
	if len(URL) < 6 {
		http.Error(w, "Incomplete URL", http.StatusBadRequest)
		return
	}

	//GET requests, URL[4] and URL [5] is APACHE and KAFKA
	json1, err := http.Get(GitHubURL + URL[4] + "/" + URL[5])
	json2, err := http.Get(GitHubURL + URL[4] + "/" + URL[5] + CommitURL)
	json3, err := http.Get(GitHubURL + URL[4] + "/" + URL[5] + LangURL)

	failsafe, err := checkNotFound(json1.Body)

	if failsafe == "Not Found" {
		http.Error(w, "Repo not found", http.StatusBadRequest)
		return
	}

	if json1.Body == nil {
		http.Error(w, "Need a JSON body", http.StatusBadRequest)
		return
	}

	//populating variables
	owner, err := getOwner(json1.Body)
	committer, commits, err := getCommitter(json2.Body)
	language, err := getLang(json3.Body)

	//error handler
	if err != nil {
		fmt.Printf("Something went wrong %s\n", err)
		http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
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
