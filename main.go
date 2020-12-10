package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
)

const JSONFILE string = "/home/sendai/37647.json"

type Rsp struct {
	Posts []Message
}

type Message struct {
	No             int    `json:"no"`
	Title          string `json:"sub"`
	Comment        string `json:"com"`
	Author         string `json:"name"`
	trip           string
	Time           int `json:"time"`
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
	filename       string
	ext            string
	tim            string
	md5            string
	Resto          int `json:"resto"`
}

func main() {
	bytes, err := ioutil.ReadFile(JSONFILE)
	if err != nil {
		panic(err)
	}

	var data Rsp
	if err := json.Unmarshal(bytes, &data); err != nil {
		panic(err)
	}
	fmt.Printf("%s\n", data.Posts[0].Comment)
}
