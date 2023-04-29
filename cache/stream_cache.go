package cache

type StreamCache struct {
	cache map[string]string
}

func NewStreamCache() *StreamCache {
	return &StreamCache{
		cache: make(map[string]string),
	}
}

func (sc *StreamCache) SetResolvedStreamURL(liquipediaUrl string, resolvedUrl string) {
	sc.cache[liquipediaUrl] = resolvedUrl
}

func (sc *StreamCache) GetResolvedStreamURL(liquipediaUrl string) (string, bool) {
	value, ok := sc.cache[liquipediaUrl]
	return value, ok
}
