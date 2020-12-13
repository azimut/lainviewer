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

	_, err := url.Parse(uri)
	if err != nil {
		return fmt.Errorf("unparsable url")
	}

	// if u.Host != "lainchan.org" {
	// 	return fmt.Errorf("invalid domain")
	// }
	if strings.HasSuffix(uri, ".html") {
		uri = strings.TrimSuffix(uri, ".html") + ".json"
	}
	if !strings.HasSuffix(uri, ".json") {
		return fmt.Errorf("invalid url??")
	}
	return nil
}

// srcFileneme returnes the path where the image can be found
func srcFilename(filename string, extension string) (string, error) {
	u, err := url.Parse(uri)
	if err != nil {
		return "", fmt.Errorf("unparsable url")
	}
	return fmt.Sprintf("%s://%s/%s/src/%s%s",
		u.Scheme,
		u.Host,
		strings.Split(u.Path, "/")[1],
		filename,
		extension), nil
}

// getUrl get JSON from url
func getUrl() ([]byte, error) {
	client := &http.Client{
		Timeout: time.Duration(timeout) * time.Second}

	req, err := http.NewRequest(http.MethodGet, uri, nil)
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
