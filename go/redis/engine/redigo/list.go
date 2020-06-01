package redigo

import (
	"fmt"

	"github.com/gomodule/redigo/redis"
)

// RPush append values to the key
func (r *Redigo) RPush(key string, values ...string) (int, error) {
	args := make([]interface{}, len(values)+1)
	args[0] = key
	for i, value := range values {
		args[i+1] = value
	}
	return redis.Int(r.do("RPUSH", args...))
}

// RPop Removes and returns the last element of the list stored at key
// return redis.ErrNil if the key is not exist
func (r *Redigo) RPop(key string) (string, error) {
	var resp string

	val, err := r.do("RPOP", key)
	if val != nil {
		resp = fmt.Sprintf("%s", val)
	} else if err == nil {
		// val is nil and err is nil too
		err = redis.ErrNil
	}
	return resp, err
}

// LLen get the length of the list
func (r *Redigo) LLen(key string) (int64, error) {
	return redis.Int64(r.do("LLEN", key))
}

// LPush prepend values to the list
func (r *Redigo) LPush(key string, values ...string) (int, error) {
	args := make([]interface{}, len(values)+1)
	args[0] = key
	for i, value := range values {
		args[i+1] = value
	}
	return redis.Int(r.do("LPUSH", args...))
}

// LPop removes and get the first element in the list
func (r *Redigo) LPop(key string) (string, error) {
	var resp string

	val, err := r.do("LPOP", key)
	if val != nil {
		resp = fmt.Sprintf("%s", val)
	} else if err == nil {
		// val is nil and err is nil too
		err = redis.ErrNil
	}
	return resp, err
}

// LRange returns the specified elements of the list stored at key
func (r *Redigo) LRange(key string, start, stop int64) ([]string, error) {
	return redis.Strings(r.do("LRANGE", key, start, stop))
}
