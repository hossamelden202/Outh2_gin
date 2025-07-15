package controllers

import (
	"net/http"
	"os"
	"outh2/config"
	"outh2/utils"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt"
)
func Logout(c *gin.Context){
	Refresh,err:=c.Cookie("refresh-token")
		if err!=nil{
		utils.SendError(c,http.StatusUnauthorized,"Unauthorized access provide refresh-token cookie")
		return
	}
	token,err:=jwt.Parse(Refresh,func (token*jwt.Token)(interface{},error)  {
		if _,ok:=token.Method.(*jwt.SigningMethodHMAC);!ok{
			return nil,jwt.ErrSignatureInvalid
		}
		return os.Getenv("jwt_secret"),nil
	})
	if err!=nil||!token.Valid{
		utils.SendError(c,http.StatusUnauthorized,"Unauthorized access")
		return
	}
	var email interface{}
	var devid interface{}
	
	if claims,ok:=token.Claims.(jwt.MapClaims);ok{
		var ok bool
		if email,ok=claims["email"];ok{
			if devid,ok=claims["devid"];ok{
				refreshtoken,err:=config.Rd.Get(config.Ctx,"refresh-token:"+email.(string)+":"+strconv.Itoa(int(devid.(float64)))).Result()
				if err!=nil{
					utils.SendError(c,http.StatusInternalServerError,"something went wrong")
					return
				}
				if refreshtoken!=Refresh{
					utils.SendError(c,http.StatusUnauthorized,"refresh-token in cookies isnot correct")
					return
				}
				if err:=config.Rd.Del(config.Ctx,"refresh-token:"+email.(string)+":"+strconv.Itoa(int(devid.(float64)))).Err();err!=nil{
					utils.SendError(c,http.StatusInternalServerError,"something went wrong")
					return
				}


			}
		}
	}
	c.SetCookie("refresh-token","",-1,"/","",true,true)
	tokenStr:=c.GetHeader("Authorization")

	if err:=config.Rd.Set(config.Ctx,"Block:"+email.(string)+":"+strconv.Itoa(int(devid.(float64))),tokenStr,15*time.Minute).Err();err!=nil{
		utils.SendError(c,http.StatusInternalServerError,"something went wrong")
		return
	}
res:=struct{
	Mess string
}{
	Mess:"Logout successfully",
}
utils.SendRes(c,res)
}
func LogoutALl(c *gin.Context){
	
}