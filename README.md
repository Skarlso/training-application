# Training Application

This application is used in the [cloud-native training catalogue](https://www.cloud-native.com/trainings/). The images can be found in [quay](https://quay.io/repository/kubermatic-labs/training-application?tab=tags).

The aim of this application is to be used in the context of trainings for Docker, Kubernetes, Istio, Helm, ...

## Functionality of the application

**Nothing**, besides possibly (depending on the application configuration) showing a cute cat image 🙀.

<img width="500" alt="Screenshot 2025-06-04 at 15 40 30" src="https://github.com/user-attachments/assets/7cafc452-4f21-4202-8379-2a983ecbf122" />


## Available Endpoints

> **_NOTE:_** The application offers the following endpoints on port **8080**

### `/`

Root endpoint, the output depends on the application configuration.

Eg the response could be delayed for a configurable amount of seconds.

### `/liveness`

Endpoint of the application to signal if the application is in a healthy state.

If everything is fine the application will respond with a 200 status code, if not the application should respond with a 503 status code.

### `/readiness`

Endpoint of the application to signal if the application is ready to receive requests.

If everything is fine the application will respond with a 200 status code, if not the application should respond with a 503 status code.

## Available Commands

> **_NOTE:_** The application offers the following commands **via stdin**

| Command             | Description                                                         |
| ------------------- | ------------------------------------------------------------------- |
| `help`              | Get info about available commands and endpoints                     |
| `init`              | Re-initialize the application                                       |
| `config`            | Print out the current application configuration                     |
| `set ready`         | Application readiness probe will be successful                      |
| `set unready`       | Application readiness probe will fail                               |
| `set alive`         | Application liveness probe will be successful                       |
| `set dead`          | Application liveness probe will fail                                |
| `leak mem`          | Leak memory                                                         |
| `leak cpu`          | Leak CPU                                                            |
| `request <url>`     | Request a URL, e.g., `request https://www.kubermatic.com/`          |
| `delay / <seconds>` | Set delay for the root endpoint (`/`) in seconds, e.g., `delay / 5` |

> **_INSIDE A CONTAINER_** If you want to send commands to the application you have to use of `docker attach my-training-application-container`. The container als has to have `tty` enabled.

> **_INSIDE A POD_** If you want to send commands to the application you have have to use of `kubectl attach -it training-application-pod` and have the following flags set in the pod manifest:

```yaml
apiVersion: v1
kind: Pod
metadata:
  name: training-application
spec:
  containers:
    - name: training-application
      image: quay.io/kubermatic-labs/training-application:3.0.0
      imagePullPolicy: Always
      tty: true # <= add those flags
      stdin: true # <= add those flags
      ports:
        - name: http
          containerPort: 8080
```

## Configuring the application

### `configFilePath`

- **Description**: Path to the config file
- **Type**: string
- **Default Value**: "./training-application.conf"
- **Usage**: application arg, you can set this via eg `./training-application --configFilePath my.conf`

### `alive`

- **Description**: Flag to indicate the applications liveness
- **Type**: bool
- **Default Value**: true
- **Usage**: configurable via the commands `set alive` and `set dead`

### `ready`

- **Description**: Flag to indicate the applications readiness
- **Type**: bool
- **Default Value**: false
- **Usage**: true after `startUpDelaySeconds` , false on graceful shutdown; configurable via the commands `set ready` and `set unready`

### `name`

- **Description**: The name of the application
- **Type**: string
- **Default Value**: "not set"
- **Usage**: via config file or via the environment variable `APP_NAME`

### `version`

- **Description**: The version of the application
- **Type**: string
- **Default Value**: "not set"
- **Usage**: via config file or via the environment variable `APP_VERSION`

### `message`

- **Description**: A message to be shown on the root endpoint
- **Type**: string
- **Default Value**: "not set
- **Usage**: via config file or via the environment variable `APP_MESSAGE`

### `color`

- **Description**: The background color of the root endpoint
- **Type**: string
- **Default Value**: "not set"
- **Usage**: via config file or via the environment variable `APP_COLOR`

### `rootDelaySeconds`

- **Description**: For delaying requests to the root endpoint
- **Type**: int
- **Default Value**: 0
- **Usage**: via config file or via the command `delay / <seconds>`, eg "delay / 10"

### `startUpDelaySeconds`

- **Description**: Time the application will take to start
- **Type**: int
- **Default Value**: 0
- **Usage**: via config file

### `tearDownDelaySeconds`

- **Description**: Time the application will take to gracefully shut down
- **Type**: int
- **Default Value**: 0
- **Usage**: via config file

### `logToFileOnly`

- **Description**: Log **only** to the file named `training-application.log`, if set to true no logging to stdout will happen
- **Type**: bool
- **Default Value**: false
- **Usage**: via config file

### `catMode`

- **Description**: Flag to get cute cat images in the root endpoint
- **Type**: bool
- **Default Value**: false
- **Usage**: via config file

## Building the application

```bash
make build
```

## Running the application natively

```bash
make run
```

## Building the Image

```bash
make docker-build
```

## Running the Image

```bash
make docker-run
```
