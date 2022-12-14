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
		SaveFailedDelegationTx(stake *models.Stake) (*models.Stake, error)
		GetFailedDelegations(pagination filters.Pagination) ([]models.FailedStakeShow, uint64, error)
		GetDelegationByTxHash(stake *models.Stake) (*models.Stake, error)
		GetInvitedUsersStakes(user *models.User, pagination filters.Pagination) ([]models.StakeShow, uint64, error)
		GetDelegationKey(user *models.User) (string, error)
		CheckDelegationKey(user *models.User, key string) (bool, error)

		CreatePayout(payout *models.Payout) (*models.Payout, error)
		UpdatePayout(payout *models.Payout) error
		GetPayouts(pagination filters.Pagination) ([]models.PayoutShow, uint64, error)

		OpenBox(user *models.User) (float64, error)
		GetAvailableBoxesByUserID(userID uint64) (*models.Box, error)

		UpdateReward(reward *models.Reward) error
		GetUserRewardsByID(user *models.User, pagination filters.Pagination) ([]models.RewardShow, uint64, error)
		GetAllRewards(pagination filters.Pagination) ([]models.RewardShow, uint64, error)
		GetTotalRewardStats() (*models.TotalReward, error)
		GetMyStakeSum(id uint64) (*models.StakeAndProgress, error)

		GetTotalStats(req filters.PeriodInfoRequest) (*models.TotalStats, error)
		GetBoxesStats(req filters.PeriodInfoRequest, pagination filters.Pagination) ([]models.BoxStats, uint64, error)
		GetTotalStakeStats(req filters.PeriodInfoRequest, pagination filters.Pagination) ([]models.TotalStakeStats, uint64, error)
		GetFriendsStakeStats(req filters.PeriodInfoRequest, pagination filters.Pagination) ([]models.FriendStakeStats, uint64, error)
		GetUsersInvitationsStats(pagination filters.Pagination) ([]models.InvitationsStats, uint64, error)

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
