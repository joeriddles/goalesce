module github.com/joeriddles/goalesce/examples/echo

go 1.20

replace github.com/joeriddles/goalesce => ../..

require gorm.io/gorm v1.25.10

require (
	github.com/jinzhu/inflection v1.0.0 // indirect
	github.com/jinzhu/now v1.1.5 // indirect
)
