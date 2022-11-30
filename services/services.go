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

		GetLinkByUserID(user *models.User) (*models.Link, error)

		SaveDelegationTx(stake *models.Stake) (*models.Stake, error)
		GetInvitedUsersStakes(user *models.User) ([]models.StakeShow, error)

		GetUserRewardsByID(user *models.User) ([]models.RewardShow, error)
		GetAllRewards() ([]models.RewardShow, error)
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
