package rss

import (
	"encoding/json"
	"os"
	"time"

	"github.com/mmcdole/gofeed"

	"news/pkg/storage"
)

type ParserMsg struct {
	Source string
	Data   []storage.Post
	Err    error
}

type config struct {
	RSS           []string `json:"rss"`
	RequestPeriod int      `json:"request_period"`
}

func LoadConf(path string) (*config, error) {
	b, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}

	var conf config
	err = json.Unmarshal(b, &conf)
	if err != nil {
		return nil, err
	}

	return &conf, nil
}

type Parser struct {
	Urls  []string
	Delay time.Duration
}

func (p *Parser) handleFeed(feed *gofeed.Feed) ([]storage.Post, error) {
	var posts []storage.Post

	for _, item := range feed.Items {
		publishedUTC, _ := convertToUTC(item.Published)
		post := storage.Post{
			Title:     item.Title,
			Content:   item.Description,
			Published: publishedUTC,
			Link:      item.Link,
		}
		posts = append(posts, post)
	}

	return posts, nil
}

func (p *Parser) Run(msgChan chan<- ParserMsg) {
	for _, url := range p.Urls {
		go func(url string) {
			msg := ParserMsg{Source: url}
			fp := gofeed.NewParser()

			feed, err := fp.ParseURL(url)
			if err != nil {
				msg.Err = err
				msgChan <- msg
				return
			}

			posts, err := p.handleFeed(feed)
			if err != nil {
				msg.Err = err
				msgChan <- msg
				return
			}

			msg.Data = posts
			msgChan <- msg

		}(url)
	}
}

func NewParser(conf config) *Parser {
	p := Parser{
		Urls:  conf.RSS,
		Delay: time.Minute * time.Duration(conf.RequestPeriod),
	}

	return &p
}

var rssTimeFormats = []string{
	"Mon, 2 Jan 2006 15:04:05 +0000", // Custom for single digit date.
	time.RFC1123,                     // "Mon, 02 Jan 2006 15:04:05 MST"
	time.RFC1123Z,                    // "Mon, 02 Jan 2006 15:04:05 -0700"
	time.RFC822,                      // "02 Jan 06 15:04 MST"
	time.RFC822Z,                     // "02 Jan 06 15:04 -0700"
}

func convertToUTC(rssTime string) (t time.Time, err error) {
	for _, format := range rssTimeFormats {
		t, err = time.Parse(format, rssTime)
		if err == nil {
			return t.UTC(), nil
		}
	}
	return
}
