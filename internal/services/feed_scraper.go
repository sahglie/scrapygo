package services

import (
	"context"
	"database/sql"
	"encoding/xml"
	"errors"
	"fmt"
	"github.com/google/uuid"
	"io"
	"net/http"
	"scrapygo/internal/database"
	"strings"
	"time"
)

var ErrPartialScrape = errors.New("partial scrape")
var ErrFailedScrape = errors.New("failed scrape")

func (cfg *Config) ScrapeFeeds() error {
	feeds, err := cfg.DB.GetNextFeedsToScrape(context.TODO())
	if err != nil {
		cfg.Logger.Error("failed to fetch feeds", "err", err)
		return err
	}

	if len(feeds) == 0 {
		cfg.Logger.Info("no feeds to scrape")
		return nil
	}

	msg := fmt.Sprintf("attempting to scrape %d feeds", len(feeds))
	cfg.Logger.Info(msg)

	for _, f := range feeds {
		err = cfg.ScrapeFeed(f)
		if err != nil {
			errMsg := fmt.Sprintf("failed to scrape feed %s: %s\n", f.ID, err)
			return errors.New(errMsg)
		}
	}

	return nil
}

func (cfg *Config) ScrapeFeed(feed database.Feed) error {
	cfg.Logger.Info("attempting to scrape feed", "url", feed.Url)

	feedData, err := fetchFeed(feed.Url)
	if err != nil {
		cfg.Logger.Error("failed to scrape feed", "err", err)
		return err
	}

	posts, err := cfg.DB.GetPostsByFeedID(context.TODO(), feed.ID)
	if err != nil {
		cfg.Logger.Error("failed to scrape feed", "err", err)
		return err
	}

	failedPostUrls := make([]string, 0)

	for _, p := range feedData.Posts {
		if postAlreadyScraped(p.Link, posts) {
			continue
		}

		_, err := cfg.DB.CreatePost(context.TODO(), database.CreatePostParams{
			ID:     uuid.New(),
			FeedID: feed.ID,
			Title:  p.Title,
			Description: sql.NullString{
				String: p.Description,
				Valid:  len(p.Description) > 0,
			},
			Url:         p.Link,
			PublishedAt: p.PubDate,
			CreatedAt:   time.Now(),
			UpdatedAt:   time.Now(),
		})

		if err != nil {
			cfg.Logger.Error("failed to create post", "url", p.Link, "feedId", feed.ID, "err", err)
			failedPostUrls = append(failedPostUrls, p.Link)
		} else {
			cfg.Logger.Info("created post", "url", p.Link)
		}

	}

	if len(failedPostUrls) == 0 {
		return nil
	}

	if len(failedPostUrls) == len(feedData.Posts) {
		return ErrFailedScrape
	}

	return ErrPartialScrape
}

func postAlreadyScraped(postURL string, posts []database.Post) bool {
	for _, p := range posts {
		if p.Url == postURL {
			return true
		}
	}
	return false
}

type feedData struct {
	Title         string
	Link          string
	Description   string
	Generator     string
	Language      string
	LastBuildDate time.Time
	Posts         []postData
}

type postData struct {
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

func fetchFeed(url string) (feedData, error) {
	resp, err := http.Get(url)
	if err != nil {
		return feedData{}, err
	}

	defer resp.Body.Close()

	xmlPayload, err := io.ReadAll(resp.Body)
	if err != nil {
		return feedData{}, err
	}

	feed, err := parseRssXML(string(xmlPayload))
	if err != nil {
		return feedData{}, err
	}

	return feed, nil
}

func parseRssXML(xmlPayload string) (feedData, error) {
	rss := rssXML{}

	xmlReader := strings.NewReader(xmlPayload)

	d := xml.NewDecoder(xmlReader)
	d.Strict = false

	err := d.Decode(&rss)
	if err != nil {
		return feedData{}, err
	}

	lastBuildDate, err := parseRssTime(rss.Channel.LastBuildDate)
	if err != nil {
		return feedData{}, err
	}

	items := make([]postData, 0)
	for _, i := range rss.Channel.Item {
		t, err := parseRssTime(i.PubDate)

		if err != nil {
			errMsg := fmt.Sprintf("failed to parse time: %s", i.PubDate)
			fmt.Println(errMsg)
		}

		items = append(items, postData{
			Title:       i.Title,
			Link:        i.Link,
			PubDate:     t,
			Guid:        i.Guid,
			Description: i.Description,
		})
	}

	feed := feedData{
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
