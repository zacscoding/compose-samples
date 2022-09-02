package main

import (
	"context"
	"log"
	"sync"
	"time"

	"github.com/go-zookeeper/zk"
)

func main() {
	conn, ch, err := zk.Connect([]string{
		"localhost:12181",
		"localhost:22181",
		"localhost:32181",
	}, time.Second)
	if err != nil {
		panic(err)
	}
	defer conn.Close()
	conn.SetLogger(&zkLogger{})
	running := sync.WaitGroup{}
	log.Printf("Conn State: %s", conn.State().String())

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// Subscribe events
	running.Add(1)
	go func() {
		defer running.Done()
		for {
			select {
			case e := <-ch:
				log.Printf("[EventLoop] EventOccur: %v", e)
			case <-ctx.Done():
				log.Printf("Stopping event loop")
				return
			}
		}
	}()

	// Check "/MyFirstZnode" and delete it if exist
	const path = "/MyFirstZnode"
	ok, stat, err := conn.Exists(path)
	if err != nil {
		panic(err)
	}
	log.Printf("'%s' exists: %v, stat: %v", path, ok, stat)
	if ok {
		if err := conn.Delete(path, stat.Version); err != nil {
			panic(err)
		}
	}

	// Create a new ZNode with "MyData" data.
	result, err := conn.Create(path, []byte("MyData"), 0, zk.WorldACL(zk.PermAll))
	if err != nil {
		log.Fatal(err)
	}
	log.Printf("'%s' create result: %s", path, result)

	cancel()
	running.Wait()
}

type zkLogger struct{}

func (l *zkLogger) Printf(format string, v ...interface{}) {
	log.Printf("[ZKClient] "+format, v...)
}
