package main

import (
	"github.com/assbomber/myzone/configs"
	"github.com/assbomber/myzone/pkg/constants"
	"github.com/assbomber/myzone/pkg/logger"
)

var (
	log *logger.Logger
)

func main() {
	configs.Init()
	log = logger.InitLogger()
	log.Info(constants.LOGO)
}
