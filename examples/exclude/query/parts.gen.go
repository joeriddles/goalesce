// Code generated by gorm.io/gen. DO NOT EDIT.
// Code generated by gorm.io/gen. DO NOT EDIT.
// Code generated by gorm.io/gen. DO NOT EDIT.

package query

import (
	"context"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"
	"gorm.io/gorm/schema"

	"gorm.io/gen"
	"gorm.io/gen/field"

	"gorm.io/plugin/dbresolver"

	"github.com/joeriddles/goalesce/examples/exclude/model"
)

func newPart(db *gorm.DB, opts ...gen.DOOption) part {
	_part := part{}

	_part.partDo.UseDB(db, opts...)
	_part.partDo.UseModel(&model.Part{})

	tableName := _part.partDo.TableName()
	_part.ALL = field.NewAsterisk(tableName)
	_part.ID = field.NewUint(tableName, "id")
	_part.CreatedAt = field.NewTime(tableName, "created_at")
	_part.UpdatedAt = field.NewTime(tableName, "updated_at")
	_part.DeletedAt = field.NewField(tableName, "deleted_at")
	_part.Name = field.NewString(tableName, "name")
	_part.Cost = field.NewInt(tableName, "cost")
	_part.Models = partManyToManyModels{
		db: db.Session(&gorm.Session{}),

		RelationField: field.NewRelation("Models", "model.VehicleModel"),
		Manufacturer: struct {
			field.RelationField
			Vehicles struct {
				field.RelationField
			}
		}{
			RelationField: field.NewRelation("Models.Manufacturer", "model.Manufacturer"),
			Vehicles: struct {
				field.RelationField
			}{
				RelationField: field.NewRelation("Models.Manufacturer.Vehicles", "model.VehicleModel"),
			},
		},
		Parts: struct {
			field.RelationField
			Models struct {
				field.RelationField
			}
		}{
			RelationField: field.NewRelation("Models.Parts", "model.Part"),
			Models: struct {
				field.RelationField
			}{
				RelationField: field.NewRelation("Models.Parts.Models", "model.VehicleModel"),
			},
		},
	}

	_part.fillFieldMap()

	return _part
}

type part struct {
	partDo

	ALL       field.Asterisk
	ID        field.Uint
	CreatedAt field.Time
	UpdatedAt field.Time
	DeletedAt field.Field
	Name      field.String
	Cost      field.Int
	Models    partManyToManyModels

	fieldMap map[string]field.Expr
}

func (p part) Table(newTableName string) *part {
	p.partDo.UseTable(newTableName)
	return p.updateTableName(newTableName)
}

func (p part) As(alias string) *part {
	p.partDo.DO = *(p.partDo.As(alias).(*gen.DO))
	return p.updateTableName(alias)
}

func (p *part) updateTableName(table string) *part {
	p.ALL = field.NewAsterisk(table)
	p.ID = field.NewUint(table, "id")
	p.CreatedAt = field.NewTime(table, "created_at")
	p.UpdatedAt = field.NewTime(table, "updated_at")
	p.DeletedAt = field.NewField(table, "deleted_at")
	p.Name = field.NewString(table, "name")
	p.Cost = field.NewInt(table, "cost")

	p.fillFieldMap()

	return p
}

func (p *part) GetFieldByName(fieldName string) (field.OrderExpr, bool) {
	_f, ok := p.fieldMap[fieldName]
	if !ok || _f == nil {
		return nil, false
	}
	_oe, ok := _f.(field.OrderExpr)
	return _oe, ok
}

func (p *part) fillFieldMap() {
	p.fieldMap = make(map[string]field.Expr, 7)
	p.fieldMap["id"] = p.ID
	p.fieldMap["created_at"] = p.CreatedAt
	p.fieldMap["updated_at"] = p.UpdatedAt
	p.fieldMap["deleted_at"] = p.DeletedAt
	p.fieldMap["name"] = p.Name
	p.fieldMap["cost"] = p.Cost

}

func (p part) clone(db *gorm.DB) part {
	p.partDo.ReplaceConnPool(db.Statement.ConnPool)
	return p
}

func (p part) replaceDB(db *gorm.DB) part {
	p.partDo.ReplaceDB(db)
	return p
}

