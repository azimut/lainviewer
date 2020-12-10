package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"
	"time"

	markdown "github.com/MichaelMure/go-term-markdown"
	"github.com/dustin/go-humanize"
	"github.com/jaytaylor/html2text"
)

func human_time(unix int64) string {
	unix_time := time.Unix(unix, 0)
	return humanize.Time(unix_time)
}

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

func init() {
	flag.StringVar(&uri, "u", "", "url")
	flag.IntVar(&timeout, "t", 5, "timeout after seconds")
}

// TODO: add domain check
// TODO: add html -> json convert
func validate_uri() error {
	if uri == "" {
		return fmt.Errorf("-u parameter not provided")
	}

	u, err := url.Parse(uri)
	if err != nil {
		return fmt.Errorf("unparsable url")
	}

	if u.Host != "lainchan.org" {
		return fmt.Errorf("invalid domain")
	}
	if strings.HasSuffix(uri, ".html") {
		uri = strings.TrimSuffix(uri, ".html") + ".json"
	}
	if !strings.HasSuffix(uri, ".json") {
		return fmt.Errorf("invalid url??")
	}
	return nil
}

// getUrl get JSON from url
func getUrl(url string) ([]byte, error) {
	client := &http.Client{
		Timeout: time.Duration(timeout) * time.Second}

	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("User-Agent", "LainViewer/0.1")

	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		r, _ := ioutil.ReadAll(resp.Body)
		fmt.Println(uri)
		fmt.Println(string(r))
		return nil, fmt.Errorf("invalid http status code %d", resp.StatusCode)
	}

	if b, err := ioutil.ReadAll(resp.Body); err == nil {
		return b, nil
	}
	return nil, fmt.Errorf("no body read")
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
