package scraper

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

var xmlPayload = `
<rss version="2.0">

	<channel>
	  <title>Boot.dev Blog</title>
	  <link>https://blog.boot.dev/</link>
	  <description>Recent content on Boot.dev Blog</description>
	  <generator>Hugo -- gohugo.io</generator>
	  <language>en-us</language>
	  <lastBuildDate>Wed, 01 May 2024 00:00:00 +0000</lastBuildDate>
	  <atom:link href="https://blog.boot.dev/index.xml" rel="self" type="application/rss+xml"/>
	  <item>
	    <title>The Boot.dev Beat. May 2024</title>
	    <link>https://blog.boot.dev/news/bootdev-beat-2024-05/</link>
	    <pubDate>Wed, 01 May 2024 00:00:00 +0000</pubDate>
	    <guid>https://blog.boot.dev/news/bootdev-beat-2024-05/</guid>
	    <description>
	      A new Pub/Sub Architecture course, lootable chests, and ThePrimeagen&rsquo;s Git course is only a couple weeks away.
	    </description>
	  </item>
	  <item>
	    <title>Trustworthy vs Trustless Apps</title>
	    <link>
	      https://blog.boot.dev/security/trustworthy-vs-trustless-apps/
	    </link>
	    <pubDate>Tue, 23 Jul 2019 00:00:00 +0000</pubDate>
	    <guid>
	      https://blog.boot.dev/security/trustworthy-vs-trustless-apps/
	    </guid>
	    <description>
	      In the wake of the hearings about Facebook&rsquo;s new Libra blockchain, it is more important than ever that we all understand the difference between trustworthy and trustless apps.
	    </description>
	  </item>
	</channel>
</rss>
`

func TestParseFeedXml(t *testing.T) {
	feed, err := parseRssXML(xmlPayload)
	assert.NoError(t, err)

	assert.Equal(t, "Boot.dev Blog", feed.Title)
	assert.Equal(t, "Recent content on Boot.dev Blog", feed.Description)
	assert.Equal(t, "2024-05-01 00:00:00 +0000 +0000", feed.LastBuildDate.String())
	assert.Equal(t, 2, len(feed.Entries))

	i1 := feed.Entries[0]
	assert.Equal(t, "The Boot.dev Beat. May 2024", i1.Title)
	assert.Equal(t, "2024-05-01 00:00:00 +0000 +0000", i1.PubDate.String())

	i2 := feed.Entries[1]
	assert.Equal(t, "Trustworthy vs Trustless Apps", i2.Title)
	assert.Equal(t, "2019-07-23 00:00:00 +0000 +0000", i2.PubDate.String())
}

func TestFetchFeed(t *testing.T) {
	url := "https://blog.boot.dev/index.xml"

	feed, err := FetchFeed(url)
	assert.NoError(t, err)

	assert.IsType(t, Feed{}, feed)
	assert.NotEmpty(t, feed.Title)
	assert.NotEmpty(t, feed.Description)
	assert.NotEmpty(t, feed.LastBuildDate)
	assert.Greater(t, len(feed.Entries), 5)
}