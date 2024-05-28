package scraper

import (
	"encoding/xml"
	"fmt"
	"net/http"
	"strings"
	"time"
)

type Feed struct {
	Title         string
	Link          string
	Description   string
	Generator     string
	Language      string
	LastBuildDate time.Time
	Items         []Item
}

type Item struct {
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

func FetchFeed(url string) (Feed, error) {
	resp, err := http.Get(url)
	if err != nil {
		errMsg := fmt.Sprintf("failed to fetch feed '%s': %s", url, err)
		fmt.Println(errMsg)
		return Feed{}, err
	}

	defer resp.Body.Close()

	xmlPayload := make([]byte, 0)
	_, err = resp.Body.Read(xmlPayload)

	if err != nil {
		errMsg := fmt.Sprintf("failed to read feed from body for: '%s'", url)
		fmt.Println(errMsg)
		return Feed{}, err
	}

	feed, err := parseRssXML(string(xmlPayload))
	if err != nil {
		errMsg := fmt.Sprintf("failed to parse feed rss xml for: '%s'", url)
		fmt.Println(errMsg)
		return Feed{}, err
	}

	return feed, nil
}

func parseRssXML(xmlPayload string) (Feed, error) {
	rss := rssXML{}

	xmlReader := strings.NewReader(xmlPayload)

	d := xml.NewDecoder(xmlReader)
	d.Strict = false

	err := d.Decode(&rss)
	if err != nil {
		errMsg := fmt.Sprintf("failed to xmlPayload unmarshal json: %s", err)
		fmt.Println(errMsg)
		return Feed{}, err
	}

	lastBuildDate, err := parseRssTime(rss.Channel.LastBuildDate)
	if err != nil {
		errMsg := fmt.Sprintf("failed to parse time: %s", rss.Channel.LastBuildDate)
		fmt.Println(errMsg)
		return Feed{}, err
	}

	items := make([]Item, 0)
	for _, i := range rss.Channel.Item {
		t, err := parseRssTime(i.PubDate)

		if err != nil {
			errMsg := fmt.Sprintf("failed to parse time: %s", i.PubDate)
			fmt.Println(errMsg)
		}

		items = append(items, Item{
			Title:       i.Title,
			Link:        i.Link,
			PubDate:     t,
			Guid:        i.Guid,
			Description: i.Description,
		})
	}

	feed := Feed{
		Title:         rss.Channel.Title,
		Link:          rss.Channel.Link.Href,
		Description:   rss.Channel.Description,
		Generator:     rss.Channel.Generator,
		Language:      rss.Channel.Language,
		LastBuildDate: lastBuildDate,
		Items:         items,
	}

	return feed, nil
}

func parseRssTime(rssTime string) (time.Time, error) {
	t, err := time.Parse(time.RFC1123Z, rssTime)
	if err != nil {
		errMsg := fmt.Sprintf("failed to parse time: %s", rssTime)
		fmt.Println(errMsg)
		return time.Time{}, err
	}

	return t, nil
}
