package main

import (
	"log"
	"net/http"
	"os"
	"fmt"
	_ "io"
	"io"
	"time"
	_ "github.com/natefinch/lumberjack"
	//"github.com/joho/godotenv"
	"github.com/gin-contrib/cors"
	"github.com/getsentry/raven-go"
	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/contrib/sessions"
	"github.com/gin-gonic/contrib/gzip"
	"github.com/gin-contrib/sentry"
	"github.com/joho/godotenv"
)

func init() {
	// Load configuration env
	err := godotenv.Load()
	if err != nil {
		log.Fatal("OrgError loading .env file", err)
	}
	//raven.SetEnvironment(os.Getenv("ENVIRONMENT"))
	//raven.SetDSN(os.Getenv("RAVEN_DSN"))
	// End
}

func main() {
	log.Print("Start PO financing Service")

	// Logger
	//log.SetOutput(&lumberjack.Logger{
	//	Filename:   "logs/_payment_service.log",
	//	MaxSize:    10, // megabytes
	//	MaxBackups: 10,
	//	MaxAge:     30,   //days
	//	Compress:   true, // disabled by default
	//})
	//log.SetFlags(log.Lshortfile | log.LstdFlags)
	// end Logger

	// Logger
	logFile, err := os.OpenFile("_sign_service_api.log", os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0600)
	if err != nil {
		panic(err)
	}
	gin.DefaultWriter = io.MultiWriter(logFile, os.Stdout)
	log.SetOutput(gin.DefaultWriter) // You may need this
	log.SetFlags(log.Lshortfile | log.LstdFlags)
	// end Logger

	// Load configuration
	redisHost := os.Getenv("REDIS_HOST")
	redisPassword := os.Getenv("REDIS_PASSWORD")
	sessionPrefix := os.Getenv("SESSION_PREFIX")
	// End

	// Setting router
	router := gin.New()
	// Define session
	store, _ := sessions.NewRedisStore(10, "tcp", redisHost, redisPassword, []byte(""))
	router.Use(sessions.Sessions(sessionPrefix, store))

	router.Use(RouterMiddleware())
	router.Use(CORSMiddleware())
	router.Use(sentry.Recovery(raven.DefaultClient, false))
	router.Use(gzip.Gzip(gzip.DefaultCompression))

	// Router Index
	index := router.Group("/")
	{
		index.GET("/", func(context *gin.Context) {
			context.JSON(http.StatusOK, gin.H{"status": 1, "message": "Sign Service API"})
		})
	}

	log.Printf(":%d", os.Getenv("SERVICE_PORT"))
	router.Run(fmt.Sprintf(":%s", os.Getenv("SERVICE_PORT")))
}

func RouterMiddleware() gin.HandlerFunc {
	return func(context *gin.Context) {
		// begin before request
		log.Print("RouterMiddleware: Before request")
		// end

		context.Next()

		// after request
		log.Print("RouterMiddleware: End request")
		// end
	}
}

func CORSMiddleware() gin.HandlerFunc {
	// Gin Cors setting
	return cors.New(cors.Config{
		//AllowOrigins:     []string{"*"},
		AllowMethods:     []string{"POST", "GET", "OPTIONS", "PUT", "DELETE", "UPDATE"},
		AllowHeaders:     []string{"Content-Type", "Origin", "Device-Type", "Device-Id", "Authorization", "*"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
		AllowOriginFunc: func(origin string) bool {
			return true
		},
		MaxAge: 12 * time.Hour,
	})
}
