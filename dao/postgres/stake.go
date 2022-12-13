package postgres

import (
	"errors"
	"fmt"
	"gorm.io/gorm"
	"new-refferal/filters"
	"new-refferal/models"
)

func (db *Postgres) SaveDelegationTx(stake *models.Stake) (*models.Stake, error) {
	result := db.db.Create(stake)
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

func (db *Postgres) GetInvitedUsersStakes(id uint64, pagination filters.Pagination) ([]models.StakeShow, error) {
	pagination.Validate()
	stakes := make([]models.StakeShow, 0)

	if err := db.db.Limit(int(pagination.Limit)).
		Offset(int(pagination.Offset())).
		Model(&models.StakeShow{}).
		Select("s.id, s.user_id, s.amount, s.status, s.tx_hash, b.available+b.opened as boxes, s.created").
		Table(fmt.Sprintf("%s it", models.InvitationsTable)).
		Joins("inner join stakes s on it.referral_id = s.user_id").
		Joins("inner join boxes b on b.user_id = s.user_id").
		Where("it.referrer_id = ? AND s.status = ?", id, true).
		Order("s.created desc").
		Scan(&stakes).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}

	return stakes, nil
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
