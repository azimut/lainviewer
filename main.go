package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"time"

	markdown "github.com/MichaelMure/go-term-markdown"
	"github.com/dustin/go-humanize"
	"github.com/jaytaylor/html2text"
)

func human_time(unix int64) string {
	unix_time := time.Unix(unix, 0)
	return humanize.Time(unix_time)
}

const JSONFILE string = "/home/sendai/37647.json"

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

func main() {
	bytes, err := ioutil.ReadFile(JSONFILE)
	if err != nil {
		panic(err)
	}

	var data Rsp
	if err := json.Unmarshal(bytes, &data); err != nil {
		panic(err)
	}
	print_op(data, JSONFILE)
	print_comments(data)
}

// print_comments prints the rest of the messages
// TODO: sanitize author for trailing spaces at least
// NOTE: Resto field is useless, it ALWAYS has the mainID
func print_comments(data Rsp) {
	for _, value := range data.Posts[1:] {
		fmt.Printf("%s - %d\n", value.Author, value.No)
		fmt.Printf("%s\n\n\n", html2console(value.Comment))
	}
}

func html2console(raw string) string {
	mdcomment, _ := html2text.FromString(raw, html2text.Options{PrettyTables: true})
	return string(markdown.Render(mdcomment, 80, 3))
}

// print_op Prints the main thread post
func print_op(resp Rsp, url string) {
	// Header
	fmt.Printf("\ntitle: %s\nurl: %s\navatar: %s\n\n",
		resp.Posts[0].Title,
		url,
		resp.Posts[0].Filename+resp.Posts[0].Ext, // TODO: add domain
	)
	// Message
	fmt.Printf("%s\n", html2console(resp.Posts[0].Comment))
	// Footer
	fmt.Printf("\n%s - %s\n\n",
		resp.Posts[0].Author,
		human_time(resp.Posts[0].Time))
}
