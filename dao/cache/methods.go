package cache

import "time"

func (c *Cache) AddAuthToken(key string, item interface{}, expiration time.Duration) (err error) {
	c.cache.Set(key, item, expiration)
	return nil
}

func (c *Cache) GetAuthToken(key string) (interface{}, bool, error) {
	item, ok := c.cache.Get(key)
	if !ok {
		return item, false, nil
	}

	return item, true, nil
}

func (c *Cache) RemoveAuthToken(key string) error {
	c.cache.Delete(key)
	return nil
}

func (c *Cache) CacheSave(key string, item interface{}, expiration time.Duration) (err error) {
	c.cache.Set(key, item, expiration)
	return nil
}

func (c *Cache) CacheGet(key string) (interface{}, bool, error) {
	item, ok := c.cache.Get(key)
	if !ok {
		return item, false, nil
	}

	return item, true, nil
}

func (c *Cache) CacheRemove(key string) error {
	c.cache.Delete(key)
	return nil
}
