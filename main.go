package main

import (
	"encoding/json"
	"flag"
)

type Rsp struct {
	Posts []Message
}

type Message struct {
	No             int    `json:"no"`
	Title          string `json:"sub"`
	Comment        string `json:"com"`
	Author         string `json:"name"`
	trip           string
	Time           int64 `json:"time"`
	omitted_posts  int
	omitted_images int
	sticky         int
	locked         int
	cyclical       string
	last_modified  int
	tn_h           int
	tn_w           int
	h              int
	w              int
	fsize          int
	Filename       string `json:"filename"`
	Ext            string `json:"ext"`
	tim            string
	md5            string
	Resto          int `json:"resto"`
}

var timeout int
var uri string
var maxWidth int

// TODO: flag color
// TODO: max width
func init() {
	flag.StringVar(&uri, "u", "", "url")
	flag.IntVar(&timeout, "t", 5, "timeout after seconds")
	flag.IntVar(&maxWidth, "w", 120, "max text width")
}

func main() {
	flag.Parse()
	if err := validate_uri(); err != nil {
		panic(err)
	}

	bytes, err := getUrl(uri)
	if err != nil {
		panic(err)
	}

	var data Rsp
	if err := json.Unmarshal(bytes, &data); err != nil {
		panic(err)
	}
	print_op(data, uri)
	print_comments(data)
}