type partManyToManyModels struct {
	db *gorm.DB

	field.RelationField

	Manufacturer struct {
		field.RelationField
		Vehicles struct {
			field.RelationField
		}
	}
	Parts struct {
		field.RelationField
		Models struct {
			field.RelationField
		}
	}
}

func (a partManyToManyModels) Where(conds ...field.Expr) *partManyToManyModels {
	if len(conds) == 0 {
		return &a
	}

	exprs := make([]clause.Expression, 0, len(conds))
	for _, cond := range conds {
		exprs = append(exprs, cond.BeCond().(clause.Expression))
	}
	a.db = a.db.Clauses(clause.Where{Exprs: exprs})
	return &a
}

func (a partManyToManyModels) WithContext(ctx context.Context) *partManyToManyModels {
	a.db = a.db.WithContext(ctx)
	return &a
}

func (a partManyToManyModels) Session(session *gorm.Session) *partManyToManyModels {
	a.db = a.db.Session(session)
	return &a
}

func (a partManyToManyModels) Model(m *model.Part) *partManyToManyModelsTx {
	return &partManyToManyModelsTx{a.db.Model(m).Association(a.Name())}
}

type partManyToManyModelsTx struct{ tx *gorm.Association }

func (a partManyToManyModelsTx) Find() (result []*model.VehicleModel, err error) {
	return result, a.tx.Find(&result)
}

func (a partManyToManyModelsTx) Append(values ...*model.VehicleModel) (err error) {
	targetValues := make([]interface{}, len(values))
	for i, v := range values {
		targetValues[i] = v
	}
	return a.tx.Append(targetValues...)
}

func (a partManyToManyModelsTx) Replace(values ...*model.VehicleModel) (err error) {
	targetValues := make([]interface{}, len(values))
	for i, v := range values {
		targetValues[i] = v
	}
	return a.tx.Replace(targetValues...)
}

func (a partManyToManyModelsTx) Delete(values ...*model.VehicleModel) (err error) {
	targetValues := make([]interface{}, len(values))
	for i, v := range values {
		targetValues[i] = v
	}
	return a.tx.Delete(targetValues...)
}

func (a partManyToManyModelsTx) Clear() error {
	return a.tx.Clear()
}

func (a partManyToManyModelsTx) Count() int64 {
	return a.tx.Count()
}

type partDo struct{ gen.DO }

type IPartDo interface {
	gen.SubQuery
	Debug() IPartDo
	WithContext(ctx context.Context) IPartDo
	WithResult(fc func(tx gen.Dao)) gen.ResultInfo
	ReplaceDB(db *gorm.DB)
	ReadDB() IPartDo
	WriteDB() IPartDo
	As(alias string) gen.Dao
	Session(config *gorm.Session) IPartDo
	Columns(cols ...field.Expr) gen.Columns
	Clauses(conds ...clause.Expression) IPartDo
	Not(conds ...gen.Condition) IPartDo
	Or(conds ...gen.Condition) IPartDo
	Select(conds ...field.Expr) IPartDo
	Where(conds ...gen.Condition) IPartDo
	Order(conds ...field.Expr) IPartDo
	Distinct(cols ...field.Expr) IPartDo
	Omit(cols ...field.Expr) IPartDo
	Join(table schema.Tabler, on ...field.Expr) IPartDo
	LeftJoin(table schema.Tabler, on ...field.Expr) IPartDo
	RightJoin(table schema.Tabler, on ...field.Expr) IPartDo
	Group(cols ...field.Expr) IPartDo
	Having(conds ...gen.Condition) IPartDo
	Limit(limit int) IPartDo
	Offset(offset int) IPartDo
	Count() (count int64, err error)
	Scopes(funcs ...func(gen.Dao) gen.Dao) IPartDo
	Unscoped() IPartDo
	Create(values ...*model.Part) error
	CreateInBatches(values []*model.Part, batchSize int) error
	Save(values ...*model.Part) error
	First() (*model.Part, error)
	Take() (*model.Part, error)
	Last() (*model.Part, error)
	Find() ([]*model.Part, error)
	FindInBatch(batchSize int, fc func(tx gen.Dao, batch int) error) (results []*model.Part, err error)
	FindInBatches(result *[]*model.Part, batchSize int, fc func(tx gen.Dao, batch int) error) error
	Pluck(column field.Expr, dest interface{}) error
	Delete(...*model.Part) (info gen.ResultInfo, err error)
	Update(column field.Expr, value interface{}) (info gen.ResultInfo, err error)
	UpdateSimple(columns ...field.AssignExpr) (info gen.ResultInfo, err error)
	Updates(value interface{}) (info gen.ResultInfo, err error)
	UpdateColumn(column field.Expr, value interface{}) (info gen.ResultInfo, err error)
	UpdateColumnSimple(columns ...field.AssignExpr) (info gen.ResultInfo, err error)
	UpdateColumns(value interface{}) (info gen.ResultInfo, err error)
	UpdateFrom(q gen.SubQuery) gen.Dao
	Attrs(attrs ...field.AssignExpr) IPartDo
	Assign(attrs ...field.AssignExpr) IPartDo
	Joins(fields ...field.RelationField) IPartDo
	Preload(fields ...field.RelationField) IPartDo
	FirstOrInit() (*model.Part, error)
	FirstOrCreate() (*model.Part, error)
	FindByPage(offset int, limit int) (result []*model.Part, count int64, err error)
	ScanByPage(result interface{}, offset int, limit int) (count int64, err error)
	Scan(result interface{}) (err error)
	Returning(value interface{}, columns ...string) IPartDo
	UnderlyingDB() *gorm.DB
	schema.Tabler
}

