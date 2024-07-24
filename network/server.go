package network

import (
	"chat_server/types"
	"github.com/gin-gonic/gin"
	"net/http"
)

type data struct {
	userMap map[string]bool
}

func registerServer(engine *gin.Engine) *data {
	d := &data{
		userMap: make(map[string]bool),
	}

	engine.POST("/login", d.login)

	r := NewRoom()
	go r.Run()

	engine.GET("/room", r.ServeHTTP)

	return d
}

func (d *data) login(c *gin.Context) {
	var req types.LoginReq

	if err := c.ShouldBind(&req); err != nil {
		response(c, http.StatusUnprocessableEntity, nil, err.Error())
	} else {
		d.userMap[req.Name] = true
		response(c, http.StatusOK, nil, req.Name)
	}
}
