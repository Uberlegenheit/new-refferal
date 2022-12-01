package services

import (
	"fmt"
	"new-refferal/models"
)

func (s *ServiceFacade) UpdateReward(reward *models.Reward) error {
	err := s.dao.UpdateReward(reward)
	if err != nil {
		return fmt.Errorf("dao.UpdateReward: %s", err.Error())
	}

	return nil
}

func (s *ServiceFacade) GetUserRewardsByID(user *models.User) ([]models.RewardShow, error) {
	rewards, err := s.dao.GetUserRewardsByID(user.ID)
	if err != nil {
		return nil, fmt.Errorf("dao.GetUserRewardsByID: %s", err.Error())
	}

	return rewards, nil
}

func (s *ServiceFacade) GetAllRewards() ([]models.RewardShow, error) {
	rewards, err := s.dao.GetAllRewards()
	if err != nil {
		return nil, fmt.Errorf("dao.GetAllRewards: %s", err.Error())
	}

	return rewards, nil
}
