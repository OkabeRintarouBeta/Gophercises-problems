package hn

import (
	"encoding/json"
	"fmt"
	"net/http"
)

const (
	apiBase = "https://hacker-news.firebaseio.com/v0"
)

type Client struct {
	apiBase string
}

type Item struct {
	By          string `json:"by"`
	Descendants int    `json:"descendants"`
	ID          int    `json:"id"`
	Kids        []int  `json:"kids"`
	Score       int    `json:"score"`
	Time        int    `json:"time"`
	Title       string `json:"title"`
	Type        string `json:"type"`

	Text string `json:"text"`
	URL  string `json:"url"`
}

func (c *Client) defaultify() {
	c.apiBase = apiBase
}

func (c *Client) TopItems() ([]int, error) {
	c.defaultify()
	queryUrl := fmt.Sprintf("%s/topstories.json", c.apiBase)

	resp, err := http.Get(queryUrl)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var target []int
	err = json.NewDecoder(resp.Body).Decode(&target)
	if err != nil {
		return nil, err
	}
	return target, nil
}

func (c *Client) GetStory(id int) (Item, error) {
	c.defaultify()
	queryUrl := fmt.Sprintf("%s/item/%d.json", c.apiBase, id)
	resp, err := http.Get(queryUrl)
	var item Item
	if err != nil {
		return Item{}, err
	}
	defer resp.Body.Close()

	err = json.NewDecoder(resp.Body).Decode(&item)
	if err != nil {
		return Item{}, err
	}
	// fmt.Println(item.Title)
	return item, nil
}
