package controllers
import (
	"encoding/json"
	"fmt"
	"net/http"
	"outh2/config"
	"outh2/model"
	"outh2/utils"
	"strings"

	"github.com/gin-gonic/gin"
)
func GetSession(c *gin.Context){
	var cursor uint64
	var sessions []model.Session
for{
	scannedkeys,cursor,err:=config.Rd.Scan(config.Ctx,cursor,"session:"+c.GetString("id")+":*",100).Result()
	if err!=nil{
		utils.SendError(c,http.StatusInternalServerError,"something went wrong")
		return
	}
	for _,key :=range scannedkeys{

		parts:=strings.Split(key,":")
		if len(parts)==3{
			
			 res,err:=config.Rd.Get(config.Ctx,key).Result()
			 if err!=nil{
				utils.SendError(c,http.StatusInternalServerError,"something went wrong")
				return

			}
			var session model.Session

			errr:=json.Unmarshal([]byte(res),&session)
			fmt.Println("one::",session)
			
		if errr!=nil{
			utils.SendError(c,http.StatusInternalServerError,"something went wrong")
			return
			}

			sessions=append(sessions,session)

		}
	}
	if cursor==0{
		break
	}

}
fmt.Println("final:::",sessions)
	utils.SendRes(c,sessions)
}
