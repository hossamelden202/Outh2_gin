package main

import (
	"net/http"
	"outh2/config"
	"outh2/routes"
	"outh2/utils"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)
func main(){
	config.Connect()
	r:=gin.Default()
	err:=godotenv.Load(".env")
	if err!=nil{
		utils.SendError(&gin.Context{},http.StatusInternalServerError,"something went wrong")
		return
	}
	routes.Routing(r)
	r.Run(":8080")


}