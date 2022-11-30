package dao

import (
	"fmt"
	"new-refferal/conf"
	"new-refferal/dao/cache"
	"new-refferal/dao/postgres"
	"new-refferal/models"
)

type (
	DAO interface {
		Postgres
		Cache
	}

	Postgres interface {
		CreateUser(user *models.User) (*models.User, error)
		CreateUserAndLink(user *models.User, code string) (*models.User, error)
		GetUserByWalletAddress(addr string) (*models.User, error)

		GetLastLink() (*models.Link, error)
	}

	Cache   interface{}
	daoImpl struct {
		*postgres.Postgres
		*cache.Cache
	}
)

func New(cfg conf.Config) (DAO, error) {
	pg, err := postgres.NewPostgres(cfg.Postgres)
	if err != nil {
		return nil, fmt.Errorf("postgres.NewPostgres: %s", err.Error())
	}
	return daoImpl{
		Postgres: pg,
		Cache:    cache.NewCache(pg),
	}, nil
}
