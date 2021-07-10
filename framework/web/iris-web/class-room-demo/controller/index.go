package controller

import (
	"class-room-demo/database"
	"class-room-demo/service"
)

var db = database.DB
var messageService = service.NewMessage()
var userService = service.NewUser()
