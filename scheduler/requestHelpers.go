package main

import (
	"encoding/json"
	"fmt"
	"sync"
	"time"

	pb "github.com/bgroupe/obey/jobscheduler"
	"github.com/gomodule/redigo/redis"
	"github.com/google/uuid"
	log "github.com/sirupsen/logrus"
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
	ID            string `redis:"id" json:"id"`
	Address       string `redis:"address" json:"address"`
	BroadcastAddr string `redis:"broadcastAddress" json:"broadcast-addr,omitempty"`
	EnvName       string `redis:"envName" json:"env"`
	EnvType       string `redis:"envType" json:"env-type"`
	ServerStatus  string `redis:"serverStatus" json:"server-status,omitempty"`
	LaunchTime    string `redis:"launchTime" json:"launch-time,omitempty"`
}

// newWorker creates a new worker instance and adds
// the new worker to the map.
// Returns:
// 		- string: worker id
func newWorker(r *pb.RegisterReq) (string, error) {
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

	// TODO: Add dynamic iteration of struct fields with `reflect` package
	// https://stackoverflow.com/questions/18926303/iterate-through-the-fields-of-a-struct-in-go
	reply, err := db.Conn.Do(
		"HMSET",
		workerKey,
		"address",
		r.Address,
		"broadcastAddress",
		r.BroadcastAddress,
		"id",
		workerID,
		"envName",
		r.EnvName,
		"envType",
		r.EnvType,
		"launchTime",
		time.Now().UTC().String(),
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

	// TODO: move worker env to archive
	return nil

}

// reportServiceData
func acceptReportServiceData(dataReq *pb.ReportServiceDataRequest) error {
	db, err := ConnectDB(config.RedisServer.Addr)

	if err != nil {
		return err
	}

	defer db.Conn.Close()

	// Full Marshaled Request
	fullMarshaledReq, err := json.Marshal(dataReq)

	values, err := redis.Strings(db.Conn.Do(
		"LRANGE",
		"workers",
		"0",
		"-1",
	))

	if err != nil {
		return err
	}
	var envUUID string

	for _, uuid := range values {
		workerKey := fmt.Sprintf("worker:%v", uuid)

		workerEnv, err := redis.String(db.Conn.Do("HGET", workerKey, "envName"))

		if err != nil {
			return err
		}

		if workerEnv == dataReq.Name {
			envUUID = uuid
			// maybe not break here...
			break
		} else {
			envUUID = ""
		}

	}

	if err != nil {
		return err
	}

	if envUUID != "" {
		workerEnvKey := fmt.Sprintf("env:%s", envUUID)
		current, err := redis.String(db.Conn.Do(
			"HGET",
			workerEnvKey,
			"json",
		))

		// 1. fetch current record
		// 2. Check if empty
		// 2. unmarshal into struct
		// 3. call smart replace on current fields
		// 4. marshal that value and send

		// If first time discovery, handle the nil response and dump the whole manifest
		if err != nil {
			if err == redis.ErrNil {
				reply, err := db.Conn.Do(
					"HMSET",
					workerEnvKey,
					"json",
					fullMarshaledReq,
				)

				if err != nil {
					return err
				}

				log.WithFields(log.Fields{
					"status": reply,
				}).Info("environment saved")
				return nil
			}
			log.Warn("DB Error")
			return err
		}

		svcDataCurrent := pb.ReportServiceDataRequest{}
		err = json.Unmarshal([]byte(current), &svcDataCurrent)

		if err != nil {
			return err
		}

		for _, sv := range dataReq.ServiceData {
			for i, svc := range svcDataCurrent.ServiceData {
				if svc.Name == sv.Name {
					log.WithFields(log.Fields{
						"service": svc.Name,
					}).Info("updating service")

					svcDataCurrent.ServiceData[i] = sv
				}
			}

		}

		// call of json.Marshal copies lock value: github.com/bgroupe/obey/jobscheduler.ReportServiceDataRequest contains google.golang.org/protobuf/internal/impl.MessageState contains sync.Mutex
		// copylocks
		updatedMarshaledReq, err := json.Marshal(svcDataCurrent)
		reply, err := db.Conn.Do(
			"HMSET",
			workerEnvKey,
			"json",
			updatedMarshaledReq,
		)

		if err != nil {
			return err
		}

		log.WithFields(log.Fields{
			"status": reply,
		}).Info("services updated")

		return nil
	}

	return fmt.Errorf("could not find worker by name: %s", dataReq.Name)
}

func listEnvs() ([]*pb.ReportServiceDataRequest, error) {
	// Step 1. Create DB Connection
	db, err := ConnectDB(config.RedisServer.Addr)

	if err != nil {
		var dummy []*pb.ReportServiceDataRequest
		return dummy, err
	}

	defer db.Conn.Close()

	// Step 2. Get UUIDs from worker list
	values, err := redis.Strings(db.Conn.Do(
		"LRANGE",
		"workers",
		"0",
		"-1",
	))

	fmt.Println(values)

	if err != nil {
		fmt.Println("Fuck!")
		var dummy []*pb.ReportServiceDataRequest
		return dummy, err
	}

	// Step 3. Get envs using worker UUIDS
	envz := []*pb.ReportServiceDataRequest{}
	for _, uuid := range values {
		envKey := fmt.Sprintf("env:%v", uuid)

		marshaledJSON, err := redis.String(db.Conn.Do("HGET", envKey, "json"))

		if err != nil {
			return envz, err
		}
		rsdr := pb.ReportServiceDataRequest{}

		fmt.Println(values)
		err = json.Unmarshal([]byte(marshaledJSON), &rsdr)

		if err != nil {
			return envz, err
		}

		envz = append(envz, &rsdr)
	}
	return envz, nil
}

// Contains does a thing
func Contains(a []string, x string) bool {
	for _, n := range a {
		if x == n {
			return true
		}
	}
	return false
}

// Remove via append
// x := []Foo{{Name: "test", Type: "fart",},{Name: "fuck", Type: "Poo",}}

// 	for i, t := range x {
// 	  if t.Name == "fuck" {
// 	    x = append(x[:i], x[i+1:]...)
// 	  }
// 	}
