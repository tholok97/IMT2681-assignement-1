package main

/*
 * how do we test something that doesn't return anything
 * difference between unmarshal and decode? decode didn't accept incomplete
	struct
 * how to test something that connects to the internet
*/

import "testing"

func TestDetermineTopCommiter(t *testing.T) {

	// test urls:
	urlValid := "https://api.github.com/repos/tholok97/the-t-files/contributors"
	urlInvaid := "this is not an url! xD"
	urlMisguided := "https://api.github.com/repos/tholok97/the-t-files/"

	// test valid url
	topValid, errValid := determineTopCommiter(urlValid)
	if errValid != nil {
		t.Error("Gave error on valid url")
	} else if topValid.Login != "tholok97" {
		t.Error("Didn't return correct contributor")
	}

	// test invalid url
	_, errInvalid := determineTopCommiter(urlInvaid)
	if errInvalid == nil {
		t.Error("Didn't give error on invalid url")
	}

	// test url pointed at wrong json code
	_, errMisguided := determineTopCommiter(urlMisguided)
	if errMisguided == nil {
		t.Error("Didn't give error on invalid json")
	}
}

func TestDetermineLanguages(t *testing.T) {

	// test urls:
	urlValid := "https://api.github.com/repos/tholok97/the-t-files/languages"
	urlInvaid := "this is not an url! xD"
	urlMisguided := "https://api.github.com/repos/tholok97/the-t-files/"

	// test valid url
	langsValid, errValid := determineLanguages(urlValid)
	if errValid != nil {
		t.Error("Gave error on valid url")
	}

	if len(langsValid) == 0 {
		t.Error("Gave empty slice on valid url")
	}

	if langsValid[0] != "C++" {
		t.Error("Didn't have correct langauge as first element (C++)")
	}

	// test invalid url
	_, errInvalid := determineLanguages(urlInvaid)
	if errInvalid == nil {
		t.Error("Didn't give error on invalid url")
	}

	// test misguided url
	_, errMisguided := determineLanguages(urlMisguided)
	if errMisguided == nil {
		t.Error("Didnt' give error on invalid json")
	}
}

func TestTopContributor(t *testing.T) {

	// make example contributor list
	contrs := []contributor{
		{
			Login:         "not this one!",
			Contributions: 15,
		},
		{
			Login:         "thomas", // this one should be returned!
			Contributions: 923,
		},
		{
			Login:         "pkbuer",
			Contributions: 10,
		},
		{
			Login:         "",
			Contributions: 0,
		},
	}

	// test with list
	top := topContributor(contrs)
	if top != contrs[1] {
		t.Error("Contributor returned was not the top one")
	}

	// test with empty list
	emptyTop := topContributor(make([]contributor, 0))
	if (emptyTop != contributor{}) {
		t.Error("Contributor should be empty")
	}
}
