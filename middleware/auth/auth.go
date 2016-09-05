package auth

import (
	"errors"
	"github.com/appleboy/gin-jwt"
	"github.com/gin-gonic/gin"
	"github.com/goinggo/tracelog"
	"github.com/vivek-yadav/UserManagementService/utils"
	"github.com/vivek-yadav/ticks24/config"
	"net/http"
	"net/rpc"
	"strconv"
	"time"
)

func ConnectAuthService(ip string, port int) (client *rpc.Client, er error) {
	typeCon := "tcp"

	// LATER to setup other type of connections
	//if this.Config.AuthService.IsHttp {
	//	typeCon = "http"
	//}
	path := ip + ":" + strconv.Itoa(port)
	//client, er = rpc.DialHTTP(typeCon, path)
	client, er = rpc.Dial(typeCon, path)
	if er != nil {
		er = errors.New("ERROR : Failed to conenct to Auth Service: (" + path + " (\n\t" + er.Error() + "\n)")
	}
	return
}

func Setup(config *config.Config) *jwt.GinJWTMiddleware {
	auth := &jwt.GinJWTMiddleware{
		Realm:      "test zone",
		Key:        []byte("secret key"),
		Timeout:    time.Hour,
		MaxRefresh: time.Hour,
		Authenticator: func(email string, password string, c *gin.Context) (string, bool) {
			user := Login{
				Email:    email,
				Password: password,
			}
			var result bool
			client, er := ConnectAuthService(config.LoginService.Ip, int(config.LoginService.Port))
			if er != nil {
				tracelog.Errorf(er, "auth", "Login", "Failed to connect to Login service")
			}
			er = client.Call("User.IsLogin", &user, &result)
			if er != nil {
				c.JSON(http.StatusBadRequest, gin.H{
					"code":    http.StatusBadRequest,
					"message": er.Error(),
				})
				return email, false
			}
			if result {
				h := c.Writer.Header()
				h.Set("email", email)
				return email, true
			}
			return email, result
		},
		Authorizator: func(email string, c *gin.Context) bool {
			req := AuthRequest{
				Email:       email,
				AppToken:    config.AppToken,
				AccessLevel: int(utils.GetAccessLevel(c.Request.Method)),
				Path:        c.Request.RequestURI,
			}
			var result bool
			client, er := ConnectAuthService(config.AuthService.Ip, int(config.AuthService.Port))
			if er != nil {
				tracelog.Errorf(er, "auth", "Authoriation", "Failed to connect to auth service")
			}
			er = client.Call("AuthRequest.IsAuth", &req, &result)
			if er != nil {
				tracelog.Errorf(er, "auth", "Authorization", "Not valid credientials.")
				return false
			}
			if result {
				h := c.Writer.Header()
				h.Set("email", email)
				return true
			}
			//h := c.Writer.Header()
			//h.Set("email", email)

			return false
		},
		Unauthorized: func(c *gin.Context, code int, message string) {
			c.JSON(code, gin.H{
				"code":    code,
				"message": message,
			})
		},
	}
	return auth

}
