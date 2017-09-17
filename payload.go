package main

import (
	"encoding/json"
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
	Contributors_url string `json: "contributors_url"` // url to contributors
	Languages_url    string `json: "languages_url"`    // url to langauges
}

// generate payload by requesting github for the info we need and then basing
// the payload off the resopnse
func generateResponsePayload(user, repo string) (Payload, error) {

	// (try to) get json from url pointing to $user, $repo
	respJson, getErr := getJson("https://api.github.com/repos/" + user +
		"/" + repo)
	if getErr != nil {
		return Payload{}, getErr
	}

	// (try to) fill ghReposResp with data from json response
	var ghReposResp githubReposResponse
	jsonErr := json.Unmarshal(respJson, &ghReposResp)
	if jsonErr != nil {
		return Payload{}, jsonErr
	}

	// initialize payload with user and project
	var pload = Payload{Owner: user, Project: repo}

	// (try to) determine top commiter. if unable, leave blank
	topCommiter, commitNum, commitErr :=
		determineTopCommiter(ghReposResp.Contributors_url)
	if commitErr == nil {
		pload.Comitter = topCommiter
		pload.Commits = commitNum
	}

	// (try to) determine langauges. if unable, leave blank
	langauges, langErr := determineLanguages(ghReposResp.Languages_url)
	if langErr == nil {
		pload.Languages = langauges
	}

	return pload, nil
}

// based on url from github payload ^, get top commiters name and num of commits
func determineTopCommiter(url string) (string, int, error) {
	return "julenissen", 923, nil
}

// try to get the languages used in the project based on url pointing to json
// info. return error if anything goes wrong
func determineLanguages(url string) ([]string, error) {

	// (try to) get json from url as []byte
	respJson, getErr := getJson(url)
	if getErr != nil {
		return nil, getErr
	}

	// (try to) unmarshal json info into map of { language: bytes }
	var data map[string]interface{}
	marshalErr := json.Unmarshal([]byte(respJson), &data)
	if marshalErr != nil {
		return nil, marshalErr
	}

	// fill string slice with indexes from map (the languages) and return it
	ret := make([]string, 0)
	for index := range data {
		ret = append(ret, index)
	}

	return ret, nil
}

// get json that url points to and return it as []byte. return error if
// anything goes wrong
func getJson(url string) ([]byte, error) {

	// (try to) get response from url
	resp, getErr := http.Get(url)
	if getErr != nil {
		return nil, getErr
	}

	// we need to close the body when we're done. defer it
	defer resp.Body.Close()

	// (try to) read []byte from response
	bodyBytes, readErr := ioutil.ReadAll(resp.Body)
	if readErr != nil {
		return nil, readErr
	}

	return bodyBytes, nil
}
