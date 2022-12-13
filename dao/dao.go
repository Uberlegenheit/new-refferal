package dao

import (
	"fmt"
	"new-refferal/conf"
	"new-refferal/dao/cache"
	"new-refferal/dao/postgres"
	"new-refferal/filters"
	"new-refferal/models"
	"time"
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
		SaveDelegationTxAndAddBoxes(stake *models.Stake) (*models.Stake, error)
		GetDelegationByTxHash(stake *models.Stake) (*models.Stake, error)
		SetUserDelegationsFalse(id uint64) error
		GetInvitedUsersStakes(id uint64, pagination filters.Pagination) ([]models.StakeShow, error)
		GetStakeAndBoxUserStatByID(id uint64) (*models.StakeAndBoxStat, error)

		CreatePayout(payout *models.Payout) (*models.Payout, error)
		UpdatePayout(payout *models.Payout) error
		GetPayouts(pagination filters.Pagination) ([]models.PayoutShow, error)

		AddBoxesByUserID(userID uint64, newBoxes int64) error
		OpenBoxByUserID(userID uint64) error
		GetAvailableBoxesByUserID(userID uint64) (*models.Box, error)

		SaveReward(reward *models.Reward) (*models.Reward, error)
		UpdateReward(reward *models.Reward) error
		SaveTXAndUpdateReward(info *models.StakeAndBoxStat, newStake, reward float64) error
		GetUserRewardsByID(id uint64, pagination filters.Pagination) ([]models.RewardShow, error)
		GetAllRewards(pagination filters.Pagination) ([]models.RewardShow, error)
		GetTotalRewardStats() (*models.TotalReward, error)
		GetMyStakeSum(id uint64) (*models.StakeAndProgress, error)

		GetTotalStats(req filters.PeriodInfoRequest, pagination filters.Pagination) (*models.TotalStats, error)
		GetBoxesStats(req filters.PeriodInfoRequest, pagination filters.Pagination) ([]models.BoxStats, error)
		GetTotalStakeStats(req filters.PeriodInfoRequest, pagination filters.Pagination) ([]models.TotalStakeStats, error)
		GetFriendsStakeStats(req filters.PeriodInfoRequest, pagination filters.Pagination) ([]models.FriendStakeStats, error)
		GetUsersInvitationsStats(pagination filters.Pagination) ([]models.InvitationsStats, error)

		GetRewardsPool() (*models.RewardsPool, error)
		UpdateRewardsPool(pool *models.RewardsPool) error
		SetDailyPoolLimit(pool *models.RewardsPool) error

		CreateAndUpdateRewardsState(pool *models.RewardsPool, user *models.User, amount float64) error
	}

	Cache interface {
		AddAuthToken(key string, item interface{}, expiration time.Duration) error
		GetAuthToken(token string) (interface{}, bool, error)
		RemoveAuthToken(key string) error

		CacheSave(key string, item interface{}, expiration time.Duration) error
		CacheGet(token string) (interface{}, bool, error)
		CacheRemove(key string) error
	}
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
