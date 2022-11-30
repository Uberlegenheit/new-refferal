package postgres

import (
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
		Update("available", gorm.Expr("available-?", 1)).Error
}
