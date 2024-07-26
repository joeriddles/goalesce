# Goalesce

`goalesce` is a command-line tool to generate OpenAPI CRUD routes from GORM models.

## Install
```shell
$ go install github.com/joeriddles/goalesce/cmd/goalesce@latest
```

## Usage
`goalesce` is largely configured using a YAML configuration file.

Example
```shell
$ goalesce ./pkg/model
```
