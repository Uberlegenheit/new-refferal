package services

import (
	"encoding/json"
	"fmt"
	"github.com/google/martian/log"
	"github.com/roylee0704/gron"
	"github.com/roylee0704/gron/xtime"
	"github.com/shopspring/decimal"
	"net/http"
	"net/url"
	"new-refferal/models"
	"os"
	"strconv"
	"strings"
	"time"
)

const (
	EverstakeCosmosAddr string = "cosmosvaloper1tflk30mq5vgqjdly92kkhhq3raev2hnz6eete3"
	RewardsPath         string = "%s/cosmos/distribution/v1beta1/delegators/%s/rewards/%s"
	StakePath           string = "%s/cosmos/staking/v1beta1/validators/%s/delegations/%s"
	TxPath              string = "%s/cosmos/tx/v1beta1/txs/%s"
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

	cron.AddFunc(gron.Every(1*xtime.Day).At("00:00"), func() {
		err := s.checkDailyPoolLimit()
		if err != nil {
			log.Errorf("setting daily pool limit: %s", err.Error())
			return
		}
	})

	_ = s.checkDailyPoolLimit()
}

func (s *ServiceFacade) parseDelegations() error {
	users, err := s.dao.GetAllUsers()
	if err != nil {
		return fmt.Errorf("dao.GetAllUsers: %s", err.Error())
	}

	for i := range users {
		u := url.URL{
			Scheme: "https",
			Host:   os.Getenv("COSMOS_API"),
			Path:   fmt.Sprintf(RewardsPath, os.Getenv("NODE_TOKEN"), users[i].WalletAddress, EverstakeCosmosAddr),
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
			Host:   os.Getenv("COSMOS_API"),
			Path:   fmt.Sprintf(StakePath, os.Getenv("NODE_TOKEN"), EverstakeCosmosAddr, users[i].WalletAddress),
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

		if resp.StatusCode == 400 || resp.StatusCode == 404 {
			info, err := s.dao.GetStakeAndBoxUserStatByID(users[i].ID)
			if err != nil {
				return fmt.Errorf("dao.GetStakeAndBoxUserStatByID: %s", err.Error())
			}

			err = s.dao.SaveTXAndUpdateReward(info, 0, 0)
			if err != nil {
				return fmt.Errorf("dao.SaveTXAndUpdateReward: %s", err.Error())
			}
			continue
		}

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

		rewardF, _ := rewardD.Float64()
		stakeF, _ := stakeD.Float64()

		info, err := s.dao.GetStakeAndBoxUserStatByID(users[i].ID)
		if err != nil {
			return fmt.Errorf("dao.GetStakeAndBoxUserStatByID: %s", err.Error())
		}

		err = s.dao.SaveTXAndUpdateReward(info, stakeF, rewardF)
		if err != nil {
			return fmt.Errorf("dao.SaveTXAndUpdateReward: %s", err.Error())
		}
	}

	return nil
}

func (s *ServiceFacade) checkDailyPoolLimit() error {
	pool, err := s.dao.GetRewardsPool()
	if err != nil {
		return fmt.Errorf("dao.GetRewardsPool: %s", err.Error())
	}

	date := strings.Split(time.Now().Format("2006-01-02"), "-")
	year, _ := strconv.Atoi(date[0])
	month, _ := strconv.Atoi(date[1])
	day, _ := strconv.Atoi(date[2])
	days := time.Date(year, time.Month(month), 0, 0, 0, 0, 0, time.UTC).Day()

	if day == days {
		days = 1
	} else {
		days = days - day
	}

	pool.DailyLimit = (pool.Available / float64(days)) * 0.9

	err = s.dao.SetDailyPoolLimit(pool)

	return nil
}
