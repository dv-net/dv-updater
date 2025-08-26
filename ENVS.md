# Environments

| Name                        | Required | Secret | Default value | Usage                                | Example        |
|-----------------------------|----------|--------|---------------|--------------------------------------|----------------|
| `APP_PROFILE`               |          |        | `dev`         |                                      |                |
| `HTTP_HOST`                 |          |        | `localhost`   |                                      |                |
| `HTTP_PORT`                 |          |        | `8081`        |                                      |                |
| `HTTP_FETCH_INTERVAL`       |          |        | `30s`         |                                      |                |
| `HTTP_CONNECT_TIMEOUT`      |          |        | `5s`          |                                      |                |
| `HTTP_READ_TIMEOUT`         |          |        | `10s`         |                                      |                |
| `HTTP_WRITE_TIMEOUT`        |          |        | `10s`         |                                      |                |
| `HTTP_MAX_HEADER_MEGABYTES` |          |        | `1`           |                                      |                |
| `HTTP_CORS_ENABLED`         |          |        | `true`        | allows to disable cors               | `true / false` |
| `HTTP_CORS_ALLOWED_ORIGINS` |          |        |               |                                      |                |
| `LOG_FORMAT`                |          |        | `json`        | allows to set custom formatting      | `json`         |
| `LOG_LEVEL`                 |          |        | `info`        | allows to set custom logger level    | `info`         |
| `LOG_CONSOLE_COLORED`       |          |        | `false`       | allows to set colored console output | `false`        |
| `LOG_TRACE`                 |          |        | `fatal`       | allows to set custom trace level     | `fatal`        |
| `LOG_WITH_CALLER`           |          |        | `false`       | allows to show caller                | `false`        |
| `LOG_WITH_STACK_TRACE`      |          |        | `false`       | allows to show stack trace           | `false`        |
| `AUTO_UPDATE_ENABLED`       |          |        | `true`        |                                      |                |
