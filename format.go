package main

import (
	"fmt"
	"log"
	"strconv"
	"strings"
	"time"

	markdown "github.com/MichaelMure/go-term-markdown"
	"github.com/PuerkitoBio/goquery"
	"github.com/dustin/go-humanize"
	"github.com/jaytaylor/html2text"
)

func Min(x, y int) int {
	if x > y {
		return y
	}
	return x
}

func Max(x, y int) int {
	if x > y {
		return x
	}
	return y
}

func human_time(unix int64) string {
	unix_time := time.Unix(unix, 0)
	return humanize.Time(unix_time)
}

func (m *Message) Parent() (int, error) {
	// Starts with the quote
	if !strings.HasPrefix(m.Comment, "<a onclick=") {
		return 0, nil
	}
	// Has only 1 quote
	doc, err := goquery.NewDocumentFromReader(strings.NewReader(m.Comment))
	if err != nil {
		return -1, err
	}
	if doc.Find("a[href*=res]").Size() != 1 {
		return 0, nil
	}
	// Extract ID of the quote
	sel, exists := doc.Find("a[href*=res]").First().Attr("href")
	if !exists {
		return 0, nil
	}
	num, err := strconv.Atoi(strings.Split(sel, "#")[1])
	if err != nil {
		return -1, err
	}
	return num, nil
}

// isOrphan if no parent exists for post, like linking elsewhere
func (c *Message) isOrphan(rsp Rsp) (bool, error) {
	found := true
	parentid, err := c.Parent()
	if err != nil {
		return false, err
	}
	for _, post := range rsp.Posts[1:] {
		if post.No == parentid {
			found = false
			break
		}
	}
	return found, nil
}

// getMedia returns a slice of media urls and original filename
func (m *Message) getMedia() (media []string) {
	if m.Filename == "" {
		return media
	}
	mainMedia, err := srcFilename(m.Tim, m.Ext)
	if err != nil {
		log.Print(err)
	}
	media = append(media, fmt.Sprintf("%s (%s)", mainMedia, m.Filename+m.Ext))
	for _, e := range m.ExtraFiles {
		extraMedia, err := srcFilename(e.Tim, e.Ext)
		if err != nil {
			log.Print(err)
		}
		media = append(media, fmt.Sprintf("%s (%s)", extraMedia, e.Filename+e.Ext))
	}
	return media
}

// print_comments prints the rest of the messages
// TODO: sanitize author for trailing spaces at least
// NOTE: Resto field is useless, it ALWAYS has the mainID
func print_comments(data Rsp) {
	for _, post := range data.Posts[1:] {
		parent, err := post.Parent()
		if err != nil {
			log.Print(err)
		}
		orphan, err := post.isOrphan(data)
		if err != nil {
			log.Print(err)
		}
		if parent == 0 || parent == data.Posts[0].No || orphan {
			print_comment(post, data, 0)
		}
	}
}

// print_comment prints the message provided as well as any children
func print_comment(msg Message, rsp Rsp, depth int) {
	parentId, err := msg.Parent()
	if err != nil {
		log.Print(err)
	}
	if parentId > 0 {
		fmt.Printf("%s", html2console(
			// TODO: properly remove the quote
			strings.Join(strings.Split(msg.Comment, "<br/>")[1:], ""),
			depth))
	} else {
		fmt.Printf("%s", html2console(msg.Comment, depth))
	}
	// Media
	for _, media := range msg.getMedia() {
		fmt.Printf(strings.Repeat(" ", Max(depth*3+1, 0))+"%s\n", media)
	}
	// Footer
	fmt.Printf(strings.Repeat(" ", Max(depth*3, 0))+">> %s - %d - %s\n\n",
		msg.Author, msg.No, human_time(msg.Time))
	// Try to find childs
	for _, othermsg := range rsp.Posts[1:] {
		otherParentId, err := othermsg.Parent()
		if err != nil {
			log.Print(err)
		}
		if msg.No == otherParentId {
			print_comment(othermsg, rsp, depth+1)
		}
	}
}

func html2console(raw string, depth int) string {
	mdcomment, _ := html2text.FromString(raw, html2text.Options{PrettyTables: true})
	// TODO: use console width
	return string(markdown.Render(mdcomment, 80, Max(depth*3+1, 0)))
}

// print_op Prints the main thread post
func print_op(resp Rsp) {
	// Header
	fmt.Printf("\ntitle: %s\nurl: %s\n", resp.Posts[0].Title, uri)
	for _, media := range resp.Posts[0].getMedia() {
		fmt.Printf("media: %s\n", media)
	}
	fmt.Println()
	// Message
	fmt.Printf("%s\n", html2console(resp.Posts[0].Comment, 1))
	// Footer
	fmt.Printf("%s - %s\n\n\n",
		resp.Posts[0].Author,
		human_time(resp.Posts[0].Time))
}
