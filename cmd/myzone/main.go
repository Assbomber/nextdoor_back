package main

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/assbomber/myzone/configs"
	"github.com/assbomber/myzone/docs"
	"github.com/assbomber/myzone/internal/auth"
	"github.com/assbomber/myzone/internal/server"
	store "github.com/assbomber/myzone/internal/store/sqlc"
	"github.com/assbomber/myzone/pkg/constants"
	"github.com/assbomber/myzone/pkg/db"
	"github.com/assbomber/myzone/pkg/logger"
	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"

	// swagger embed files
	swaggerFiles "github.com/swaggo/files"
	// gin-swagger middleware
	ginSwagger "github.com/swaggo/gin-swagger"
)

var (
	log     *logger.Logger
	queries *store.Queries
	svc     *server.Server
	redisIn *redis.Client
)

func main() {
	// Initializing swagger info
	docs.SwaggerInfo.BasePath = "/api"
	docs.SwaggerInfo.Title = "MyZone"
	docs.SwaggerInfo.Version = "1.0"

	configs.Init()
	if configs.GetString("RUNTIME_ENV") == constants.Environments.PRODUCTION {
		gin.SetMode(gin.ReleaseMode)
	} else {
		gin.SetMode(gin.DebugMode)
	}

	// logger
	log = logger.InitLogger()
	// postgres
	postgresDB := db.ConnectPostgres(log, configs.GetString("POSTGRES_URL"))
	queries = store.New(postgresDB)
	// redis
	redisIn = db.ConnectRedis(log, configs.GetString("redisHost"))

	//server
	svc = server.New(log, configs.GetString("PORT"), configs.GetInt("server.readTimeout"), configs.GetInt("server.writeTimeout"))
	go svc.Start()
	time.Sleep(2 * time.Second)
	fmt.Println(constants.LOGO)

	// registering routes
	registerRoutes()

	// Shutdown logic
	channel := make(chan os.Signal, 1)
	signal.Notify(channel, syscall.SIGINT, syscall.SIGTERM)
	<-channel
	svc.Shutdown()
	postgresDB.Close()
}

func registerRoutes() {

	// * swagger route --------------------------------
	svc.Router.GET(fmt.Sprint("/swagger"), func(ctx *gin.Context) { ctx.Redirect(301, fmt.Sprint("/docs/index.html")) })
	svc.Router.GET(fmt.Sprint("/docs/*any"), ginSwagger.WrapHandler(swaggerFiles.Handler))

	baseRoute := svc.Router.Group("/api")

	// Auth
	authService := auth.NewService(log, configs.GetString(constants.JWT_SECRET), queries, redisIn)
	authHandler := auth.NewHandler(log, authService)
	authHandler.RegisterRoutes(baseRoute)
}
