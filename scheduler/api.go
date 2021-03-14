package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"

	"github.com/julienschmidt/httprouter"
)

const (
	contentTypeHeader     = "Content-Type"
	applicationJSONHeader = "application/json"
	corsHeader            = "Access-Control-Allow-Origin"
)

func apiStartJob(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	startJobReq := apiStartJobReq{}

	w.Header().Set(contentTypeHeader, applicationJSONHeader)

	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(apiError{Error: err.Error()})
		return
	}

	err = json.Unmarshal(body, &startJobReq)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(apiError{Error: err.Error()})
		return
	}

	jobID, err := startJobOnWorker(startJobReq)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(apiError{Error: err.Error()})
		return
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(apiStartJobRes{JobID: jobID})
}

func apiStopJob(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	stopJobReq := apiStopJobReq{}

	w.Header().Set(contentTypeHeader, applicationJSONHeader)

	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(apiError{Error: err.Error()})
		return
	}

	err = json.Unmarshal(body, &stopJobReq)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(apiError{Error: err.Error()})
		return
	}

	if err := stopJobOnWorker(stopJobReq); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(apiError{Error: err.Error()})
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(apiStopJobRes{Success: true})
}

func apiQueryJob(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	queryJobReq := apiQueryJobReq{}

	w.Header().Set(contentTypeHeader, applicationJSONHeader)

	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(apiError{Error: err.Error()})
		return
	}

	err = json.Unmarshal(body, &queryJobReq)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(apiError{Error: err.Error()})
		return
	}

	jobDone, jobError, jobErrorText, err := queryJobOnWorker(queryJobReq)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(apiError{Error: err.Error()})
		return
	}

	queryJobRes := apiQueryJobRes{
		Done:      jobDone,
		Error:     jobError,
		ErrorText: jobErrorText,
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(queryJobRes)
}

func apiQueryServiceVersionOnWorker(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	serviceName := ps.ByName("service")
	workerID := ps.ByName("worker")

	fmt.Println(serviceName)
	fmt.Println(workerID)

	w.Header().Set(contentTypeHeader, applicationJSONHeader)

	queryServiceReq := apiServiceVersionReq{
		ServiceName: serviceName,
		WorkerID:    workerID,
	}

	serviceVersion, err := queryServiceVersionOnWorker(queryServiceReq)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(apiError{Error: err.Error()})
		return
	}

	serviceVersionRes := apiServiceVersionRes{
		Service: serviceName,
		Version: serviceVersion,
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(serviceVersionRes)
}

func apiListWorkers(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	workers, err := listWorkers()
	w.Header().Set(corsHeader, "*")
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(apiError{Error: err.Error()})
		return
	}
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(workers)
}

func apiListEnv(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	envs, err := listEnvs()
	w.Header().Set(corsHeader, "*")
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(apiError{Error: err.Error()})
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(envs)
}

func createRouter() *httprouter.Router {
	router := httprouter.New()

	router.POST("/start", apiStartJob)
	router.POST("/stop", apiStopJob)
	router.POST("/query", apiQueryJob)

	router.GET("/workers/list", apiListWorkers)
	router.GET("/env/list", apiListEnv)
	router.GET("/version/:service/:worker", apiQueryServiceVersionOnWorker)

	return router
}

func api() {
	srv := &http.Server{
		Addr:    config.HTTPServer.Addr,
		Handler: createRouter(),
	}

	log.Println("HTTP Server listening on", config.HTTPServer.Addr)
	if err := srv.ListenAndServe(); err != nil {
		log.Fatal(err)
	}
}
