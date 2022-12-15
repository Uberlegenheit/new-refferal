package services

import (
	"encoding/json"
	"fmt"
	"github.com/google/uuid"
	"github.com/shopspring/decimal"
	"math"
	"net/http"
	"net/url"
	"new-refferal/filters"
	"new-refferal/models"
	"os"
	"strconv"
	"time"
)

func roundFloat(val float64, precision uint) float64 {
	ratio := math.Pow(10, float64(precision))
	return math.Round(val*ratio) / ratio
}

func checkHashAndSum(stake *models.Stake, addr string) (bool, error) {
	u := url.URL{
		Scheme: "https",
		Host:   os.Getenv("COSMOS_API"),
		Path:   fmt.Sprintf(TxPath, os.Getenv("NODE_TOKEN"), stake.Hash),
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

	if time.Now().Unix()-tx.TxFetchTime.Timestamp.Unix() > 120 {
		return false, fmt.Errorf("old transaction")
	}

	if len(tx.Tx.Body.Body) != 0 {
		if tx.Tx.Body.Body[0].DelegatorAddr != addr {
			return false, fmt.Errorf("wrong delegator address")
		}

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

func (s *ServiceFacade) SaveDelegationTx(stake *models.Stake, user *models.User) (*models.Stake, error) {
	ok, err := checkHashAndSum(stake, user.WalletAddress)
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

	stb := os.Getenv("STAKE_TO_BOX")
	stakeToBox, err := strconv.ParseFloat(stb, 64)
	if err != nil {
		return nil, fmt.Errorf("ParseFloat: %s", err.Error())
	}
	boxesOnSum := int64(roundFloat((stats.TotalStake)/stakeToBox, 5))
	boxesAvailable := int64(roundFloat((stats.TotalStake+stake.Amount)/stakeToBox, 5))
	newBoxes := boxesAvailable - boxesOnSum
	if newBoxes > 0 {
		err := s.dao.AddBoxesByUserID(stake.UserID, newBoxes)
		if err != nil {
			return nil, fmt.Errorf("dao.AddBoxesByUserID: %s", err.Error())
		}
	}

	stake.BoxesGiven = uint64(newBoxes)
	stake, err = s.dao.SaveDelegationTxAndCreateReward(stake)
	if err != nil {
		return nil, fmt.Errorf("dao.SaveDelegationTx: %s", err.Error())
	}

	return stake, nil
}

func (s *ServiceFacade) SaveFailedDelegationTx(stake *models.Stake) (*models.Stake, error) {
	stake.Status = false
	stake, err := s.dao.SaveFailedDelegationTx(stake)
	if err != nil {
		return nil, fmt.Errorf("dao.SaveFailedDelegationTx: %s", err.Error())
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

func (s *ServiceFacade) GetFailedDelegations(pagination filters.Pagination) ([]models.FailedStakeShow, uint64, error) {
	delegations, length, err := s.dao.GetFailedDelegations(pagination)
	if err != nil {
		return nil, length, fmt.Errorf("dao.GetFailedDelegations: %s", err.Error())
	}

	return delegations, length, nil
}

func (s *ServiceFacade) GetInvitedUsersStakes(user *models.User, pagination filters.Pagination) ([]models.StakeShow, uint64, error) {
	stakes, length, err := s.dao.GetInvitedUsersStakes(user.ID, pagination)
	if err != nil {
		return nil, length, fmt.Errorf("dao.GetInvitedUsersStakes: %s", err.Error())
	}

	return stakes, length, nil
}

func (s *ServiceFacade) GetDelegationKey(user *models.User) (string, error) {
	key := uuid.New()

	err := s.dao.CacheSave(user.WalletAddress, key.String(), time.Minute*2)
	if err != nil {
		return "", fmt.Errorf("dao.CacheSave: %s", err.Error())
	}

	return key.String(), nil
}

func (s *ServiceFacade) CheckDelegationKey(user *models.User, key string) (bool, error) {
	defer s.dao.CacheRemove(user.WalletAddress)

	sKey, ok, err := s.dao.CacheGet(user.WalletAddress)
	if err != nil {
		return false, fmt.Errorf("dao.CacheGet: %s", err.Error())
	}

	if !ok {
		return false, fmt.Errorf("no key generated")
	}

	if key != sKey.(string) {
		return false, fmt.Errorf("key is invalid")
	}

	return true, nil
}
