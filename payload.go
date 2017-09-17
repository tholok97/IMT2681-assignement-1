package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
)

// Payload contains the information sent back to the user
type Payload struct {
	Project   string   `json:"project"`
	Owner     string   `json:"owner"`
	Comitter  string   `json:"comitter"`
	Commits   int      `json:"commits"`
	Languages []string `json:"languages"`
}

// Contains information we wanted from the initial github request
type githubReposResponse struct {
	ContributorsURL string `json:"contributors_url"` // url to contributors
	LanguagesURL    string `json:"languages_url"`    // url to langauges
}

// contains info we gain from the "/contributors" request on github
type contributor struct {
	Login         string `json:"login"`
	Contributions int    `json:"contributionsi"`
}

// returns top contributor in slice of contributors. nil if empty
func topContributor(contrs []contributor) contributor {

	// empty slice: give up
	if len(contrs) < 1 {
		return contributor{}
	}

	// find top contributor
	top := contrs[0]
	for _, value := range contrs {
		if value.Contributions > top.Contributions {
			top = value
		}
	}

	return top
}

// generate payload by requesting github for the info we need and then basing
// the payload off the resopnse
func generateResponsePayload(user, repo string) (Payload, error) {

	// (try to) get json from url pointing to $user, $repo
	respJSON, getErr := getJSON("https://api.github.com/repos/" + user +
		"/" + repo)
	if getErr != nil {
		return Payload{}, getErr
	}

	// (try to) fill ghReposResp with data from json response
	var ghReposResp githubReposResponse
	jsonErr := json.Unmarshal(respJSON, &ghReposResp)
	if jsonErr != nil {
		return Payload{}, jsonErr
	}

	// initialize payload with user and project
	var pload = Payload{Owner: user, Project: repo}

	// (try to) determine top commiter. if unable, leave blank
	top, commitErr := determineTopCommiter(ghReposResp.ContributorsURL)
	if commitErr == nil {
		pload.Comitter = top.Login
		pload.Commits = top.Contributions
	}

	// (try to) determine langauges. if unable, leave blank
	langauges, langErr := determineLanguages(ghReposResp.LanguagesURL)
	if langErr == nil {
		pload.Languages = langauges
	}

	return pload, nil
}

// based on url from github payload ^, get top commiters name and num of commits
func determineTopCommiter(url string) (contributor, error) {

	// (try to) get json ffrom url as []byte
	respJSON, getErr := getJSON(url)
	if getErr != nil {
		return contributor{}, getErr
	}

	// (try to) parse json into list of contributors
	contrs := make([]contributor, 0)
	marshalErr := json.Unmarshal(respJSON, &contrs)
	if marshalErr != nil {
		return contributor{}, marshalErr
	}

	// find top contributor. (if contrs is empty, give up...)
	top := topContributor(contrs)
	if top.Login == "" {
		return contributor{}, fmt.Errorf("Empty contributors slice. Something is wrong")
	}

	return top, nil
}

// try to get the languages used in the project based on url pointing to json
// info. return error if anything goes wrong
func determineLanguages(url string) ([]string, error) {

	// (try to) get json from url as []byte
	respJSON, getErr := getJSON(url)
	if getErr != nil {
		return nil, getErr
	}

	// (try to) unmarshal json info into map of { language: bytes }
	var data map[string]interface{}
	marshalErr := json.Unmarshal([]byte(respJSON), &data)
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
func getJSON(url string) ([]byte, error) {

	// (try to) get response from url
	resp, getErr := http.Get(url)
	if getErr != nil || resp.StatusCode != http.StatusOK {
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
