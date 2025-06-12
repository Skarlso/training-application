# todos

## verify leak cpu method

## add timestamp when the process was started

## add ip info to application

## make root.html configurable (parse on runtime) for CF

## pass in version into docker build?

## write unit tests to ensure training usecases work (KFD, KSM)

## k8s meta info without downward api

- serviceaccount
- nodeName
- podNamespace
- podName
- podIP

## linting in github action

## packaging

- linux service (for LF training)
- compose
- k8s
- helm

## bugs

### leak cpu

sometimes results in a container restart
=> does log.Errorf("error on opening /dev/null: %s", err) at least give info why?
