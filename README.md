# `gorm-oapi-codegen`

`gorm-oapi-codegen` is a command-line tool to generate OpenAPI CRUD routes from GORM models.

## Install
```shell
$ go install github.com/joeriddles/gorm-oapi-codegen/cmd/gorm-oapi-codegen@latest
```

## Usage
`gorm-oapi-codegen` is largely configured using a YAML configuration file.

Example
```shell
$ gorm-oapi-codegen ./pkg/model
```
