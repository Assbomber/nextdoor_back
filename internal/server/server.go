package server

import (
	"context"
	"net/http"
	"time"

	"github.com/assbomber/myzone/pkg/constants"
	"github.com/assbomber/myzone/pkg/logger"
	"github.com/gin-gonic/gin"
	"golang.org/x/net/http2"
	"golang.org/x/net/http2/h2c"
)

type Server struct {
	logger       *logger.Logger
	port         string
	readTimeout  int
	writeTimeout int
	Router       *gin.Engine
	server       *http.Server
}

// Returns new instance of server.
func New(logger *logger.Logger, port string, readTimeout, writeTimeout int) *Server {
	return &Server{
		logger:       logger,
		port:         port,
		readTimeout:  readTimeout,
		writeTimeout: writeTimeout,
	}
}

// Start the server
func (s *Server) Start() {
	s.logger.Info(constants.PENDING + " Starting Server...")

	router := gin.Default()
	s.Router = router
	router.UseH2C = true

	router.HandleMethodNotAllowed = true

	router.NoMethod(func(c *gin.Context) {
		c.JSON(http.StatusMethodNotAllowed, gin.H{"status": "method not allowed"})
	})
	router.NoRoute(func(c *gin.Context) {
		c.JSON(http.StatusNotFound, gin.H{"status": "not found"})
	})

	bindingAddress := ":" + s.port
	h2s := &http2.Server{}
	s.server = &http.Server{
		Addr:           bindingAddress,
		Handler:        h2c.NewHandler(router, h2s),
		ReadTimeout:    time.Duration(s.readTimeout) * time.Second,
		WriteTimeout:   time.Duration(s.writeTimeout) * time.Second,
		MaxHeaderBytes: 1 << 20,
	}

	s.logger.Info(constants.SUCCESS + " Started Server Successfuly on port: " + s.port)
	err := s.server.ListenAndServe()
	if err != nil && err != http.ErrServerClosed {
		s.logger.Fatal(constants.FAILURE + " Error starting Server")
	}
}

// Shutdowns the server gracefully within 20 sec
func (s *Server) Shutdown() {
	s.logger.Info(constants.PENDING + " Graceful Shutdown Started...")
	if s.server == nil {
		s.logger.Error(constants.WARNING+"no server found for shutdown", nil)
		return
	}

	ctx, cancel := context.WithTimeout(context.TODO(), 20*time.Second)
	defer cancel()
	err := s.server.Shutdown(ctx)
	if err != nil {
		s.logger.Fatal(constants.FAILURE + " Error Shutting download server")
	}
	s.logger.Info(constants.STOP + " Server stopped successfuly")
}
