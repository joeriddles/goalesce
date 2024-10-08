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

	"github.com/joeriddles/goalesce/examples/cars/model"
)

func newVehicleForSale(db *gorm.DB, opts ...gen.DOOption) vehicleForSale {
	_vehicleForSale := vehicleForSale{}

	_vehicleForSale.vehicleForSaleDo.UseDB(db, opts...)
	_vehicleForSale.vehicleForSaleDo.UseModel(&model.VehicleForSale{})

	tableName := _vehicleForSale.vehicleForSaleDo.TableName()
	_vehicleForSale.ALL = field.NewAsterisk(tableName)
	_vehicleForSale.ID = field.NewUint(tableName, "id")
	_vehicleForSale.CreatedAt = field.NewTime(tableName, "created_at")
	_vehicleForSale.UpdatedAt = field.NewTime(tableName, "updated_at")
	_vehicleForSale.DeletedAt = field.NewField(tableName, "deleted_at")
	_vehicleForSale.VehicleID = field.NewUint(tableName, "vehicle_id")
	_vehicleForSale.Amount = field.NewField(tableName, "amount")
	_vehicleForSale.Duration = field.NewInt64(tableName, "duration")
	_vehicleForSale.Vehicle = vehicleForSaleBelongsToVehicle{
		db: db.Session(&gorm.Session{}),

		RelationField: field.NewRelation("Vehicle", "model.Vehicle"),
		VehicleModel: struct {
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
		}{
			RelationField: field.NewRelation("Vehicle.VehicleModel", "model.VehicleModel"),
			Manufacturer: struct {
				field.RelationField
				Vehicles struct {
					field.RelationField
				}
			}{
				RelationField: field.NewRelation("Vehicle.VehicleModel.Manufacturer", "model.Manufacturer"),
				Vehicles: struct {
					field.RelationField
				}{
					RelationField: field.NewRelation("Vehicle.VehicleModel.Manufacturer.Vehicles", "model.VehicleModel"),
				},
			},
			Parts: struct {
				field.RelationField
				Models struct {
					field.RelationField
				}
			}{
				RelationField: field.NewRelation("Vehicle.VehicleModel.Parts", "model.Part"),
				Models: struct {
					field.RelationField
				}{
					RelationField: field.NewRelation("Vehicle.VehicleModel.Parts.Models", "model.VehicleModel"),
				},
			},
		},
		Person: struct {
			field.RelationField
		}{
			RelationField: field.NewRelation("Vehicle.Person", "model.Person"),
		},
	}

	_vehicleForSale.fillFieldMap()

	return _vehicleForSale
}

type vehicleForSale struct {
	vehicleForSaleDo

	ALL       field.Asterisk
	ID        field.Uint
	CreatedAt field.Time
	UpdatedAt field.Time
	DeletedAt field.Field
	VehicleID field.Uint
	Amount    field.Field
	Duration  field.Int64
	Vehicle   vehicleForSaleBelongsToVehicle

	fieldMap map[string]field.Expr
}

func (v vehicleForSale) Table(newTableName string) *vehicleForSale {
	v.vehicleForSaleDo.UseTable(newTableName)
	return v.updateTableName(newTableName)
}

func (v vehicleForSale) As(alias string) *vehicleForSale {
	v.vehicleForSaleDo.DO = *(v.vehicleForSaleDo.As(alias).(*gen.DO))
	return v.updateTableName(alias)
}

func (v *vehicleForSale) updateTableName(table string) *vehicleForSale {
	v.ALL = field.NewAsterisk(table)
	v.ID = field.NewUint(table, "id")
	v.CreatedAt = field.NewTime(table, "created_at")
	v.UpdatedAt = field.NewTime(table, "updated_at")
	v.DeletedAt = field.NewField(table, "deleted_at")
	v.VehicleID = field.NewUint(table, "vehicle_id")
	v.Amount = field.NewField(table, "amount")
	v.Duration = field.NewInt64(table, "duration")

	v.fillFieldMap()

	return v
}

