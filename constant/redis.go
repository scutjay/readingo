package constant

var (
	SupportedReadCommands = map[string]string{
		// Key
		"TTL": "*Key",
		// String
		"GET": "*Key",
		// Hash
		"HGET":    "*Key *Field",
		"HEXISTS": "*Key *Field",
		"HGETALL": "*Key",
		// List
		"LLEN": "*Key",
		// Set
		"SMEMBERS":  "*Key",
		"SISMEMBER": "*Key *Value",
		"SCARD":     "*Key",
		// ZSET
		"ZCARD":         "*Key",
		"ZCOUNT":        "*Key *Min *Max",
		"ZRANK":         "*Key *Value",
		"ZREVANK":       "*Key *Value",
		"ZSCORE":        "*Key *Value",
		"ZRANGE":        "*Key *Min *Max",
		"ZRANGEBYSCORE": "*Key *Min *Max",
	}

	SupportedWriteCommands = map[string]string{
		// Key
		"DEL": "*Key",
		// String
		"SET":    "*Key *Value",
		"GETSET": "*Key *Value",
		"APPEND": "*Key *Value",
		// Hash
		"HSET": "*Key *Field *Value",
		"HDEL": "*Key *Field",
		// List
		"LPUSH": "*Key *Value",
		"RPUSH": "*Key *Value",
		// Set
		"SADD": "*Key *Value",
		"SREM": "*Key *Value",
		// ZSet
		"ZADD":    "*Key *Score *Value",
		"ZINCRBY": "*Key *Score *Value",
		"ZREM":    "*Key *Value",
	}
)
