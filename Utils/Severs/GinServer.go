package Severs

import (
	BatchNodeInstall2 "Funnull/Interface/BatchNodeInstall"
	"github.com/gin-gonic/gin"
)

func ServerRoute(server *gin.Engine) {
	server.Any("/BatchNodeInstall", BatchNodeInstall2.BatchNodeInstall)
}

func ServerRun() {
	Server := gin.Default()
	ServerRoute(Server)
	Server.Run()
}
