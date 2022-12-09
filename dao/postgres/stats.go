package postgres

import (
	"errors"
	"gorm.io/gorm"
	"new-refferal/filters"
	"new-refferal/models"
)

func (db *Postgres) GetTotalStats(req filters.PeriodInfoRequest) ([]models.TotalStats, error) {
	stats := make([]models.TotalStats, 0)

	if err := db.db.Raw(`select (select coalesce(round(CAST(sum(s.amount) as numeric), 8), 0)
						from stakes s
						where s.status = true and s.type_id = 1) as stake_sum,
						(select coalesce(round(CAST(sum(s.amount) as numeric), 8), 0)
						from stakes s
						where s.status = true and s.type_id = 2) as redelegation_sum,
						(select coalesce(round(CAST(sum(s.amount) as numeric), 8), 0)
						 from invitations i
						  inner join stakes s on s.user_id = i.referral_id
						 where s.status = true) as invited_sum,
						(select coalesce(sum(b.opened+b.available), 0)
						 from boxes b) as boxes_given,
						(select coalesce(sum(b.opened), 0)
						 from boxes b) as boxes_opened,
						(select coalesce(round(CAST(sum(r.amount) as numeric), 8), 0)
						 from rewards r
						 where r.type_id = 2 and status = 'paid') as boxes_rewards,
						(select coalesce(round(CAST(sum(r.amount) as numeric), 8), 0)
						 from rewards r
						 where r.type_id = 2 and status = 'pending') as boxes_unpaid`).
		Scan(&stats).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}

	return stats, nil
}

func (db *Postgres) GetTotalStakeStats(req filters.PeriodInfoRequest) ([]models.TotalStakeStats, error) {
	stats := make([]models.TotalStakeStats, 0)

	if err := db.db.Model(&models.TotalStakeStats{}).
		Select(`select u.id as user_id,
							   u.wallet_name,
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
		Where("s.created >= ? AND s.created <= ?", req.Start, req.End).
		Scan(&stats).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}

	return stats, nil
}

func (db *Postgres) GetFriendsStakeStats(req filters.PeriodInfoRequest) ([]models.FriendStakeStats, error) {
	stats := make([]models.FriendStakeStats, 0)

	if err := db.db.Model(&models.FriendStakeStats{}).
		Select(`select i.referrer_id,
							   i.referral_id,
							   u.wallet_name,
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
		Where("s.created >= ? AND s.created <= ?", req.Start, req.End).
		Scan(&stats).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}

	return stats, nil
}

func (db *Postgres) GetRewardPaymentStats(req filters.PeriodInfoRequest) ([]models.RewardPaymentsStats, error) {
	stats := make([]models.RewardPaymentsStats, 0)

	if err := db.db.Model(&models.RewardPaymentsStats{}).
		Select(`select r.id,
							   r.user_id,
							   u.wallet_name,
							   r.status,
							   r.amount,
							   r.tx_hash,
							   r.created
						from rewards r
						inner join reward_types rt on rt.id = r.type_id
						inner join users u on r.user_id = u.id`).
		Where("r.type_id = 2 AND (r.created >= ? AND r.created <= ?)", req.Start, req.End).
		Scan(&stats).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}

	return stats, nil
}
