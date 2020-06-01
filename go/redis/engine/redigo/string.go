package redigo

import (
	"fmt"

	"github.com/gomodule/redigo/redis"

	"github.com/boxofimagination/bxdk/go/redis/engine"
)

// Set key and value
func (r *Redigo) Set(key string, value interface{}) error {
	ok, err := redis.String(r.do("SET", key, value))
	if ok != "OK" && err == nil {
		return engine.ErrNotOK
	}
	return err
}

// SetNX do SETNX (only set if not exist) with SET's NX & EX args.
// It sets the key which will expired in `expire` seconds
func (r *Redigo) SetNX(key string, value interface{}, expire int) (string, error) {
	var resp string

	val, err := r.do("SET", key, value, "NX", "EX", expire)
	if val != nil {
		resp = fmt.Sprintf("%s", val)
	}
	return resp, err
}

// SetEX key and value
// It sets the key wich will expired in `expire` seconds
func (r *Redigo) SetEX(key string, value interface{}, expire int) (string, error) {
	return redis.String(r.do("SETEX", key, expire, value))
}

// Get string value
func (r *Redigo) Get(key string) (string, error) {
	return redis.String(r.do("GET", key))
}

// MSet keys and values
// please use basic types only (no struct, array, or map) for arguments
func (r *Redigo) MSet(pairs ...interface{}) error {
	ok, err := redis.String(r.do("MSET", pairs...))
	if ok != "OK" && err == nil {
		return engine.ErrNotOK
	}
	return err
}

// MGet keys
func (r *Redigo) MGet(keys ...string) ([]string, error) {
	args := make([]interface{}, len(keys))
	for i, key := range keys {
		args[i] = key
	}
	return redis.Strings(r.do("MGET", args...))
}

// HSetEX key and value and sets the expiration to the given `expire` seconds
func (r *Redigo) HSetEX(key, field string, value interface{}, expire int) (int, error) {
	// we don't use r.do here because we do two commands
	conn, err := r.getConn()
	if err != nil {
		return 0, err
	}
	defer conn.Close()

	_, err = redis.Int(conn.Do("HSET", key, field, value))
	if err != nil {
		return 0, err
	}

	return redis.Int(conn.Do("EXPIRE", key, expire))
}

// HGet key and value
func (r *Redigo) HGet(key, field string) (string, error) {
	return redis.String(r.do("HGET", key, field))
}

// HMSet function
// please use basic types only (no struct, array, or map) for kv value
func (r *Redigo) HMSet(key string, kv map[string]interface{}) (string, error) {
	var (
		args = make([]interface{}, 1+(len(kv)*2))
		idx  = 1
	)
	args[0] = key
	for k, v := range kv {
		args[idx] = k
		args[idx+1] = v
		idx += 2
	}
	return redis.String(r.do("HMSET", args...))
}

// HMGet keys and value
func (r *Redigo) HMGet(key string, fields ...string) ([]string, error) {
	args := make([]interface{}, len(fields)+1)
	args[0] = key
	for i, field := range fields {
		args[i+1] = field
	}
	return redis.Strings(r.do("HMGET", args...))
}

// HDel fields of a key
func (r *Redigo) HDel(key string, fields ...string) (int, error) {
	args := make([]interface{}, len(fields)+1)
	args[0] = key
	for i, field := range fields {
		args[i+1] = field
	}
	return redis.Int(r.do("HDEL", args...))
}

// Append string to existing value in the key
func (r *Redigo) Append(key, value string) (int, error) {
	return redis.Int(r.do("APPEND", key, value))
}
