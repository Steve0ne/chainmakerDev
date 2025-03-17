package mock

type MemCache map[string][]byte

func NewMemCache() MemCache {
	m := make(MemCache)
	return m
}

func (m *MemCache) Put(key string, value []byte) {
	(*m)[key] = value
}

func (m *MemCache) Get(key string) []byte {
	return (*m)[key]
}
