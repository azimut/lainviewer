package main

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"
	"time"
)

// TODO: remove url arguments?
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
	// TODO: set UA through a flag
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
