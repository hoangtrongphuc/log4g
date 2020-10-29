## *log4g* - Logger interface 
Some logger instance:
* Logrus ( [github.com/sirupsen/logrus](github.com/sirupsen/logrus) )
* Zap ([go.uber.org/zap](go.uber.org/zap))

### 1. Setup & install
#### 1.1: Using for any project
Example with woker
```go
package main

import (
	"context"
	"fmt"
	"logger/log4g"
)

func main() {
	ctxLogger, err := log4g.NewContextLogger(context.Background(), log4g.Configuration{
    		LogLevel:        log4g.INFO,
    		JSONFormat:      true,
    		TimestampFormat: "2006-01-02 15:04:05",
    		FieldMap: log4g.FieldMap{
    			log4g.FieldKeyTime:  "_datetime",
    			log4g.FieldKeyMsg:   "message",
    			log4g.FieldKeyLevel: "severity",
    			log4g.FieldKeyFunc:  "func",
    		},
    	}, log4g.InstanceZapLogger)
    	if err != nil {
    		fmt.Println("Error", err)
    		return
    	}
    	ctxLogger = ctxLogger.WithFields(log4g.Fields{
    		log4g.AppName:    "XXX",
    		log4g.XRequestID: "xRequestId",
    	})
    	ctxLogger.Info("%v", "This is info")
    	ctxLogger.Error("%v", "BUGGG")
    	ctxLogger.Msg("show me")
}
```
Result with zaplog:
```text
{"severity":"info","_datetime":"2020-10-27 14:41:06","caller":"logger-interface/main.go:42","message":"show me","app_name":"XXX","x_request_id":"xRequestId","_err_0":"BUGGG","_info_0":"This is info"}
```
Result with logrus:
```text
{"_datetime":"2020-10-27 14:40:29","_err_0":"BUGGG","_info_0":"This is info","app_name":"XXX","message":"show me","severity":"info","x_request_id":"xRequestId"}
```
Example with Go Gin
```go
// Logger Middleware
func LoggerMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		logger, err := log4g.NewContextLogger(c, log4g.Configuration{
			LogLevel:        log4g.INFO,
			JSONFormat:      true,
			TimestampFormat: "2006-01-02 15:04:05",
			FieldMap: log4g.FieldMap{
				log4g.FieldKeyTime:  "_datetime",
				log4g.FieldKeyMsg:   "message",
				log4g.FieldKeyLevel: "severity",
				log4g.FieldKeyFunc:  "func",
			},
		}, log4g.InstanceZapLogger)
		if err != nil {
			fmt.Println("Error", err)
			return
		}
		ctxLogger := logger.WithFields(log4g.Fields{
			log4g.AppName:    "sale-app",
			log4g.XRequestID: "x12n123",
		})
		c.Next()
		ctxLogger.Msg("show me")
	}
}
```
```go
// Main
package main

import (
	"context"
	"fmt"
	"github.com/gin-gonic/gin"
	"logger/log4g"
)

func main() {
	r := gin.Default()
    	r.Use(LoggerMiddleware())
    	r.GET("/ping", func(c *gin.Context) {
    		log4g.Get(c).Info("%v", "Ping Handler")
    		log4g.Get(c).Info("%v", "Everything OK")
    		c.JSON(200, gin.H{
    			"message": "pong",
    		})
    	})
    	r.GET("/pong", func(c *gin.Context) {
    		log4g.Get(c).Error("%v", "Pong Handler")
    		c.JSON(200, gin.H{
    			"message": "pong",
    		})
    	})
    	r.Run()
}
```
Result:
```text
{"severity":"info","_datetime":"2020-10-28 15:41:01","caller":"logger-interface/main.go:88","func":"main.LoggerMiddleware.func1","message":"show me","x_request_id":"x12n123","app_name":"logger-interface","_err_0":"Pong Handler"}
```

### Referer document
* https://cloud.google.com/error-reporting/docs/setup/
* https://github.com/uber-go/zap
* github.com/sirupsen/logrus