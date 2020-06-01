package redis

import (
	"errors"
	"fmt"

	"github.com/boxofimagination/bxdk/go/defaults"
	"github.com/boxofimagination/bxdk/go/redis/engine"
	"github.com/boxofimagination/bxdk/go/redis/engine/redigo"
)

// Config defines configuration for the redis library
type Config = engine.Config

var (
	// ErrNilFiller returned when `GetWithSingleFiller` called with nil filler
	ErrNilFiltre = errors.New("empty filter")
)


// Redis defines interface for TDK redis library
type Redis = engine.Redis

// Pipeliner alias of engine.Pipeliner, the caller don't have to import engine
type Pipeliner = engine.Pipeliner


// Client defines a redis client
type Client struct {
	Redis
}

// New creates new redis tdk library from the given config.
// The default engine used is `redigo` based engine
func New(cfg Config) (*Client, error ) {
	err := defaults.SetDefault(&cfg)
	if err != nil {
		return nil, err
	}

	var eng Redis

	switch cfg.EngineType {
	case engine.Redigo:
		eng = redigo.New(cfg)
	default:
		return nil, fmt.Errorf("invalid engine type: %v", cfg.EngineType)
	}

	if !cfg.NoPingOnCreate {
		if _, err = eng.Ping(); err != nil {
			return nil, err
		}
	}

	return &Client{
		Redis: eng,
	}, nil
}