func (v *vehicleForSale) GetFieldByName(fieldName string) (field.OrderExpr, bool) {
	_f, ok := v.fieldMap[fieldName]
	if !ok || _f == nil {
		return nil, false
	}
	_oe, ok := _f.(field.OrderExpr)
	return _oe, ok
}

func (v *vehicleForSale) fillFieldMap() {
	v.fieldMap = make(map[string]field.Expr, 8)
	v.fieldMap["id"] = v.ID
	v.fieldMap["created_at"] = v.CreatedAt
	v.fieldMap["updated_at"] = v.UpdatedAt
	v.fieldMap["deleted_at"] = v.DeletedAt
	v.fieldMap["vehicle_id"] = v.VehicleID
	v.fieldMap["amount"] = v.Amount
	v.fieldMap["duration"] = v.Duration

}

func (v vehicleForSale) clone(db *gorm.DB) vehicleForSale {
	v.vehicleForSaleDo.ReplaceConnPool(db.Statement.ConnPool)
	return v
}

func (v vehicleForSale) replaceDB(db *gorm.DB) vehicleForSale {
	v.vehicleForSaleDo.ReplaceDB(db)
	return v
}

type vehicleForSaleBelongsToVehicle struct {
	db *gorm.DB

	field.RelationField

	VehicleModel struct {
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
	Person struct {
		field.RelationField
	}
}

func (a vehicleForSaleBelongsToVehicle) Where(conds ...field.Expr) *vehicleForSaleBelongsToVehicle {
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

func (a vehicleForSaleBelongsToVehicle) WithContext(ctx context.Context) *vehicleForSaleBelongsToVehicle {
	a.db = a.db.WithContext(ctx)
	return &a
}

func (a vehicleForSaleBelongsToVehicle) Session(session *gorm.Session) *vehicleForSaleBelongsToVehicle {
	a.db = a.db.Session(session)
	return &a
}

func (a vehicleForSaleBelongsToVehicle) Model(m *model.VehicleForSale) *vehicleForSaleBelongsToVehicleTx {
	return &vehicleForSaleBelongsToVehicleTx{a.db.Model(m).Association(a.Name())}
}

type vehicleForSaleBelongsToVehicleTx struct{ tx *gorm.Association }

func (a vehicleForSaleBelongsToVehicleTx) Find() (result *model.Vehicle, err error) {
	return result, a.tx.Find(&result)
}

func (a vehicleForSaleBelongsToVehicleTx) Append(values ...*model.Vehicle) (err error) {
	targetValues := make([]interface{}, len(values))
	for i, v := range values {
		targetValues[i] = v
	}
	return a.tx.Append(targetValues...)
}

func (a vehicleForSaleBelongsToVehicleTx) Replace(values ...*model.Vehicle) (err error) {
	targetValues := make([]interface{}, len(values))
	for i, v := range values {
		targetValues[i] = v
	}
	return a.tx.Replace(targetValues...)
}

func (a vehicleForSaleBelongsToVehicleTx) Delete(values ...*model.Vehicle) (err error) {
	targetValues := make([]interface{}, len(values))
	for i, v := range values {
		targetValues[i] = v
	}
	return a.tx.Delete(targetValues...)
}

func (a vehicleForSaleBelongsToVehicleTx) Clear() error {
	return a.tx.Clear()
}

func (a vehicleForSaleBelongsToVehicleTx) Count() int64 {
	return a.tx.Count()
}

type vehicleForSaleDo struct{ gen.DO }

type IVehicleForSaleDo interface {
	gen.SubQuery
	Debug() IVehicleForSaleDo
	WithContext(ctx context.Context) IVehicleForSaleDo
	WithResult(fc func(tx gen.Dao)) gen.ResultInfo
	ReplaceDB(db *gorm.DB)
	ReadDB() IVehicleForSaleDo
	WriteDB() IVehicleForSaleDo
	As(alias string) gen.Dao
	Session(config *gorm.Session) IVehicleForSaleDo
	Columns(cols ...field.Expr) gen.Columns
	Clauses(conds ...clause.Expression) IVehicleForSaleDo
	Not(conds ...gen.Condition) IVehicleForSaleDo
	Or(conds ...gen.Condition) IVehicleForSaleDo
	Select(conds ...field.Expr) IVehicleForSaleDo
	Where(conds ...gen.Condition) IVehicleForSaleDo
	Order(conds ...field.Expr) IVehicleForSaleDo
	Distinct(cols ...field.Expr) IVehicleForSaleDo
	Omit(cols ...field.Expr) IVehicleForSaleDo
	Join(table schema.Tabler, on ...field.Expr) IVehicleForSaleDo
	LeftJoin(table schema.Tabler, on ...field.Expr) IVehicleForSaleDo
	RightJoin(table schema.Tabler, on ...field.Expr) IVehicleForSaleDo
	Group(cols ...field.Expr) IVehicleForSaleDo
	Having(conds ...gen.Condition) IVehicleForSaleDo
	Limit(limit int) IVehicleForSaleDo
	Offset(offset int) IVehicleForSaleDo
	Count() (count int64, err error)
	Scopes(funcs ...func(gen.Dao) gen.Dao) IVehicleForSaleDo
	Unscoped() IVehicleForSaleDo
	Create(values ...*model.VehicleForSale) error
	CreateInBatches(values []*model.VehicleForSale, batchSize int) error
	Save(values ...*model.VehicleForSale) error
	First() (*model.VehicleForSale, error)
	Take() (*model.VehicleForSale, error)
	Last() (*model.VehicleForSale, error)
	Find() ([]*model.VehicleForSale, error)
	FindInBatch(batchSize int, fc func(tx gen.Dao, batch int) error) (results []*model.VehicleForSale, err error)
	FindInBatches(result *[]*model.VehicleForSale, batchSize int, fc func(tx gen.Dao, batch int) error) error
	Pluck(column field.Expr, dest interface{}) error
	Delete(...*model.VehicleForSale) (info gen.ResultInfo, err error)
	Update(column field.Expr, value interface{}) (info gen.ResultInfo, err error)
	UpdateSimple(columns ...field.AssignExpr) (info gen.ResultInfo, err error)
	Updates(value interface{}) (info gen.ResultInfo, err error)
	UpdateColumn(column field.Expr, value interface{}) (info gen.ResultInfo, err error)
	UpdateColumnSimple(columns ...field.AssignExpr) (info gen.ResultInfo, err error)
	UpdateColumns(value interface{}) (info gen.ResultInfo, err error)
	UpdateFrom(q gen.SubQuery) gen.Dao
	Attrs(attrs ...field.AssignExpr) IVehicleForSaleDo
	Assign(attrs ...field.AssignExpr) IVehicleForSaleDo
	Joins(fields ...field.RelationField) IVehicleForSaleDo
	Preload(fields ...field.RelationField) IVehicleForSaleDo
	FirstOrInit() (*model.VehicleForSale, error)
	FirstOrCreate() (*model.VehicleForSale, error)
	FindByPage(offset int, limit int) (result []*model.VehicleForSale, count int64, err error)
	ScanByPage(result interface{}, offset int, limit int) (count int64, err error)
	Scan(result interface{}) (err error)
	Returning(value interface{}, columns ...string) IVehicleForSaleDo
	UnderlyingDB() *gorm.DB
	schema.Tabler
}

func (v vehicleForSaleDo) Debug() IVehicleForSaleDo {
	return v.withDO(v.DO.Debug())
}

func (v vehicleForSaleDo) WithContext(ctx context.Context) IVehicleForSaleDo {
	return v.withDO(v.DO.WithContext(ctx))
}

func (v vehicleForSaleDo) ReadDB() IVehicleForSaleDo {
	return v.Clauses(dbresolver.Read)
}

func (v vehicleForSaleDo) WriteDB() IVehicleForSaleDo {
	return v.Clauses(dbresolver.Write)
}

func (v vehicleForSaleDo) Session(config *gorm.Session) IVehicleForSaleDo {
	return v.withDO(v.DO.Session(config))
}

func (v vehicleForSaleDo) Clauses(conds ...clause.Expression) IVehicleForSaleDo {
	return v.withDO(v.DO.Clauses(conds...))
}

func (v vehicleForSaleDo) Returning(value interface{}, columns ...string) IVehicleForSaleDo {
	return v.withDO(v.DO.Returning(value, columns...))
}

func (v vehicleForSaleDo) Not(conds ...gen.Condition) IVehicleForSaleDo {
	return v.withDO(v.DO.Not(conds...))
}

func (v vehicleForSaleDo) Or(conds ...gen.Condition) IVehicleForSaleDo {
	return v.withDO(v.DO.Or(conds...))
}

func (v vehicleForSaleDo) Select(conds ...field.Expr) IVehicleForSaleDo {
	return v.withDO(v.DO.Select(conds...))
}

func (v vehicleForSaleDo) Where(conds ...gen.Condition) IVehicleForSaleDo {
	return v.withDO(v.DO.Where(conds...))
}

func (v vehicleForSaleDo) Order(conds ...field.Expr) IVehicleForSaleDo {
	return v.withDO(v.DO.Order(conds...))
}

func (v vehicleForSaleDo) Distinct(cols ...field.Expr) IVehicleForSaleDo {
	return v.withDO(v.DO.Distinct(cols...))
}

func (v vehicleForSaleDo) Omit(cols ...field.Expr) IVehicleForSaleDo {
	return v.withDO(v.DO.Omit(cols...))
}

func (v vehicleForSaleDo) Join(table schema.Tabler, on ...field.Expr) IVehicleForSaleDo {
	return v.withDO(v.DO.Join(table, on...))
}

func (v vehicleForSaleDo) LeftJoin(table schema.Tabler, on ...field.Expr) IVehicleForSaleDo {
	return v.withDO(v.DO.LeftJoin(table, on...))
}

func (v vehicleForSaleDo) RightJoin(table schema.Tabler, on ...field.Expr) IVehicleForSaleDo {
	return v.withDO(v.DO.RightJoin(table, on...))
}

func (v vehicleForSaleDo) Group(cols ...field.Expr) IVehicleForSaleDo {
	return v.withDO(v.DO.Group(cols...))
}

func (v vehicleForSaleDo) Having(conds ...gen.Condition) IVehicleForSaleDo {
	return v.withDO(v.DO.Having(conds...))
}

func (v vehicleForSaleDo) Limit(limit int) IVehicleForSaleDo {
	return v.withDO(v.DO.Limit(limit))
}

func (v vehicleForSaleDo) Offset(offset int) IVehicleForSaleDo {
	return v.withDO(v.DO.Offset(offset))
}

func (v vehicleForSaleDo) Scopes(funcs ...func(gen.Dao) gen.Dao) IVehicleForSaleDo {
	return v.withDO(v.DO.Scopes(funcs...))
}

func (v vehicleForSaleDo) Unscoped() IVehicleForSaleDo {
	return v.withDO(v.DO.Unscoped())
}

func (v vehicleForSaleDo) Create(values ...*model.VehicleForSale) error {
	if len(values) == 0 {
		return nil
	}
	return v.DO.Create(values)
}

func (v vehicleForSaleDo) CreateInBatches(values []*model.VehicleForSale, batchSize int) error {
	return v.DO.CreateInBatches(values, batchSize)
}

// Save : !!! underlying implementation is different with GORM
// The method is equivalent to executing the statement: db.Clauses(clause.OnConflict{UpdateAll: true}).Create(values)
func (v vehicleForSaleDo) Save(values ...*model.VehicleForSale) error {
	if len(values) == 0 {
		return nil
	}
	return v.DO.Save(values)
}

func (v vehicleForSaleDo) First() (*model.VehicleForSale, error) {
	if result, err := v.DO.First(); err != nil {
		return nil, err
	} else {
		return result.(*model.VehicleForSale), nil
	}
}

func (v vehicleForSaleDo) Take() (*model.VehicleForSale, error) {
	if result, err := v.DO.Take(); err != nil {
		return nil, err
	} else {
		return result.(*model.VehicleForSale), nil
	}
}

func (v vehicleForSaleDo) Last() (*model.VehicleForSale, error) {
	if result, err := v.DO.Last(); err != nil {
		return nil, err
	} else {
		return result.(*model.VehicleForSale), nil
	}
}

func (v vehicleForSaleDo) Find() ([]*model.VehicleForSale, error) {
	result, err := v.DO.Find()
	return result.([]*model.VehicleForSale), err
}

func (v vehicleForSaleDo) FindInBatch(batchSize int, fc func(tx gen.Dao, batch int) error) (results []*model.VehicleForSale, err error) {
	buf := make([]*model.VehicleForSale, 0, batchSize)
	err = v.DO.FindInBatches(&buf, batchSize, func(tx gen.Dao, batch int) error {
		defer func() { results = append(results, buf...) }()
		return fc(tx, batch)
	})
	return results, err
}

func (v vehicleForSaleDo) FindInBatches(result *[]*model.VehicleForSale, batchSize int, fc func(tx gen.Dao, batch int) error) error {
	return v.DO.FindInBatches(result, batchSize, fc)
}

func (v vehicleForSaleDo) Attrs(attrs ...field.AssignExpr) IVehicleForSaleDo {
	return v.withDO(v.DO.Attrs(attrs...))
}

func (v vehicleForSaleDo) Assign(attrs ...field.AssignExpr) IVehicleForSaleDo {
	return v.withDO(v.DO.Assign(attrs...))
}

func (v vehicleForSaleDo) Joins(fields ...field.RelationField) IVehicleForSaleDo {
	for _, _f := range fields {
		v = *v.withDO(v.DO.Joins(_f))
	}
	return &v
}

func (v vehicleForSaleDo) Preload(fields ...field.RelationField) IVehicleForSaleDo {
	for _, _f := range fields {
		v = *v.withDO(v.DO.Preload(_f))
	}
	return &v
}

func (v vehicleForSaleDo) FirstOrInit() (*model.VehicleForSale, error) {
	if result, err := v.DO.FirstOrInit(); err != nil {
		return nil, err
	} else {
		return result.(*model.VehicleForSale), nil
	}
}

func (v vehicleForSaleDo) FirstOrCreate() (*model.VehicleForSale, error) {
	if result, err := v.DO.FirstOrCreate(); err != nil {
		return nil, err
	} else {
		return result.(*model.VehicleForSale), nil
	}
}

func (v vehicleForSaleDo) FindByPage(offset int, limit int) (result []*model.VehicleForSale, count int64, err error) {
	result, err = v.Offset(offset).Limit(limit).Find()
	if err != nil {
		return
	}

	if size := len(result); 0 < limit && 0 < size && size < limit {
		count = int64(size + offset)
		return
	}

	count, err = v.Offset(-1).Limit(-1).Count()
	return
}

func (v vehicleForSaleDo) ScanByPage(result interface{}, offset int, limit int) (count int64, err error) {
	count, err = v.Count()
	if err != nil {
		return
	}

	err = v.Offset(offset).Limit(limit).Scan(result)
	return
}

func (v vehicleForSaleDo) Scan(result interface{}) (err error) {
	return v.DO.Scan(result)
}

func (v vehicleForSaleDo) Delete(models ...*model.VehicleForSale) (result gen.ResultInfo, err error) {
	return v.DO.Delete(models)
}

func (v *vehicleForSaleDo) withDO(do gen.Dao) *vehicleForSaleDo {
	v.DO = *do.(*gen.DO)
	return v
}
