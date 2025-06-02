github action => is it save to put quay credentials in public repo

write unit tests to ensure training usecases work (KFD, KSM)

# leak cpu

sometimes results in a container restart
=> does log.Errorf("Error on opening /dev/null: %s", err) at least give info why?

# istio

recheck if istio really calls the readiness endpoint => WTF?!?!?!
liveness and readiness probes logging is info too much => go smaller?
make log level customizable
do I need the caller? eg if the request is coming from k8s or istio?

# TODO helm chart

# TODO add configuration possibilities into readme

- app.conf
- env vars

# loggging issue

- not to stdout only in some file

# pass in version into docker build?

## Training Application

show pid in root
sigterm state in root
