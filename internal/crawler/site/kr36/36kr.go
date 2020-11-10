package kr36

import (
	"context"
	"fmt"
	"hotrss/internal/util"
	"net/http"
	"net/url"
	"time"

	"github.com/PuerkitoBio/goquery"
	"github.com/antlabs/pcurl"
	"github.com/gorilla/feeds"
	"github.com/rs/zerolog/log"
)

// copy from chrome
const cURLData = `curl 'https://36kr.com/hot-list/catalog' \
-X 'GET' \
-H 'Cookie: Hm_lpvt_713123c60a0e86982326bae1a51083e1=1605005386; Hm_lvt_713123c60a0e86982326bae1a51083e1=1605005386; UM_distinctid=175b1c6a1e5912-0ee75313990e4a-5c465d7b-1fa400-175b1c6a1e6b02; CNZZDATA1256793290=250613029-1605004590-%7C1605004590; SERVERID=6eb0a1872728d69c244094a636b7db3b|1605005386|1605005385; Hm_lpvt_1684191ccae0314c6254306a8333d090=1605005386; Hm_lvt_1684191ccae0314c6254306a8333d090=1605005386; sajssdk_2015_cross_new_user=1; sensorsdata2015jssdkcross=%7B%22distinct_id%22%3A%22175b1c69e761216-0ba3749c32ee31-5c465d7b-2073600-175b1c69e77ac0%22%2C%22%24device_id%22%3A%22175b1c69e761216-0ba3749c32ee31-5c465d7b-2073600-175b1c69e77ac0%22%2C%22props%22%3A%7B%22%24latest_traffic_source_type%22%3A%22%E7%9B%B4%E6%8E%A5%E6%B5%81%E9%87%8F%22%2C%22%24latest_referrer%22%3A%22%22%2C%22%24latest_referrer_host%22%3A%22%22%2C%22%24latest_search_keyword%22%3A%22%E6%9C%AA%E5%8F%96%E5%88%B0%E5%80%BC_%E7%9B%B4%E6%8E%A5%E6%89%93%E5%BC%80%22%7D%7D; acw_tc=2760828316050053847227222ea59567b47f2154951c7c3186cf1c9cf53ec3' \
-H 'Accept: text/html,application/xhtml+xml,application/xml;q=0.9,*/*;q=0.8' \
-H 'Host: 36kr.com' \
-H 'User-Agent: Mozilla/5.0 (Macintosh; Intel Mac OS X 11_0) AppleWebKit/605.1.15 (KHTML, like Gecko) Version/14.0.1 Safari/605.1.15' \
-H 'Accept-Language: zh-cn' \
-H 'Accept-Encoding: gzip, deflate, br' \
-H 'Connection: keep-alive'`

// NewKr36 new Kr36 crawler
func NewKr36() *Kr36 {
	return &Kr36{
		cURLData: cURLData,
		Client:   http.DefaultClient,
	}
}

// Kr36 crawler url https://bbs.hupu.com/all-gambia
type Kr36 struct {
	Client   *http.Client
	cURLData string
	req      *http.Request
	baseURL  string
}

// GenRssFeed impl interface Crawler
func (c *Kr36) GenRssFeed(ctx context.Context) (*feeds.Feed, error) {
	// 解析curl
	req, err := pcurl.ParseAndRequest(c.cURLData)
	if err != nil {
		return nil, fmt.Errorf("parse request failed %w", err)
	}
	// req
	// 保存解析后的req请求
	c.req = req
	// 保存baseURL，用来生成每一页的url
	c.baseURL = fmt.Sprintf("%s://%s", c.req.URL.Scheme, c.req.URL.Host)
	// 请求热榜url数据
	res, err := util.Request(c.req, c.Client)
	if err != nil {
		return nil, fmt.Errorf("request url %s failed %w", c.req.URL, err)
	}
	// 设置referer
	c.req.Header.Set("Referer", c.req.URL.String())

	// 解析请求到的数据
	document, err := goquery.NewDocumentFromReader(res)
	if err != nil {
		return nil, fmt.Errorf("new document from resp failed %w", err)
	}

	now := time.Now()
	// 获取请求后的数据
	feed := feeds.Feed{
		Title:       document.Find("#app > div > div.kr-layout-main.clearfloat > div.main-right > div > div > div.main-wrapper > div.list-wrapper > div:nth-child(1) > div.list-title > div").Text(),
		Link:        &feeds.Link{Href: c.req.URL.String()},
		Description: document.Find("#app > div > div.kr-layout-main.clearfloat > div.main-right > div > div > div.main-wrapper > div.list-wrapper > div:nth-child(1) > div.list-title > div").Text(),
		Items:       make([]*feeds.Item, 0),
		Created:     now,
	}
	// 用来保存每一页的url
	pageurls := []string{}
	// 如果页面的指定内容
	// 打开开发者模式，选取所需要的html数据，然后右键，选择Copy->Copy selector,获取selector
	document.Find("#app > div > div.kr-layout-main.clearfloat > div.main-right > div > div > div.main-wrapper > div.list-wrapper > div:nth-child(1) > div.article-list > div").Each(func(i int, s *goquery.Selection) {
		if i >= 10 {
			return
		}
		href, ok := s.Find("div.kr-shadow-content > div.article-item-pic-wrapper > a").Attr("href")
		if !ok {
			return
		}
		pageurls = append(pageurls, c.baseURL+href)
	})

	failedURL := []string{}
	for _, url := range pageurls {
		item, err := c.getPage(url)
		if err != nil {
			failedURL = append(failedURL, url)
			continue
		}
		item.Created = now
		feed.Items = append(feed.Items, item)
	}

	if len(failedURL) > 0 {
		log.Error().Msgf("failed get url count %d urls: %v", len(failedURL), failedURL)
	}
	return &feed, nil
}

// getPage get url page body
func (c *Kr36) getPage(pageURL string) (*feeds.Item, error) {
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
	// 标题
	title := document.Find("#app > div > div.box-kr-article-new-y > div > div.kr-layout-main.clearfloat > div.main-right > div > div > div > div.article-detail-wrapper-box > div > div.article-left-container > div.article-content > div > div > div:nth-child(1) > div > h1").Text()
	// 作者名称
	name := document.Find("#app > div > div.box-kr-article-new-y > div > div.kr-layout-main.clearfloat > div.main-right > div > div > div > div.article-detail-wrapper-box > div > div.article-left-container > div.article-content > div > div > div:nth-child(1) > div > div.article-title-icon.common-width.margin-bottom-40 > a").Text()
	// 文章内容
	ret, err := document.Find("#app > div > div.box-kr-article-new-y > div > div.kr-layout-main.clearfloat > div.main-right > div > div > div > div.article-detail-wrapper-box > div > div.article-left-container > div.article-content > div > div > div.common-width.margin-bottom-20 > div").Html()
	if err != nil {
		return nil, fmt.Errorf("parse html '#tpc > div > div.floor_box > table.case > tbody > tr > td > div.quote-content' failed %w", err)
	}
	// 文章创建时间
	// 没有文章创建时间 使用现在的时间
	// 最讨厌xx小时前了 sb
	createdAt := time.Now()
	item := &feeds.Item{
		Title:   title,
		Link:    &feeds.Link{Href: pageURL},
		Author:  &feeds.Author{Name: name},
		Content: ret,
		Id:      pageURL,
		Created: createdAt,
	}
	return item, nil
}
