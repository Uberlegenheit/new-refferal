package services

import (
	"encoding/json"
	"fmt"
	"github.com/google/martian/log"
	"github.com/roylee0704/gron"
	"github.com/shopspring/decimal"
	"net/http"
	"net/url"
	"new-refferal/models"
	"strconv"
	"time"
)

const (
	CosmosAPI           string = "api.cosmos.network"
	RewardsPath         string = "cosmos/distribution/v1beta1/delegators/%s/rewards/%s"
	EverstakeCosmosAddr string = "cosmosvaloper1tflk30mq5vgqjdly92kkhhq3raev2hnz6eete3"
	StakePath           string = "cosmos/staking/v1beta1/validators/%s/delegations/%s"
)

func (s *ServiceFacade) InitCron(cron *gron.Cron) {
	dur := time.Hour * 6
	log.Infof("Scheduled delegations parse every %s", dur)
	cron.AddFunc(gron.Every(dur), func() {
		err := s.parseDelegations()
		if err != nil {
			log.Errorf("delegations parsing failed: %s", err.Error())
			return
		}
	})
}

func (s *ServiceFacade) parseDelegations() error {
	users, err := s.dao.GetAllUsers()
	if err != nil {
		return fmt.Errorf("dao.GetAllUsers: %s", err.Error())
	}

	for i := range users {
		u := url.URL{
			Scheme: "https",
			Host:   CosmosAPI,
			Path:   fmt.Sprintf(RewardsPath, users[i].WalletAddress, EverstakeCosmosAddr),
		}

		resp, err := http.Get(u.String())
		if err != nil {
			return fmt.Errorf("http.Get 1: %s", err.Error())
		}

		var sar models.StakeAndReward
		err = json.NewDecoder(resp.Body).Decode(&sar)
		if err != nil {
			return err
		}
		resp.Body.Close()

		u = url.URL{
			Scheme: "https",
			Host:   CosmosAPI,
			Path:   fmt.Sprintf(StakePath, EverstakeCosmosAddr, users[i].WalletAddress),
		}

		resp, err = http.Get(u.String())
		if err != nil {
			return fmt.Errorf("http.Get 2: %s", err.Error())
		}

		err = json.NewDecoder(resp.Body).Decode(&sar)
		if err != nil {
			return err
		}
		resp.Body.Close()

		stake, err := strconv.ParseInt(sar.Stake.Balance.Amount, 10, 64)
		if err != nil {
			return err
		}
		reward, err := strconv.ParseFloat(sar.Rewards[0].Amount, 64)
		if err != nil {
			return err
		}

		rewardD := decimal.NewFromFloat(reward).Div(decimal.New(1, 6))
		stakeD := decimal.New(stake, -6)

		rewardF, ok := rewardD.Float64()
		if !ok {
			return fmt.Errorf("rewardD.Float64()")
		}

		stakeF, ok := stakeD.Float64()
		if !ok {
			return fmt.Errorf("stakeD.Float64()")
		}

		info, err := s.dao.GetStakeAndBoxUserStatByID(users[i].ID)
		if err != nil {
			return fmt.Errorf("dao.GetStakeAndBoxUserStatByID: %s", err.Error())
		}

		if info.TotalStake != stakeF {
			err := s.dao.SetUserDelegationsFalse(users[i].ID)
			if err != nil {
				return fmt.Errorf("dao.SetUserDelegationsFalse: %s", err.Error())
			}
			_, err = s.dao.SaveDelegationTx(&models.Stake{
				UserID:  users[i].ID,
				Amount:  stakeF,
				Status:  true,
				Hash:    "updated delegation balance",
				Created: time.Now(),
			})
			if err != nil {
				return fmt.Errorf("dao.SaveDelegationTx: %s", err.Error())
			}
		}

		err = s.dao.UpdateReward(&models.Reward{
			UserID:  users[i].ID,
			Status:  "updated",
			TypeID:  1,
			Amount:  rewardF,
			Hash:    "updated rewards",
			Created: time.Now(),
		})
		if err != nil {
			return fmt.Errorf("dao.UpdateReward: %s", err.Error())
		}
	}

	return nil
}
