package services

import (
	"encoding/json"
	"fmt"
	"github.com/shopspring/decimal"
	"net/http"
	"net/url"
	"new-refferal/models"
	"strconv"
)

func checkHashAndSum(stake *models.Stake) (bool, error) {
	u := url.URL{
		Scheme: "https",
		Host:   CosmosAPI,
		Path:   fmt.Sprintf(TxPath, NodeToken, stake.Hash),
	}

	resp, err := http.Get(u.String())
	if err != nil {
		return false, fmt.Errorf("http.Get: %s", err.Error())
	}
	defer resp.Body.Close()
	if resp.StatusCode != 200 {
		return false, nil
	}

	var tx models.TxFetch
	err = json.NewDecoder(resp.Body).Decode(&tx)
	if err != nil {
		return false, err
	}

	if len(tx.Tx.Body.Body) != 0 {
		stakeStr, err := strconv.ParseInt(tx.Tx.Body.Body[0].Amount.Amount, 10, 64)
		if err != nil {
			return false, err
		}

		stakeD := decimal.New(stakeStr, -6)
		stakeF, _ := stakeD.Float64()

		if stakeF != stake.Amount {
			return false, nil
		}
	} else {
		return false, fmt.Errorf("msg length is 0")
	}

	return true, nil
}

func (s *ServiceFacade) SaveDelegationTx(stake *models.Stake) (*models.Stake, error) {
	ok, err := checkHashAndSum(stake)
	if err != nil {
		return nil, fmt.Errorf("checkHashAndSum: %s", err.Error())
	}
	if !ok {
		return nil, fmt.Errorf("checkHashAndSum: transaction differs from imput data")
	}

	stats, err := s.dao.GetStakeAndBoxUserStatByID(stake.UserID)
	if err != nil {
		return nil, fmt.Errorf("dao.GetStakeAndBoxUserStatByID: %s", err.Error())
	}

	boxesAvailable := int64((stats.TotalStake + stake.Amount) / 0.001 /*10.0*/)
	newBoxes := boxesAvailable - stats.TotalBoxes
	if newBoxes != 0 {
		err := s.dao.AddBoxesByUserID(stake.UserID, newBoxes)
		if err != nil {
			return nil, fmt.Errorf("dao.AddBoxesByUserID: %s", err.Error())
		}
	}

	stake.BoxesGiven = uint64(newBoxes)
	stake, err = s.dao.SaveDelegationTx(stake)
	if err != nil {
		return nil, fmt.Errorf("dao.SaveDelegationTx: %s", err.Error())
	}

	return stake, nil
}

func (s *ServiceFacade) GetDelegationByTxHash(stake *models.Stake) (*models.Stake, error) {
	dbStake, err := s.dao.GetDelegationByTxHash(stake)
	if err != nil {
		return nil, fmt.Errorf("dao.GetDelegationByTxHash: %s", err.Error())
	}

	if dbStake != nil {
		if dbStake.ID != 0 {
			return dbStake, nil
		}
	}

	return nil, nil
}

func (s *ServiceFacade) GetInvitedUsersStakes(user *models.User) ([]models.StakeShow, error) {
	stakes, err := s.dao.GetInvitedUsersStakes(user.ID)
	if err != nil {
		return nil, fmt.Errorf("dao.GetInvitedUsersStakes: %s", err.Error())
	}

	return stakes, nil
}
