package main

import "github.com/kataras/iris/v12"

type (
	request struct {
		Firstname string `json:"firstname"`
		Lastname  string `json:"lastname"`
	}

	response struct {
		ID      uint64 `json:"id"`
		Message string `json:"message"`
	}
)

func main() {
	app := iris.New()
	app.Handle("GET", "/", func(context iris.Context) {
		context.Text("hello world")
	})
	app.Handle("PUT", "/users/{id:uint64}", updateUser)
	app.Listen(":8080")
}

func updateUser(context iris.Context) {
	id, _ := context.Params().GetUint64("id")

	var req request
	if err := context.ReadJSON(&req); err != nil {
		context.StopWithError(iris.StatusBadRequest, err)
		return
	}

	resp := response{
		ID:      id,
		Message: req.Firstname + " updated successfully",
	}
	context.JSON(resp)
}
