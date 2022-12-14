package services

import (
	"fmt"
	"new-refferal/filters"
	"new-refferal/models"
)

func (s *ServiceFacade) UpdateReward(reward *models.Reward) error {
	err := s.dao.UpdateReward(reward)
	if err != nil {
		return fmt.Errorf("dao.UpdateReward: %s", err.Error())
	}

	return nil
}

func (s *ServiceFacade) GetUserRewardsByID(user *models.User, pagination filters.Pagination) ([]models.RewardShow, uint64, error) {
	rewards, length, err := s.dao.GetUserRewardsByID(user.ID, pagination)
	if err != nil {
		return nil, length, fmt.Errorf("dao.GetUserRewardsByID: %s", err.Error())
	}

	return rewards, length, nil
}

func (s *ServiceFacade) GetAllRewards(pagination filters.Pagination) ([]models.RewardShow, uint64, error) {
	rewards, length, err := s.dao.GetAllRewards(pagination)
	if err != nil {
		return nil, length, fmt.Errorf("dao.GetAllRewards: %s", err.Error())
	}

	return rewards, length, nil
}
