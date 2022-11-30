package dao

import (
	"fmt"
	"strings"
	"time"

	"github.com/patrickmn/go-cache"
)

// BasicModelCache is used to store models only by id-keys
// and code-keys in memory cache
type BasicModelCache struct {
	ttl        time.Duration
	id2model   *cache.Cache
	code2model *cache.Cache
}

func NewCache(ttl time.Duration) *BasicModelCache {
	return &BasicModelCache{ttl: ttl, id2model: cache.New(ttl, ttl), code2model: cache.New(ttl, ttl)}
}

func (mc *BasicModelCache) Store(id uint64, code string, obj interface{}) {
	code = strings.ToLower(code)

	mc.id2model.Set(fmt.Sprintf("%d", id), obj, mc.ttl)
	mc.code2model.Set(code, obj, mc.ttl)
}

func (mc *BasicModelCache) GetByCode(code string) interface{} {
	code = strings.ToLower(code)

	model, _ := mc.code2model.Get(code)
	return model
}

func (mc *BasicModelCache) GetById(id uint64) interface{} {
	model, _ := mc.id2model.Get(fmt.Sprintf("%d", id))
	return model
}

func (mc *BasicModelCache) GetAll() interface{} {
	items := mc.id2model.Items()
	var res []interface{}
	for _, it := range items {
		res = append(res, it.Object)
	}
	return res
}
