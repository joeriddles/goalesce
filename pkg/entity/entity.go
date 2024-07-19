package entity

type GormModelMetadata struct {
	Name   string
	Fields []*GormModelField
}

type GormModelField struct {
	Name string
	Type string
	Tag  *string
}
