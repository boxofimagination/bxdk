package redigo

import (
	"github.com/gomodule/redigo/redis"
)

// SAdd Add the specified members to the set stored at key.
// Specified members that are already a member of this set are ignored.
// If key does not exist, a new set is created before adding the specified members.
// An error is returned when the value stored at key is not a set.
func (r *Redigo) SAdd(key string, members ...interface{}) (int64, error) {
	args := append([]interface{}{key}, members...)
	return redis.Int64(r.do("SADD", args...))
}

// SRem Remove the specified members from the set stored at key.
// Specified members that are not a member of this set are ignored.
// If key does not exist, it is treated as an empty set and this command returns 0.
// An error is returned when the value stored at key is not a set.
func (r *Redigo) SRem(key string, members ...interface{}) (int64, error) {
	args := append([]interface{}{key}, members...)
	return redis.Int64(r.do("SREM", args...))
}

// SMembers Returns all the members of the set value stored at key.
func (r *Redigo) SMembers(key string) ([]string, error) {
	return redis.Strings(r.do("SMEMBERS", key))
}
