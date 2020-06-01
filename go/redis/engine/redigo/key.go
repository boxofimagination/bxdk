package redigo

import (
	"github.com/gomodule/redigo/redis"
)

// Expire set expiration time for a key
// `expiry` is in seconds
func (r *Redigo) Expire(key string, expiry int) (int, error) {
	return redis.Int(r.do("EXPIRE", key, expiry))
}

// TTL return remaining ttl of a key
// from: https://redis.io/commands/ttl
// The command returns -2 if the key does not exist.
// The command returns -1 if the key exists but has no associated expire.
func (r *Redigo) TTL(key string) (int, error) {
	return redis.Int(r.do("TTL", key))
}

// Exists check key existence
func (r *Redigo) Exists(key string) (bool, error) {
	return redis.Bool(r.do("EXISTS", key))
}

// Sort ordered the value lexicographically
// Temporarily unexported, cannot test using miniredis
/*func (r *Redigo) sort(key string, alpha bool, asc bool) ([]string, error) {
	conn, err := r.getConn()
	if err != nil {
		return nil, err
	}
	defer conn.Close()

	args := make([]interface{}, 3)
	args[0] = key
	if alpha {
		args[1] = "ALPHA"
	}
	if asc {
		args[2] = "ASC"
	} else {
		args[2] = "DESC"
	}
	return redis.Strings(conn.Do("SORT", args...))
}*/

// Delete function
func (r *Redigo) Delete(keys ...string) (int, error) {
	// copy string to array of interface
	args := make([]interface{}, len(keys))
	for i, key := range keys {
		args[i] = key
	}
	return redis.Int(r.do("DEL", args...))
}
