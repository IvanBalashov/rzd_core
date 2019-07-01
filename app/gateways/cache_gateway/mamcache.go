package cache_gateway

import (
	"errors"
	"fmt"
	"github.com/bradfitz/gomemcache/memcache"
)

type Memcache struct {
	CLI memcache.Client
	Ttl int32
}

//TODO: rewrire errors

func NewMemcache(client memcache.Client, ttl int32) *Memcache {
	return &Memcache{CLI: client, Ttl: ttl}
}

func (m *Memcache) Get(key string) ([]byte, error) {
	data, err := m.CLI.Get(key)
	if err != nil {
		return nil, errors.New(fmt.Sprintf("Gateways->Cache_Gateway->Get: Error in get data by key: %s, err: %s", key, err))
	}

	return data.Value, nil
}

func (m *Memcache) Set(key string, data []byte) error {
	err := m.CLI.Set(&memcache.Item{Key: key, Value: data, Expiration: m.Ttl})
	if err != nil {
		return errors.New(fmt.Sprintf("Gateways->Cache_Gateway->Get: Error in set data by key: %s, err: %s", key, err))
	}

	return nil
}

func (m *Memcache) Del(key string) error {
	err := m.CLI.Delete(key)
	if err != nil {
		return errors.New(fmt.Sprintf("Gateways->Cache_Gateway->Get: Error in delete data by key: %s, err: %s", key, err))
	}

	return nil
}

func (m *Memcache) GetMulti(keys []string) (map[string][]byte, error) {
	answer := map[string][]byte{}

	data, err := m.CLI.GetMulti(keys)
	if err != nil {
		return nil, errors.New(fmt.Sprintf("Gateways->Cache_Gateway->Get: Error in get data by keys: %v, err: %s", keys, err))
	}

	for key, val := range data {
		answer[key] = val.Value
	}

	return answer, nil
}
