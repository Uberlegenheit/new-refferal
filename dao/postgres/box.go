package postgres

import (
	"errors"
	"fmt"
	"gorm.io/gorm"
	"new-refferal/models"
)

func (db *Postgres) AddBoxesByUserID(userID uint64, newBoxes int64) error {
	return db.db.Model(&models.Box{}).
		Where("user_id = ?", userID).
		Update("available", gorm.Expr("available+?", newBoxes)).Error
}

func (db *Postgres) OpenBoxByUserID(userID uint64) error {
	return db.db.Model(&models.Box{}).
		Where("user_id = ?", userID).
		Update("available", gorm.Expr("available-?", 1)).
		Update("opened", gorm.Expr("opened+?", 1)).Error
}

func (db *Postgres) GetAvailableBoxesByUserID(userID uint64) (*models.Box, error) {
	box := new(models.Box)
	err := db.db.First(box, "user_id = ?", userID).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, fmt.Errorf("no rows in db")
		}

		return nil, err
	}

	return box, nil
}
