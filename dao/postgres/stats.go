package postgres

import (
	"errors"
	"fmt"
	"gorm.io/gorm"
	"new-refferal/filters"
	"new-refferal/models"
)

func (db *Postgres) GetTotalStats(req filters.PeriodInfoRequest) (*models.TotalStats, error) {
	stats := new(models.TotalStats)

	if err := db.db.Raw(`select (select coalesce(round(CAST(sum(s.amount) as numeric), 8), 0)
						from stakes s
						where s.status = true and s.type_id = 1 and (s.created >= ? and s.created <= ?)) as stake_sum,
						(select coalesce(round(CAST(sum(s.amount) as numeric), 8), 0)
						from stakes s
						where s.status = true and s.type_id = 2 and (s.created >= ? and s.created <= ?)) as redelegation_sum,
						(select coalesce(round(CAST(sum(s.amount) as numeric), 8), 0)
						 from invitations i
						 inner join stakes s on s.user_id = i.referral_id
						 where s.status = true and (s.created >= ? and s.created <= ?)) as invited_sum,
						(select coalesce(sum(s.boxes_given), 0)
						 from stakes s
						 where s.status = true and (s.created >= ? and s.created <= ?)) as boxes_given,
						(select coalesce(count(r.id), 0)
						 from rewards r
						 where r.type_id = 2 and (r.created >= ? and r.created <= ?)) as boxes_opened,
						(select coalesce(round(CAST(sum(r.amount) as numeric), 8), 0)
						 from rewards r
						 where r.type_id = 2 and status = 'paid' and (r.created >= ? and r.created <= ?)) as boxes_rewards,
						(select coalesce(round(CAST(sum(r.amount) as numeric), 8), 0)
						 from rewards r
						 where r.type_id = 2 and status = 'pending' and (r.created >= ? and r.created <= ?)) as boxes_unpaid`,
		req.Start, req.End,
		req.Start, req.End,
		req.Start, req.End,
		req.Start, req.End,
		req.Start, req.End,
		req.Start, req.End,
		req.Start, req.End).
		Scan(&stats).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}

	return stats, nil
}

func (db *Postgres) GetMyStakeSum(id uint64) (*models.StakeAndProgress, error) {
	stake := new(models.StakeAndProgress)

	if err := db.db.Model(&models.StakeAndProgress{}).
		Select("s.user_id, sum(s.amount) as total_stake").
		Table(fmt.Sprintf("%s s", models.StakesTable)).
		Where("s.user_id = ?", id).
		Group("s.user_id").
		Scan(&stake).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}

	return stake, nil
}

func (db *Postgres) GetTotalStakeStats(req filters.PeriodInfoRequest, pagination filters.Pagination) ([]models.TotalStakeStats, uint64, error) {
	pagination.Validate()
	stats := make([]models.TotalStakeStats, 0)

	if err := db.db.Limit(int(pagination.Limit)).
		Offset(int(pagination.Offset())).
		Model(&models.TotalStakeStats{}).
		Select(`select u.id as user_id,
							   u.wallet_name,
							   u.wallet_address,
							   st.id as type_id,
							   st.name as type,
							   s.amount,
							   s.boxes_given as boxes,
							   s.created,
							   s.tx_hash
						from users u
						inner join stakes s on u.id = s.user_id
						inner join stake_types st on s.type_id = st.id
						inner join boxes b on u.id = b.user_id`).
		Where("s.status = true AND (s.created >= ? AND s.created <= ?)", req.Start, req.End).
		Scan(&stats).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, 0, nil
		}
		return nil, 0, err
	}

	length := uint64(len(stats))
	offset := pagination.Offset()
	limit := pagination.Limit
	if offset > length {
		return nil, length, nil
	} else if limit > length || offset+limit > length {
		stats = stats[offset:length]
	} else {
		stats = stats[offset : offset+limit]
	}

	return stats, length, nil
}

func (db *Postgres) GetBoxesStats(req filters.PeriodInfoRequest, pagination filters.Pagination) ([]models.BoxStats, uint64, error) {
	pagination.Validate()
	stats := make([]models.BoxStats, 0)

	if err := db.db.Limit(int(pagination.Limit)).
		Offset(int(pagination.Offset())).
		Model(&models.TotalStakeStats{}).
		Select(`select u.id as user_id,
							   u.wallet_name,
							   u.wallet_address,
							   r.amount,
							   r.status,
							   r.created,
							   r.tx_hash
						from users u
						inner join rewards r on u.id = r.user_id`).
		Where("r.type_id = 2 AND (r.created >= ? AND r.created <= ?)", req.Start, req.End).
		Scan(&stats).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, 0, nil
		}
		return nil, 0, err
	}

	length := uint64(len(stats))
	offset := pagination.Offset()
	limit := pagination.Limit
	if offset > length {
		return nil, length, nil
	} else if limit > length || offset+limit > length {
		stats = stats[offset:length]
	} else {
		stats = stats[offset : offset+limit]
	}

	return stats, length, nil
}

func (db *Postgres) GetFriendsStakeStats(req filters.PeriodInfoRequest, pagination filters.Pagination) ([]models.FriendStakeStats, uint64, error) {
	pagination.Validate()
	stats := make([]models.FriendStakeStats, 0)

	if err := db.db.Limit(int(pagination.Limit)).
		Offset(int(pagination.Offset())).
		Model(&models.FriendStakeStats{}).
		Select(`select i.referrer_id,
							   i.referral_id,
							   u.wallet_name,
							   u.wallet_address,
							   st.id as type_id,
							   st.name as type,
							   s.amount,
							   s.boxes_given as boxes,
							   s.created,
							   s.tx_hash
						from invitations i
						inner join users u on i.referral_id = u.id
						inner join stakes s on u.id = s.user_id
						inner join stake_types st on s.type_id = st.id
						inner join boxes b on u.id = b.user_id`).
		Where("s.status = true AND (s.created >= ? AND s.created <= ?)", req.Start, req.End).
		Scan(&stats).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, 0, nil
		}
		return nil, 0, err
	}

	length := uint64(len(stats))
	offset := pagination.Offset()
	limit := pagination.Limit
	if offset > length {
		return nil, length, nil
	} else if limit > length || offset+limit > length {
		stats = stats[offset:length]
	} else {
		stats = stats[offset : offset+limit]
	}

	return stats, length, nil
}
