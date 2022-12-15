package postgres

import (
	"errors"
	"gorm.io/gorm"
	"new-refferal/models"
)

func (db *Postgres) GetInviterByUserID(id uint64) (*models.Invitation, error) {
	info := new(models.Invitation)
	if err := db.db.First(info, "referral = ?", id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}

		return nil, err
	}

	return info, nil
}
