package dao

import (
	"fmt"
	"new-refferal/conf"
	"new-refferal/dao/cache"
	"new-refferal/dao/postgres"
	"new-refferal/models"
)

type (
	DAO interface {
		Postgres
		Cache
	}

	Postgres interface {
		CreateUser(user *models.User) (*models.User, error)
		CreateUserAndLink(user *models.User, code string) (*models.User, error)
		GetUserByWalletAddress(addr string) (*models.User, error)
		GetAllUsers() ([]models.User, error)

		GetLastLink() (*models.Link, error)
		GetLinkByUserID(id uint64) (*models.Link, error)

		SaveDelegationTx(stake *models.Stake) (*models.Stake, error)
		SetUserDelegationsFalse(id uint64) error
		GetInvitedUsersStakes(id uint64) ([]models.StakeShow, error)
		GetStakeAndBoxUserStatByID(id uint64) (*models.StakeAndBoxStat, error)

		AddBoxesByUserID(userID uint64, newBoxes int64) error
		OpenBoxByUserID(userID uint64) error

		SaveReward(reward *models.Reward) (*models.Reward, error)
		UpdateReward(reward *models.Reward) error
		GetUserRewardsByID(id uint64) ([]models.RewardShow, error)
		GetAllRewards() ([]models.RewardShow, error)
		GetTotalRewardStats() ([]models.TotalReward, error)
		GetUsersInvitationsStats() ([]models.InvitationsStats, error)
	}

	Cache   interface{}
	daoImpl struct {
		*postgres.Postgres
		*cache.Cache
	}
)

func New(cfg conf.Config) (DAO, error) {
	pg, err := postgres.NewPostgres(cfg.Postgres)
	if err != nil {
		return nil, fmt.Errorf("postgres.NewPostgres: %s", err.Error())
	}
	return daoImpl{
		Postgres: pg,
		Cache:    cache.NewCache(pg),
	}, nil
}
