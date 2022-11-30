package services

import (
	"new-refferal/models"
	"time"

	"new-refferal/conf"
	"new-refferal/dao"
	"new-refferal/helpers/scheduler"
)

type (
	Service interface {
		LogInOrRegister(user *models.User) (*models.User, error)
	}
	Scheduler interface {
		AddProcessWithInterval(process scheduler.Process, interval time.Duration)
		AddProcessWithPeriod(process scheduler.Process, period time.Duration)
		EveryDayAt(process scheduler.Process, hour int, minutes int)
		EveryMonthAt(process scheduler.Process, day int, hours int, minutes int)
	}

	ServiceFacade struct {
		cfg conf.Config
		dao dao.DAO
	}
)

func NewService(cfg conf.Config, dao dao.DAO) (*ServiceFacade, error) {

	return &ServiceFacade{
		cfg: cfg,
		dao: dao,
	}, nil
}
