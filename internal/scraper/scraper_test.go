package scraper

import (
	"fmt"
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
	feed, _ := parseRssXML(xmlPayload)
	fmt.Printf("%+v\n", feed)
}

func TestFetchFeed(t *testing.T) {
	url := "https://blog.boot.dev/index.xml"
	feed, _ := FetchFeed(url)

	fmt.Printf("%v\n", feed)
}
