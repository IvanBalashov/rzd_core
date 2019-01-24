package cache_gateway

type CacheGateway interface {
	Get(key string) ([]byte, error)
	Set(key string, data []byte) error
	Del(key string) error
	GetMulti(keys []string) (map[string][]byte, error)
}
