package postgres

import (
	"errors"
	"fmt"
	"gorm.io/gorm"
	"new-refferal/filters"
	"new-refferal/models"
)

func (db *Postgres) CreatePayout(payout *models.Payout) (*models.Payout, error) {
	result := db.db.Create(payout)
	if result.Error != nil {
		return nil, result.Error
	}

	return payout, nil
}

func (db *Postgres) UpdatePayout(payout *models.Payout) error {
	changes := make(map[string]interface{})

	if payout.UserID != 0 {
		changes["user_id"] = payout.UserID
	}
	if payout.Amount != 0 {
		changes["amount"] = payout.Amount
	}
	if payout.Fee != 0 {
		changes["fee"] = payout.Fee
	}
	if payout.TxHash != "" {
		changes["tx_hash"] = payout.TxHash
	}

	result := db.db.Model(&models.Payout{}).
		Updates(changes).
		Where("id = ?", payout.ID)
	if result.Error != nil {
		return result.Error
	}

	return nil
}

func (db *Postgres) GetPayouts(pagination filters.Pagination) ([]models.PayoutShow, uint64, error) {
	pagination.Validate()
	payouts := make([]models.PayoutShow, 0)

	if err := db.db.Limit(int(pagination.Limit)).
		Offset(int(pagination.Offset())).
		Model(&models.Payout{}).
		Select(`p.id, p.user_id, u.wallet_name, u.wallet_address, p.amount, p.fee, p.tx_hash, p.created`).
		Table(fmt.Sprintf("%s p", models.PayoutsTable)).
		Joins("inner join users u on u.id = p.user_id").
		Order("p.created desc").
		Scan(&payouts).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, 0, nil
		}
		return nil, 0, err
	}

	length := uint64(len(payouts))
	offset := pagination.Offset()
	limit := pagination.Limit
	if offset > length {
		return nil, length, nil
	} else if limit > length {
		payouts = payouts[offset:length]
	} else {
		payouts = payouts[offset : offset+limit]
	}

	return payouts, length, nil
}
