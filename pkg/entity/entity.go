package entity

type GormModelMetadata struct {
	Name                string
	IsGormModelEmbedded bool
	Fields              []*GormModelField
}

type GormModelField struct {
	Name                string
	Type                string
	IsGormModelEmbedded bool
	Tag                 *string
}
