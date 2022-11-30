package cache

import (
	"time"

	"github.com/patrickmn/go-cache"
	"new-refferal/dao/postgres"
	"new-refferal/helpers/dao"
)

const (
	defaultExpiration      = time.Minute * 5
	defaultCleanUpDuration = time.Minute * 8
)

type Cache struct {
	cache         *cache.Cache
	currencyCache *dao.BasicModelCache
	db            *postgres.Postgres
}

func NewCache(db *postgres.Postgres) *Cache {
	return &Cache{
		cache:         cache.New(defaultExpiration, defaultCleanUpDuration),
		currencyCache: dao.NewCache(defaultExpiration),
		db:            db,
	}
}
