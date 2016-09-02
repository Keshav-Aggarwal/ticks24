package utils

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/goinggo/tracelog"
	"time"
)

func Logger() gin.HandlerFunc {
	return func(c *gin.Context) {
		var logMsg string
		t := time.Now()

		// before request

		c.Next()

		// after request
		latency := time.Since(t)

		// access the status we are sending
		status := c.Writer.Status()
		username := c.Writer.Header().Get("username")
		email := c.Writer.Header().Get("email")
		if username == "" && email == "" {
			username = "Not Logged In"
		}
		if email != "" && username == "" {
			username = email
		}
		dataB := c.Writer.Size()
		if dataB == -1 {
			dataB = 0
		}
		dataK := dataB / 1024
		dataB %= 1024
		dataM := dataK / 1024
		dataK %= 1024
		dataSize := fmt.Sprint(dataM, "m ", dataK, "k ", dataB, "B   ", c.Writer.Size(), "B")
		logMsg = fmt.Sprint("| ", latency, " | ", status, " | ", dataSize, " | ",
			c.ClientIP(), " | ", username, " | ", c.Request.Host, " | ", c.Request.URL, " |")

		tracelog.Trace(c.Request.Method, "Logger", logMsg)
	}
}
