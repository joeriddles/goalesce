set -eux
go run . examples/basic
cd generated
oapi-codegen --config ../types.cfg.yaml ./openapi.yaml
