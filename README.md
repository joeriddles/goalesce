# Goalesce

`goalesce` is a command-line tool to generate OpenAPI CRUD routes from GORM models.

---
![Test workflow](https://github.com/joeriddles/goalesce/actions/workflows/test.yaml/badge.svg) ![Release workflow](https://github.com/joeriddles/goalesce/actions/workflows/release.yaml/badge.svg)

## Features
- Generate [OpenAPI YAML](https://swagger.io/specification/) files from [GORM](https://gorm.io/) model types
- Generate CRUD paths for each GORM model
- Generate controllers, mappers, and repositories for each GORM model
- `goalesce` uses [`oapi-codegen`](https://github.com/oapi-codegen/oapi-codegen/) for generating server and controller interfaces

## Install
```shell
$ go install github.com/joeriddles/goalesce/cmd/goalesce@latest
```

Goalesce requires additional tooling to work correctly. 
- Ensure a modern version of [Node](https://nodejs.org) is installed (LTS+).
- Ensure the `goimports` tool is also installed:
	- `go install golang.org/x/tools/cmd/goimports@latest`

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

## Releasing

To release a new change, simply run `rev-tag.sh` with the desired [semver](https://semver.org/) update: major, minor, or patch:

```shell
$ ./rev-tag.sh patch
./rev-tag.sh patch
Old tag: v1.0.1
New tag: v1.0.2
Created new tag v1.0.2
Do you want to push this tag? y
Enumerating objects: 4, done.
Counting objects: 100% (4/4), done.
Delta compression using up to 8 threads
Compressing objects: 100% (3/3), done.
Writing objects: 100% (3/3), 769 bytes | 769.00 KiB/s, done.
Total 3 (delta 1), reused 0 (delta 0), pack-reused 0 (from 0)
remote: Resolving deltas: 100% (1/1), completed with 1 local object.
To https://github.com/joeriddles/goalesce.git
 * [new tag]         v1.0.2 -> v1.0.2
```

The shell script will create a new tag and push it to GitHub, which the [release](https://github.com/joeriddles/goalesce/blob/main/.github/workflows/release.yaml) workflow will detect and create a new release for.
