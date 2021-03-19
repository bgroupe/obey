<img src="./ui/src/assets/obey-logo.svg" height="520px" width="520px"/>

# Environment Discovery Service

### Goal
### General Overview

### Scheduler Overview
### Worker Overview


### TODO:
---
- Finish streaming output of a job
- Change the insecure dialing for grpc

[x] Delete worker from Redis store

[] Database functions

[x] Serialize store to marshaled json

[] Create shared models

[x] Surface errors for deregistering workers

[] Use redis.Pool

[] remove workers.Mutex()

[] potentially look into redlock https://github.com/go-redsync/redsync

[] change marshaling function to account for service poll
[] move module functions out of worker/server
    [] separate grpc requests
    [] separate config
