package postgres

import (
	"errors"
	"fmt"
	"gorm.io/gorm"
	"new-refferal/models"
)

func (db *Postgres) CreateUser(user *models.User) (*models.User, error) {
	result := db.db.Create(user)
	if result.Error != nil {
		return nil, result.Error
	}

	return user, nil
}

func (db *Postgres) CreateUserAndLink(user *models.User, code string) (*models.User, error) {
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

func (db *Postgres) GetAllUsers() ([]models.User, error) {
	users := make([]models.User, 0)

	if err := db.db.Model(&models.User{}).
		Select("distinct u.id, u.wallet_name, u.wallet_address, u.created, u.role").
		Table(fmt.Sprintf("%s u", models.UsersTable)).
		Joins("inner join stakes s on u.id = s.user_id").
		Scan(&users).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}

	return users, nil
}
