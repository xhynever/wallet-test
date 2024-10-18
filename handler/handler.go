package handler

import (
	"github.com/gin-gonic/gin"
	"github.com/xhynever/wallet-test/services"
)

type Handler struct {
	services *services.Service
}

func NewHandler(services *services.Service) *Handler {
	return &Handler{services: services}
}

func (h *Handler) InitRouter(app *gin.Engine) *gin.Engine {

	accounts := app.Group("/accounts")
	accounts.POST("/creat", h.CreateAccount)

	// 查看账户信息
	accounts.GET("/:id", h.GetAccount)

	// 查询owner拥有多少个账户,或者分页查询账户
	accounts.GET("/owners", h.ListAccounts)

	accounts.DELETE("/:id", h.DeleteAccount)

	funds := app.Group("/funds")

	//存取款
	funds.POST("/business", h.Business)
	// 取款
	// 存款和取款逻辑调用相同，可以只用一个接口，也可以拆分。这样传参只用传账户和金额。
	// business.POST("/draw", h.WithDraw)

	funds.POST("/transfers", h.CreateTransfer)

	// 获得交易记录
	funds.GET("/tx", h.ListEntries)

	return app
}
