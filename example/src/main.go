package main

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/vivek-yadav/ticks24"
	"net/http"
)

func main() {
	fmt.Println("LocoLink Server Running....")
	service, er := ticks24.GetInstance()
	if er != nil {
		fmt.Printf("ERROR : In creating server instance : ( %v )\n", er.Error())
	}

	// Load Config
	configFilePath := "example/config/config.toml"
	if r, err := service.SetConfigFile(configFilePath); r == false && err != nil {
		configFilePath := "config/config.toml"
		r, err = service.SetConfigFile(configFilePath)
		if r == false && err != nil {
			fmt.Printf("Error in setting configurations from file of ums service : (%v)  : ( %v )", configFilePath, err.Error())
			return
		}
	}

	// Load Command Line Args
	if r, err := service.SetCmdArgs(); r == false && err != nil {
		fmt.Printf("ERROR : In setting configurations from file of Loco Logic service : ( %v )", err.Error())
		return
	}

	router, er := service.GetRootRouter()
	if er != nil {
		fmt.Printf("Error in setting up Root Router :  ( %v )", er.Error())
		return
	}

	router.GET("/user/:name", HelloUser)

	// Authorization group
	authRouter := router.Group("/auth")
	authRouter.Use(service.AuthMiddleware.MiddlewareFunc())
	{
		authRouter.GET("/user/:name", HelloUser)
	}

	var stop chan bool
	go func() {
		service.Start(false)
		stop <- true
	}()
	//go func(){
	//	authService.StartService(":7001")
	stop <- true
	//}()
	for i := 0; i < 2; i++ {
		<-stop
	}

}

func HelloUser(c *gin.Context) {
	name := c.Param("name")
	msg := "How have you been " + name + "?"
	//c.String(http.StatusOK, "Hello %s", name)
	c.HTML(http.StatusOK, "main/index.html", gin.H{
		"title":    "Welcome " + name,
		"userName": name,
		"msg":      msg,
	})
}
