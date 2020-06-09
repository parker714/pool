package pool

import (
	"errors"
	"testing"
)

var (
	errDial  = errors.New("mock: dial error")
	errPing  = errors.New("mock: ping error")
	pingFlag bool
)

type connection struct{}

func mockDialSuccess() (interface{}, error) {
	return connection{}, nil
}

func mockDialFail() (interface{}, error) {
	return nil, errDial
}

func mockPing(interface{}) error {
	if pingFlag {
		return nil
	}
	return errPing
}

func TestConn_Full(t *testing.T) {
	t.Run("normal", func(t *testing.T) {
		_, err := NewConn(2, mockDialSuccess, nil)
		if err != nil {
			t.Fatalf("want ``, got `%s`", err)
		}
	})

	t.Run("conn dial fail", func(t *testing.T) {
		_, err := NewConn(2, mockDialFail, nil)
		if err != errDial {
			t.Fatalf("want `%s`, got `%s`", errDial, err)
		}
	})

	t.Run("pool put fail", func(t *testing.T) {
		cp := conn{
			pool: pool{
				resources: make(chan interface{}, 1),
				size:      1,
			},
			dial: mockDialSuccess,
			ping: nil,
		}

		if err := cp.Full(1); err != nil {
			t.Fatalf("want ``, got `%s`", err)
		}

		if err := cp.Full(1); err != errFull {
			t.Fatalf("want `%s`, got `%s`", errFull, err)
		}
	})
}

func TestConn_Get(t *testing.T) {
	t.Run("normal", func(t *testing.T) {
		cp, err := NewConn(1, mockDialSuccess, nil)
		if err != nil {
			t.Fatalf("want ``, got `%s`", err)
		}
		_, err = cp.Get()
		if err != nil {
			t.Fatalf("want ``, got `%s`", err)
		}
	})

	t.Run("no available", func(t *testing.T) {
		cp, err := NewConn(0, mockDialSuccess, nil)
		if err != nil {
			t.Fatalf("want `%s`, got `%s`", errDial, err)
		}
		_, err = cp.Get()
		if err != errNoAvailable {
			t.Fatalf("want `%s`, got `%s`", errNoAvailable, err)
		}
	})

	t.Run("ping fail", func(t *testing.T) {
		cp, err := NewConn(1, mockDialSuccess, mockPing)
		if err != nil {
			t.Fatalf("want `%s`, got `%s`", errDial, err)
		}
		pingFlag = false
		_, err = cp.Get()
		if err != errPing {
			t.Fatalf("want `%s`, got `%s`", errPing, err)
		}
	})
}

func BenchmarkConn_Get(b *testing.B) {
	pingFlag = true
	cp, err := NewConn(10, mockDialSuccess, mockPing)
	if err != nil {
		b.Errorf("want ``, got `%s`", err)
	}

	b.SetParallelism(1)
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			x, err := cp.Get()
			if err != nil {
				b.Fatalf("want ``, got `%s`", err)
			}

			if err := cp.Put(x); err != nil {
				b.Fatalf("want ``, got `%s`", err)
			}
		}
	})
}
