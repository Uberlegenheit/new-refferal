package services

import (
	"new-refferal/filters"
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

		OpenBox(user *models.User) error

		UpdateReward(reward *models.Reward) error
		GetUserRewardsByID(user *models.User) ([]models.RewardShow, error)
		GetAllRewards() ([]models.RewardShow, error)
		GetTotalRewardStats() ([]models.TotalReward, error)
		GetTotalStats(req filters.PeriodInfoRequest) ([]models.TotalStats, error)
		GetTotalStakeStats(req filters.PeriodInfoRequest) ([]models.TotalStakeStats, error)
		GetFriendsStakeStats(req filters.PeriodInfoRequest) ([]models.FriendStakeStats, error)
		GetRewardPaymentStats(req filters.PeriodInfoRequest) ([]models.RewardPaymentsStats, error)
		GetUsersInvitationsStats() ([]models.InvitationsStats, error)
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
