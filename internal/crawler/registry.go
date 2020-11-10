package crawler

import (
	"context"
	"encoding/xml"
	"hotrss/internal/crawler/site/hupu"
	"hotrss/internal/crawler/site/kr36"
	"hotrss/internal/crawler/site/xueqiu"
	"hotrss/internal/crawler/site/zdm_hot1"
	"time"
)

var crawle = NewCrawler()

// RegistryCrawlers registry all crawler to Crawle
func RegistryCrawlers(ctx context.Context) {
	baseurl := ctx.Value("baseurl").(string)

	crawle.BaseURL = baseurl
	// hupu
	crawle.Registry("hupu", hupu.NewBXJ(), time.Hour*2)
	// 雪球热榜
	crawle.Registry("xueqiu", xueqiu.NewXueqiu(), time.Hour*2)
	// 36kr
	crawle.Registry("36kr", kr36.NewKr36(), time.Hour*2)
	// 什么值得买
	crawle.Registry("zdm_hot1", zdm_hot1.NewZdmHot1(), time.Hour*2)
	crawle.Start(ctx)
}

// GetAllFeeds get all feeds
func GetAllFeeds(feedtype string) []*CrawleInfo {
	return crawle.Feeds(feedtype)
}

// FeedIndex feed index
type FeedIndex struct {
	JSONRss string `json:"json_rss"`
	XMLRss  string `json:"xml_rss"`
	OmplURL string `json:"opml"`
}

// GetFeedIndex get feed index
func GetFeedIndex() FeedIndex {
	feedindex := FeedIndex{
		JSONRss: crawle.BaseURL + "/feeds/json",
		XMLRss:  crawle.BaseURL + "/feeds/xml",
		OmplURL: crawle.BaseURL + "/opml",
	}
	return feedindex
}

type RssOpml struct {
	XMLName xml.Name `xml:"opml"`
	Version string   `xml:"version,attr"`
	Outline Outline  `xml:"body>outline"`
}

type Outline struct {
	XMLName     xml.Name  `xml:"outline"`
	Title       string    `xml:"title,attr,omitempty"`
	Text        string    `xml:"text,attr,omitempty"`
	Description string    `xml:"description,attr,omitempty"`
	Type        string    `xml:"type,attr,omitempty"`
	Version     string    `xml:"version,attr,omitempty"`
	HTMLURL     string    `xml:"htmlUrl,attr,omitempty"`
	XMLURL      string    `xml:"xmlUrl,attr,omitempty"`
	Outline     []Outline `xml:"outline"`
}

// GetFeedOpml generate rss opml file from rss
func GetFeedOpml() *RssOpml {
	return crawle.GetFeedOpml()
}
