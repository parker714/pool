Pool [![Release](https://img.shields.io/github/release/parker714/pool.svg)](https://github.com/parker714/pool/releases)
=====================

## Overview
```go
package main

import (
	"errors"
	"github.com/parker714/pool"
	"github.com/streadway/amqp"
	"log"
)

func dial() (interface{}, error) {
	return amqp.Dial("amqp://guest:guest@127.0.0.1:5672//")
}

func ping(x interface{}) (err error) {
	if x.(*amqp.Connection).IsClosed() {
		err = errors.New("example: amqp connection was closed")
	}
	return
}

func main() {
	cp, err := pool.NewConn(2, dial, ping)
	if err != nil {
		log.Fatalf("main: new conn pool err, %s\n", err)
	}

	x, err := cp.Get()
	if err != nil {
		log.Fatalf("example: cp get conn err, %s\n", err)
	}
	log.Printf("example: get %t\n", x)

	if err := cp.Put(x); err != nil {
		log.Printf("example: cp put err, %s\n", err)
	}
}
```

## License

This project is under the MIT License. See the [LICENSE](https://github.com/parker714/pool/blob/master/LICENSE) file for the full license text.
