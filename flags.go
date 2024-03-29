package main

import (
	"errors"
	"flag"
	"fmt"
	"net/url"
	"os"
	"strings"

	"github.com/nathan-fiscaletti/consolesize-go"
)

func initializeWidth() {
	if width == 0 {
		cols, _ := consolesize.GetConsoleSize()
		width = uint(cols)
	}
}

func validateFlags() error {
	initializeWidth()
	if err := validateUri(); err != nil {
		return err
	}
	return nil
}

// TODO: remove url arguments?
func validateUri() error {
	_, err := url.Parse(uri)
	if err != nil {
		return err
	}

	if strings.HasSuffix(uri, ".html") {
		uri = strings.TrimSuffix(uri, ".html") + ".json"
	}

	if !strings.HasSuffix(uri, ".json") {
		return errors.New("invalid url??")

	}
	return nil
}

func myUsage() {
	fmt.Printf("Usage: %s [OPTIONS] URL ...\n", os.Args[0])
	flag.PrintDefaults()
}

func processFlags() error {
	flag.Parse()
	flag.Usage = myUsage
	if flag.NArg() != 1 {
		flag.Usage()
		return errors.New("Error: Missing URL argument.")
	}
	uri = flag.Args()[0]
	return nil
}
