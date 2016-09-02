package ticks24

import (
	"errors"
	"fmt"
	"github.com/appleboy/gin-jwt"
	"github.com/braintree/manners"
	"github.com/gin-gonic/contrib/gzip"
	"github.com/gin-gonic/contrib/static"
	"github.com/gin-gonic/gin"
	"github.com/goinggo/tracelog"
	"github.com/itsjamie/gin-cors"
	"github.com/vivek-yadav/ticks24/config"
	"github.com/vivek-yadav/ticks24/db/mongo"
	"github.com/vivek-yadav/ticks24/middleware/auth"
	"github.com/vivek-yadav/ticks24/utils"
	"html/template"
	"net/http"
	"net/rpc"
	"strconv"
	"strings"
	"time"
)

type Service struct {
	Config         config.Config `json:"Config"`
	Engine         *gin.Engine
	RootRouter     *gin.RouterGroup
	AppDB          mongo.AppDB
	AuthMiddleware *jwt.GinJWTMiddleware
}

var serviceI *Service

// This returns a new Instance of User Management Service
func GetInstance() (*Service, error) {
	if serviceI == nil {
		service := Service{}
		result, err := service.Config.SetEnvArgs()
		if result == false && err != nil {
			return nil, errors.New("ERROR : Environment Variables were not proper ( " + err.Error() + " )")
		}
		serviceI = &service
	}
	return serviceI, nil
}

// This SetConfig function takes filePath of the config file
// and loads the User Management Service Instance with specified settings
// if some error occurs it throws error.
// if no file is sent in filePath param then default settings are loaded
func (this *Service) SetConfigFile(filePath string) (bool, error) {
	return this.Config.SetFromFile(filePath)
}

// This sets configuration from command line arguments.
// Use this when you think your users might want to give command line arguments.
// Call this after SetConfig if you want it to have more priority.
func (this *Service) SetCmdArgs() (bool, error) {
	return this.Config.SetFromCmdArgs()
}

// This function is used to start-up the service with given settings or default settings
// If you send isblocking true then the system waits for the server to end first before return
// Else the call starts the server and returns, then it is up to you to hold the system to keep the
// service running.
func (this *Service) Start(isBlocking bool) {
	var paths gin.RoutesInfo
	this.RootRouter.GET("/_routes", func(c *gin.Context) {
		c.JSON(http.StatusOK, paths)
	})
	paths = this.Engine.Routes()
	this.Config.Show()
	if isBlocking {
		r := make(chan bool)
		go func(v chan bool) {
			serverPort := fmt.Sprintf(":%v", this.Config.WebServer.Port)
			manners.ListenAndServe(serverPort, this.Engine)
			v <- true
		}(r)
		<-r
	} else {
		go func() {
			serverPort := fmt.Sprintf(":%v", this.Config.WebServer.Port)
			manners.ListenAndServe(serverPort, this.Engine)
		}()
	}

}

func (this *Service) GetRootRouter() (*gin.RouterGroup, error) {
	if r, err := this.InitService(); r == false && err != nil {
		return nil, err
	}
	if r, err := this.setupRootRouter(); r == false && err != nil {
		return nil, err
	}
	return this.RootRouter, nil
}

// This function sets up the root routing
func (this *Service) setupRootRouter() (bool, error) {
	this.RootRouter = this.Engine.Group("/")
	//routes.Setup(this.RootRouter)
	return true, nil
}

func (this *Service) ConnectAuthService() (client *rpc.Client, er error) {
	typeCon := "tcp"

	// TODO : to setup other type of connections
	//if this.Config.AuthService.IsHttp {
	//	typeCon = "http"
	//}
	path := this.Config.AuthService + ":" + strconv.Itoa(this.Config.AuthService.Port)
	client, er = rpc.DialHTTP(typeCon, path)
	if er != nil {
		er = errors.New("ERROR : Failed to conenct to Auth Service: (" + path + " (\n\t" + er.Error() + "\n)")
	}
	return
}

func (this *Service) InitService() (bool, error) {
	if len(this.Config.AppDatabases) > 0 {
		this.AppDB.Config = &this.Config.AppDatabases[0]
		var er error
		if this.AppDB, er = this.AppDB.Setup(); er != nil {
			return false, errors.New("ERROR : Failed to conenct to AppDatabase[0] (\n\t" + er.Error() + "\n)")
		}
	}

	_, er := this.ConnectAuthService()
	if er != nil {
		return false, errors.New("ERROR : Failed to conenct to Auth Service (\n\t" + er.Error() + "\n)")
	}
	this.AuthMiddleware = auth.Setup(this)

	router := gin.New()
	if this.Config.LogConfig.Path != "" {
		var level int32
		switch strings.ToUpper(this.Config.LogConfig.Level) {
		case "TRACE":
			level = tracelog.LevelTrace
		case "INFO":
			level = tracelog.LevelInfo
		case "WARN":
			level = tracelog.LevelWarn
		case "ERROR":
			level = tracelog.LevelError
		default:
			this.Config.LogConfig.Level = "TRACE"
			level = tracelog.LevelTrace
		}
		tracelog.StartFile(level, this.Config.LogConfig.Path, int(this.Config.LogConfig.Days))
		router.Use(utils.Logger())
	} else {
		router.Use(gin.Logger())
	}
	router.Use(gin.Recovery())
	router.Use(gzip.Gzip(gzip.DefaultCompression))
	router.Use(static.Serve("/", static.LocalFile(this.Config.FrontEnd.ViewsPath, true)))
	html, err := template.New("").Delims(this.Config.FrontEnd.TemplateDelimiterStart, this.Config.FrontEnd.TemplateDelimiterEnd).ParseGlob(this.Config.FrontEnd.TemplatesPath + "/**/*")
	if err != nil {
		return false, errors.New("ERROR : Failed to set Templates Path for Server : ( " + err.Error() + " )")
	}
	router.SetHTMLTemplate(html)

	// Apply the middleware to the router (works with groups too)
	router.Use(cors.Middleware(cors.Config{
		Origins:         "*",
		Methods:         "GET, PUT, POST, DELETE",
		RequestHeaders:  "Origin, Authorization, Content-Type",
		ExposedHeaders:  "",
		MaxAge:          50 * time.Second,
		Credentials:     true,
		ValidateHeaders: false,
	}))

	this.Engine = router

	return true, nil
}

// This function is used to stop the service
func (this *Service) Stop() (bool, error) {
	return true, nil
}

// This function is used to Re-Start the service
func (this *Service) ReStart() (bool, error) {
	return true, nil
}
