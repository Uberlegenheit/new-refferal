package services

import (
	"fmt"
	"new-refferal/models"
)

func (s *ServiceFacade) FetchAuth(authD *models.AccessDetails) (string, error) {
	if authD == nil {
		return "", fmt.Errorf("empty token details")
	}

	walletAddr, ok, err := s.dao.GetAuthToken(authD.AccessUuid)
	if err != nil || !ok {
		return "", err
	}

	return walletAddr.(string), nil
}

func (s *ServiceFacade) DeleteAuth(UUID ...string) error {
	for i := range UUID {
		err := s.dao.RemoveAuthToken(UUID[i])
		if err != nil {
			return err
		}
	}

	return nil
}
