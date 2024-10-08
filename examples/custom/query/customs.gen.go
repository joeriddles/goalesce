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

	"github.com/joeriddles/goalesce/examples/custom/model"
)

func newCustom(db *gorm.DB, opts ...gen.DOOption) custom {
	_custom := custom{}

	_custom.customDo.UseDB(db, opts...)
	_custom.customDo.UseModel(&model.Custom{})

	tableName := _custom.customDo.TableName()
	_custom.ALL = field.NewAsterisk(tableName)
	_custom.ID = field.NewInt64(tableName, "id")
	_custom.CreatedAt = field.NewTime(tableName, "created_at")
	_custom.UpdatedAt = field.NewTime(tableName, "updated_at")
	_custom.DeletedAt = field.NewField(tableName, "deleted_at")
	_custom.Name = field.NewString(tableName, "name")

	_custom.fillFieldMap()

	return _custom
}

type custom struct {
	customDo

	ALL       field.Asterisk
	ID        field.Int64
	CreatedAt field.Time
	UpdatedAt field.Time
	DeletedAt field.Field
	Name      field.String

	fieldMap map[string]field.Expr
}

func (c custom) Table(newTableName string) *custom {
	c.customDo.UseTable(newTableName)
	return c.updateTableName(newTableName)
}

func (c custom) As(alias string) *custom {
	c.customDo.DO = *(c.customDo.As(alias).(*gen.DO))
	return c.updateTableName(alias)
}

func (c *custom) updateTableName(table string) *custom {
	c.ALL = field.NewAsterisk(table)
	c.ID = field.NewInt64(table, "id")
	c.CreatedAt = field.NewTime(table, "created_at")
	c.UpdatedAt = field.NewTime(table, "updated_at")
	c.DeletedAt = field.NewField(table, "deleted_at")
	c.Name = field.NewString(table, "name")

	c.fillFieldMap()

	return c
}

func (c *custom) GetFieldByName(fieldName string) (field.OrderExpr, bool) {
	_f, ok := c.fieldMap[fieldName]
	if !ok || _f == nil {
		return nil, false
	}
	_oe, ok := _f.(field.OrderExpr)
	return _oe, ok
}

func (c *custom) fillFieldMap() {
	c.fieldMap = make(map[string]field.Expr, 5)
	c.fieldMap["id"] = c.ID
	c.fieldMap["created_at"] = c.CreatedAt
	c.fieldMap["updated_at"] = c.UpdatedAt
	c.fieldMap["deleted_at"] = c.DeletedAt
	c.fieldMap["name"] = c.Name
}

func (c custom) clone(db *gorm.DB) custom {
	c.customDo.ReplaceConnPool(db.Statement.ConnPool)
	return c
}

func (c custom) replaceDB(db *gorm.DB) custom {
	c.customDo.ReplaceDB(db)
	return c
}

type customDo struct{ gen.DO }

