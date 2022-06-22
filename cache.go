package cache

import (
	"time"
)

type Value struct {
	key, value string
	deadline   *time.Time
}

func NewValue(key, value string, deadline *time.Time) Value {
	return Value{key: key, value: value, deadline: deadline}
}

func (v Value) IsExpired() bool {
	if v.deadline != nil && v.deadline.Before(time.Now()) {
		return true
	}

	return false
}

type Cache struct {
	values []Value
}

func NewCache() Cache {
	return Cache{}
}

func (c *Cache) remove(key string) {
	var values []Value

	for _, v := range c.values {
		if v.key != key {
			values = append(values, v)
		}
	}

	c.values = values
}

func (c *Cache) Get(key string) (string, bool) {
	c.Cleanup()

	for _, v := range c.values {
		if v.key == key {
			return v.value, true
		}
	}

	return "", false
}

func (c *Cache) Set(key, value string, deadline *time.Time) {
	for _, v := range c.values {
		if v.key == key {
			v.value = value
			v.deadline = deadline
			return
		}
	}

	c.values = append(c.values, NewValue(key, value, deadline))
}

func (c *Cache) Put(key, value string) {
	c.Set(key, value, nil)
}

func (c *Cache) PutTill(key, value string, deadline time.Time) {
	c.Set(key, value, &deadline)
}

func (c *Cache) Keys() []string {
	c.Cleanup()

	var keys []string

	for _, v := range c.values {
		keys = append(keys, v.key)
	}

	return keys
}

func (c *Cache) Cleanup() {
	for _, v := range c.values {
		if v.IsExpired() {
			c.remove(v.key)
		}
	}
}
