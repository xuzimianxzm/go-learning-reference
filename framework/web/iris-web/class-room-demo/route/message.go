package route

import (
	"class-room-demo/controller"
	"class-room-demo/middleware"
	"github.com/kataras/iris/v12/core/router"
)

func routeMessage(party router.Party) {
	party.Post("/message", middleware.JWT.Serve, middleware.Logined, controller.PostMessage)
	party.Get("/message", middleware.JWT.Serve, middleware.Logined, controller.GetMessage)
}
