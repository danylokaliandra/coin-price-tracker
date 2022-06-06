package main

import (
	"fmt"
	"os"

	"github.com/eliasbokreta/tracklet/cmd"
	"github.com/eliasbokreta/tracklet/pkg/utils"
	log "github.com/sirupsen/logrus"
)

func init() {
	log.SetFormatter(&log.TextFormatter{
		DisableColors: false,
		FullTimestamp: true,
	})

	if err := utils.LoadConfig(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func main() {
	if err := cmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
