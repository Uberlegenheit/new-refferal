package services

import (
	"fmt"
	"new-refferal/models"
	"strconv"
	"strings"
	"time"
)

func (s *ServiceFacade) OpenBox(user *models.User) (*models.RewardShow, error) {
	info, err := s.dao.GetAvailableBoxesByUserID(user.ID)
	if err != nil {
		return nil, fmt.Errorf("dao.GetAvailableBoxesByUserID: %s", err.Error())
	}

	if info.Available <= 0 {
		return nil, fmt.Errorf("you don't have boxes to open")
	}
	date := strings.Split(time.Now().Format("2006-01-02"), "-")
	year, _ := strconv.Atoi(date[0])
	month, _ := strconv.Atoi(date[1])
	day, _ := strconv.Atoi(date[2])
	days := time.Date(year, time.Month(month), 0, 0, 0, 0, 0, time.UTC).Day()

	if day == days {
		days = 1
	} else {
		days = days - day
	}

	pool, err := s.dao.GetRewardsPool()
	if err != nil {
		return nil, fmt.Errorf("dao.GetRewardsPool: %s", err.Error())
	}

	/* formulas */

	err = s.dao.CreateAndUpdateRewardsState(pool, user, 5.5 /********************/)
	if err != nil {
		return nil, fmt.Errorf("dao.CreateAndUpdateRewardsState: %s", err.Error())
	}

	return nil, nil
}
