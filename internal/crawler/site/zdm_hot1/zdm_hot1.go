package zdm_hot1

import (
	"context"
	"fmt"
	"hotrss/internal/util"
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/PuerkitoBio/goquery"
	"github.com/antlabs/pcurl"
	"github.com/gorilla/feeds"
	"github.com/rs/zerolog/log"
)

// copy from chrome
const cURLData = `curl 'https://post.smzdm.com/hot_1/' \
-X 'GET' \
-H 'Cookie: _ga=GA1.2.1283885670.1603952429; _gid=GA1.2.1062520729.1604983437; CNZZDATA1256793290=1927633017-1604637381-https%253A%252F%252Fwww.smzdm.com%252F%7C1605015390; zdm_qd=%7B%22referrer%22%3A%22https%3A%2F%2Fmo.fish%2F%3Fclass_id%3D%25E5%2585%25A8%25E9%2583%25A8%26hot_id%3D117%22%7D; sensorsdata2015jssdkcross=%7B%22distinct_id%22%3A%22175730215d53a6-08f2903f24c906-3e63694b-2073600-175730215d6e0e%22%2C%22first_id%22%3A%22%22%2C%22props%22%3A%7B%22%24latest_traffic_source_type%22%3A%22%E7%9B%B4%E6%8E%A5%E6%B5%81%E9%87%8F%22%2C%22%24latest_search_keyword%22%3A%22%E6%9C%AA%E5%8F%96%E5%88%B0%E5%80%BC_%E7%9B%B4%E6%8E%A5%E6%89%93%E5%BC%80%22%2C%22%24latest_referrer%22%3A%22%22%2C%22%24latest_landing_page%22%3A%22https%3A%2F%2Fpost.smzdm.com%2Fhot_1%2F%22%7D%2C%22%24device_id%22%3A%22175730215d53a6-08f2903f24c906-3e63694b-2073600-175730215d6e0e%22%7D; Hm_lpvt_9b7ac3d38f30fe89ff0b8a0546904e58=1605018470; Hm_lvt_9b7ac3d38f30fe89ff0b8a0546904e58=1604329024,1604640356,1604640393,1604983436; _zdmA.time=1605018464680.970.https%3A%2F%2Fwww.smzdm.com%2F; _zdmA.uid=ZDMA.2hKfMB-Bp.1605018351.2419200; _zdmA.vid=*; homepage_sug=a; r_sort_type=score; __jsluid_h=3461464f04326e85da32c2dab04c12d1; smidV2=20201106133845196cac922453566168d2ffe9963da44f00bc8ea24246d38f0; wt3_eid=%3B999768690672041%7C2160464098400742548%232160464112200953206; smzdm_user_source=D17885CC32EBF26A3530D565B9D991E4; UM_distinctid=17589753d4812b0-01fdf9653622bf8-5c465d7b-13c680-17589753d49d58; device_id=10326674771603952317737249ee7b6f66905c485fbfdac37c3c50e038; __ckguid=bIUQ2bgL4ynrJb2G2wV252; __jsluid_s=93243d08094757115b4b544713890d57' \
-H 'Accept: text/html,application/xhtml+xml,application/xml;q=0.9,*/*;q=0.8' \
-H 'Upgrade-Insecure-Requests: 1' \
-H 'Host: post.smzdm.com' \
-H 'User-Agent: Mozilla/5.0 (Macintosh; Intel Mac OS X 11_0) AppleWebKit/605.1.15 (KHTML, like Gecko) Version/14.0.1 Safari/605.1.15' \
-H 'Accept-Language: zh-cn' \
-H 'Accept-Encoding: gzip, deflate'`

// NewZdmHot1 new ZdmHot1 crawler
func NewZdmHot1() *ZdmHot1 {
	return &ZdmHot1{
		cURLData: cURLData,
		Client:   http.DefaultClient,
	}
}

// ZdmHot1 crawler url https://bbs.hupu.com/all-gambia
type ZdmHot1 struct {
	Client   *http.Client
	cURLData string
	req      *http.Request
	baseURL  string
}

// GenRssFeed impl interface Crawler
func (c *ZdmHot1) GenRssFeed(ctx context.Context) (*feeds.Feed, error) {
	// 解析curl
	req, err := pcurl.ParseAndRequest(c.cURLData)
	if err != nil {
		return nil, fmt.Errorf("parse request failed %w", err)
	}
	req.WithContext(ctx)
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
		Title:       "什么值得买_今日热门",
		Link:        &feeds.Link{Href: c.req.URL.String()},
		Description: "什么值得买_今日热门",
		Items:       make([]*feeds.Item, 0),
		Created:     now,
	}
	// 用来保存每一页的url
	pageurls := []string{}
	// 如果页面的指定内容
	// 打开开发者模式，选取所需要的html数据，然后右键，选择Copy->Copy selector,获取selector
	document.Find("#feed-main-list > li").Each(func(i int, s *goquery.Selection) {
		if i >= 10 {
			return
		}
		href, exist := s.Find("div > div.z-feed-content > div.feed-block-describe > a").Attr("href")
		if !exist {
			return
		}
		pageurls = append(pageurls, href)
	})

	failedURL := []string{}
	for _, url := range pageurls {
		item, err := c.getPage(url)
		if err != nil {
			log.Error().Msgf("get page failed: %v", err)
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
func (c *ZdmHot1) getPage(pageURL string) (*feeds.Item, error) {
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
	// 文章标题
	title := strings.TrimSpace(document.Find("#articleId > h1").Text()) 
	// 文章作者
	name := document.Find("#feed-side > div:nth-child(2) > div.user_tx > div > div > h2 > a").Text()

	// 文章发布时间
	// 如果文章没有时间用time.Now()
	timeStr := document.Find("#articleId > div.recommend-tab.z-clearfix.item-preferential > span > span:nth-child(1)").Text()
	var createAt time.Time
	createAt, err = time.Parse("2006-01-02 15:04:05", timeStr)
	if err != nil {
		log.Error().Msgf("parse time %s failed: %v", timeStr, err)
		createAt = time.Now()
	}

	// 文章内容

	selection := document.Find("#articleId")
	selection.Find("h1").Remove()
	selection.Find("div").Remove()
	ret, err := selection.Html()
	if err != nil {
		return nil, fmt.Errorf("parse html '#tpc > div > div.floor_box > table.case > tbody > tr > td > div.quote-content' failed %w", err)
	}
	item := &feeds.Item{
		Title:   title,
		Link:    &feeds.Link{Href: pageURL},
		Author:  &feeds.Author{Name: name},
		Content: strings.TrimSpace(ret), // for json
		Id:      pageURL,
		Created: createAt,
	}
	return item, nil
}
