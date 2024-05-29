package scraper

import (
	"encoding/xml"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"
)

type FeedData struct {
	Title         string
	Link          string
	Description   string
	Generator     string
	Language      string
	LastBuildDate time.Time
	Posts         []PostData
}

type PostData struct {
	Title       string
	Link        string
	PubDate     time.Time
	Guid        string
	Description string
}

type rssXML struct {
	XMLName xml.Name `xml:"rss"`
	Text    string   `xml:",chardata"`
	Version string   `xml:"version,attr"`
	Channel struct {
		Text  string `xml:",chardata"`
		Title string `xml:"title"`
		Link  struct {
			Text string `xml:",chardata"`
			Href string `xml:"href,attr"`
			Rel  string `xml:"rel,attr"`
			Type string `xml:"type,attr"`
		} `xml:"link"`
		Description   string `xml:"description"`
		Generator     string `xml:"generator"`
		Language      string `xml:"language"`
		LastBuildDate string `xml:"lastBuildDate"`
		Item          []struct {
			Text        string `xml:",chardata"`
			Title       string `xml:"title"`
			Link        string `xml:"link"`
			PubDate     string `xml:"pubDate"`
			Guid        string `xml:"guid"`
			Description string `xml:"description"`
		} `xml:"item"`
	} `xml:"channel"`
}

func FetchFeed(url string) (FeedData, error) {
	resp, err := http.Get(url)
	if err != nil {
		return FeedData{}, err
	}

	defer resp.Body.Close()

	xmlPayload, err := io.ReadAll(resp.Body)
	if err != nil {
		return FeedData{}, err
	}

	feed, err := parseRssXML(string(xmlPayload))
	if err != nil {
		return FeedData{}, err
	}

	return feed, nil
}

func parseRssXML(xmlPayload string) (FeedData, error) {
	rss := rssXML{}

	xmlReader := strings.NewReader(xmlPayload)

	d := xml.NewDecoder(xmlReader)
	d.Strict = false

	err := d.Decode(&rss)
	if err != nil {
		return FeedData{}, err
	}

	lastBuildDate, err := parseRssTime(rss.Channel.LastBuildDate)
	if err != nil {
		return FeedData{}, err
	}

	items := make([]PostData, 0)
	for _, i := range rss.Channel.Item {
		t, err := parseRssTime(i.PubDate)

		if err != nil {
			errMsg := fmt.Sprintf("failed to parse time: %s", i.PubDate)
			fmt.Println(errMsg)
		}

		items = append(items, PostData{
			Title:       i.Title,
			Link:        i.Link,
			PubDate:     t,
			Guid:        i.Guid,
			Description: i.Description,
		})
	}

	feed := FeedData{
		Title:         rss.Channel.Title,
		Link:          rss.Channel.Link.Href,
		Description:   rss.Channel.Description,
		Generator:     rss.Channel.Generator,
		Language:      rss.Channel.Language,
		LastBuildDate: lastBuildDate,
		Posts:         items,
	}

	return feed, nil
}

func parseRssTime(rssTime string) (time.Time, error) {
	t, err := time.Parse(time.RFC1123Z, rssTime)
	if err != nil {
		return time.Time{}, err
	}

	return t, nil
}
