package utils

import (
	"encoding/json"
	"fmt"
	"io"
	"net"
	"net/http"
	"net/smtp"
	"os"
	"outh2/config"
	"outh2/model"
	"regexp"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt"
	"github.com/google/uuid"
)
type GeoData struct {
	Query        string `json:"query"`
	Country      string `json:"country"`
	RegionName   string `json:"regionName"`
	City         string `json:"city"`
	ISP          string `json:"isp"`
	Timezone     string `json:"timezone"`
	Org          string `json:"org"`
	Status       string `json:"status"`
	Lat          float64 `json:"lat"`
	Lon          float64 `json:"lon"`
	Zip          string     `json:"zip"`
	CountryCode  string      `json:"countryCode"`
}
func GenerteJwt(c *gin.Context,Username string,Email string ,id int,role string,Time time.Duration,version int,devid int)(string,error){
	
jti:=uuid.New().String()	
exp:=time.Now().Add(Time).Unix()
token:=jwt.NewWithClaims(jwt.SigningMethodHS256,jwt.MapClaims{
	"Username":Username,
	"email":Email,
	"role":role,
	"id":id,
	"exp":exp,
	"iat":time.Now().Unix(),
	"version":version,
	"jti":jti,
	"devid":devid,


})
c.Set("Username",Username)
c.Set("email",Email)
c.Set("role",role)
c.Set("id",id)
c.Set("exp",exp)
c.Set("jti",jti)
c.Set("version",version)





return token.SignedString([]byte(os.Getenv("jwt_secret")))

}
func Sendlocation(ip string)(*GeoData,error) {
	// fmt.Println("first",ip)
	
ip,_,err2:=net.SplitHostPort(ip)
if err2!=nil{return nil,err2}
// fmt.Println("sec",ip)
resp,err:=http.Get("http://ip-api.com/json/"+"8.8.8.8")



if err!=nil{return nil,err}
body,error:=io.ReadAll(resp.Body)
defer resp.Body.Close()
if error!=nil{return nil,error}
var data *GeoData
fmt.Println(resp)
fmt.Println("body:",string(body))
//can use parseAndsend written funcion  
error=json.Unmarshal(body,&data)
if error!=nil{return nil,error}
if data.Status!="success"{return nil,fmt.Errorf("failed to get geo info")}
return data,nil
}
func GenterUserName(name string)string{
sub1:=uuid.New().String()

sub:=sub1[:8]
strings.ToLower(name)
strings.ReplaceAll(name," ","")
re:=regexp.MustCompile(`[^a-zA-Z]+`)
clean := re.ReplaceAllString(name, "")
username:=clean+"_"+ string(sub)
return username
}
func SetDeviceInfo(c *gin.Context,email string)(int,model.DeviceRecord){
	
var user model.Users
if err := config.DB.Where("email = ?", email).First(&user).Error; err != nil {
	SendError(c, http.StatusInternalServerError, "failed to load user")
	return 0,model.DeviceRecord{}
}

	ip:=c.Request.RemoteAddr

browser:=c.Request.UserAgent()
geo,errw:=Sendlocation(ip)

	if errw!=nil{
		SendError(c,http.StatusInternalServerError,"something went wrong")
		fmt.Println("ashd",errw)
		return 0,model.DeviceRecord{}
	}
	dev:=model.DeviceRecord{
    	UserID:user.ID,
		City:geo.City,
		Region:geo.RegionName,
		Country:geo.Country,
		Lat:geo.Lat,
		Lon:geo.Lon,
		ZipCode:geo.Zip,
		Locale:"en"+"-"+geo.CountryCode,
		Browser:browser,
		LastLogin:time.Now(),
	}
erre := config.DB.Model(&model.DeviceRecord{}).
	Where("userID=?", user.ID).
	Assign(dev).
	FirstOrCreate(&dev).Error


if erre!=nil{
	SendError(c,http.StatusInternalServerError,"something went wrong")
	fmt.Println("sja",erre)
	return 0,model.DeviceRecord{}

}

fmt.Println("this is now the device id been set:",int(dev.ID))

c.Set("devid",int(dev.ID))
return int(dev.ID),dev
}
func SendError(c *gin.Context,status int,err string){
if c.Writer.Written(){
	return
}
	switch c.ContentType(){
		case "application/xml":
			if err!=""{
			c.XML(status,gin.H{"error":err})
			return
			}
		case "application/x-yaml":
			if err!=""{
			c.YAML(status,gin.H{"error":err})
			return
			}
		default:
			if err!=""{
				c.JSON(status,gin.H{"error":err})
			}
	}
}
func SendRes(c *gin.Context,res interface{}){
	if c.Writer.Written(){
		return
	}
	switch c.ContentType(){
case "application/xml":

	c.XML(http.StatusOK,res)


case "application/x-yaml":
c.YAML(http.StatusOK,res)
default:
	c.JSON(http.StatusOK,res)
}

}
func SetSession(c *gin.Context)bool{

session:=model.Session{
	
	Jti:c.GetString("jti"),
	UserID: c.GetInt("id"),
	IsActive: true,
	IssuedAT: time.Now(),
	DeviceInfoId: c.GetInt("devid"),
	ExpireAt:time.Now().Add(time.Minute*15),
}
fmt.Println("iam in setSessions")
fmt.Println("id:", c.GetInt("id"))
fmt.Println("jti:", c.GetString("jti"))
fmt.Println("devid:", c.GetInt("devid"))


bytes,err:=json.Marshal(session)
if err!=nil{
	return false
}
if err:=config.Rd.Set(config.Ctx,"session:"+ c.GetString("id")+":"+c.GetString("jti"),bytes,time.Minute*15).Err();err!=nil{
fmt.Println(err.Error())
return false
}
return true


}
func SendEmail(c *gin.Context,email string){
	Addr:=c.Request.RemoteAddr
geo,err:=Sendlocation(Addr)
if err!=nil{SendError(c,http.StatusInternalServerError,"something went wrong")
fmt.Println(err.Error())
return}
fmt.Println("hello gere is addr",Addr)
 var user model.Users
// email:="kc334844@gmail.com"
config.DB.Where("email=?",email).First(&user)
Message := fmt.Sprintf(
`Subject: 🎉 Welcome to Racist Team, %s!

Hello %s 👋,

Welcome aboard! Your account with the email: %s has just been created successfully.

🗺️ Location at Signup:
- IP Address: %s
- Country:%s
- Region: %s
- City: %s
- ISP: %s
- Organization: %s
- Timezone: %s

🕒 Signup Time: %v

We’re thrilled to have you in the Racist Team family.  
Feel free to explore, connect, and enjoy everything we’ve built for you 💥

If you didn’t create this account, please contact us immediately.

Cheers,  
The Racist Team 🛡️

`,
	user.Username,
	user.Username,         //  — username
	user.Email,            // %s — email
	geo.Query,             // %s — IP
	geo.Country,           // %s — country
	geo.RegionName,        // %s — region
	geo.City,              // %s — city
	geo.ISP,               // %s — ISP
	geo.Org,               // %s — organization
	geo.Timezone,          // %s — timezone
	time.Now(),            // %v — timestamp
)

SendEmailSmtp(c,email,Message)
}
func SendEmailSmtp(c *gin.Context,email string,Message string ){
	
// 	ms:=mailersend.NewMailersend(os.Getenv("Ms_API_KEY"))
// 	msg:=ms.Email.NewMessage()
// 	From:=mailersend.From{Name:"team",Email: os.Getenv("Mail_email")}
// 	to:=mailersend.Recipient{Name:userName,Email: email}
// 	msg.SetFrom(From)
// 	msg.SetRecipients([]mailersend.Recipient{to})
// 	msg.SetSubject("Alert: Login issue")
// 	// msg.SetText(Message)
// 	msg.SetHTML("<pre>"+html.EscapeString(Message)+"</pre>")
// fmt.Println("EMAIL MESSAGE:\n", Message)

// 	_,_,err:=ms.BulkEmail.Send(config.Ctx,[]*mailersend.Message{msg})

auth:=smtp.PlainAuth("",os.Getenv("Mail_email"),os.Getenv("Mail_password"),"smtp.gmail.com")
    Addr:="smtp.gmail.com"+":"+"587"
	err:=smtp.SendMail(Addr,auth,os.Getenv("Mail_email"),[]string{email},[]byte(Message))
	

	if err!=nil{
		
		SendError(c,http.StatusInternalServerError,fmt.Sprintf("something went wrong:%s",err))
		return
	}
}
