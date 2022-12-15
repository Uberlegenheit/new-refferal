package postgres

import (
	"errors"
	"fmt"
	"gorm.io/gorm"
	"new-refferal/models"
)

func (db *Postgres) GetRewardsPool() (*models.RewardPool, error) {
	pool := new(models.RewardPool)
	err := db.db.Table(models.RewardsPoolTable).
		Last(pool).
		Order("id desc").Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, fmt.Errorf("no rewards pool")
		}
		return nil, err
	}

	return pool, nil
}

func (db *Postgres) UpdateRewardsPool(pool *models.RewardPool) error {
	result := db.db.Table(models.RewardsPoolTable).
		Model(&models.RewardPool{}).
		Where("id = ?", pool.ID).
		Updates(map[string]interface{}{
			"available": pool.Available,
			"sent":      pool.Sent,
		})
	if result.Error != nil {
		return result.Error
	}

	return nil
}

func (db *Postgres) SetDailyPoolLimit(pool *models.RewardPool) error {
	result := db.db.Table(models.RewardsPoolTable).
		Model(&models.RewardPool{}).
		Where("id = ?", pool.ID).
		Updates(map[string]interface{}{
			"daily_limit": pool.DailyLimit,
		})
	if result.Error != nil {
		return result.Error
	}

	return nil
}
