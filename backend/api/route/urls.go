package route

import (
	"github.com/Kotzyk/go-short/api/controller"
	"github.com/Kotzyk/go-short/api/db"
	"github.com/gin-gonic/gin"
)

func SetUrlsRouter(r *gin.Engine) {
	mw := db.RateLimitMiddleware(db.CreateClient(0))
	urls := r.Group("/urls")
	{
		urls.POST("/shorten", mw, controller.ShortenUrl)
	}
	r.GET("/:slug", mw, controller.ResolveUrl)
}
