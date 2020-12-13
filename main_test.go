package main

import (
	"encoding/json"
	"io/ioutil"
	"testing"
)

func TestMain(t *testing.T) {

	uri = "https://lainchan.org/music/res/6019.html"
	showAuthors = false
	if err := validate_uri(); err != nil {
		panic(err)
	}

	bytes, err := ioutil.ReadFile("testdata/6019.json")
	if err != nil {
		panic(err)
	}

	var data Rsp
	if err := json.Unmarshal(bytes, &data); err != nil {
		panic(err)
	}
	print_op(data)
	print_comments(data)
}