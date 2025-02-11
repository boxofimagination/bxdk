// Package redigo provide reliable redigo wrapper.
//
// It has these features:
//  - helper to a lot of redis commands
//  - safe access to redis pool, waiting for certain amount of time when the pool is exhausted.
//	  The default redigo will simply return invalid connection when the pool is exhausted.

package redigo

import (
	"context"
	"fmt"
	"time"

	"github.com/gomodule/redigo/redis"

	"github.com/boxofimagination/bxdk/go/redis/engine"
)


const (
	networkTCP = "tcp"
)

type (
	// Redigo defines the redis wrapper
	Redigo struct {
		pool *redis.Pool // redis pool
		poolWaitTime time.Duration // duration to wait when the pool exhausted
	}
)

// New creates new Redigo object from the given config
func New(cfg engine.Config) *Redigo {
	if cfg.MaxIdle == 0 {
		cfg.MaxIdle = cfg.MaxActive
	}

	// creates the pool
	pool := &redis.Pool{
		MaxIdle: cfg.MaxIdle,
		MaxActive: cfg.MaxActive,
		IdleTimeout: time.Duration(cfg.Timeout) * time.Second,
		TestOnBorrow: func(c redis.Conn, t time.Time) error {
			if time.Since(t) < time.Duration(cfg.IdlePingPeriod)*time.Second {
				return nil
			}

			_, err := c.Do("PING")
			return err
		},
		Dial: func() (redis.Conn, error) {
			return redis.Dial(networkTCP , cfg.Address)
		},
		Wait: true,
	}

	return &Redigo{
		pool: pool,
		poolWaitTime: time.Duration(cfg.PoolWaitMs) * time.Millisecond,
	}
}

// get connection from the pool with some timeout
func (r *Redigo) getConn() (redis.Conn, error) {
	ctx, cancel := context.WithTimeout(context.Background(), r.poolWaitTime)
	defer cancel()
	return r.pool.GetContext(ctx)
}

// Ping command to redis
func (r *Redigo) Ping() (string, error) {
	val, err  := r.Do("PING")
	return fmt.Sprint(val), err
}


// Do Command
func (r *Redigo) Do(cmd string, args ...interface{}) (interface{}, error) {
	return r.do(cmd, args...)
}

func (r *Redigo) do(cmd string, args ...interface{}) (interface{}, error) {
	conn, err := r.getConn()
	if err != nil {
		return nil, err
	}

	resp, err := conn.Do(cmd, args...)

	// we don't defer the Close because:
	// - defer has some performance hit
	// - it is in hot path
	// - the func is quite short & simple, we won't miss the defer here
	conn.Close()

	return resp, err
}

// GetConn get connection from the redis pool.
// Notes:
// - Please only use it for the pipelining feature.
func (r *Redigo) GetConn() (redis.Conn, error) {
	return r.getConn()
}

// IsErrNil returns true if the err given is ErrNil value.
// in case of redigo: it is redigo.ErrNil.
// Please use this func instead of comparing to redigo.ErrNil directly
// because each library has its own ErrNil definition.
func (r *Redigo) IsErrNil(err error) bool {
	return err == redis.ErrNil
}