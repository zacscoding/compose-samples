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
		"localhost:2181",
	}, time.Second)
	if err != nil {
		panic(err)
	}
	defer conn.Close()
	conn.SetLogger(&zkLogger{})
	running := sync.WaitGroup{}

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

	//2022/09/03 01:28:44 [EventLoop] EventOccur: {EventSession StateConnecting  <nil> [::1]:2181}
	//2022/09/03 01:28:44 [ZKClient]connected to [::1]:2181
	//2022/09/03 01:28:44 [EventLoop] EventOccur: {EventSession StateConnected  <nil> [::1]:2181}
	//2022/09/03 01:28:44 [ZKClient]authenticated: id=72057868293898243, timeout=4000
	//2022/09/03 01:28:44 [EventLoop] EventOccur: {EventSession StateHasSession  <nil> [::1]:2181}
	//2022/09/03 01:28:44 [ZKClient]re-submitting `0` credentials after reconnect
	//2022/09/03 01:28:44 '/MyFirstZnode' exists: true, stat: &{10 10 1662136114471 1662136114471 0 0 0 0 6 0 10}
	//2022/09/03 01:28:44 '/MyFirstZnode' create result: /MyFirstZnode
	//2022/09/03 01:28:44 Stopping event loop
	//2022/09/03 01:28:44 [ZKClient]recv loop terminated: EOF
	//2022/09/03 01:28:44 [ZKClient]send loop terminated: <nil>
}

type zkLogger struct{}

func (l *zkLogger) Printf(format string, v ...interface{}) {
	log.Printf("[ZKClient]"+format, v...)
}
