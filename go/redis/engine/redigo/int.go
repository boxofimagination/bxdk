package redigo

import "github.com/gomodule/redigo/redis"

// Incr function
func (r *Redigo) Incr(key string) (int64, error) {
	return redis.Int64(r.do("INCR", key))
}

// IncrBy function
func (r *Redigo) IncrBy(key string, value int64) (int64, error) {
	return redis.Int64(r.do("INCRBY", key, value))
}

// Decr function
func (r *Redigo) Decr(key string) (int64, error) {
	return redis.Int64(r.do("DECR", key))
}

// DecrBy function
func (r *Redigo) DecrBy(key string, value int64) (int64, error) {
	return redis.Int64(r.do("DECRBY", key, value))
}
