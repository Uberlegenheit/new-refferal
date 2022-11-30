package postgres

import (
	"gorm.io/gorm"
	"new-refferal/models"
	"time"
)

func (db *Postgres) CreateUser(user *models.User) (*models.User, error) {
	result := db.db.Create(user)
	if result.Error != nil {
		return nil, result.Error
	}

	return user, nil
}

func (db *Postgres) CreateUserAndLink(user *models.User, code string) (*models.User, error) {
	user.Created = time.Now()

	err := db.db.Transaction(func(tx *gorm.DB) error {
		if err := tx.Create(user).Error; err != nil {
			return err
		}

		if err := tx.Create(&models.Link{
			UserID: user.ID,
			Code:   code,
		}).Error; err != nil {
			return err
		}

		if err := tx.Create(&models.Box{
			UserID:    user.ID,
			Available: 0,
			Opened:    0,
		}).Error; err != nil {
			return err
		}

		return nil
	})
	if err != nil {
		return nil, err
	}

	return user, err
}

func (db *Postgres) GetUserByWalletAddress(addr string) (*models.User, error) {
	user := new(models.User)
	result := db.db.First(user, "wallet_address = ?", addr)
	if result.Error != nil {
		return nil, result.Error
	}

	return user, nil
}
