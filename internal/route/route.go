package route

import (
	"hotrss/internal/middle"
	"hotrss/internal/route/api"

	"github.com/gin-contrib/pprof"
	"github.com/gin-gonic/gin"
)

// InitRoute init route
func InitRoute(logfile string) *gin.Engine {
	r := gin.New()

	r.Use(middle.HTPPLog(logfile))
	r.Use(middle.RecoverLog(logfile))
	r.ForwardedByClientIP = true
	pprof.Register(r)

	r.GET("/feed/:name_feedtype", api.GetFeedInfo)
	r.GET("/feeds/:feedtype", api.GetFeeds)
	r.GET("/opml", api.DownOpml)
	// r.GET("/status")
	r.GET("/", api.GetFeedIndex)
	return r
}
