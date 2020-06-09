package pool

import (
	"errors"
)

var (
	errFull        = errors.New("pool: resource is full")
	errNoAvailable = errors.New("pool: no resources available")
)

// Pool interface.
type Pool interface {
	Put(x interface{}) error
	Get() (interface{}, error)
}

type pool struct {
	resources chan interface{}
	size      int
}

// New returns new pool instance.
func New(size int) Pool {
	return &pool{
		resources: make(chan interface{}, size),
		size:      size,
	}
}

// Put adds x to the pool.
func (p *pool) Put(x interface{}) error {
	if p.size == len(p.resources) {
		return errFull
	}
	p.resources <- x
	return nil
}

// Get returns x.
func (p *pool) Get() (interface{}, error) {
	select {
	case resource := <-p.resources:
		return resource, nil
	default:
		return nil, errNoAvailable
	}
}
