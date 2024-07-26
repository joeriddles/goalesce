# Goalesce

`goalesce` is a command-line tool to generate OpenAPI CRUD routes from GORM models.

## Features
- Generate [OpenAPI YAML](https://swagger.io/specification/) files from [GORM](https://gorm.io/) model types

## Install
```shell
$ go install github.com/joeriddles/goalesce/cmd/goalesce@latest
```

## Usage
`goalesce` is largely configured using a YAML configuration file. Check out the GoDoc for [`Config`](https://pkg.go.dev/github.com/joeriddles/goalesce/pkg/config#Config) for more detail.


Example
```shell
$ goalesce ./pkg/model
```
