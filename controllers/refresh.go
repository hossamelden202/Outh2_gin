package controllers

import (
	"net/http"
	"os"
	"outh2/config"
	"outh2/model"
	"outh2/utils"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt"
)
func Refresh(c *gin.Context){

	RefreshToken,errh:=c.Cookie("refresh-token")
		if errh!=nil{
		utils.SendError(c,http.StatusInternalServerError,"something went wrong")
		return
	}
	
	token,err:=jwt.Parse(RefreshToken,func (token *jwt.Token) (interface{},error){
		if _,ok:=token.Method.(*jwt.SigningMethodHMAC);!ok{
		return nil,jwt.ErrSignatureInvalid
		}
		return os.Getenv("jwt_secret"),nil
	})
	if err!=nil||!token.Valid{
		utils.SendError(c,http.StatusUnauthorized,"Refreshtoken is ivalid")
		return
	}
if claims,ok:=token.Claims.(jwt.MapClaims);ok{
	if email,ok:=claims["email"];ok{
		if devid:=claims["devid"];ok{
		RereshtokenSent,errr:=config.Rd.Get(config.Ctx,"refresh-token:"+email.(string)+":"+strconv.Itoa(int(devid.(float64)))).Result()
		if errr!=nil{
		utils.SendError(c,http.StatusInternalServerError,"something went wrong")
		return
		}
		if RereshtokenSent!=RefreshToken{
		utils.SendError(c,http.StatusUnauthorized,"unAuthorized accesss sent RefreshToken")
		return
		}
		var user model.Users
		if err:=config.DB.Model(&model.Users{}).Where("email=?",email.(string)).First(&user).Error;err!=nil{
		utils.SendError(c,http.StatusInternalServerError,"something went wrong")
		return
		}
		devid3,_:=utils.SetDeviceInfo(c,user.Email)
		token,errs:=utils.GenerteJwt(c,user.Username,user.Email,int(user.ID),user.Role,15*time.Minute,user.Token_version,devid3)
		if errs!=nil{
		utils.SendError(c,http.StatusInternalServerError,"something went wrong")
		return
		}
		refresh,errs3:=utils.GenerteJwt(c,user.Username,user.Email,int(user.ID),user.Role,30*time.Hour*24,user.Token_version,devid3)
		if errs3!=nil{
		utils.SendError(c,http.StatusInternalServerError,"something went wrong")
		return
		}
		err4:=config.Rd.Del(config.Ctx,"refresh-token:"+email.(string)+":"+strconv.Itoa(int(devid.(float64)))).Err()
			if err4!=nil{
		utils.SendError(c,http.StatusInternalServerError,"something went wrong")
		return
		}
		err5:=config.Rd.Set(config.Ctx,"refresh-token:"+email.(string)+":"+strconv.Itoa(devid3),refresh,30*time.Hour*24).Err()
		if err5!=nil{
		utils.SendError(c,http.StatusInternalServerError,"something went wrong")
		return
		}

		res:=struct{
			Token string
			RefreshToken string
		}{
			Token: token,
			RefreshToken:refresh ,
		}
		utils.SendRes(c,res)
	}
	}
}


}