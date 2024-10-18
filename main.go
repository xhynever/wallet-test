package main

import (
	"github.com/gin-gonic/gin"
	_ "github.com/lib/pq"
	"github.com/sirupsen/logrus"
	"github.com/xhynever/wallet-test/handler"
	repository "github.com/xhynever/wallet-test/repository/sqlc"
	"github.com/xhynever/wallet-test/services"
	"github.com/xhynever/wallet-test/util"
)

func main() {
	// 读配置
	config := util.InitConfig()
	// 初始化db
	conn, err := util.InitDB()
	if err != nil {
		logrus.Fatal("cannot connect to db", err)
	}
	// 设置db连接池
	conn.SetMaxOpenConns(config.PgPool)
	conn.SetMaxIdleConns(config.PgPool)
	//initialize store, services and server
	store := repository.NewStore(conn)
	services := services.NewService(store)
	handlers := handler.NewHandler(services)
	r := gin.Default()
	// api接口，主要重心
	handlers.InitRouter(r)
	logrus.Fatal(r.Run(config.ServersAddress))
}
