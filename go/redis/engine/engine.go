// Package engine contains various engines for the redis client
package engine

import (
	"errors"
)

// Type defines engine type.
// It currently only support `redigo` engine
type Type string

// ScanAllResult struct define response for redis scan all
type ScanAllResult struct {
	Keys []string
	Err  error
}

const (
	// Redigo is redigo engine
	Redigo Type = "redigo"
)

// ErrNotOK returned if redis not respond with OK but error is nil
var ErrNotOK = errors.New("not ok")

// Config of redis engine
type Config struct {
	// EngineType defines engine's/library type.
	// The supported values : redigo,goredis
	EngineType Type `yaml:"engine_type" default:"redigo"`

	// Redis server address
	Address string `yaml:"address"`

	// Maximum number of idle connections in the pool.
	// Only for redigo engine
	MaxIdle int `yaml:"maxidle"`

	// Maximum number of connections allocated by the pool at a given time.
	// When zero, there is no limit on the number of connections in the pool.
	MaxActive int `yaml:"maxactive" default:"50"`

	// Close connections after remaining idle for this duration. If the value
	// is zero, then idle connections are not closed. Applications should set
	// the timeout to a value less than the server's timeout.
	Timeout int `yaml:"timeout"`

	// IdlePingPeriod defines period in seconds after which the connection need
	// to be check(PING) before being used.
	// Default is 10 seconds
	IdlePingPeriod int `yaml:"idle_ping_period" default:"10"`

	// pool wait time in millisecond.
	// If > 0 and the pool is at the MaxActive limit, then Get() waits
	// for a connection to be returned to the pool before returning. It only waits
	// for PoolWaitMs millisecond
	PoolWaitMs int `yaml:"pool_wait_ms" default:"1000"`

	// NoPingOnCreate is a flag to indicate whether it will be do ping check on `New` or not.
	// If true: client will do redis PING on `New`, make sure that the server is up.
	NoPingOnCreate bool `yaml:"no_ping_on_create"`
}

// Redis defines interface for BXDK redis library
type Redis interface {
	// Ping command to redis
	Ping() (string, error)

	// Do command
	Do(cmd string, args ...interface{}) (interface{}, error)

	// IsErrNil returns true if the err given is ErrNil value.
	// in case of redis: it is redis.ErrNil.
	// Please use this func instead of comparing to redis.ErrNil directly
	// because each library has its own ErrNil definition.
	IsErrNil(err error) bool

	// Set key and value
	Set(key string, value interface{}) error

	// SetNX do SETNX (only set if not exist) with SET's NX & EX args.
	// It sets the key which will expired in `expire` seconds
	SetNX(key string, value interface{}, expire int) (string, error)

	// SetEX key and value
	// It sets the key wich will expired in `expire` seconds
	SetEX(key string, value interface{}, expire int) (string, error)

	// Get string value
	Get(key string) (string, error)

	// Delete delete keys from the server
	Delete(keys ...string) (int, error)

	// MSet keys and values
	// please use basic types only (no struct, array, or map) for arguments
	MSet(pairs ...interface{}) error

	// MGet keys
	MGet(keys ...string) ([]string, error)

	// HSetEX key and value and sets the expiration to the given `expire` seconds
	HSetEX(key, field string, value interface{}, expire int) (int, error)

	// HGet key and value
	HGet(key, field string) (string, error)

	// HMSet function
	// please use basic types only (no struct, array, or map) for kv value
	HMSet(key string, kv map[string]interface{}) (string, error)

	// HMGet keys and value
	HMGet(key string, fields ...string) ([]string, error)

	// HDel fields of a key
	HDel(key string, fields ...string) (int, error)

	// Incr function
	Incr(key string) (int64, error)

	// IncrBy function
	IncrBy(key string, value int64) (int64, error)

	// Decr function
	Decr(key string) (int64, error)

	// DecrBy function
	DecrBy(key string, value int64) (int64, error)

	// Expire set expiration time for a key
	// `expiry` is in seconds
	Expire(key string, expiry int) (int, error)

	// TTL return remaining ttl of a key
	// from: https://redis.io/commands/ttl
	// The command returns -2 if the key does not exist.
	// The command returns -1 if the key exists but has no associated expire.
	TTL(key string) (int, error)

	// Exists checks if a key exists.
	// return true if exists.
	Exists(key string) (bool, error)

	// LLen get the length of the list
	LLen(key string) (int64, error)

	// LPush prepends values to the list and returns the length of the list
	LPush(key string, values ...string) (int, error)

	// LPop removes and get the first element in the list
	LPop(key string) (string, error)

	// LRange returns the specified elements of the list stored at key
	LRange(key string, start, stop int64) ([]string, error)

	// RPush append values to the list and return the length of the list
	RPush(key string, values ...string) (int, error)

	// RPop Removes and returns the last element of the list stored at key
	// return redigo.ErrNil if the key is not exist
	RPop(key string) (string, error)

	// Scan will do SCAN command to get keys by given pattern
	// returning keys, cursor, and error
	Scan(pattern string, cursor uint64, count int64) ([]string, uint64, error)

	Pipeline(retry, numCmdHint int) Pipeliner

	// SAdd Add the specified members to the set stored at key.
	// Specified members that are already a member of this set are ignored.
	// If key does not exist, a new set is created before adding the specified members.
	// An error is returned when the value stored at key is not a set.
	SAdd(key string, members ...interface{}) (int64, error)

	// SRem Remove the specified members from the set stored at key.
	// Specified members that are not a member of this set are ignored.
	// If key does not exist, it is treated as an empty set and this command returns 0.
	// An error is returned when the value stored at key is not a set.
	SRem(key string, members ...interface{}) (int64, error)

	// SMembers Returns all the members of the set value stored at key.
	SMembers(key string) ([]string, error)

	// Append string to existing value in the key
	Append(key, value string) (int, error)
}

// CmdErr is redis command, args, and error
type CmdErr interface {
	Name() string
	Args() []interface{}
	Err() error
}

// Pipeliner is interface for pipeline object
type Pipeliner interface {
	// AddRawCmd adds/queues raw redis command to the pipeline
	AddRawCmd(cmd string, args ...interface{})

	// Exec executes all queued commands using redis pipeline.
	// - cmdErrs is command & error of each of the queued commands.
	// - firstErr is first index of the errored commands.
	// 		there is no error if firstErr < 0
	// - err is not nil if the whole pipeline execution failed
	Exec() (cmdErrs []CmdErr, firstErr int, err error)

	// Discard resets the pipeline and discards queued commands
	Discard() error

	// Close closes the pipeline, releasing any open resources.
	Close() error

	Incr(string)
	IncrBy(string, int64)
	Decr(string)
	DecrBy(string, int64)
	Expire(string, int)
	Delete(keys ...string)
	HMSet(key string, kv map[string]interface{})
	HDel(key string, fields ...string)
}
