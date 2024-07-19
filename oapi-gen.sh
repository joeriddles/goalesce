set -eux
rm -rf ./generated/**
go run . examples/basic

cd ./generated
npx @redocly/openapi-cli@latest bundle openapi_base.gen.yaml > openapi.yaml

oapi-codegen --config ../types.cfg.yaml openapi.yaml
oapi-codegen --config ../server.cfg.yaml openapi.yaml
