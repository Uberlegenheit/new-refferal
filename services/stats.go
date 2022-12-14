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

func (s *ServiceFacade) GetTotalStats(req filters.PeriodInfoRequest) (*models.TotalStats, error) {
	stats, err := s.dao.GetTotalStats(req)
	if err != nil {
		return nil, fmt.Errorf("dao.GetTotalStats: %s", err.Error())
	}

	return stats, nil
}

func (s *ServiceFacade) GetTotalStakeStats(req filters.PeriodInfoRequest, pagination filters.Pagination) ([]models.TotalStakeStats, uint64, error) {
	stats, length, err := s.dao.GetTotalStakeStats(req, pagination)
	if err != nil {
		return nil, length, fmt.Errorf("dao.GetTotalStakeStats: %s", err.Error())
	}

	return stats, length, nil
}

func (s *ServiceFacade) GetBoxesStats(req filters.PeriodInfoRequest, pagination filters.Pagination) ([]models.BoxStats, uint64, error) {
	stats, length, err := s.dao.GetBoxesStats(req, pagination)
	if err != nil {
		return nil, length, fmt.Errorf("dao.GetBoxesStats: %s", err.Error())
	}

	return stats, length, nil
}

func (s *ServiceFacade) GetFriendsStakeStats(req filters.PeriodInfoRequest, pagination filters.Pagination) ([]models.FriendStakeStats, uint64, error) {
	stats, length, err := s.dao.GetFriendsStakeStats(req, pagination)
	if err != nil {
		return nil, length, fmt.Errorf("dao.GetFriendsStakeStats: %s", err.Error())
	}

	return stats, length, nil
}

func (s *ServiceFacade) GetUsersInvitationsStats(pagination filters.Pagination) ([]models.InvitationsStats, uint64, error) {
	stats, length, err := s.dao.GetUsersInvitationsStats(pagination)
	if err != nil {
		return nil, length, fmt.Errorf("dao.GetUsersInvitationsStats: %s", err.Error())
	}

	return stats, length, nil
}
