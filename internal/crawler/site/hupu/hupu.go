package hupu

import (
	"context"
	"fmt"
	"net/http"
	"net/url"
	"rss_hot/internal/util"
	"strings"
	"time"

	"github.com/PuerkitoBio/goquery"
	"github.com/antlabs/pcurl"
	"github.com/gorilla/feeds"
	"github.com/rs/zerolog/log"
)

// copy from chrome
const cURLData = `curl 'https://bbs.hupu.com/all-gambia' \
-H 'authority: bbs.hupu.com' \
-H 'cache-control: max-age=0' \
-H 'upgrade-insecure-requests: 1' \
-H 'user-agent: Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_6) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/85.0.4183.121 Safari/537.36' \
-H 'accept: text/html,application/xhtml+xml,application/xml;q=0.9,image/avif,image/webp,image/apng,*/*;q=0.8,application/signed-exchange;v=b3;q=0.9' \
-H 'sec-fetch-site: same-site' \
-H 'sec-fetch-mode: navigate' \
-H 'sec-fetch-user: ?1' \
-H 'sec-fetch-dest: document' \
-H 'referer: https://www.hupu.com/' \
-H 'accept-language: zh-CN,zh;q=0.9,en;q=0.8' \
-H 'cookie: _dacevid3=fea0749f.2c8f.e76c.f688.60dcb06d67f0; sensorsdata2015jssdkcross=%7B%22distinct_id%22%3A%2216e72a1b1095a5-0615fee4b372d4-1d3e6a5a-1296000-16e72a1b10a3ef%22%2C%22%24device_id%22%3A%2216e72a1b1095a5-0615fee4b372d4-1d3e6a5a-1296000-16e72a1b10a3ef%22%2C%22props%22%3A%7B%22%24latest_traffic_source_type%22%3A%22%E5%BC%95%E8%8D%90%E6%B5%81%E9%87%8F%22%2C%22%24latest_referrer%22%3A%22https%3A%2F%2Fmo.fish%2F%22%2C%22%24latest_search_keyword%22%3A%22%E6%9C%AA%E5%8F%96%E5%88%B0%E5%80%BC%22%7D%7D; Hm_lvt_39fc58a7ab8a311f2f6ca4dc1222a96e=1596450051,1596450053,1596450134,1596450280; PHPSESSID=b5c0b8ea3a6b1aa26eb90e592a7157c2; _cnzz_CV30020080=buzi_cookie%7Cfea0749f.2c8f.e76c.f688.60dcb06d67f0%7C-1; Hm_lvt_4fac77ceccb0cd4ad5ef1be46d740615=1602138091; Hm_lvt_b241fb65ecc2ccf4e7e3b9601c7a50de=1602138091; _fmdata=QnBISsVgHYRsOvkwwh%2BRdzPkSJyLnYHpyk66AyihuOIQm8zzAUp4sgIBgnZrI0V9Q6rFGrOPBHlxs1D1SmcbAyN2H9kLpL8GixM3bmJIPE8%3D; acw_tc=781bad2916021775896718696e5f90cae6d74ed5abca9a65353fed90d05554; __dacevst=ef98d509.150beb6c|1602179408862; Hm_lpvt_b241fb65ecc2ccf4e7e3b9601c7a50de=1602177610; Hm_lpvt_4fac77ceccb0cd4ad5ef1be46d740615=1602177610; Hm_lvt_c324100ace03a4c61826ef5494c44048=1602138092,1602177591,1602177610; Hm_lpvt_c324100ace03a4c61826ef5494c44048=1602177610' \
--compressed`

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
	req, err := pcurl.ParseAndRequest(c.cURLData)
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
	c.req.Header.Set("Referer", c.req.URL.Host)

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
		href, exist := s.Find("span.textSpan > a").Attr("href")
		if !exist {
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
	item := &feeds.Item{
		Title:   title,
		Link:    &feeds.Link{Href: pageURL},
		Author:  &feeds.Author{Name: name},
		Content: ret + browse, // for json
		Id:      pageURL,
	}
	return item, nil
}
