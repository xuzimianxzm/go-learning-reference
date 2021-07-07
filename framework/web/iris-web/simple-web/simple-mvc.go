package main

import (
	"github.com/kataras/iris/v12"
	"github.com/kataras/iris/v12/mvc"
)

func main() {
	app := iris.New()
	mvc.Configure(app.Party("/users"), configureMVC)
	mvc.Configure(app.Party("/"), func(application *mvc.Application) {
		application.Handle(new(indexController))
	})
	app.Listen(":8080")
}

func configureMVC(app *mvc.Application) {
	app.Handle(new(userController))
}

type indexController struct {
}

type userController struct {
}

func (i *indexController) Get() string {
	return "hello world"
}

func (c *userController) PutBy(id uint64, req request) response {
	return response{
		ID:      id,
		Message: req.Firstname + " updated successfully",
	}
}
