# Training Application

This application is used in the [cloud-native training catalogue](https://www.cloud-native.com/trainings/). The images can be found in [quay](https://quay.io/repository/kubermatic-labs/training-application?tab=tags).

## Available Endpoints

| Endpoint     | Description                                                        |
| ------------ | ------------------------------------------------------------------ |
| `/`          | Root endpoint, the output depends on the application configuration |
| `/liveness`  | Liveness probe                                                     |
| `/readiness` | Readiness probe                                                    |

## Available Commands

| Command             | Description                                                         |
| ------------------- | ------------------------------------------------------------------- |
| `help`              | Get info about available commands and endpoints                     |
| `init`              | Set readiness true, liveness true and delay 0                       |
| `config`            | Print out the current application configuration                     |
| `set ready`         | Application readiness probe will be successful                      |
| `set unready`       | Application readiness probe will fail                               |
| `set alive`         | Application liveness probe will be successful                       |
| `set dead`          | Application liveness probe will fail                                |
| `leak mem`          | Leak memory                                                         |
| `leak cpu`          | Leak CPU                                                            |
| `request <url>`     | Request a URL, e.g., `request https://www.google.com`               |
| `delay / <seconds>` | Set delay for the root endpoint (`/`) in seconds, e.g., `delay / 5` |
