package route

import (
	"rss_hot/internal/route/api"

	"github.com/gin-gonic/gin"
)

// InitRoute init route
func InitRoute() *gin.Engine {
	r := gin.Default()
	r.GET("/feed/:name_feedtype", api.GetFeedInfo)
	r.GET("/feeds/:feedtype", api.GetFeeds)
	r.GET("/opml", api.DownOpml)
	r.GET("/", api.GetFeedIndex)
	return r
}
