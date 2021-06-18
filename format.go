package main

import (
	"bytes"
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

// printComments prints the rest of the messages
// TODO: sanitize author for trailing spaces at least
// NOTE: Resto field is useless, it ALWAYS has the mainID
func printComments(data Thread) string {
	var bytes bytes.Buffer
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
			printComment(&bytes, post, data, 0)
		}
	}
	return bytes.String()
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

// printComment prints the message provided as well as any children
func printComment(bytes *bytes.Buffer, msg Message, rsp Thread, depth int) {
	parentId, err := msg.Parent()
	if err != nil {
		log.Print(err)
	}
	if parentId > 0 {
		msg, err := remove_citations(msg.Comment)
		if err != nil {
			log.Print(err)
		}
		fmt.Fprintf(bytes, "%s", html2console(
			msg,
			depth))
	} else {
		fmt.Fprintf(bytes, "%s", html2console(msg.Comment, depth))
	}
	// Media
	for _, media := range msg.getMedia() {
		fmt.Fprintf(bytes, strings.Repeat(" ", Max(depth*3+1, 0))+"%s\n", media)
	}
	// Footer
	if showAuthors {
		fmt.Fprintf(bytes, strings.Repeat(" ", Max(depth*3, 0))+">> %s - %d - %s\n\n",
			msg.Author,
			msg.No,
			humanTime(msg.Time))
	} else {
		fmt.Fprintf(bytes, strings.Repeat(" ", Max(depth*3, 0))+">> %d - %s\n\n",
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
			printComment(bytes, othermsg, rsp, depth+1)
		}
	}
}

func html2console(raw string, depth int) string {
	md, _ := html2text.FromString(raw, html2text.Options{PrettyTables: true})
	s, _ := text.WrapLeftPadded(md, int(width), depth*3+1)
	return s + fmt.Sprintln()
}

func (t Thread) String() string {
	return printOp(t.Posts[0]) + printComments(t)
}

// printOp Prints the main thread post
func printOp(msg Message) string {
	var bytes bytes.Buffer
	// Header
	fmt.Fprintf(&bytes, "\ntitle: %s\nurl: %s\n", msg.Title, uri)
	for _, media := range msg.getMedia() {
		fmt.Fprintf(&bytes, "media: %s\n", media)
	}
	// Message
	fmt.Fprintf(&bytes, "\n%s\n", html2console(msg.Comment, 1))
	// Footer
	author := msg.Author
	date := humanTime(msg.Time)
	id := msg.No
	if showAuthors {
		fmt.Fprintf(&bytes, "%s - %d -  %s\n\n\n", author, id, date)
	} else {
		fmt.Fprintf(&bytes, "%d - %s\n\n\n", id, date)
	}
	return bytes.String()
}
