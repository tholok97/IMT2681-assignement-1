package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
)

// associated user account/org, indication of prog lang(s) used, account name
// of top committer (comitter with largest number of commits to the project)
type Payload struct {
	Project   string   `json: "project"`
	Owner     string   `json: "owner"`
	Comitter  string   `json "comitter"`
	Commits   int      `json: "commits"`
	Languages []string `json: "languages"`
}

// Contains information we wanted from the initial github request
type githubReposResponse struct {
	contributors_url string `json: "contributors_url"` // url to contributors
	languages_url    string `json: "languages_url"`    // url to langauges
}

// generate payload by requesting github for the info we need and then basing
// the payload off the resopnse
func generateResponsePayload(user, repo string) Payload {

	// make request
	resp, err := http.Get("https://api.github.com/repos/" + user + "/" + repo)
	defer resp.Body.Close() // we need to close it when we're done

	// error cheks:

	if err != nil {
		return Payload{}
	}

	if resp.StatusCode != http.StatusOK {
		return Payload{}
	}

	// get body as bytes, and return
	body, err := ioutil.ReadAll(resp.Body)

	var ghReposResp githubReposResponse

	jsonErr := json.Unmarshal(body, &ghReposResp)
	if jsonErr != nil {
		return Payload{}
	}

	var pload = Payload{Owner: user, Project: repo}

	topCommiter, commitNum := determineTopCommiter(ghReposResp.contributors_url)

	pload.Comitter = topCommiter
	pload.Commits = commitNum

	langauges := determineLanguages(ghReposResp.languages_url)

	pload.Languages = langauges

	return pload
}

// based on url from github payload ^, get top commiters name and num of commits
func determineTopCommiter(url string) (string, int) {
	return "julenissen", 923
}

func determineLanguages(url string) []string {

	resp, err := http.Get("http://localhost:8080/projectinfo/v1/github.com/tholok97/the-t-files/languages")

	if err != nil {
		fmt.Println("koko1")
		return []string{}
	}
	defer resp.Body.Close()

	responseBody, err2 := ioutil.ReadAll(resp.Body)

	if err2 != nil {
		fmt.Println("koko1")
		return []string{}
	}

	var data map[string]interface{}
	err = json.Unmarshal([]byte(responseBody), &data)
	if err != nil {
		panic(err)
	}

	//ret := make([]string, 0)
	for index, value := range data {
		fmt.Println(index, ": ", value)
	}

	return []string{"C++", "Vimscript"}
}