type ICustomDo interface {
	gen.SubQuery
	Debug() ICustomDo
	WithContext(ctx context.Context) ICustomDo
	WithResult(fc func(tx gen.Dao)) gen.ResultInfo
	ReplaceDB(db *gorm.DB)
	ReadDB() ICustomDo
	WriteDB() ICustomDo
	As(alias string) gen.Dao
	Session(config *gorm.Session) ICustomDo
	Columns(cols ...field.Expr) gen.Columns
	Clauses(conds ...clause.Expression) ICustomDo
	Not(conds ...gen.Condition) ICustomDo
	Or(conds ...gen.Condition) ICustomDo
	Select(conds ...field.Expr) ICustomDo
	Where(conds ...gen.Condition) ICustomDo
	Order(conds ...field.Expr) ICustomDo
	Distinct(cols ...field.Expr) ICustomDo
	Omit(cols ...field.Expr) ICustomDo
	Join(table schema.Tabler, on ...field.Expr) ICustomDo
	LeftJoin(table schema.Tabler, on ...field.Expr) ICustomDo
	RightJoin(table schema.Tabler, on ...field.Expr) ICustomDo
	Group(cols ...field.Expr) ICustomDo
	Having(conds ...gen.Condition) ICustomDo
	Limit(limit int) ICustomDo
	Offset(offset int) ICustomDo
	Count() (count int64, err error)
	Scopes(funcs ...func(gen.Dao) gen.Dao) ICustomDo
	Unscoped() ICustomDo
	Create(values ...*model.Custom) error
	CreateInBatches(values []*model.Custom, batchSize int) error
	Save(values ...*model.Custom) error
	First() (*model.Custom, error)
	Take() (*model.Custom, error)
	Last() (*model.Custom, error)
	Find() ([]*model.Custom, error)
	FindInBatch(batchSize int, fc func(tx gen.Dao, batch int) error) (results []*model.Custom, err error)
	FindInBatches(result *[]*model.Custom, batchSize int, fc func(tx gen.Dao, batch int) error) error
	Pluck(column field.Expr, dest interface{}) error
	Delete(...*model.Custom) (info gen.ResultInfo, err error)
	Update(column field.Expr, value interface{}) (info gen.ResultInfo, err error)
	UpdateSimple(columns ...field.AssignExpr) (info gen.ResultInfo, err error)
	Updates(value interface{}) (info gen.ResultInfo, err error)
	UpdateColumn(column field.Expr, value interface{}) (info gen.ResultInfo, err error)
	UpdateColumnSimple(columns ...field.AssignExpr) (info gen.ResultInfo, err error)
	UpdateColumns(value interface{}) (info gen.ResultInfo, err error)
	UpdateFrom(q gen.SubQuery) gen.Dao
	Attrs(attrs ...field.AssignExpr) ICustomDo
	Assign(attrs ...field.AssignExpr) ICustomDo
	Joins(fields ...field.RelationField) ICustomDo
	Preload(fields ...field.RelationField) ICustomDo
	FirstOrInit() (*model.Custom, error)
	FirstOrCreate() (*model.Custom, error)
	FindByPage(offset int, limit int) (result []*model.Custom, count int64, err error)
	ScanByPage(result interface{}, offset int, limit int) (count int64, err error)
	Scan(result interface{}) (err error)
	Returning(value interface{}, columns ...string) ICustomDo
	UnderlyingDB() *gorm.DB
	schema.Tabler
}

func (c customDo) Debug() ICustomDo {
	return c.withDO(c.DO.Debug())
}

func (c customDo) WithContext(ctx context.Context) ICustomDo {
	return c.withDO(c.DO.WithContext(ctx))
}

func (c customDo) ReadDB() ICustomDo {
	return c.Clauses(dbresolver.Read)
}

func (c customDo) WriteDB() ICustomDo {
	return c.Clauses(dbresolver.Write)
}

func (c customDo) Session(config *gorm.Session) ICustomDo {
	return c.withDO(c.DO.Session(config))
}

func (c customDo) Clauses(conds ...clause.Expression) ICustomDo {
	return c.withDO(c.DO.Clauses(conds...))
}

func (c customDo) Returning(value interface{}, columns ...string) ICustomDo {
	return c.withDO(c.DO.Returning(value, columns...))
}

func (c customDo) Not(conds ...gen.Condition) ICustomDo {
	return c.withDO(c.DO.Not(conds...))
}

func (c customDo) Or(conds ...gen.Condition) ICustomDo {
	return c.withDO(c.DO.Or(conds...))
}

func (c customDo) Select(conds ...field.Expr) ICustomDo {
	return c.withDO(c.DO.Select(conds...))
}

func (c customDo) Where(conds ...gen.Condition) ICustomDo {
	return c.withDO(c.DO.Where(conds...))
}

func (c customDo) Order(conds ...field.Expr) ICustomDo {
	return c.withDO(c.DO.Order(conds...))
}

func (c customDo) Distinct(cols ...field.Expr) ICustomDo {
	return c.withDO(c.DO.Distinct(cols...))
}

func (c customDo) Omit(cols ...field.Expr) ICustomDo {
	return c.withDO(c.DO.Omit(cols...))
}

func (c customDo) Join(table schema.Tabler, on ...field.Expr) ICustomDo {
	return c.withDO(c.DO.Join(table, on...))
}

func (c customDo) LeftJoin(table schema.Tabler, on ...field.Expr) ICustomDo {
	return c.withDO(c.DO.LeftJoin(table, on...))
}

func (c customDo) RightJoin(table schema.Tabler, on ...field.Expr) ICustomDo {
	return c.withDO(c.DO.RightJoin(table, on...))
}

func (c customDo) Group(cols ...field.Expr) ICustomDo {
	return c.withDO(c.DO.Group(cols...))
}

func (c customDo) Having(conds ...gen.Condition) ICustomDo {
	return c.withDO(c.DO.Having(conds...))
}

func (c customDo) Limit(limit int) ICustomDo {
	return c.withDO(c.DO.Limit(limit))
}