func (p partDo) Debug() IPartDo {
	return p.withDO(p.DO.Debug())
}

func (p partDo) WithContext(ctx context.Context) IPartDo {
	return p.withDO(p.DO.WithContext(ctx))
}

func (p partDo) ReadDB() IPartDo {
	return p.Clauses(dbresolver.Read)
}

func (p partDo) WriteDB() IPartDo {
	return p.Clauses(dbresolver.Write)
}

func (p partDo) Session(config *gorm.Session) IPartDo {
	return p.withDO(p.DO.Session(config))
}

func (p partDo) Clauses(conds ...clause.Expression) IPartDo {
	return p.withDO(p.DO.Clauses(conds...))
}

func (p partDo) Returning(value interface{}, columns ...string) IPartDo {
	return p.withDO(p.DO.Returning(value, columns...))
}

func (p partDo) Not(conds ...gen.Condition) IPartDo {
	return p.withDO(p.DO.Not(conds...))
}

func (p partDo) Or(conds ...gen.Condition) IPartDo {
	return p.withDO(p.DO.Or(conds...))
}

func (p partDo) Select(conds ...field.Expr) IPartDo {
	return p.withDO(p.DO.Select(conds...))
}

func (p partDo) Where(conds ...gen.Condition) IPartDo {
	return p.withDO(p.DO.Where(conds...))
}

func (p partDo) Order(conds ...field.Expr) IPartDo {
	return p.withDO(p.DO.Order(conds...))
}

func (p partDo) Distinct(cols ...field.Expr) IPartDo {
	return p.withDO(p.DO.Distinct(cols...))
}

func (p partDo) Omit(cols ...field.Expr) IPartDo {
	return p.withDO(p.DO.Omit(cols...))
}

func (p partDo) Join(table schema.Tabler, on ...field.Expr) IPartDo {
	return p.withDO(p.DO.Join(table, on...))
}

func (p partDo) LeftJoin(table schema.Tabler, on ...field.Expr) IPartDo {
	return p.withDO(p.DO.LeftJoin(table, on...))
}

func (p partDo) RightJoin(table schema.Tabler, on ...field.Expr) IPartDo {
	return p.withDO(p.DO.RightJoin(table, on...))
}

func (p partDo) Group(cols ...field.Expr) IPartDo {
	return p.withDO(p.DO.Group(cols...))
}

func (p partDo) Having(conds ...gen.Condition) IPartDo {
	return p.withDO(p.DO.Having(conds...))
}

func (p partDo) Limit(limit int) IPartDo {
	return p.withDO(p.DO.Limit(limit))
}

func (p partDo) Offset(offset int) IPartDo {
	return p.withDO(p.DO.Offset(offset))
}

