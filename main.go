package main

import (
	"github.com/joho/godotenv"
	"github.com/roylee0704/gron"
	"go.uber.org/zap"
	"new-refferal/api"
	"new-refferal/conf"
	"new-refferal/dao"
	"new-refferal/helpers/modules"
	"new-refferal/log"
	"new-refferal/services"
	"os"
	"os/signal"
	"syscall"
)

func init() {
	if err := godotenv.Load(); err != nil {
		log.Error("No .env file found")
	}
}

func main() {
	cfg, err := conf.GetNewConfig()
	if err != nil {
		log.Fatal("can`t read config from file", zap.Error(err))
	}

	d, err := dao.New(cfg)
	if err != nil {
		log.Fatal("dao.New", zap.Error(err))
	}

	s, err := services.NewService(cfg, d)
	if err != nil {
		log.Fatal("services.NewService", zap.Error(err))
	}

	cron := gron.New()
	s.InitCron(cron)
	cron.Start()
	defer cron.Stop()

	a, err := api.NewAPI(cfg, s)
	if err != nil {
		log.Fatal("api.NewAPI", zap.Error(err))
	}

	mds := []modules.Module{a}

	modules.Run(mds)

	var gracefulStop = make(chan os.Signal, 1)
	signal.Notify(gracefulStop, syscall.SIGTERM, syscall.SIGINT)

	<-gracefulStop
	modules.Stop(mds)
}
