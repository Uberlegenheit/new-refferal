package services

import (
	"fmt"
	"math"
	"math/rand"
	"new-refferal/models"
	"strconv"
	"strings"
	"time"
)

const (
	usersToParticipate = 10.0
)

func (s *ServiceFacade) OpenBox(user *models.User) (float64, error) {
	info, err := s.dao.GetAvailableBoxesByUserID(user.ID)
	if err != nil {
		return 0, fmt.Errorf("dao.GetAvailableBoxesByUserID: %s", err.Error())
	}

	if info.Available <= 0 {
		return 0, fmt.Errorf("you don't have boxes to open")
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
		return 0, fmt.Errorf("dao.GetRewardsPool: %s", err.Error())
	}

	x1, x2 := 1.0, usersToParticipate
	y1, y2 := 1.0, (pool.Available/float64(days))*0.9

	a := math.Log(y1) - (math.Log(y1/y2)*x1)/(x1-x2)
	b := math.Log(y1/y2) / (x1 - x2)

	rand.Seed(time.Now().UnixNano())
	min := 1
	max := usersToParticipate
	r := rand.Intn(int(max)-min+1) + min

	winAmount := -(a * math.Exp(b*float64(r)))

	pool.Available = pool.Available - winAmount
	pool.Sent = pool.Sent + winAmount

	err = s.dao.CreateAndUpdateRewardsState(pool, user, winAmount)
	if err != nil {
		return 0, fmt.Errorf("dao.CreateAndUpdateRewardsState: %s", err.Error())
	}

	return winAmount, nil
}

func (s *ServiceFacade) GetAvailableBoxesByUserID(userID uint64) (*models.Box, error) {
	box, err := s.dao.GetAvailableBoxesByUserID(userID)
	if err != nil {
		return nil, fmt.Errorf("dao.GetAvailableBoxesByUserID: %s", err.Error())
	}

	return box, nil
}
