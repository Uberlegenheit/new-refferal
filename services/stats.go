package services

import (
	"fmt"
	"math"
	"new-refferal/filters"
	"new-refferal/models"
)

func (s *ServiceFacade) GetTotalRewardStats() (*models.TotalReward, error) {
	rewards, err := s.dao.GetTotalRewardStats()
	if err != nil {
		return nil, fmt.Errorf("dao.GetTotalRewardStats: %s", err.Error())
	}

	return rewards, nil
}

func (s *ServiceFacade) GetMyStakeSum(id uint64) (*models.StakeAndProgress, error) {
	stake, err := s.dao.GetMyStakeSum(id)
	if err != nil {
		return nil, fmt.Errorf("dao.GetMyStakeSum: %s", err.Error())
	}

	toNewBox := stake.TotalStake / StakeToBox
	prcToOne := toNewBox - float64(int(toNewBox))
	stake.Progress = math.Floor(prcToOne*100000.0) / 100000.0

	return stake, nil
}

func (s *ServiceFacade) GetTotalStats(req filters.PeriodInfoRequest) ([]models.TotalStats, error) {
	stats, err := s.dao.GetTotalStats(req)
	if err != nil {
		return nil, fmt.Errorf("dao.GetTotalStats: %s", err.Error())
	}

	return stats, nil
}

func (s *ServiceFacade) GetTotalStakeStats(req filters.PeriodInfoRequest) ([]models.TotalStakeStats, error) {
	stats, err := s.dao.GetTotalStakeStats(req)
	if err != nil {
		return nil, fmt.Errorf("dao.GetTotalStakeStats: %s", err.Error())
	}

	return stats, nil
}

func (s *ServiceFacade) GetFriendsStakeStats(req filters.PeriodInfoRequest) ([]models.FriendStakeStats, error) {
	stats, err := s.dao.GetFriendsStakeStats(req)
	if err != nil {
		return nil, fmt.Errorf("dao.GetFriendsStakeStats: %s", err.Error())
	}

	return stats, nil
}

func (s *ServiceFacade) GetRewardPaymentStats(req filters.PeriodInfoRequest) ([]models.RewardPaymentsStats, error) {
	stats, err := s.dao.GetRewardPaymentStats(req)
	if err != nil {
		return nil, fmt.Errorf("dao.GetRewardPaymentStats: %s", err.Error())
	}

	return stats, nil
}

func (s *ServiceFacade) GetUsersInvitationsStats() ([]models.InvitationsStats, error) {
	stats, err := s.dao.GetUsersInvitationsStats()
	if err != nil {
		return nil, fmt.Errorf("dao.GetUsersInvitationsStats: %s", err.Error())
	}

	return stats, nil
}
