package main

import (
	"encoding/json"
	"flag"
	"time"
)

var (
	showAuthors bool
	timeout     time.Duration
	uri         string
	userAgent   string
	width       uint
)

func init() {
	flag.StringVar(&userAgent, "A", "LainViewer/0.1", "user agent")
	flag.UintVar(&width, "w", 0, "width")
	flag.DurationVar(&timeout, "t", time.Second*5, "timeout after")
	flag.BoolVar(&showAuthors, "a", false, "show author comments")
}

func main() {
	if err := processFlags(); err != nil {
		panic(err)
	}
	if err := validateFlags(); err != nil {
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

	printOp(thread.Posts[0])
	print_comments(thread)
}
