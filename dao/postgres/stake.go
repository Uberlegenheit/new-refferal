package postgres

import (
	"errors"
	"fmt"
	"gorm.io/gorm"
	"new-refferal/filters"
	"new-refferal/models"
	"time"
)

func (db *Postgres) SaveDelegationTx(stake *models.Stake) (*models.Stake, error) {
	result := db.db.Table(models.StakesTable).Create(stake)
	if result.Error != nil {
		return nil, result.Error
	}

	return stake, nil
}

func (db *Postgres) SaveDelegationTxAndCreateReward(stake *models.Stake) (*models.Stake, error) {
	err := db.db.Transaction(func(tx *gorm.DB) error {
		if err := tx.Table(models.StakesTable).
			Create(stake).Error; err != nil {
			return err
		}

		reward := new(models.Reward)
		if err := tx.First(reward, "user_id = ? AND type_id = 1", stake.UserID).Error; err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				reward = nil
			} else {
				return err
			}
		}

		if reward == nil {
			if err := tx.Create(&models.Reward{
				UserID:  stake.UserID,
				Status:  "get reward in wallet",
				TypeID:  1,
				Amount:  0,
				Hash:    "",
				Created: time.Now(),
			}).Error; err != nil {
				return err
			}
		}

		return nil
	})
	if err != nil {
		return nil, err
	}

	return stake, err
}

func (db *Postgres) SaveFailedDelegationTx(stake *models.Stake) (*models.Stake, error) {
	result := db.db.Table(models.FailedStakesTable).Create(stake)
	if result.Error != nil {
		return nil, result.Error
	}

	return stake, nil
}

func (db *Postgres) SaveDelegationTxAndAddBoxes(stake *models.Stake) (*models.Stake, error) {
	err := db.db.Transaction(func(tx *gorm.DB) error {
		if err := tx.Create(stake).Error; err != nil {
			return err
		}

		if err := tx.Model(&models.Box{}).
			Where("user_id = ?", stake.UserID).
			Update("available", gorm.Expr("available+?", stake.BoxesGiven)).
			Error; err != nil {
			return err
		}

		return nil
	})
	if err != nil {
		return nil, err
	}

	return stake, err
}

func (db *Postgres) SetUserDelegationsFalse(id uint64) error {
	result := db.db.Model(&models.Stake{}).Update("status", false).Where("user_id = ?", id)
	if result.Error != nil {
		return result.Error
	}

	return nil
}

func (db *Postgres) GetInvitedUsersStakes(id uint64, pagination filters.Pagination) ([]models.StakeShow, uint64, error) {
	pagination.Validate()
	stakes := make([]models.StakeShow, 0)

	if err := db.db.Model(&models.StakeShow{}).
		Select("s.id, s.user_id, s.amount, s.status, s.tx_hash, b.available+b.opened as boxes, s.created").
		Table(fmt.Sprintf("%s it", models.InvitationsTable)).
		Joins("inner join stakes s on it.referral_id = s.user_id").
		Joins("inner join boxes b on b.user_id = s.user_id").
		Where("it.referrer_id = ? AND s.status = ?", id, true).
		Order("s.created desc").
		Scan(&stakes).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, 0, nil
		}
		return nil, 0, err
	}

	length := uint64(len(stakes))
	offset := pagination.Offset()
	limit := pagination.Limit
	if offset > length {
		return nil, length, nil
	} else if limit > length {
		stakes = stakes[offset:length]
	} else {
		stakes = stakes[offset : offset+limit]
	}

	return stakes, length, nil
}

func (db *Postgres) GetDelegationByTxHash(stake *models.Stake) (*models.Stake, error) {
	dbStake := new(models.Stake)

	if err := db.db.Model(&models.Stake{}).
		Select("*").
		Table(models.StakesTable).
		Where("tx_hash = ?", stake.Hash).
		Scan(&dbStake).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}

	return dbStake, nil
}

func (db *Postgres) GetStakeAndBoxUserStatByID(id uint64) (*models.StakeAndBoxStat, error) {
	stats := new(models.StakeAndBoxStat)

	if err := db.db.Model(&models.StakeAndBoxStat{}).
		Select("s.user_id, coalesce(round(CAST(sum(s.amount) as numeric), 8), 0) as total_stake, b.total_boxes").
		Table(fmt.Sprintf("%s s", models.StakesTable)).
		Joins("join (select b.user_id, sum(b.available+b.opened) as total_boxes from boxes b group by b.user_id) b on b.user_id = s.user_id").
		Where("s.user_id = ? AND s.status = true", id).
		Group("s.user_id, b.total_boxes").
		Scan(&stats).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}

	return stats, nil
}

func (db *Postgres) GetFailedDelegations(pagination filters.Pagination) ([]models.FailedStakeShow, uint64, error) {
	pagination.Validate()
	stakes := make([]models.FailedStakeShow, 0)

	if err := db.db.Model(&models.FailedStakeShow{}).
		Select("fs.id, fs.user_id, u.wallet_name, u.wallet_address, fs.amount, fs.status, fs.type_id, fs.boxes_given, fs.tx_hash, fs.created").
		Table(fmt.Sprintf("%s fs", models.FailedStakesTable)).
		Joins("inner join stake_types st on st.id = fs.type_id").
		Joins("inner join users u on u.id = fs.user_id").
		Order("fs.created desc").
		Scan(&stakes).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, 0, nil
		}
		return nil, 0, err
	}

	length := uint64(len(stakes))
	offset := pagination.Offset()
	limit := pagination.Limit
	if offset > length {
		return nil, length, nil
	} else if limit > length {
		stakes = stakes[offset:length]
	} else {
		stakes = stakes[offset : offset+limit]
	}

	return stakes, length, nil
}

func (db *Postgres) SaveTXAndUpdateReward(info *models.StakeAndBoxStat, stake, reward float64) error {
	err := db.db.Transaction(func(tx *gorm.DB) error {
		if info.TotalStake != stake {
			if err := tx.Model(&models.Stake{}).
				Table(models.StakesTable).
				Where("user_id = ?", info.UserID).
				Update("status", false).Error; err != nil {
				return err
			}

			if err := tx.Create(&models.Stake{
				UserID:  info.UserID,
				Amount:  stake,
				Status:  true,
				TypeID:  3,
				Hash:    "updated delegation balance",
				Created: time.Now(),
			}).Error; err != nil {
				return err
			}
		}

		if err := db.db.Model(&models.Reward{}).
			Where("user_id = ? and type_id = 1", info.UserID).
			Updates(&models.Reward{
				UserID:  info.UserID,
				Status:  "get reward in wallet",
				TypeID:  1,
				Amount:  reward,
				Hash:    "updated rewards",
				Created: time.Now(),
			}).Error; err != nil {
			return err
		}

		return nil
	})
	if err != nil {
		return err
	}

	return nil
}
