package controllers

import (
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"os"
	"outh2/config"
	"outh2/model"
	"outh2/utils"
	"strconv"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)
type CallbackToken struct{
	AccessToken string `json:"access_token"`
	ExpiresIn int      `json:"expires_in"`
	RefreshToken string `json:"refresh_token"`
	Scope string         `json:"scope"`
	TokenType string      `json:"token_type"`
	IdToken string       `json:"id_token"`


}
func CallbackGoogle( c *gin.Context){
	//var body io.Reader

	
	
	var claims map[string]interface{}
	
	fmt.Println("callback returned with code i hope")
code:=c.Query("code")
massege:=fmt.Sprintf("code=%s&client_id=%s&client_secret=%s&redirect_uri=%s&grant_type=%s",code,os.Getenv("CLIENT_ID"),os.Getenv("CLIENT_SECRET"),os.Getenv("RED_URL"),"authorization_code")	

//fmt.Println("here is code :",code)

reader:=strings.NewReader(massege)

    res,err:=http.Post("https://oauth2.googleapis.com/token","application/x-www-form-urlencoded",reader)
	if err!=nil{
		c.JSON(http.StatusInternalServerError,gin.H{"error":"something went wrong"})
		return
	}
	body,errBody:=io.ReadAll(res.Body)
	if errBody!=nil{
		c.JSON(http.StatusInternalServerError,gin.H{"error":"something went wrong"})
		return
	}
	//fmt.Println(body,"AND",res.Body)
	var Obj CallbackToken
	errJson:=json.Unmarshal(body,&Obj)
	if errJson!=nil{
		c.JSON(http.StatusInternalServerError,gin.H{"error":"something went wrong"})
		return
	}
	// token,error:=jwt.Parse(obj.id_token,func(token *jwt.Token)(interface{},error){
	// 	if _,ok:=token.Method.(*jwt.SigningMethodHMAC);!ok{
	// 		return nil,jwt.ErrSignatureInvalid
	// 	}
	// 	return os.Getenv("jwt_secret"),nil
	// })
	// if error!=nil||!token.Valid{
	// 	c.JSON(http.StatusUnauthorized,"you cannot pass to get user info from google try again or manul logging")
	// 	return
	// }
	parts:=strings.Split(Obj.IdToken,".")
	if len(parts)<2{
	utils.SendError(c,http.StatusUnauthorized,"parts of openID token must be more than  2")
	return
	}
// 	if len(parts)>=3{
// 	fmt.Println("token signuture",parts[2])
// }
// 	fmt.Println("token Header",parts[0])
// 	fmt.Println("token encoded",parts[1])
	
	 payload,errp:=base64.RawURLEncoding.DecodeString(parts[1])
	//  fmt.Println("token decoded",payload)
 	if errp!=nil{
		utils.SendError(c,http.StatusInternalServerError,"something went wrong")
		return
	 }
	 erre:=json.Unmarshal(payload,&claims)
	//  fmt.Println("Claims",claims)
	 if erre!=nil{
		utils.SendError(c,http.StatusInternalServerError,"something went wrong")
		return
	 }   	
	    var name,email interface{}
	 	var ok bool
		if name,ok=claims["given_name"];!ok{
			fmt.Println("given_name not found")
			if name,ok=claims["name"];!ok{ 
             fmt.Println("name not found")
			}
			

		}
		if email,ok=claims["email"];!ok{
				fmt.Println("email not found")
		}
		username:=utils.GenterUserName(name.(string))
		var user model.Users
		errq:=config.DB.Model(&model.Users{}).Where("email=?",email.(string)).First(&user).Error
		if errq!=nil{
			if errors.Is(errq,gorm.ErrRecordNotFound){
				user=model.Users{
					Username :username,
					Name:name.(string),
					Email:email.(string),
					Role:"user",
					Token_version: 0,
					Provider: "google",
				
				}
				if errw:=config.DB.Create(&user).Error;errw!=nil{
					utils.SendError(c,http.StatusInternalServerError,"something went wrong")
					return
				}
			

			}else{
            utils.SendError(c,http.StatusInternalServerError,"something went wrong")
			return
			}
			
		}
		
		devid,devinfo:=utils.SetDeviceInfo(c,email.(string))
				var version int
		if user.Token_version==0{
version=1
		}else{
			version=user.Token_version
		}
token,error:=utils.GenerteJwt(c,username,email.(string),int(user.ID),"user",time.Minute*15,version,devid)
if error !=nil{
	utils.SendError(c,http.StatusInternalServerError,"something went wrong")
	return
}
Rereshtoken,errorr:=utils.GenerteJwt(c,username,email.(string),int(user.ID),"user",time.Hour*30*24,version,devid)
if errorr !=nil{
	utils.SendError(c,http.StatusInternalServerError,"something went wrong")
	return
}

if er:=config.Rd.Set(config.Ctx,"refresh-token:"+user.Email+":"+strconv.Itoa(devid),Rereshtoken,24*time.Hour*30).Err();er!=nil{
			utils.SendError(c,http.StatusInternalServerError,"something went wrong seting refresh-token")
			return
}
c.SetCookie("refresh-token",Rereshtoken,int(24*time.Hour*30),"/","localhost",true,true)
if err:=config.DB.Model(&model.Users{}).Where("email=?",email).Updates(map[string]interface{}{
"token_version":user.Token_version+1,
"provider":"google",
}).Error;err!=nil{
			utils.SendError(c,http.StatusInternalServerError,"something went wrong updating token vrsion")
			return
		}
		errq2:=config.DB.Model(&model.Users{}).Where("email=?",email).First(&user).Error
		if errq2!=nil{
		utils.SendError(c,http.StatusInternalServerError,"something went wrong")
			return
		}
		
rese:=struct{
	Token string
	User model.Users
	Devinfo model.DeviceRecord
}{
	User:user,
	Token:token,
	Devinfo: devinfo,
}
utils.SendRes(c,rese)
	
}
		
	



