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
func CallbackLinkedin(c *gin.Context){
	state:=c.Query("state")
	fmt.Println("my state :",os.Getenv("LINKDIN_STATE"),"url state:",state)
	if state!=os.Getenv("LINKDIN_STATE"){
		fmt.Println("my state :",os.Getenv("LINKDIN_STATE"),"url state:",state)
		utils.SendError(c,http.StatusUnauthorized,"you r redirect url hit unauthorized point")
		return
	}
	error:=c.Query("error")
	if error=="user_cancelled_login"{
		utils.SendError(c,http.StatusInternalServerError,"user canceled login")
		return
	}else if error=="user_cancelled_authorize"{
		utils.SendError(c,http.StatusInternalServerError,"user refused to authorize login")
		return
	}
	error_description:=c.Query("error_description")
	if error_description!=""{
		utils.SendError(c,http.StatusInternalServerError,error_description)
		return
	}
code :=c.Query("code")
body:=fmt.Sprintf("code=%s&client_id=%s&client_secret=%s&redirect_uri=%s&grant_type=authorization_code",code,os.Getenv("LINKDIN_CLIENT_ID"),os.Getenv("LINKDIN_CLIENT_SECRET"),os.Getenv("LINKDIN_RED_URL"))
BodyReader:=strings.NewReader(body)
req,err:=http.NewRequest("POST", "https://www.linkedin.com/oauth/v2/accessToken",BodyReader)
if err!=nil{
	utils.SendError(c,http.StatusInternalServerError,err.Error())
	fmt.Println(err.Error())
	return
}
req.Header.Set("Content-type","application/x-www-form-urlencoded")
res,err1:=http.DefaultClient.Do(req)
if err1!=nil{
	utils.SendError(c,http.StatusInternalServerError,err1.Error())
	fmt.Println(err1.Error())
	return
}
defer res.Body.Close()
bodyres,err2:=io.ReadAll(res.Body)
if err2!=nil{
	utils.SendError(c,http.StatusInternalServerError,err2.Error())
	fmt.Println(err2.Error())
	return
}
var resObj map[string]interface{}
json.Unmarshal(bodyres,&resObj)
//name
var test io.Reader
req2,err3:=http.NewRequest("GET","https://api.linkedin.com/v2/me",test)
if err3!=nil{
	utils.SendError(c,http.StatusInternalServerError,err3.Error())
	fmt.Println(err3.Error())
	return
}
req2.Header.Set("Authorization","Bearer"+" "+resObj["access_token"].(string))
res2,err4:=http.DefaultClient.Do(req2)
if err4!=nil{
	utils.SendError(c,http.StatusInternalServerError,err4.Error())
	fmt.Println(err4.Error())
	return
}
defer res2.Body.Close()
body3,err5:=io.ReadAll(res2.Body)
if err5!=nil{
	utils.SendError(c,http.StatusInternalServerError,err5.Error())
	fmt.Println(err5.Error())
	return
}
var US map[string]interface{}
json.Unmarshal(body3,&US)
///email
req3,err6:=http.NewRequest("GET","https://api.linkedin.com/v2/emailAddress?q=members&projection=(elements*(handle~))",test)
if err6!=nil{
	utils.SendError(c,http.StatusInternalServerError,err6.Error())
	fmt.Println(err6.Error())
	return
}
req3.Header.Set("Authorization","Bearer"+" "+resObj["access_token"].(string))
res3,err7:=http.DefaultClient.Do(req3)
if err7!=nil{
	utils.SendError(c,http.StatusInternalServerError,err7.Error())
	fmt.Println(err7.Error())
	return
}
defer res3.Body.Close()
body4,err8:=io.ReadAll(res3.Body)
if err8!=nil{
	utils.SendError(c,http.StatusInternalServerError,err8.Error())
	fmt.Println(err8.Error())
	return
}
var emailResponse struct {
    Elements []struct {
        Handle struct {
            Email string `json:"emailAddress"`
        } `json:"handle~"`
    } `json:"elements"`
}

json.Unmarshal(body4,&emailResponse)
email := emailResponse.Elements[0].Handle.Email
//
name:=US["localizedFirstName"].(string)+" "+US["localizedLastName"].(string)
username:=utils.GenterUserName(name)
		var user model.Users
		errq:=config.DB.Model(&model.Users{}).Where("email=?",email).First(&user).Error
		if errq!=nil{
			if errors.Is(errq,gorm.ErrRecordNotFound){
				user=model.Users{
					Username :username,
					Name:name,
					Email:email,
					Role:"user",
					Token_version: 0,
					Provider: "linkedin",
				
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
		
		devid,devinfo:=utils.SetDeviceInfo(c,email)
				var version int
		if user.Token_version==0{
version=1
		}else{
			version=user.Token_version
		}
token,errorr:=utils.GenerteJwt(c,username,email,int(user.ID),"user",time.Minute*15,version,devid)
if errorr !=nil{
	utils.SendError(c,http.StatusInternalServerError,"something went wrong")
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
"provider":"linkedin",
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