package main

import (
	"fmt"
	"log"
	"strconv"
	"strings"
	"time"

	text "github.com/MichaelMure/go-term-text"
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

func humanTime(unix int64) string {
	return humanize.Time(time.Unix(unix, 0))
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
func (m *Message) isOrphan(rsp Thread) (bool, error) {
	found := true
	parentid, err := m.Parent()
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
	mediaUrl, err := srcFilename(m.Tim, m.Ext)
	if err != nil {
		log.Print(err)
		return media
	}
	mediaSize := humanize.Bytes(m.Fsize)
	media = append(media, fmt.Sprintf("%s (%s) %s", mediaUrl, m.Filename+m.Ext, mediaSize))
	for _, e := range m.ExtraFiles {
		extraUrl, err := srcFilename(e.Tim, e.Ext)
		if err != nil {
			log.Print(err)
			break
		}
		extraSize := humanize.Bytes(e.Fsize)
		media = append(media, fmt.Sprintf("%s (%s) %s", extraUrl, e.Filename+e.Ext, extraSize))
	}
	return media
}

// print_comments prints the rest of the messages
// TODO: sanitize author for trailing spaces at least
// NOTE: Resto field is useless, it ALWAYS has the mainID
func print_comments(data Thread) {
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

// remove_citations removes ALL citations to the original post AND link to parent post. Might be too much for normal post, but we limit the usage to single parent posts. Otherwise parsers gets confused and print a long line. Like this one.
func remove_citations(m string) (string, error) {
	doc, err := goquery.NewDocumentFromReader(strings.NewReader(m))
	if err != nil {
		return "", err
	}
	doc.Find("span.quote").Each(func(_ int, s *goquery.Selection) {
		s.Remove()
	})
	doc.Find("a").First().Remove()
	return doc.Html()
}

// print_comment prints the message provided as well as any children
func print_comment(msg Message, rsp Thread, depth int) {
	parentId, err := msg.Parent()
	if err != nil {
		log.Print(err)
	}
	if parentId > 0 {
		msg, err := remove_citations(msg.Comment)
		if err != nil {
			log.Print(err)
		}
		fmt.Printf("%s", html2console(
			msg,
			depth))
	} else {
		fmt.Printf("%s", html2console(msg.Comment, depth))
	}
	// Media
	for _, media := range msg.getMedia() {
		fmt.Printf(strings.Repeat(" ", Max(depth*3+1, 0))+"%s\n", media)
	}
	// Footer
	if showAuthors {
		fmt.Printf(strings.Repeat(" ", Max(depth*3, 0))+">> %s - %d - %s\n\n",
			msg.Author,
			msg.No,
			humanTime(msg.Time))
	} else {
		fmt.Printf(strings.Repeat(" ", Max(depth*3, 0))+">> %d - %s\n\n",
			msg.No,
			humanTime(msg.Time))
	}
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
	md, _ := html2text.FromString(raw, html2text.Options{PrettyTables: true})
	// TODO: use console width
	s, _ := text.WrapLeftPadded(md, 80, depth*3+1)
	return s + fmt.Sprintln()
}

// printOp Prints the main thread post
func printOp(resp Thread) {
	// Header
	fmt.Printf("\ntitle: %s\nurl: %s\n", resp.Posts[0].Title, uri)
	for _, media := range resp.Posts[0].getMedia() {
		fmt.Printf("media: %s\n", media)
	}
	fmt.Println()
	// Message
	fmt.Printf("%s\n", html2console(resp.Posts[0].Comment, 1))
	// Footer
	author := resp.Posts[0].Author
	date := humanTime(resp.Posts[0].Time)
	id := resp.Posts[0].No
	if showAuthors == true {
		fmt.Printf("%s - %d -  %s\n\n\n", author, id, date)
	} else {
		fmt.Printf("%d - %s\n\n\n", id, date)
	}
}
