package crawler

import (
	"context"
	"fmt"
	"hotrss/internal/storage"
	"hotrss/internal/util"
	"sync"
	"time"

	"github.com/gorilla/feeds"
	"github.com/rs/zerolog/log"
)

// Crawler get page msg and gnenerate feed
type Crawler interface {
	GenRssFeed(ctx context.Context) (*feeds.Feed, error)
}

// CrawleInfo  crawle info
type CrawleInfo struct {
	Name        string        `json:"-"`
	Cwer        Crawler       `json:"-"`
	Interval    time.Duration `json:"-"`
	IntervalStr string        `json:"interval"`
	URL         string        `json:"url"`
	LastUpdate  util.JSONTime `json:"last_update"`
	Title       string        `json:"title"`
	Status      string        `json:"status"`
	Count       int           `json:"count"`
}

// Run start run crawle url
func (ci *CrawleInfo) Run(ctx context.Context) {
	var interval time.Duration
	log.Info().Msgf("start run crawle task %s ", ci.Name)
	for {
		select {
		case <-ctx.Done():
			log.Info().Msgf("start stop crawle %s", ci.Name)
			return
		case <-time.After(interval):
			go ci.run(ctx)
			if interval == 0 {
				interval = ci.Interval // 首次运行后就执行抓取
			}
		}
	}
}

func (ci *CrawleInfo) run(ctx context.Context) {
	ci.Status = "crawling"
	log.Info().Msgf("start crawle %s data", ci.Name)
	ci.LastUpdate = util.JSONTime(time.Now())
	feed, err := ci.Cwer.GenRssFeed(ctx)
	if err != nil {
		log.Error().Msgf("gen rss feed failed %v", err)
		ci.Status = "fail"
		return
	}
	ci.Count = len(feed.Items)
	ci.Status = "finish"
	ci.Title = feed.Title

	xmlfeed, err := feed.ToRss()
	if err != nil {
		log.Error().Msgf("feed to xml rss failed %w", err)
	}
	if xmlfeed != "" {
		key := fmt.Sprintf("%s.%s", ci.Name, "xml")
		storage.DefaultStorage.SaveFeedData(key, util.StringToByte(xmlfeed))
	}

	jsonfeed, err := feed.ToJSON()
	if err != nil {
		log.Error().Msgf("feed to json rss failed %w", err)
	}
	if jsonfeed != "" {
		key := fmt.Sprintf("%s.%s", ci.Name, "json")
		storage.DefaultStorage.SaveFeedData(key, util.StringToByte(jsonfeed))
	}
	log.Info().Msgf("%s crawler feed total %d", ci.Name, len(feed.Items))
}

// CrawleManager manager crawler
type CrawleManager struct {
	sync.Mutex
	BaseURL     string
	crawleInfos []*CrawleInfo
}

// Registry registry crawler
func (cm *CrawleManager) Registry(name string,
	cwawler Crawler,
	interval time.Duration) {

	crwaleinfo := CrawleInfo{
		Name:        name,
		Cwer:        cwawler,
		Interval:    interval,
		IntervalStr: fmt.Sprintf("%s", interval),
	}

	cm.crawleInfos = append(cm.crawleInfos, &crwaleinfo)
	log.Info().Msgf("registry crawler %s interval: %s ", name, interval)
}

// Start start run crawlemanager
func (cm *CrawleManager) Start(ctx context.Context) {
	for _, crawleinfo := range cm.crawleInfos {
		crawleinfocopy := crawleinfo
		go crawleinfocopy.Run(ctx)
	}
}

// Feeds get all feeds
func (cm *CrawleManager) Feeds(feedtype string) []*CrawleInfo {
	for i := range cm.crawleInfos {
		cm.crawleInfos[i].URL = fmt.Sprintf("%s/feed/%s.%s", cm.BaseURL, cm.crawleInfos[i].Name, feedtype)
	}
	return cm.crawleInfos
}

// GetFeedOpml generate opml file from rss
func (cm *CrawleManager) GetFeedOpml() *RssOpml {
	feedopml := &RssOpml{Version: "1.1",
		Outline: Outline{
			Title:   "实时热榜",
			Text:    "实时热榜",
			Outline: []Outline{},
		},
	}
	for _, crawlerinfo := range cm.crawleInfos {
		feedopml.Outline.Outline = append(feedopml.Outline.Outline, Outline{
			Title:       crawlerinfo.Title,
			Text:        crawlerinfo.Title,
			Description: crawlerinfo.Title,
			Type:        "rss",
			Version:     "RSS",
			XMLURL:      crawlerinfo.URL,
		})
	}
	return feedopml
}

// NewCrawler init crawler
func NewCrawler() *CrawleManager {
	return &CrawleManager{
		crawleInfos: make([]*CrawleInfo, 0),
	}
}
