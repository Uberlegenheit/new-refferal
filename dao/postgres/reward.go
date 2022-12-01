package postgres

import (
	"errors"
	"fmt"
	"gorm.io/gorm"
	"new-refferal/models"
)

func (db *Postgres) SaveReward(reward *models.Reward) (*models.Reward, error) {
	result := db.db.Create(reward)
	if result.Error != nil {
		return nil, result.Error
	}

	return reward, nil
}

func (db *Postgres) UpdateReward(reward *models.Reward) error {
	changes := make(map[string]interface{})

	if reward.Hash != "" {
		changes["tx_hash"] = reward.Hash
	}
	if reward.Status != "" {
		changes["status"] = reward.Status
	}
	if reward.Amount != 0 {
		changes["amount"] = reward.Amount
	}

	result := db.db.Model(&models.Reward{}).Updates(changes)
	if result.Error != nil {
		return result.Error
	}

	return nil
}

func (db *Postgres) GetUserRewardsByID(id uint64) ([]models.RewardShow, error) {
	rewards := make([]models.RewardShow, 0)

	if err := db.db.Model(&models.RewardShow{}).
		Select(`r.id, u.wallet_name, r.status, rt."name" as type, r.amount, r.tx_hash, r.created`).
		Table(fmt.Sprintf("%s r", models.RewardsTable)).
		Joins("inner join reward_types rt on rt.id = r.type_id").
		Joins("inner join users u on u.id = r.user_id").
		Where("r.user_id = ?", id).
		Order("r.created desc").
		Scan(&rewards).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}

	return rewards, nil
}

func (db *Postgres) GetAllRewards() ([]models.RewardShow, error) {
	rewards := make([]models.RewardShow, 0)

	if err := db.db.Model(&models.RewardShow{}).
		Select(`r.id, u.wallet_name, r.status, rt."name" as type, r.amount, r.tx_hash, r.created`).
		Table(fmt.Sprintf("%s r", models.RewardsTable)).
		Joins("inner join reward_types rt on rt.id = r.type_id").
		Joins("inner join users u on u.id = r.user_id").
		Order("r.created desc").
		Scan(&rewards).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}

	return rewards, nil
}
