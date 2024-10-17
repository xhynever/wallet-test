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
	//Accounts route
	// 可以添加一个账户登陆的中间层 accounts := app.Group("/accounts").Use(Middleware())

	accounts := app.Group("/accounts")
	accounts.POST("/creat", h.CreateAccount)

	// 查看账户信息
	accounts.GET("/:id", h.ListAccounts)

	// 查看owner拥有多少个账户
	accounts.GET("/:owner", h.GetAccount)

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
