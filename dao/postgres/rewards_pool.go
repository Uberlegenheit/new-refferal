package postgres

import (
	"errors"
	"fmt"
	"gorm.io/gorm"
	"new-refferal/models"
)

func (db *Postgres) GetRewardsPool() (*models.RewardsPool, error) {
	pool := new(models.RewardsPool)
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

func (db *Postgres) UpdateRewardsPool(pool *models.RewardsPool) error {
	result := db.db.Table(models.RewardsPoolTable).
		Model(&models.RewardsPool{}).
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
