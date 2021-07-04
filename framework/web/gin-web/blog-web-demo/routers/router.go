package routers

import (
	"github.com/gin-gonic/gin"
	"github.com/xuzimian/blog-web-demo/pkg/setting"
	"github.com/xuzimian/blog-web-demo/routers/api"
)

func InitRouter() *gin.Engine {
	engine := gin.New()

	engine.Use(gin.Logger())

	engine.Use(gin.Recovery())

	gin.SetMode(setting.RunMode)

	apis := engine.Group("/api")
	{
		//获取标签列表
		apis.GET("/tags", api.GetTags)
		//新建标签
		apis.POST("/tags", api.AddTag)
		//更新指定标签
		apis.PUT("/tags/:id", api.EditTag)
		//删除指定标签
		apis.DELETE("/tags/:id", api.DeleteTag)
	}

	{
		//获取文章列表
		apis.GET("/articles", api.GetArticles)
		//获取指定文章
		apis.GET("/articles/:id", api.GetArticle)
		//新建文章
		apis.POST("/articles", api.AddArticle)
		//更新指定文章
		apis.PUT("/articles/:id", api.EditArticle)
		//删除指定文章
		apis.DELETE("/articles/:id", api.DeleteArticle)
	}

	return engine
}
