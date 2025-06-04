# todos

## later

- pass in version into docker build?
- github action => is it save to put quay credentials in public repo
- write unit tests to ensure training usecases work (KFD, KSM)

## lint

### go code

- fix remaining linter issues and add dep from build step in makefile

### docker build

- fix remaining linter issues and add dep from build step in makefile

## bugs

### leak cpu

sometimes results in a container restart
=> does log.Errorf("Error on opening /dev/null: %s", err) at least give info why?

### istio

recheck if istio really calls the readiness endpoint => WTF?!?!?!
liveness and readiness probes logging is info too much => go smaller?
make log level customizable
do I need the caller? eg if the request is coming from k8s or istio?

## now

### TODO add configuration possibilities into readme

- app.conf
- env vars
- bin params
