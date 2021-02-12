package main

// apiStartJobReq expected API payload for `/start`
type apiStartJobReq struct {
	Command  string `json:"command"`
	Path     string `json:"path"`
	WorkerID string `json:"worker_id"`
}

// apiStartJobRes returned API payload for `/start`
type apiStartJobRes struct {
	JobID string `json:"job_id"`
}

// apiStopJobReq expected API payload for `/stop`
type apiStopJobReq struct {
	JobID    string `json:"job_id"`
	WorkerID string `json:"worker_id"`
}

// apiStopJobRes returned API payload for `/stop`
type apiStopJobRes struct {
	Success bool `json:"success"`
}

// apiQueryJobReq expected API payload for `/query`
type apiQueryJobReq struct {
	JobID    string `json:"job_id"`
	WorkerID string `json:"worker_id"`
}

// apiQueryJobRes returned API payload for `/query`
type apiQueryJobRes struct {
	Done      bool   `json:"done"`
	Error     bool   `json:"error"`
	ErrorText string `json:"error_text"`
}

// apiServiceVersionResponse for `/version/:service`
type apiServiceVersionReq struct {
	ServiceName string `json:"service"`
	WorkerID    string `json:"worker_id"`
}

// apiServiceVersionResponse for `/version/:service`
type apiServiceVersionRes struct {
	Version string `json:"version"`
	Service string `json:"service"`
}

// apiError is used as a generic api response error
type apiError struct {
	Error string `json:"error"`
}
