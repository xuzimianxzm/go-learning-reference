package route

import (
	"class-room-demo/controller"
	"class-room-demo/middleware"
	"github.com/kataras/iris/v12/core/router"
)

func routeToken(party router.Party) {
	party.Get("/token/info", middleware.JWT.Serve, middleware.Logined, controller.GetTokenInfo)
}
