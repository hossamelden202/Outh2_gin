package routes

import (
	"outh2/controllers"
	"outh2/middlerware"

	"github.com/gin-gonic/gin"
) 
func Routing(r *gin.Engine){
 
	userRoutes:=r.Group("/oauth2")
	userRoutes.GET("/authorize/:provider",controllers.Authorize)
	userRoutes.GET("/callback/google",controllers.CallbackGoogle)
	userRoutes.GET("/callback/github",controllers.CallbackGitub)
	userRoutes.GET("/callback/facebook",controllers.CallbackFacebook)
	userRoutes.GET("/callback/linkedin",controllers.CallbackLinkedin)
	userRoutes.POST("/refresh",controllers.Refresh)
	userRoutes.POST("/logout",controllers.Logout)
	userRoutes.POST("/logoutall",controllers.LogoutALl)
	userRoutes.GET("/get-sessions",middlerware.Auth(),controllers.GetSession)
//http://localhost:8080/oauth2/authorize/linkedin

}