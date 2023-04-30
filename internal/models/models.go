package models

type RedisJSON struct {
	Key   string `json:"key"`
	Value int64  `json:"value"`
}

type CryptPair struct {
	Text string `json:"text"`
	Key  string `json:"key"`
}
