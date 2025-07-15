package controllers

import (
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
func CallbackGitub(c *gin.Context){
		code:=c.Query("code")
		params:=fmt.Sprintf("client_id=%s&client_secret=%s&code=%s&redirect_uri=%s",os.Getenv("GITHUB_CLIENT_ID"),os.Getenv("GITHUB_CLIENT_SECRET"),code,os.Getenv("GITHUB_RED_URL"))
		pram:=strings.NewReader(params)
		req,errReq:=http.NewRequest("POST","https://github.com/login/oauth/access_token",pram)
		if errReq!=nil{
			utils.SendError(c,http.StatusInternalServerError,"Github fetching token of api failded")
			return
		}
		req.Header.Set("Content-type","application/x-www-form-urlencoded")
		req.Header.Set("Accept","application/json")
		res,err:=http.DefaultClient.Do(req)
		if err!=nil{
			utils.SendError(c,http.StatusInternalServerError,err.Error())
			return
		}
		byte,errp:=io.ReadAll(res.Body)
			if errp!=nil{
			utils.SendError(c,http.StatusInternalServerError,errp.Error())
			return
		}
	
		var Res map[string]interface{}
		errj:=json.Unmarshal(byte,&Res)
		if errj!=nil{
			utils.SendError(c,http.StatusInternalServerError,errj.Error())
			return
		}
	



		var ty io.Reader
		req2,_:=http.NewRequest("GET","https://api.github.com/user/emails",ty)
		req2.Header.Set("Authorization","Bearer"+" "+Res["access_token"].(string))
		req2.Header.Set("Accept","application/json")
		res,errh:=http.DefaultClient.Do(req2)
		if errh!=nil{
			utils.SendError(c,http.StatusInternalServerError,errh.Error())
			return
		}
		byte3,errp3:=io.ReadAll(res.Body)
			if errp3!=nil{
			utils.SendError(c,http.StatusInternalServerError,errp3.Error())
			return
		}
		var emails []map[string]interface{}
	  errj3:=json.Unmarshal(byte3,&emails)
		if errj3!=nil{
			utils.SendError(c,http.StatusInternalServerError,errj3.Error())
			return
		}
		var email string
		for _,emailList:=range emails{
			if emailList["primary"].(bool)&&emailList["verified"].(bool){
				email=emailList["email"].(string)
			}
		}
			if email == "" {
				utils.SendError(c, http.StatusBadRequest, "No verified email found")
			return
		}
		  
	 reqName,_:=http.NewRequest("GET","https://api.github.com/user",ty)
	 reqName.Header.Set("Authorization","Bearer"+" "+Res["access_token"].(string))
	 reqName.Header.Set("Accept","application/json")
	 resuser,_:=http.DefaultClient.Do(reqName)
	 if resuser.StatusCode != 200 {
		utils.SendError(c, http.StatusInternalServerError, "GitHub user info failed")
		return
	}

	  var resU map[string]interface{}
		bytes,_:=io.ReadAll(resuser.Body)
		errJson:=json.Unmarshal(bytes,&resU)
			if errJson!=nil{
			utils.SendError(c,http.StatusInternalServerError,errJson.Error())
			return
		}
		name:=resU["name"]
		// if name,ok=claims["given_name"];!ok{
		// 	fmt.Println("given_name not found")
		// 	if name,ok=claims["name"];!ok{ 
        //      
		// 	}
			

		// }
		 fmt.Println("email not found",email)
	if name==nil{
			name = resU["login"]
		fmt.Println("name not found ")
	
	}
	
		username:=utils.GenterUserName(name.(string))
		var user model.Users
		errq:=config.DB.Model(&model.Users{}).Where("email=?",email).First(&user).Error
		if errq!=nil{
			if errors.Is(errq,gorm.ErrRecordNotFound){
				user=model.Users{
					Username :username,
					Name:name.(string),
					Email:email,
					Role:"user",
					Token_version: 0,
					Provider: "github",
				
				}
				if errw:=config.DB.Create(&user).Error;errw!=nil{
					utils.SendError(c,http.StatusInternalServerError,"something went wrong")
					return
				}
				utils.SendEmail(c,email)
			

			}else{
				
            utils.SendError(c,http.StatusInternalServerError,"something went wrong")
			return
			}
			
		}
		
		devid,devinfo:=utils.SetDeviceInfo(c,email)
		var version int
		if user.Token_version==0{
version=1
		}else{
			version=user.Token_version
		}
token,error:=utils.GenerteJwt(c,username,email,int(user.ID),"user",time.Minute*15,version,devid)
if error !=nil{
	utils.SendError(c,http.StatusInternalServerError,"something went wrong")
	return
}

done:=utils.SetSession(c)
if !done{
	utils.SendError(c,http.StatusInternalServerError,"something went wrong seting session")
	return
}
Rereshtoken,errorr:=utils.GenerteJwt(c,username,email,int(user.ID),"user",time.Hour*30*24,version,devid)
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
"provider":"github",
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
	
	

	