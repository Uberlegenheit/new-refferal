package services

import (
	"fmt"
	"new-refferal/models"
)

func (s *ServiceFacade) SaveDelegationTx(stake *models.Stake) (*models.Stake, error) {
	stake, err := s.dao.SaveDelegationTx(stake)
	if err != nil {
		return nil, fmt.Errorf("dao.SaveDelegationTx: %s", err.Error())
	}

	stats, err := s.dao.GetStakeAndBoxUserStatByID(stake.UserID)
	if err != nil {
		return nil, fmt.Errorf("dao.GetStakeAndBoxUserStatByID: %s", err.Error())
	}

	boxesAvailable := int64(stats.TotalStake / 10.0)
	newBoxes := boxesAvailable - stats.TotalBoxes
	if newBoxes != 0 {
		err := s.dao.AddBoxesByUserID(stake.UserID, newBoxes)
		if err != nil {
			return nil, fmt.Errorf("dao.AddBoxesByUserID: %s", err.Error())
		}
	}

	return stake, nil
}

func (s *ServiceFacade) GetInvitedUsersStakes(user *models.User) ([]models.StakeShow, error) {
	stakes, err := s.dao.GetInvitedUsersStakes(user.ID)
	if err != nil {
		return nil, fmt.Errorf("dao.GetInvitedUsersStakes: %s", err.Error())
	}

	return stakes, nil
}
