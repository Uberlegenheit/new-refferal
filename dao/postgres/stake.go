package postgres

import (
	"errors"
	"fmt"
	"gorm.io/gorm"
	"new-refferal/models"
)

func (db *Postgres) SaveDelegationTx(stake *models.Stake) (*models.Stake, error) {
	result := db.db.Create(stake)
	if result.Error != nil {
		return nil, result.Error
	}

	return stake, nil
}

func (db *Postgres) GetInvitedUsersStakes(id uint64) ([]models.StakeShow, error) {
	stakes := make([]models.StakeShow, 0)

	if err := db.db.Model(&models.StakeShow{}).
		Select("s.id, s.user_id, s.amount, s.status, s.tx_hash, b.available+b.opened as boxes, s.created").
		Table(fmt.Sprintf("%s it", models.InvitationsTable)).
		Joins("inner join stakes s on it.referral_id = s.user_id").
		Joins("inner join boxes b on b.user_id = s.user_id").
		Where("it.referrer_id = ?", id).
		Order("s.created desc").
		Scan(&stakes).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}

	return stakes, nil
}

func (db *Postgres) GetStakeAndBoxUserStatByID(id uint64) (*models.StakeAndBoxStat, error) {
	stats := new(models.StakeAndBoxStat)

	if err := db.db.Model(&models.StakeAndBoxStat{}).
		Select("s.user_id, sum(s.amount) as total_stake, b.total_boxes").
		Table(fmt.Sprintf("%s s", models.StakesTable)).
		Joins("inner join (select b.user_id, sum(b.available+b.opened) as total_boxes from boxes b group by b.user_id) b on b.user_id = s.user_id").
		Where("s.user_id = ?", id).
		Group("s.user_id, b.total_boxes").
		Scan(&stats).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}

	return stats, nil
}
