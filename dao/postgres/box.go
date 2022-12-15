package postgres

import (
	"errors"
	"fmt"
	"gorm.io/gorm"
	"new-refferal/models"
)

func (db *Postgres) AddBoxesByUserID(userID uint64, newBoxes int64) error {
	err := db.db.Transaction(func(tx *gorm.DB) error {
		err := tx.Model(&models.Box{}).
			Where("user_id = ?", userID).
			Update("available", gorm.Expr("available+?", newBoxes)).Error
		if err != nil {
			return err
		}

		info := new(models.Invitation)
		if err := db.db.First(info, "referral_id = ?", userID).Error; err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				info = nil
			} else {
				return err
			}
		}

		if info != nil {
			err = tx.Model(&models.Box{}).
				Where("user_id = ?", info.ReferrerID).
				Update("available", gorm.Expr("available+?", newBoxes)).Error
			if err != nil {
				return err
			}
		}

		return nil
	})
	if err != nil {
		return err
	}

	return nil
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
