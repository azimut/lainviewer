package main

import (
	"encoding/json"
	"flag"
)

type Thread struct {
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
	Fsize          uint64 `json:"fsize"`
	Filename       string `json:"filename"`
	Ext            string `json:"ext"`
	Tim            string `json:"tim"`
	md5            string
	Resto          int     `json:"resto"`
	ExtraFiles     []Extra `json:"extra_files"`
}

type Extra struct {
	tn_h     int
	tn_w     int
	h        int
	w        int
	Fsize    uint64 `json:"fsize"`
	Filename string `json:"filename"`
	Ext      string `json:"ext"`
	Tim      string `json:"tim"`
	md5      string
}

var timeout int
var uri string
var maxWidth int
var showAuthors bool

// TODO: flag color
// TODO: max width
func init() {
	flag.StringVar(&uri, "u", "", "url")
	flag.IntVar(&timeout, "t", 5, "timeout after seconds")
	flag.IntVar(&maxWidth, "w", 120, "max text width")
	flag.BoolVar(&showAuthors, "A", true, "show comment authors")
}

func main() {
	flag.Parse()
	if err := validate_uri(); err != nil {
		panic(err)
	}

	bytes, err := getUrl()
	if err != nil {
		panic(err)
	}

	var thread Thread
	if err := json.Unmarshal(bytes, &thread); err != nil {
		panic(err)
	}
	printOp(thread)
	print_comments(thread)
}
