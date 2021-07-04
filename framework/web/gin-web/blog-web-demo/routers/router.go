package routers

import (
	"github.com/gin-gonic/gin"
	"github.com/swaggo/gin-swagger"
	"github.com/swaggo/gin-swagger/swaggerFiles"
	_ "github.com/xuzimian/blog-web-demo/docs"
	"github.com/xuzimian/blog-web-demo/middleware/jwt"
	"github.com/xuzimian/blog-web-demo/pkg/setting"
	"github.com/xuzimian/blog-web-demo/routers/api"
)

func InitRouter() *gin.Engine {
	engine := gin.New()

	engine.Use(gin.Logger())

	engine.Use(gin.Recovery())

	gin.SetMode(setting.ServerSetting.RunMode)

	engine.GET("/auth", api.GetAuth)
	engine.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	routerGroup := engine.Group("/api")
	routerGroup.Use(jwt.JWT())
	configRouters(routerGroup)

	return engine
}

func configRouters(routerGroup *gin.RouterGroup) {
	{
		//获取标签列表
		routerGroup.GET("/tags", api.GetTags)
		//新建标签
		routerGroup.POST("/tags", api.AddTag)
		//更新指定标签
		routerGroup.PUT("/tags/:id", api.EditTag)
		//删除指定标签
		routerGroup.DELETE("/tags/:id", api.DeleteTag)
	}

	{
		//获取文章列表
		routerGroup.GET("/articles", api.GetArticles)
		//获取指定文章
		routerGroup.GET("/articles/:id", api.GetArticle)
		//新建文章
		routerGroup.POST("/articles", api.AddArticle)
		//更新指定文章
		routerGroup.PUT("/articles/:id", api.EditArticle)
		//删除指定文章
		routerGroup.DELETE("/articles/:id", api.DeleteArticle)
	}
}
