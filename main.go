package main

import (
	"fmt"

	"github.com/assbomber/myzone/configs"
	"github.com/assbomber/myzone/constants"
	"github.com/assbomber/myzone/utils"
)

var (
	logger *utils.Logger
)

func main() {
	configs.Init()
	logger = utils.InitLogger()
	fmt.Println(constants.LOGO)
}
