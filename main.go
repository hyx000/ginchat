package main

import (
	"ginchat/router"
	"ginchat/utils"
	_ "github.com/gin-gonic/gin"
	_ "net/http"
)

func main() {
	utils.InitConfig()
	utils.InitMySQL()
	utils.InitRedis()
	r := router.Router()
	r.Run(":8080")

}
