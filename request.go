package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"
)

// srcFileneme returnes the path where the image can be found
func srcFilename(filename string, extension string) (string, error) {
	u, err := url.Parse(uri)
	if err != nil {
		return "", err
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
	req, err := http.NewRequest(http.MethodGet, uri, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("User-Agent", userAgent)

	client := &http.Client{Timeout: timeout}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	if resp.StatusCode != http.StatusOK {
		fmt.Println(uri)
		fmt.Println(string(body))
		return nil, fmt.Errorf("invalid http status code %d", resp.StatusCode)
	}

	return body, nil
}

func doRequest() (*Thread, error) {
	bytes, err := getUrl()
	if err != nil {
		return nil, err
	}
	var thread Thread
	if err := json.Unmarshal(bytes, &thread); err != nil {
		return nil, err
	}
	return &thread, err
}
