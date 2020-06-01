package redigo

import (
	"errors"

	"github.com/gomodule/redigo/redis"
)

// Scan function return keys that match the pattern
func (r *Redigo) Scan(pattern string, cursor uint64, count int64) ([]string, uint64, error) {
	var (
		result       []string
		err          error
		rawResult    []interface{}
		newCursor    uint64
		rawFoundKeys []string
	)

	rawResult, err = redis.Values(r.do("SCAN", cursor, "MATCH", pattern, "COUNT", count))
	if err != nil {
		return result, newCursor, err
	}
	// raw result from redis will give us two index array, 1st is new cursor, 2nd is found keys
	if len(rawResult) < 2 {
		return result, newCursor, errors.New("fail to scan")
	}
	newCursor, _ = redis.Uint64(rawResult[0], nil)     // err ignored since redis.Uint64() function already handle it and will return 0
	rawFoundKeys, _ = redis.Strings(rawResult[1], nil) // err ignored since redis.Strings() function already handle it and will return ""

	result = append(result, rawFoundKeys...)

	return result, newCursor, err
}
