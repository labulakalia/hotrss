![GitHub Workflow Status](https://img.shields.io/github/workflow/status/labulaka521/hotrss/Build_release)
[![Downloads](https://img.shields.io/github/downloads/labulaka521/hotrss/total.svg)](https://github.com/labulaka521/hotrss/releases)
[![license](https://img.shields.io/github/license/mashape/apistatus.svg?maxAge=2592000)](https://github.com/labulaka521/hotrss/blob/master/LICENSE)
[![Release](https://img.shields.io/github/release/labulaka521/hotrss.svg?label=Release)](https://github.com/labulaka521/hotrss/releases)
![GitHub issues](https://img.shields.io/github/issues/labulaka521/hotrss)
![GitHub pull requests](https://img.shields.io/github/issues-pr-raw/labulaka521/hotrss)

## HotRSS
一款专注于网站热榜的RSS服务

## 使用
- [主页](https://hotrss.top)
    - [RssByJson](https://hotrss.top/feeds/json)
    - [RssByXml](https://hotrss.top/feeds/json)
    - [下载opml](https://hotrss.top/opml)
    - [新热榜Rss通知(TODO)](https://hotrss.top)

### 已支持网站
- [虎扑步行街热帖](https://bbs.hupu.com/all-gambia)
- [36氪人气榜](https://36kr.com/hot-list/catalog)
- [雪球热帖](https://xueqiu.com/?category=snb_article)
- [什么值得买_今日热门](https://post.smzdm.com/hot_1/)

## 部署
```shell
./hotrss -baseurl http://yourip:8080 -port 8080
# baseurl为访问rss服务时需要的IP或域名
# port为web服务的端口
```

### 访问Rss
- 主页   
    `baseurl`
- 所有的Json格式rss    
    `baseurl/feeds/json`
- 所有的Xml格式rss  
    `baseurl/feeds/xml`
- rss热榜opml文件下载  
    `baseurl/opml`
> baseurl即为启动服务的baseurl参数


## 如何添加新的站点
### 方法1
- 首先找到所需要的热榜url链接
- 然后在`internal/crawler/site`下新建一个新的目录
- 然后新建一个结构体(struct)并实现`Crawler`接口，可以参考`internal/crawler/site/exmaple`
- 然后在`internal/crawler/registry`注册，注册时的name为访问rss时的url，然后设置抓取周期
- 然后启动服务等首次运行抓取完成即可访问

### 方法2
- 新建Issue，并提供所需要的热榜url



## 说明
- 数据目前存储在内存中，如果需要存储到专业的缓存请实现`FeedStorager`接口



## 版权声明
本服务提供的信息资料、图片、视频等均来自于公开网络，如有侵权，请与我们联系，我们会尽快处理，并删除侵权内容
