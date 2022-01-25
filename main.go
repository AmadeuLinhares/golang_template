package main

import (
	"syscall"

	pocConfig "github.com/tradersclub/TCTemplateBack/config"
	pocServer "github.com/tradersclub/TCTemplateBack/server"
	conf "github.com/tradersclub/TCUtils/config"
	confCommon "github.com/tradersclub/TCUtils/config/common"
	"github.com/tradersclub/TCUtils/logger"
	"github.com/tradersclub/TCUtils/ossignal"
)

var (
	server pocServer.Server
	cfg    *confCommon.Cfg
)

func init() {
	cfg = conf.New()
	if err := cfg.Config.LoadConfig(&pocConfig.ConfigGlobal); err != nil {
		logger.Fatal(err.Error())
	}
	logger.Info("Config loaded")
	server = pocServer.New()
}

func main() {
	go server.Start()
	defer server.Stop()

	go ossignal.On(ossignal.Funcs{
		syscall.SIGTERM: func(exit chan struct{}) bool {
			logger.Info("TERMINATED BY k8s")
			server.Stop()
			exit <- struct{}{}
			return true
		},
		syscall.SIGHUP: func(exit chan struct{}) bool {
			if err := cfg.Config.LoadConfig(&pocConfig.ConfigGlobal); err != nil {
				logger.Error("cannot reload config because: ", err.Error())
			}
			logger.Info("Reloading conf by signal")
			server.ReloadConnections()
			return false
		},
	})

	<-ossignal.Exit
}
