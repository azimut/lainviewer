package main

import (
	"fmt"
	"time"

	markdown "github.com/MichaelMure/go-term-markdown"
	"github.com/dustin/go-humanize"
	"github.com/jaytaylor/html2text"
)

func human_time(unix int64) string {
	unix_time := time.Unix(unix, 0)
	return humanize.Time(unix_time)
}

// print_comments prints the rest of the messages
// TODO: sanitize author for trailing spaces at least
// NOTE: Resto field is useless, it ALWAYS has the mainID
func print_comments(data Rsp) {
	for _, value := range data.Posts[1:] {
		fmt.Printf("%s - %d - %s\n", value.Author, value.No, human_time(value.Time))
		fmt.Printf("%s\n\n\n", html2console(value.Comment))
	}
}

func html2console(raw string) string {
	mdcomment, _ := html2text.FromString(raw, html2text.Options{PrettyTables: true})
	// TODO: use console width
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
	fmt.Printf("%s - %s\n\n\n",
		resp.Posts[0].Author,
		human_time(resp.Posts[0].Time))
}
