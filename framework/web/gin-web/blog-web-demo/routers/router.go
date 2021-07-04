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

	return engine
}
