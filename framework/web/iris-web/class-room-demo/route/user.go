package route

import (
	"class-room-demo/controller"
	"class-room-demo/middleware"
	"github.com/kataras/iris/v12/core/router"
)

func routeUser(party router.Party) {
	party.Post("/login", controller.PostLogin)

	party.Post("/user", controller.PostUser)
	party.Get("/user", controller.GetUser)
	party.Put("/user", middleware.JWT.Serve, middleware.Logined, controller.PutUser)
}
