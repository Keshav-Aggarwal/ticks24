package auth

import (
	"github.com/appleboy/gin-jwt"
	"github.com/gin-gonic/gin"
	"github.com/goinggo/tracelog"
	"github.com/vivek-yadav/ticks24"
	"net/http"
	"time"
)

func Setup(service *ticks24.Service) *jwt.GinJWTMiddleware {
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
			client, er := service.ConnectAuthService()
			if er != nil {
				tracelog.Errorf(er, "auth", "Login", "Failed to connect to auth serive")
			}
			er = client.Call("AuthRequest.IsAuth", &user, &result)
			if er != nil {
				c.JSON(http.StatusBadRequest, gin.H{
					"code":    http.StatusBadRequest,
					"message": er.Error(),
				})
				return
			}
			//
			//result,err := user.Get()
			//if err != nil {
			//	return email,false
			//}
			//if result {
			//	h := c.Writer.Header()
			//	h.Set("email", email)
			//	return email, true
			//}
			return email, false
		},
		Authorizator: func(email string, c *gin.Context) bool {
			//user := auth.User{
			//	Email:email,
			//}
			//level := auth.GetAccessLevel(c.Request.Method)
			//result,err := user.Auth(initConfig.ServerConfig.AppName,level,c.Request.URL.Path)
			//if err != nil {
			//	return false
			//}
			//if result {
			//	h := c.Writer.Header()
			//	h.Set("email", email)
			//	return true
			//}
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
