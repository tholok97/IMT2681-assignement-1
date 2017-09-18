package main

/*
 * TESTS
	test -v (more detailed)
	test -cover	 (percentage)
	test -coverprofile=coverage.out (store coverage in file)
	teset cover -html=coverage.out (show colorcoded in browser)
 * empty languages, empty contributor list
 * mariusj is using ->>> url . Errorf
 * synmaxcols
*/

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
	Contributions int    `json:"contributions"`
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

	// TODO DEBUG
	if url == "https://api.github.com/repos/tholok97/the-t-files/languages" {
		return []byte(`{ "C++": 85890, "Vim script": 15368 }`), nil
	}

	switch url {
	case "https://api.github.com/repos/tholok97/the-t-files":
		return []byte(`{
			"id": 71177175,
			"name": "the-t-files",
			"full_name": "tholok97/the-t-files",
			"owner": {
				"login": "tholok97",
				"id": 22896227,
				"avatar_url": "https://avatars0.githubusercontent.com/u/22896227?v=4",
				"gravatar_id": "",
				"url": "https://api.github.com/users/tholok97",
				"html_url": "https://github.com/tholok97",
				"followers_url": "https://api.github.com/users/tholok97/followers",
				"following_url": "https://api.github.com/users/tholok97/following{/other_user}",
				"gists_url": "https://api.github.com/users/tholok97/gists{/gist_id}",
				"starred_url": "https://api.github.com/users/tholok97/starred{/owner}{/repo}",
				"subscriptions_url": "https://api.github.com/users/tholok97/subscriptions",
				"organizations_url": "https://api.github.com/users/tholok97/orgs",
				"repos_url": "https://api.github.com/users/tholok97/repos",
				"events_url": "https://api.github.com/users/tholok97/events{/privacy}",
				"received_events_url": "https://api.github.com/users/tholok97/received_events",
				"type": "User",
				"site_admin": false
			},
			"private": false,
			"html_url": "https://github.com/tholok97/the-t-files",
			"description": "Div. C++ prosjekter",
			"fork": false,
			"url": "https://api.github.com/repos/tholok97/the-t-files",
			"forks_url": "https://api.github.com/repos/tholok97/the-t-files/forks",
			"keys_url": "https://api.github.com/repos/tholok97/the-t-files/keys{/key_id}",
			"collaborators_url": "https://api.github.com/repos/tholok97/the-t-files/collaborators{/collaborator}",
			"teams_url": "https://api.github.com/repos/tholok97/the-t-files/teams",
			"hooks_url": "https://api.github.com/repos/tholok97/the-t-files/hooks",
			"issue_events_url": "https://api.github.com/repos/tholok97/the-t-files/issues/events{/number}",
			"events_url": "https://api.github.com/repos/tholok97/the-t-files/events",
			"assignees_url": "https://api.github.com/repos/tholok97/the-t-files/assignees{/user}",
			"branches_url": "https://api.github.com/repos/tholok97/the-t-files/branches{/branch}",
			"tags_url": "https://api.github.com/repos/tholok97/the-t-files/tags",
			"blobs_url": "https://api.github.com/repos/tholok97/the-t-files/git/blobs{/sha}",
			"git_tags_url": "https://api.github.com/repos/tholok97/the-t-files/git/tags{/sha}",
			"git_refs_url": "https://api.github.com/repos/tholok97/the-t-files/git/refs{/sha}",
			"trees_url": "https://api.github.com/repos/tholok97/the-t-files/git/trees{/sha}",
			"statuses_url": "https://api.github.com/repos/tholok97/the-t-files/statuses/{sha}",
			"languages_url": "https://api.github.com/repos/tholok97/the-t-files/languages",
			"stargazers_url": "https://api.github.com/repos/tholok97/the-t-files/stargazers",
			"contributors_url": "https://api.github.com/repos/tholok97/the-t-files/contributors",
			"subscribers_url": "https://api.github.com/repos/tholok97/the-t-files/subscribers",
			"subscription_url": "https://api.github.com/repos/tholok97/the-t-files/subscription",
			"commits_url": "https://api.github.com/repos/tholok97/the-t-files/commits{/sha}",
			"git_commits_url": "https://api.github.com/repos/tholok97/the-t-files/git/commits{/sha}",
			"comments_url": "https://api.github.com/repos/tholok97/the-t-files/comments{/number}",
			"issue_comment_url": "https://api.github.com/repos/tholok97/the-t-files/issues/comments{/number}",
			"contents_url": "https://api.github.com/repos/tholok97/the-t-files/contents/{+path}",
			"compare_url": "https://api.github.com/repos/tholok97/the-t-files/compare/{base}...{head}",
			"merges_url": "https://api.github.com/repos/tholok97/the-t-files/merges",
			"archive_url": "https://api.github.com/repos/tholok97/the-t-files/{archive_format}{/ref}",
			"downloads_url": "https://api.github.com/repos/tholok97/the-t-files/downloads",
			"issues_url": "https://api.github.com/repos/tholok97/the-t-files/issues{/number}",
			"pulls_url": "https://api.github.com/repos/tholok97/the-t-files/pulls{/number}",
			"milestones_url": "https://api.github.com/repos/tholok97/the-t-files/milestones{/number}",
			"notifications_url": "https://api.github.com/repos/tholok97/the-t-files/notifications{?since,all,participating}",
			"labels_url": "https://api.github.com/repos/tholok97/the-t-files/labels{/name}",
			"releases_url": "https://api.github.com/repos/tholok97/the-t-files/releases{/id}",
			"deployments_url": "https://api.github.com/repos/tholok97/the-t-files/deployments",
			"created_at": "2016-10-17T20:13:35Z",
			"updated_at": "2016-11-03T16:35:36Z",
			"pushed_at": "2017-04-24T21:50:48Z",
			"git_url": "git://github.com/tholok97/the-t-files.git",
			"ssh_url": "git@github.com:tholok97/the-t-files.git",
			"clone_url": "https://github.com/tholok97/the-t-files.git",
			"svn_url": "https://github.com/tholok97/the-t-files",
			"homepage": "",
			"size": 3511,
			"stargazers_count": 0,
			"watchers_count": 0,
			"language": "C++",
			"has_issues": true,
			"has_projects": true,
			"has_downloads": true,
			"has_wiki": true,
			"has_pages": false,
			"forks_count": 0,
			"mirror_url": null,
			"open_issues_count": 0,
			"forks": 0,
			"open_issues": 0,
			"watchers": 0,
			"default_branch": "master",
			"network_count": 0,
			"subscribers_count": 1
		}
		`), nil
	case "https://api.github.com/repos/tholok97/the-t-files/languages":
		return []byte(`{ "C++": 85890, "Vim script": 15368 }`), nil
	case "https://api.github.com/repos/tholok97/the-t-files/contributors":
		return []byte(`[ { "login": "tholok97", "id": 22896227, "avatar_url": "https://avatars0.githubusercontent.com/u/22896227?v=4", "gravatar_id": "", "url": "https://api.github.com/users/tholok97", "html_url": "https://github.com/tholok97", "followers_url": "https://api.github.com/users/tholok97/followers", "following_url": "https://api.github.com/users/tholok97/following{/other_user}", "gists_url": "https://api.github.com/users/tholok97/gists{/gist_id}", "starred_url": "https://api.github.com/users/tholok97/starred{/owner}{/repo}", "subscriptions_url": "https://api.github.com/users/tholok97/subscriptions", "organizations_url": "https://api.github.com/users/tholok97/orgs", "repos_url": "https://api.github.com/users/tholok97/repos", "events_url": "https://api.github.com/users/tholok97/events{/privacy}", "received_events_url": "https://api.github.com/users/tholok97/received_events", "type": "User", "site_admin": false, "contributions": 120 } ] `), nil
	}

	fmt.Println("requesting github.... (bad?)")

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
