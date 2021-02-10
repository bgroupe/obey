package main

import (
	"fmt"
	"sync"

	"github.com/gomodule/redigo/redis"
	"github.com/google/uuid"
)

var (
	workersMutex = &sync.Mutex{}
	workers      = make(map[string]*Worker)
	db           DB
)

// worker holds the information about registered workers
// 		- id: uuid assigned when the worker first register.
// 		- addr: workers network address, later used to create
// 				grpc client to the worker
type worker struct {
	id   string
	addr string
}

// Worker holds information about registered worker
type Worker struct {
	ID      string `redis:"id"`
	Address string `redis:"address"`
}

// newWorker creates a new worker instance and adds
// the new worker to the map.
// Returns:
// 		- string: worker id
func newWorker(address string) (string, error) {
	workersMutex.Lock()
	defer workersMutex.Unlock()

	db, err := ConnectDB(config.RedisServer.Addr)

	if err != nil {
		return "", err
	}
	defer db.Conn.Close()

	workerID := uuid.New().String()
	// `worker:<uuid>`
	workerKey := fmt.Sprintf("worker:%s", workerID)

	reply, err := db.Conn.Do(
		"HMSET",
		workerKey,
		"address",
		address,
		"id",
		workerID,
	)

	if err != nil {
		return "", err
	}

	fmt.Print(reply)

	var w Worker
	values, err := redis.Values(db.Conn.Do("HGETALL", workerKey))

	if err != nil {
		return "", err
	}
	err = redis.ScanStruct(values, &w)

	if err != nil {
		return "", err
	}

	reply, err = db.Conn.Do(
		"LPUSH",
		"workers",
		workerID,
	)

	if err != nil {
		return "", err
	}

	workers[workerID] = &w

	return workerID, nil
}

func listWorkers() (map[string]*Worker, error) {
	wrkrs := make(map[string]*Worker)
	db, err := ConnectDB(config.RedisServer.Addr)

	if err != nil {
		var dummy map[string]*Worker
		return dummy, err
	}

	defer db.Conn.Close()

	values, err := redis.Strings(db.Conn.Do(
		"LRANGE",
		"workers",
		"0",
		"-1",
	))

	if err != nil {
		dummy := make(map[string]*Worker)
		return dummy, err
	}

	for _, uuid := range values {
		workerKey := fmt.Sprintf("worker:%v", uuid)
		var w Worker
		values, err := redis.Values(db.Conn.Do("HGETALL", workerKey))

		if err != nil {
			return wrkrs, err
		}
		err = redis.ScanStruct(values, &w)

		if err != nil {
			return wrkrs, err
		}

		wrkrs[uuid] = &w
	}

	return wrkrs, nil
}

// delWorker removes the worker with the given id
// from known workers map.
func delWorker(id string) error {
	workersMutex.Lock()
	defer workersMutex.Unlock()

	db, err := ConnectDB(config.RedisServer.Addr)

	if err != nil {
		return err
	}
	defer db.Conn.Close()

	workerKey := fmt.Sprintf("worker:%v", id)

	// DEL returns `1` or `0` if successful
	// Remove worker
	_, err = redis.Int(db.Conn.Do("DEL", workerKey))

	if err != nil {
		return err
	}
	// Remove id from worker list
	_, err = redis.Int(db.Conn.Do("LREM", "workers", 0, id))
	if err != nil {
		return err
	}
	return nil

}

// TODO:
// [x] Delete worker from Redis store
// [] Database functions
// [] Serialize store to marshaled json
// [] Create shared models
// [x] Surface errors for deregistering workers
// [] Use redis.Pool
// [] remove workers.Mutex()
// [] potentially look into redlock https://github.com/go-redsync/redsync
