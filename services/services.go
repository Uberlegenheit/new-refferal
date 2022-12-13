package services

import (
	"github.com/gin-gonic/gin"
	"net/http"
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
		GetUserByWalletAddress(addr string) (*models.User, error)

		GetLinkByUserID(user *models.User) (*models.Link, error)

		SaveDelegationTx(stake *models.Stake, user *models.User) (*models.Stake, error)
		GetDelegationByTxHash(stake *models.Stake) (*models.Stake, error)
		GetInvitedUsersStakes(user *models.User, pagination filters.Pagination) ([]models.StakeShow, error)
		GetDelegationKey(user *models.User) (string, error)
		CheckDelegationKey(user *models.User, key string) (bool, error)

		CreatePayout(payout *models.Payout) (*models.Payout, error)
		UpdatePayout(payout *models.Payout) error
		GetPayouts(pagination filters.Pagination) ([]models.PayoutShow, error)

		OpenBox(user *models.User) error
		GetAvailableBoxesByUserID(userID uint64) (*models.Box, error)

		UpdateReward(reward *models.Reward) error
		GetUserRewardsByID(user *models.User, pagination filters.Pagination) ([]models.RewardShow, error)
		GetAllRewards(pagination filters.Pagination) ([]models.RewardShow, error)
		GetTotalRewardStats() (*models.TotalReward, error)
		GetMyStakeSum(id uint64) (*models.StakeAndProgress, error)

		GetTotalStats(req filters.PeriodInfoRequest, pagination filters.Pagination) (*models.TotalStats, error)
		GetBoxesStats(req filters.PeriodInfoRequest, pagination filters.Pagination) ([]models.BoxStats, error)
		GetTotalStakeStats(req filters.PeriodInfoRequest, pagination filters.Pagination) ([]models.TotalStakeStats, error)
		GetFriendsStakeStats(req filters.PeriodInfoRequest, pagination filters.Pagination) ([]models.FriendStakeStats, error)
		GetUsersInvitationsStats(pagination filters.Pagination) ([]models.InvitationsStats, error)

		CreateToken(walletAddr string) (*models.TokenDetails, error)
		CreateAuth(walletAddr string, td *models.TokenDetails) error
		ExtractTokenMetadata(c *gin.Context) (*models.AccessDetails, error)
		Refresh(r *http.Request) (*models.TokenDetails, error)

		FetchAuth(authD *models.AccessDetails) (string, error)
		DeleteAuth(UUID ...string) error
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
