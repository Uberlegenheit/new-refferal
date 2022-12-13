package services

import (
	"fmt"
	"new-refferal/filters"
	"new-refferal/models"
)

func (s *ServiceFacade) CreatePayout(payout *models.Payout) (*models.Payout, error) {
	saved, err := s.dao.CreatePayout(payout)
	if err != nil {
		return nil, fmt.Errorf("dao.CreatePayout: %s", err.Error())
	}

	return saved, nil
}

func (s *ServiceFacade) UpdatePayout(payout *models.Payout) error {
	err := s.dao.UpdatePayout(payout)
	if err != nil {
		return fmt.Errorf("dao.UpdatePayout: %s", err.Error())
	}

	return nil
}

func (s *ServiceFacade) GetPayouts(pagination filters.Pagination) ([]models.PayoutShow, error) {
	payouts, err := s.dao.GetPayouts(pagination)
	if err != nil {
		return nil, fmt.Errorf("dao.GetPayouts: %s", err.Error())
	}

	return payouts, nil
}
