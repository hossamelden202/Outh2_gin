package controllers

import (
	"fmt"
	"net/url"
	"os"

	"github.com/gin-gonic/gin"
)
func Authorize( c *gin.Context){
	
provider:=c.Param("provider")

if provider=="google"{

url := fmt.Sprintf(
  "https://accounts.google.com/o/oauth2/v2/auth?client_id=%s&redirect_uri=%s&response_type=code&scope=%s&state=%s",
  os.Getenv("CLIENT_ID"),
  os.Getenv("RED_URL"),
  "openid email profile",
  os.Getenv("GOOGLE_STATE"),
)
fmt.Println("now we are calling the google to provide us with code")

c.Redirect(302,url)

}
if provider=="github"{
	url:=fmt.Sprintf("https://github.com/login/oauth/authorize?client_id=%s&redirect_uri=%s&scope=%s&state=%s&prompt=select_account&response_type=code",os.Getenv("GITHUB_CLIENT_ID"),os.Getenv("GITHUB_RED_URL"),"user user:email",os.Getenv("GITHUB_STATE"))
fmt.Println("now we are calling the github to provide us with code")

	c.Redirect(302,url)

}
if provider=="linkedin"{
	fmt.Println("LINKDIN_CLIENT_ID:", os.Getenv("LINKDIN_CLIENT_ID"))
fmt.Println("LINKDIN_RED_URL:", os.Getenv("LINKDIN_RED_URL"))
fmt.Println("LINKDIN_STATE:", os.Getenv("LINKDIN_STATE"))

url := fmt.Sprintf("https://www.linkedin.com/oauth/v2/authorization?client_id=%s&redirect_uri=%s&state=%s&scope=%s&response_type=code",
    os.Getenv("LINKDIN_CLIENT_ID"),
    url.QueryEscape(os.Getenv("LINKDIN_RED_URL")),
    os.Getenv("LINKDIN_STATE"),
    url.QueryEscape("r_liteprofile r_emailaddress"),
)

fmt.Println("now u hit linkdin url")
c.Redirect(302,url)

}
}