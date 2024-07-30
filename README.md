# Goalesce

`goalesce` is a command-line tool to generate OpenAPI CRUD routes from GORM models.


## Features
- Generate [OpenAPI YAML](https://swagger.io/specification/) files from [GORM](https://gorm.io/) model types
- Generate CRUD paths for each GORM model
- Generate controllers, mappers, and repositories for each GORM model
- `goalesce` uses [`oapi-codegen`](https://github.com/oapi-codegen/oapi-codegen/) for generating server and controller interfaces

## Install
```shell
$ go install github.com/joeriddles/goalesce/cmd/goalesce@latest
```

## Usage
`goalesce` is largely configured using a YAML configuration file. Check out the GoDoc for [`Config`](https://pkg.go.dev/github.com/joeriddles/goalesce/pkg/config#Config) for more detail.


Example:
```yaml
# ./config.yaml
input_folder_path: ./model
output_file_path: ./generated
module_name: github.com/joeriddles/goalesce/examples/basic
models_package: github.com/joeriddles/goalesce/examples/basic/model
query_package: github.com/joeriddles/goalesce/examples/basic/query
clear_output_dir: true
```

```go
// ./model/model.go
package model

import "gorm.io/gorm"

type User struct {
	gorm.Model
	Name string `gorm:"column:name;"`
}
```

```go
// ./main.go
package main

import (
	"github.com/joeriddles/goalesce/examples/basic/model"
	"gorm.io/gen"
)

func main() {
	g := gen.NewGenerator(gen.Config{
		OutPath: "query",
		Mode:    gen.WithoutContext | gen.WithQueryInterface,
	})
	g.ApplyBasic(model.User{})
	g.Execute()
}
```

```shell
# Run GORM gen
$ go run .
# Run Goalesce gen
$ goalesce -config config.yaml
```
