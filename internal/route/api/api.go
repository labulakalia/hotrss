package api

import (
	"encoding/xml"
	"fmt"
	"net/http"
	"rss_hot/internal/crawler"
	"rss_hot/internal/storage"

	"github.com/gin-gonic/gin"
)

// GetFeedInfo get rss feed
func GetFeedInfo(c *gin.Context) {
	nameType := c.Param("name_feedtype")
	data, err := storage.DefaultStorage.GetFeedData(nameType)
	if err != nil {
		c.Status(http.StatusBadRequest)
		c.Writer.WriteString(err.Error())
		return
	}
	c.Writer.WriteHeaderNow()
	c.Writer.Write(data)
}

// GetFeeds get all feeds
func GetFeeds(c *gin.Context) {
	feedtype := c.Param("feedtype")
	if feedtype != "xml" && feedtype != "json" {
		c.JSON(http.StatusBadRequest, "only allow /feeds/[xml|json]")
		return
	}
	feeds := crawler.GetAllFeeds(feedtype)
	c.JSON(http.StatusOK, feeds)
}

// DownOpml download rss opml
func DownOpml(c *gin.Context) {
	feedopml := crawler.GetFeedOpml()
	out, err := xml.MarshalIndent(feedopml, " ", "  ")
	if err != nil {
		c.Writer.WriteHeader(http.StatusInternalServerError)
		c.Writer.WriteString(err.Error())
		return
	}
	c.Writer.WriteHeader(http.StatusOK)
	c.Header("Content-Disposition", "attachment; filename=hot_rss.opml")
	c.Header("Content-Type", "application/text/plain")
	c.Header("Accept-Length", fmt.Sprintf("%d", len(out)))
	c.Writer.Write(out)
}

// GetFeedIndex get feed index
func GetFeedIndex(c *gin.Context) {
	feedindex := crawler.GetFeedIndex()
	c.JSON(http.StatusOK, feedindex)
}
