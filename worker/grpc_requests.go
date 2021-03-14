package main

import (
	"context"
	"log"
	"time"

	pb "github.com/bgroupe/obey/jobscheduler"
	"github.com/bgroupe/obey/services"
	"google.golang.org/grpc"
)

// registerWorker dials the scheduler GRPC server and registers
// the calling worker with the worker's GRPC server address.
// Worker's GRPC server address is later used by the scheduler to dial
// worker to start/stop/query jobs.
func registerWorker() {
	conn, err := grpc.Dial(config.Scheduler.Addr, grpc.WithInsecure())
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	defer conn.Close()
	c := pb.NewSchedulerClient(conn)

	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	registerReq := pb.RegisterReq{
		Address:          config.GRPCServer.Addr,
		BroadcastAddress: config.WorkerEnvConfig.BroadcastAddr,
		EnvName:          config.WorkerEnvConfig.Name,
		EnvType:          config.WorkerEnvConfig.Type,
		// https://pkg.go.dev/github.com/golang/protobuf/ptypes#Timestamp
		// LaunchTime:       time.Unix(time.Now()),
	}
	r, err := c.RegisterWorker(ctx, &registerReq)
	if err != nil {
		log.Fatalf("could not register: %v", err)
	}

	workerID = r.WorkerID
	log.Printf("Registered ID: %s", r.WorkerID)
}

// deregisterWorker deregisters the calling worker from the scheduler.
// Scheduler will remove the worker from the known workers. Any nonpanic
// exit by the worker application should be calling deregister function
// before termination.
func deregisterWorker() {
	conn, err := grpc.Dial(config.Scheduler.Addr, grpc.WithInsecure())
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	defer conn.Close()
	c := pb.NewSchedulerClient(conn)

	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	deregisterReq := pb.DeregisterReq{
		WorkerID: workerID,
	}
	r, err := c.DeregisterWorker(ctx, &deregisterReq)
	if err != nil {
		log.Fatalf("could not deregister: %v", err)
	}

	log.Printf("Deregistered OK: %t", r.Success)
}

// reportServiceData sends reports to the scheduler in
func reportServiceData(sc []services.ScrapedContainer) {
	conn, err := grpc.Dial(config.Scheduler.Addr, grpc.WithInsecure())
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	defer conn.Close()
	c := pb.NewSchedulerClient(conn)
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	// Need to append pointers here
	var serviceData []*pb.ServiceData

	for _, container := range sc {
		sd := pb.ServiceData{
			Name:     container.Name,
			Version:  container.Version,
			State:    container.State,
			Status:   container.Status,
			Created:  container.Created,
			Revision: container.ImageSha,
		}

		serviceData = append(serviceData, &sd)
	}

	reportReq := pb.ReportServiceDataRequest{
		Name:        config.WorkerEnvConfig.Name,
		ServiceData: serviceData,
	}

	r, err := c.ReportServiceData(ctx, &reportReq)

	if err != nil {
		log.Fatalf("could not report service data: %v", err)
	}

	log.Printf("Scrape Report OK: %t", r.Success)
}
