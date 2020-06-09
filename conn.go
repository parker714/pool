package pool

// Conn pool interface.
type Conn interface {
	Pool
	Full(int) error
}

type conn struct {
	pool
	// Dial method for creating net resources.
	dial func() (interface{}, error)
	// Ping method for check whether the resource is valid.
	ping func(interface{}) error
}

// NewConn returns net conn pool.
func NewConn(size int, dial func() (interface{}, error), ping func(interface{}) error) (Conn, error) {
	c := &conn{
		pool: pool{
			resources: make(chan interface{}, size),
			size:      size,
		},
		dial: dial,
		ping: ping,
	}

	return c, c.Full(size)
}

// Full num net conn.
func (c *conn) Full(num int) error {
	for i := 0; i < num; i++ {
		x, err := c.dial()
		if err != nil {
			return err
		}
		if err := c.pool.Put(x); err != nil {
			return err
		}
	}
	return nil
}

// Get returns net conn.
func (c *conn) Get() (interface{}, error) {
	x, err := c.pool.Get()
	if err != nil {
		return nil, err
	}
	if c.ping != nil {
		if err := c.ping(x); err != nil {
			return nil, err
		}
	}
	return x, nil
}