func (c customDo) Offset(offset int) ICustomDo {
	return c.withDO(c.DO.Offset(offset))
}

func (c customDo) Scopes(funcs ...func(gen.Dao) gen.Dao) ICustomDo {
	return c.withDO(c.DO.Scopes(funcs...))
}

func (c customDo) Unscoped() ICustomDo {
	return c.withDO(c.DO.Unscoped())
}

func (c customDo) Create(values ...*model.Custom) error {
	if len(values) == 0 {
		return nil
	}
	return c.DO.Create(values)
}

func (c customDo) CreateInBatches(values []*model.Custom, batchSize int) error {
	return c.DO.CreateInBatches(values, batchSize)
}

// Save : !!! underlying implementation is different with GORM
// The method is equivalent to executing the statement: db.Clauses(clause.OnConflict{UpdateAll: true}).Create(values)
func (c customDo) Save(values ...*model.Custom) error {
	if len(values) == 0 {
		return nil
	}
	return c.DO.Save(values)
}

func (c customDo) First() (*model.Custom, error) {
	if result, err := c.DO.First(); err != nil {
		return nil, err
	} else {
		return result.(*model.Custom), nil
	}
}

func (c customDo) Take() (*model.Custom, error) {
	if result, err := c.DO.Take(); err != nil {
		return nil, err
	} else {
		return result.(*model.Custom), nil
	}
}

func (c customDo) Last() (*model.Custom, error) {
	if result, err := c.DO.Last(); err != nil {
		return nil, err
	} else {
		return result.(*model.Custom), nil
	}
}

func (c customDo) Find() ([]*model.Custom, error) {
	result, err := c.DO.Find()
	return result.([]*model.Custom), err
}

func (c customDo) FindInBatch(batchSize int, fc func(tx gen.Dao, batch int) error) (results []*model.Custom, err error) {
	buf := make([]*model.Custom, 0, batchSize)
	err = c.DO.FindInBatches(&buf, batchSize, func(tx gen.Dao, batch int) error {
		defer func() { results = append(results, buf...) }()
		return fc(tx, batch)
	})
	return results, err
}

func (c customDo) FindInBatches(result *[]*model.Custom, batchSize int, fc func(tx gen.Dao, batch int) error) error {
	return c.DO.FindInBatches(result, batchSize, fc)
}

func (c customDo) Attrs(attrs ...field.AssignExpr) ICustomDo {
	return c.withDO(c.DO.Attrs(attrs...))
}

func (c customDo) Assign(attrs ...field.AssignExpr) ICustomDo {
	return c.withDO(c.DO.Assign(attrs...))
}

func (c customDo) Joins(fields ...field.RelationField) ICustomDo {
	for _, _f := range fields {
		c = *c.withDO(c.DO.Joins(_f))
	}
	return &c
}

func (c customDo) Preload(fields ...field.RelationField) ICustomDo {
	for _, _f := range fields {
		c = *c.withDO(c.DO.Preload(_f))
	}
	return &c
}

func (c customDo) FirstOrInit() (*model.Custom, error) {
	if result, err := c.DO.FirstOrInit(); err != nil {
		return nil, err
	} else {
		return result.(*model.Custom), nil
	}
}

func (c customDo) FirstOrCreate() (*model.Custom, error) {
	if result, err := c.DO.FirstOrCreate(); err != nil {
		return nil, err
	} else {
		return result.(*model.Custom), nil
	}
}

func (c customDo) FindByPage(offset int, limit int) (result []*model.Custom, count int64, err error) {
	result, err = c.Offset(offset).Limit(limit).Find()
	if err != nil {
		return
	}

	if size := len(result); 0 < limit && 0 < size && size < limit {
		count = int64(size + offset)
		return
	}

	count, err = c.Offset(-1).Limit(-1).Count()
	return
}

func (c customDo) ScanByPage(result interface{}, offset int, limit int) (count int64, err error) {
	count, err = c.Count()
	if err != nil {
		return
	}

	err = c.Offset(offset).Limit(limit).Scan(result)
	return
}

func (c customDo) Scan(result interface{}) (err error) {
	return c.DO.Scan(result)
}

func (c customDo) Delete(models ...*model.Custom) (result gen.ResultInfo, err error) {
	return c.DO.Delete(models)
}

func (c *customDo) withDO(do gen.Dao) *customDo {
	c.DO = *do.(*gen.DO)
	return c
}
