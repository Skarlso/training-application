# Training Application

This application is used in the [cloud-native training catalogue](https://www.cloud-native.com/trainings/). The images can be found in [quay](https://quay.io/repository/kubermatic-labs/training-application?tab=tags).

The aim of this application is to be used in the context of trainings for Docker, Kubernetes, Istio, Helm, ...

## Functionality of the application

**Nothing**, besides possibly (depending on the application configuration) showing a cute cat image 🙀.

## Available Endpoints

> The application offers the following endpoints on **port 8080**

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

> The application offers the following commands **via stdin**

TODO tty in docker and k8s

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

## Configuring the application

### `configFilePath`

- **Description**:
- **Type**:
- **Default Value**:
- **Usage**:
  | Path to the config file | string | "./training-application.conf" | | application arg, you can set this via eg `./training-application --configFilePath my.conf` |

### `alive`

- **Description**:
- **Type**:
- **Default Value**:
- **Usage**:

| Flag to indicate the applications liveness | bool | true | | configurable via the commands `set alive` and `set dead` |

### `ready`

- **Description**:
- **Type**:
- **Default Value**:
- **Usage**:

| Flag to indicate the applications readiness | bool | false | | true after `startUpDelaySeconds` , false on graceful shutdown; configurable via the commands `set ready` and `set unready` |

### `name`

- **Description**:
- **Type**:
- **Default Value**:
- **Usage**:

| The name of the application | string | "not set" | | via config file or via the environment variable `APP_NAME` |

### `version`

- **Description**:
- **Type**:
- **Default Value**:
- **Usage**:

| The version of the application | string | "not set" | | via config file or via the environment variable `APP_VERSION` |

### `message`

- **Description**:
- **Type**:
- **Default Value**:
- **Usage**:

| A message to be shown on the root endpoint | string | "not set" | | via config file or via the environment variable `APP_MESSAGE` |

### `color`

- **Description**:
- **Type**:
- **Default Value**:
- **Usage**:

| The background color of the root endpoint | string | "not set" | | via config file or via the environment variable `APP_COLOR` |

### `rootDelaySeconds`

- **Description**:
- **Type**:
- **Default Value**:
- **Usage**:

| For delaying requests to the root endpoint | int | 0 | | via config file or via the environment variable `APP_ROOT_DELAY_SECONDS` |

### `startUpDelaySeconds`

- **Description**:
- **Type**:
- **Default Value**:
- **Usage**:

|Time the application will take to start | int | 0 | | via config file or via the environment variable `APP_START_UP_DELAY_SECONDS` |

### `tearDownDelaySeconds`

- **Description**:
- **Type**:
- **Default Value**:
- **Usage**:

| Time the application will take to gracefully shut down | int | 0 | | via config file or via the environment variable `APP_TEAR_DOWN_DELAY_SECONDS` |

### `logToFileOnly`

- **Description**:
- **Type**:
- **Default Value**:
- **Usage**:

| Log **only** to the file `training-application.log`, if set to true no logging to stdout will happen | bool | false | | via config file |

### `catMode`

- **Description**:
- **Type**:
- **Default Value**:
- **Usage**:

| Flag to get cute cat images in the rood endpoint | bool | false | | via config file |

## Running the application

## Building the application