func (p partDo) Scopes(funcs ...func(gen.Dao) gen.Dao) IPartDo {
	return p.withDO(p.DO.Scopes(funcs...))
}

func (p partDo) Unscoped() IPartDo {
	return p.withDO(p.DO.Unscoped())
}

func (p partDo) Create(values ...*model.Part) error {
	if len(values) == 0 {
		return nil
	}
	return p.DO.Create(values)
}

func (p partDo) CreateInBatches(values []*model.Part, batchSize int) error {
	return p.DO.CreateInBatches(values, batchSize)
}

// Save : !!! underlying implementation is different with GORM
// The method is equivalent to executing the statement: db.Clauses(clause.OnConflict{UpdateAll: true}).Create(values)
func (p partDo) Save(values ...*model.Part) error {
	if len(values) == 0 {
		return nil
	}
	return p.DO.Save(values)
}

func (p partDo) First() (*model.Part, error) {
	if result, err := p.DO.First(); err != nil {
		return nil, err
	} else {
		return result.(*model.Part), nil
	}
}

func (p partDo) Take() (*model.Part, error) {
	if result, err := p.DO.Take(); err != nil {
		return nil, err
	} else {
		return result.(*model.Part), nil
	}
}

func (p partDo) Last() (*model.Part, error) {
	if result, err := p.DO.Last(); err != nil {
		return nil, err
	} else {
		return result.(*model.Part), nil
	}
}

func (p partDo) Find() ([]*model.Part, error) {
	result, err := p.DO.Find()
	return result.([]*model.Part), err
}

func (p partDo) FindInBatch(batchSize int, fc func(tx gen.Dao, batch int) error) (results []*model.Part, err error) {
	buf := make([]*model.Part, 0, batchSize)
	err = p.DO.FindInBatches(&buf, batchSize, func(tx gen.Dao, batch int) error {
		defer func() { results = append(results, buf...) }()
		return fc(tx, batch)
	})
	return results, err
}

func (p partDo) FindInBatches(result *[]*model.Part, batchSize int, fc func(tx gen.Dao, batch int) error) error {
	return p.DO.FindInBatches(result, batchSize, fc)
}

func (p partDo) Attrs(attrs ...field.AssignExpr) IPartDo {
	return p.withDO(p.DO.Attrs(attrs...))
}

func (p partDo) Assign(attrs ...field.AssignExpr) IPartDo {
	return p.withDO(p.DO.Assign(attrs...))
}

func (p partDo) Joins(fields ...field.RelationField) IPartDo {
	for _, _f := range fields {
		p = *p.withDO(p.DO.Joins(_f))
	}
	return &p
}

func (p partDo) Preload(fields ...field.RelationField) IPartDo {
	for _, _f := range fields {
		p = *p.withDO(p.DO.Preload(_f))
	}
	return &p
}

func (p partDo) FirstOrInit() (*model.Part, error) {
	if result, err := p.DO.FirstOrInit(); err != nil {
		return nil, err
	} else {
		return result.(*model.Part), nil
	}
}

func (p partDo) FirstOrCreate() (*model.Part, error) {
	if result, err := p.DO.FirstOrCreate(); err != nil {
		return nil, err
	} else {
		return result.(*model.Part), nil
	}
}

func (p partDo) FindByPage(offset int, limit int) (result []*model.Part, count int64, err error) {
	result, err = p.Offset(offset).Limit(limit).Find()
	if err != nil {
		return
	}

	if size := len(result); 0 < limit && 0 < size && size < limit {
		count = int64(size + offset)
		return
	}

	count, err = p.Offset(-1).Limit(-1).Count()
	return
}

func (p partDo) ScanByPage(result interface{}, offset int, limit int) (count int64, err error) {
	count, err = p.Count()
	if err != nil {
		return
	}

	err = p.Offset(offset).Limit(limit).Scan(result)
	return
}

func (p partDo) Scan(result interface{}) (err error) {
	return p.DO.Scan(result)
}

func (p partDo) Delete(models ...*model.Part) (result gen.ResultInfo, err error) {
	return p.DO.Delete(models)
}

func (p *partDo) withDO(do gen.Dao) *partDo {
	p.DO = *do.(*gen.DO)
	return p
}
