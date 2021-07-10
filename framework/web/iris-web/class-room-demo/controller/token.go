package controller

import (
	"class-room-demo/model"
	"class-room-demo/model/response"
	"github.com/kataras/iris/v12"
)

// GetTokenInfo 验证token是否有效，如果有效则返回token携带的信息
func GetTokenInfo(ctx iris.Context) {
	logined := ctx.Values().Get("logined").(model.Logined)

	res := response.GetTokenInfo{
		ID:       logined.ID,
		Username: logined.Username,
	}
	ctx.JSON(res)
}
