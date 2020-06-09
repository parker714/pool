package pool

import (
	"runtime"
	"testing"
)

func TestPool_Put(t *testing.T) {
	p := New(1)
	if err := p.Put(1); err != nil {
		t.Fatalf("want ``, got `%s`", err)
	}

	if err := p.Put(2); err != errFull {
		t.Fatalf("want `%s`, got `%s`", errFull, err)
	}
}

func TestPool_Get(t *testing.T) {
	p := New(1)
	if err := p.Put(1); err != nil {
		t.Fatalf("want ``, got `%s`", err)
	}

	x, err := p.Get()
	if err != nil {
		t.Fatalf("want ``, got `%s`", err)
	}
	if x.(int) != 1 {
		t.Fatalf("want `1`, got `%v`", x)
	}

	_, err = p.Get()
	if err != errNoAvailable {
		t.Fatalf("want `%s`, got `%s`", errNoAvailable, err)
	}
}

func BenchmarkPool_Get(b *testing.B) {
	num := runtime.GOMAXPROCS(0)
	p := New(num)

	for i := 0; i < num; i++ {
		_ = p.Put(1)
	}

	b.SetParallelism(1)
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			x, err := p.Get()
			if err != nil {
				b.Fatalf("want ``, got `%s`", err)
			}

			if err = p.Put(x); err != nil {
				b.Fatalf("want ``, got `%s`", err)
			}
		}
	})
}
