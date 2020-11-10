package hupu

import (
	"context"
	"fmt"
	"hotrss/internal/util"
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/PuerkitoBio/goquery"
	"github.com/gorilla/feeds"
	"github.com/rs/zerolog/log"
)

// copy from chrome
const cURLData = `curl 'https://bbs.hupu.com/all-gambia' \
-X 'GET' \
-H 'Accept: text/html,application/xhtml+xml,application/xml;q=0.9,*/*;q=0.8' \
-H 'Host: bbs.hupu.com' \
-H 'User-Agent: Mozilla/5.0 (Macintosh; Intel Mac OS X 11_0) AppleWebKit/605.1.15 (KHTML, like Gecko) Version/14.0.1 Safari/605.1.15'`

// NewBXJ new hupu crawler
func NewBXJ() *BXJ {
	return &BXJ{
		cURLData: cURLData,
		Client:   http.DefaultClient,
	}
}

// BXJ crawler url https://bbs.hupu.com/all-gambia
type BXJ struct {
	Client   *http.Client
	cURLData string
	req      *http.Request
	baseURL  string
}

// GenRssFeed impl interface Crawler
func (c *BXJ) GenRssFeed(ctx context.Context) (*feeds.Feed, error) {
	req, err := util.ParseAndRequest(c.cURLData)
	if err != nil {
		return nil, fmt.Errorf("parse request failed %w", err)
	}
	req.WithContext(ctx)
	// req
	c.req = req

	res, err := util.Request(c.req, c.Client)
	if err != nil {
		return nil, fmt.Errorf("request url %s failed %w", c.req.URL, err)
	}

	c.baseURL = fmt.Sprintf("%s://%s", c.req.URL.Scheme, c.req.URL.Host)
	c.req.Header.Set("Referer", c.req.URL.String())

	document, err := goquery.NewDocumentFromReader(res)
	if err != nil {
		return nil, fmt.Errorf("new document from resp failed %w", err)
	}

	now := time.Now()
	feed := feeds.Feed{
		Title:       document.Find("div.bbsHotPit h1").Text(),
		Link:        &feeds.Link{Href: c.req.URL.String()},
		Description: document.Find("div.bbsHotPit h1").Text(),
		Items:       make([]*feeds.Item, 0),
		Created:     now,
	}

	pageurls := []string{}

	document.Find("#container > div > div.bbsHotPit > div:nth-child(2)").First().Find("ul > li").Each(func(i int, s *goquery.Selection) {
		if i >= 10 {
			return
		}
		href, exist := s.Find("span.textSpan > a").Attr("href")
		if !exist {

			return
		}
		pageurls = append(pageurls, c.baseURL+href)
	})

	failedURL := []string{}
	fmt.Println(pageurls)
	for _, url := range pageurls {
		item, err := c.getPage(url)
		if err != nil {
			failedURL = append(failedURL, url)
			continue
		}
		feed.Items = append(feed.Items, item)
	}

	if len(failedURL) > 0 {
		log.Error().Msgf("failed get url count %d urls: %v", len(failedURL), failedURL)
	}
	return &feed, nil
}

// getPage get url page body
func (c *BXJ) getPage(pageURL string) (*feeds.Item, error) {
	newurl, err := url.Parse(pageURL)
	if err != nil {
		return nil, fmt.Errorf("parse url failed %w", err)
	}
	c.req.URL = newurl

	res, err := util.Request(c.req, c.Client)
	if err != nil {
		return nil, fmt.Errorf("request failed %w", err)
	}

	document, err := goquery.NewDocumentFromReader(res)
	if err != nil {
		return nil, fmt.Errorf("new document from resp failed %w", err)
	}

	title := document.Find("#tpc > div > div.floor_box > table.case > tbody > tr > td > div.subhead > span").Text()
	name := document.Find("#tpc > div > div.floor_box > div.author > div.left > a").Text()

	ret, err := document.Find("#tpc > div > div.floor_box > table.case > tbody > tr > td > div.quote-content").RemoveClass("small").Html()
	if err != nil {
		return nil, fmt.Errorf("parse html '#tpc > div > div.floor_box > table.case > tbody > tr > td > div.quote-content' failed %w", err)
	}
	selection := document.Find("#t_main > div.bbs_head > div.bbs-hd-h1 > span")
	browse := selection.Find("span:nth-child(1)").Text()
	browse = browse + selection.Find("span:nth-child(2)").Text()
	// NetNewsWire uses WebKit, which is Apple’s HTML rendering system. WebKit does not support WebP, so it can’t display that image.
	ret = strings.Replace(ret, "?x-oss-process=image/resize,w_800/format,webp", "", -1)
	// 图片超过三张的会被投毒，处理下
	ret = strings.Replace(ret, `https://b1.hoopchina.com.cn/web/sns/bbs/images/placeholder.png" data-original="`, "", -1)


	createAtStr := document.Find("#tpc > div > div.floor_box > div.author > div.left > span.stime").Text()
	createdAt, err := time.Parse("2006-01-02 15:04:05", createAtStr)
	item := &feeds.Item{
		Title:   title,
		Link:    &feeds.Link{Href: pageURL},
		Author:  &feeds.Author{Name: name},
		Content: ret + browse, // for json
		Id:      pageURL,
		Created: createdAt,
	}
	return item, nil
}
