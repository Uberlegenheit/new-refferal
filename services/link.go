package services

import (
	"fmt"
	"new-refferal/models"
)

func (s *ServiceFacade) GetLinkByUserID(user *models.User) (*models.Link, error) {
	link, err := s.dao.GetLinkByUserID(user.ID)
	if err != nil {
		return nil, fmt.Errorf("dao.GetLinkByUserID: %s", err.Error())
	}

	return link, nil
}
