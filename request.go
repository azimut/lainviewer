package main

import (
	"errors"
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
	client := &http.Client{
		Timeout: timeout}

	req, err := http.NewRequest(http.MethodGet, uri, nil)
	if err != nil {
		return nil, err
	}
	// TODO: set UA through a flag
	req.Header.Set("User-Agent", userAgent)

	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		r, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			return nil, err
		}
		fmt.Println(uri)
		fmt.Println(string(r))
		return nil, fmt.Errorf("invalid http status code %d", resp.StatusCode)
	}

	if b, err := ioutil.ReadAll(resp.Body); err == nil {
		return b, nil
	}
	return nil, errors.New("no body read")
}
