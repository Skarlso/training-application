# todos

## later

- pass in version into docker build?
- github action => is it save to put quay credentials in public repo
- write unit tests to ensure training usecases work (KFD, KSM)

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

### naming of the thing

- app vs application vs training-application
- conf/app.conf... really? => naming... folder...

### TEST

- config file
- env vars
