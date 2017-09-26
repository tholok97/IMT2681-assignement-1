package main

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
)

// Payload contains the information sent back to the user
type Payload struct {
	Project   string   `json:"project"`
	Owner     string   `json:"owner"`
	Committer string   `json:"comitter"`
	Commits   int      `json:"commits"`
	Language  []string `json:"language"`
}

// Contains information we wanted from the initial github request
type githubReposResponse struct {
	LanguagesURL string `json:"languages_url"` // url to langauges
}

// contains sub-object found in contributors reponse
type author struct {
	Login string `json:"login"`
}

// contains info we gain from the ".../stats/contributors" request on github
type contributor struct {
	Total  int    `json:"total"`
	Author author `json:"author"`
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
		if value.Total > top.Total {
			top = value
		}
	}

	return top
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
	top, commitErr := determineTopCommiter("https://api.github.com/repos/" + user + "/" + repo + "/stats/contributors")
	if commitErr == nil {
		pload.Committer = top.Author.Login
		pload.Commits = top.Total
	}

	// (try to) determine langauges. if unable, leave blank
	langauges, langErr := determineLanguages(ghReposResp.LanguagesURL)
	if langErr == nil {
		pload.Language = langauges
	}

	return pload, nil
}

// get json that url points to and return it as []byte. return error if
// anything goes wrong
func getJSON(url string) ([]byte, error) {

	// switch to return hardcoded results during testing. Copy-and-pasted
	//	in for simplicity

	/*
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
		case "https://api.github.com/repos/tholok97/the-t-files/stats/contributors":
			return []byte(`[
			  {
				  "total": 117,
				  "weeks": [
				  {
					  "w": 1476576000,
					  "a": 698,
					  "d": 42,
					  "c": 8
				  },
				  {
					  "w": 1477180800,
					  "a": 163,
					  "d": 25,
					  "c": 3
				  },
				  {
					  "w": 1477785600,
					  "a": 601,
					  "d": 268,
					  "c": 31
				  },
				  {
					  "w": 1478390400,
					  "a": 0,
					  "d": 0,
					  "c": 0
				  },
				  {
					  "w": 1478995200,
					  "a": 96,
					  "d": 41,
					  "c": 9
				  },
				  {
					  "w": 1479600000,
					  "a": 334,
					  "d": 335,
					  "c": 12
				  },
				  {
					  "w": 1480204800,
					  "a": 14,
					  "d": 9,
					  "c": 4
				  },
				  {
					  "w": 1480809600,
					  "a": 234,
					  "d": 19,
					  "c": 8
				  },
				  {
					  "w": 1481414400,
					  "a": 35,
					  "d": 33,
					  "c": 2
				  },
				  {
					  "w": 1482019200,
					  "a": 0,
					  "d": 0,
					  "c": 0
				  },
				  {
					  "w": 1482624000,
					  "a": 196,
					  "d": 11,
					  "c": 8
				  },
				  {
					  "w": 1483228800,
					  "a": 0,
					  "d": 0,
					  "c": 0
				  },
				  {
					  "w": 1483833600,
					  "a": 0,
					  "d": 0,
					  "c": 0
				  },
				  {
					  "w": 1484438400,
					  "a": 369,
					  "d": 82,
					  "c": 7
				  },
				  {
					  "w": 1485043200,
					  "a": 0,
					  "d": 0,
					  "c": 0
				  },
				  {
					  "w": 1485648000,
					  "a": 284,
					  "d": 135,
					  "c": 14
				  },
				  {
					  "w": 1486252800,
					  "a": 0,
					  "d": 0,
					  "c": 0
				  },
				  {
					  "w": 1486857600,
					  "a": 0,
					  "d": 0,
					  "c": 0
				  },
				  {
					  "w": 1487462400,
					  "a": 274,
					  "d": 23,
					  "c": 5
				  },
				  {
					  "w": 1488067200,
					  "a": 85,
					  "d": 28,
					  "c": 1
				  },
				  {
					  "w": 1488672000,
					  "a": 0,
					  "d": 0,
					  "c": 0
				  },
				  {
					  "w": 1489276800,
					  "a": 166,
					  "d": 120,
					  "c": 2
				  },
				  {
					  "w": 1489881600,
					  "a": 109,
					  "d": 50,
					  "c": 1
				  },
				  {
					  "w": 1490486400,
					  "a": 0,
					  "d": 0,
					  "c": 0
				  },
				  {
					  "w": 1491091200,
					  "a": 0,
					  "d": 0,
					  "c": 0
				  },
				  {
					  "w": 1491696000,
					  "a": 0,
					  "d": 0,
					  "c": 0
				  },
				  {
					  "w": 1492300800,
					  "a": 0,
					  "d": 0,
					  "c": 0
				  },
				  {
					  "w": 1492905600,
					  "a": 784,
					  "d": 162,
					  "c": 2
				  },
				  {
					  "w": 1493510400,
					  "a": 0,
					  "d": 0,
					  "c": 0
				  },
				  {
					  "w": 1494115200,
					  "a": 0,
					  "d": 0,
					  "c": 0
				  },
				  {
					  "w": 1494720000,
					  "a": 0,
					  "d": 0,
					  "c": 0
				  },
				  {
					  "w": 1495324800,
					  "a": 0,
					  "d": 0,
					  "c": 0
				  },
				  {
					  "w": 1495929600,
					  "a": 0,
					  "d": 0,
					  "c": 0
				  },
				  {
					  "w": 1496534400,
					  "a": 0,
					  "d": 0,
					  "c": 0
				  },
				  {
					  "w": 1497139200,
					  "a": 0,
					  "d": 0,
					  "c": 0
				  },
				  {
					  "w": 1497744000,
					  "a": 0,
					  "d": 0,
					  "c": 0
				  },
				  {
					  "w": 1498348800,
					  "a": 0,
					  "d": 0,
					  "c": 0
				  },
				  {
					  "w": 1498953600,
					  "a": 0,
					  "d": 0,
					  "c": 0
				  },
				  {
					  "w": 1499558400,
					  "a": 0,
					  "d": 0,
					  "c": 0
				  },
				  {
					  "w": 1500163200,
					  "a": 0,
					  "d": 0,
					  "c": 0
				  },
				  {
					  "w": 1500768000,
					  "a": 0,
					  "d": 0,
					  "c": 0
				  },
				  {
					  "w": 1501372800,
					  "a": 0,
					  "d": 0,
					  "c": 0
				  },
				  {
					  "w": 1501977600,
					  "a": 0,
					  "d": 0,
					  "c": 0
				  },
				  {
					  "w": 1502582400,
					  "a": 0,
					  "d": 0,
					  "c": 0
				  },
				  {
					  "w": 1503187200,
					  "a": 0,
					  "d": 0,
					  "c": 0
				  },
				  {
					  "w": 1503792000,
					  "a": 0,
					  "d": 0,
					  "c": 0
				  },
				  {
					  "w": 1504396800,
					  "a": 0,
					  "d": 0,
					  "c": 0
				  },
				  {
					  "w": 1505001600,
					  "a": 0,
					  "d": 0,
					  "c": 0
				  },
				  {
					  "w": 1505606400,
					  "a": 0,
					  "d": 0,
					  "c": 0
				  },
				  {
					  "w": 1506211200,
					  "a": 0,
					  "d": 0,
					  "c": 0
				  }
				  ],
				  "author": {
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
				  }
			  }
			  ] `), nil
		default:
			return []byte(""), nil
		}

		fmt.Println("requesting github.... (bad?)")
	*/

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
