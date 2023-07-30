package main

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/assbomber/myzone/configs"
	"github.com/assbomber/myzone/internal/auth"
	"github.com/assbomber/myzone/internal/server"
	store "github.com/assbomber/myzone/internal/store/sqlc"
	"github.com/assbomber/myzone/pkg/constants"
	"github.com/assbomber/myzone/pkg/db/postgres"
	"github.com/assbomber/myzone/pkg/logger"
	"github.com/gin-gonic/gin"
)

var (
	log     *logger.Logger
	queries *store.Queries
	svc     *server.Server
)

func main() {
	configs.Init()
	if configs.GetString("RUNTIME_ENV") == constants.Environments.PRODUCTION {
		gin.SetMode(gin.ReleaseMode)
	} else {
		gin.SetMode(gin.DebugMode)
	}

	// logger
	log = logger.InitLogger()
	// postgres
	postgresDB := postgres.Connect(log, configs.GetString("POSTGRES_URL"))
	queries = store.New(postgresDB)

	//server
	svc = server.New(log, configs.GetString("server.port"), configs.GetInt("server.readTimeout"), configs.GetInt("server.writeTimeout"))
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

	baseRoute := svc.Router.Group("/api")

	// Auth
	authService := auth.NewService(log, configs.GetString(constants.JWT_SECRET), queries)
	authHandler := auth.NewHandler(log, authService)
	authHandler.RegisterRoutes(baseRoute)
}
