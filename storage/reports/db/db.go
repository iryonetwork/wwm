package db

//go:generate ../../../bin/mockgen.sh storage/reports/db DB $GOFILE

import (
	"database/sql"

	"github.com/jinzhu/gorm"
)

// DB interface is used to ease mocking/stubbing gorm's DB object
type DB interface {
	AutoMigrate(...interface{}) DB
	Begin() DB
	Close() error
	Commit() DB
	Create(interface{}) DB
	Delete(interface{}, ...interface{}) DB
	Exec(string, ...interface{}) DB
	Find(interface{}, ...interface{}) DB
	First(interface{}, ...interface{}) DB
	GetError() error
	GetErrors() []error
	Model(interface{}) DB
	Order(interface{}, ...bool) DB
	Preload(string, ...interface{}) DB
	RecordNotFound() bool
	Related(interface{}, ...string) DB
	Rollback() DB
	Rows() (*sql.Rows, error)
	Scan(interface{}) DB
	Select(interface{}, ...interface{}) DB
	Save(interface{}) DB
	Set(string, interface{}) DB
	Update(...interface{}) DB
	Where(interface{}, ...interface{}) DB
}

type db struct {
	db *gorm.DB
}

func (d *db) AutoMigrate(in ...interface{}) DB {
	return &db{d.db.AutoMigrate(in)}
}

func (d *db) Begin() DB {
	return &db{d.db.Begin()}
}

func (d *db) Close() error {
	return d.db.Close()
}

func (d *db) Commit() DB {
	return &db{d.db.Commit()}
}

func (d *db) Create(in interface{}) DB {
	return &db{d.db.Create(in)}
}

func (d *db) Delete(in interface{}, where ...interface{}) DB {
	return &db{d.db.Delete(in, where...)}
}

func (d *db) Exec(sql string, values ...interface{}) DB {
	return &db{d.db.Exec(sql, values...)}
}

func (d *db) Find(in interface{}, args ...interface{}) DB {
	return &db{d.db.Find(in, args...)}
}

func (d *db) First(in interface{}, args ...interface{}) DB {
	return &db{d.db.First(in, args...)}
}

func (d *db) GetError() error {
	return d.db.Error
}

func (d *db) GetErrors() []error {
	return d.db.GetErrors()
}

func (d *db) Select(query interface{}, args ...interface{}) DB {
	return &db{d.db.Select(query, args...)}
}

func (d *db) Model(value interface{}) DB {
	return &db{d.db.Model(value)}
}

func (d *db) Order(value interface{}, reorder ...bool) DB {
	return &db{d.db.Order(value, reorder...)}
}

func (d *db) Preload(column string, conditions ...interface{}) DB {
	return &db{d.db.Preload(column, conditions...)}
}

func (d *db) RecordNotFound() bool {
	return d.db.RecordNotFound()
}

func (d *db) Related(value interface{}, foreignKeys ...string) DB {
	return &db{d.db.Related(value, foreignKeys...)}
}

func (d *db) Rollback() DB {
	return &db{d.db.Rollback()}
}

func (d *db) Rows() (*sql.Rows, error) {
	return d.db.Rows()
}

func (d *db) Scan(dest interface{}) DB {
	return &db{d.db.Scan(dest)}
}

func (d *db) Set(name string, value interface{}) DB {
	return &db{d.db.Set(name, value)}
}

func (d *db) Save(value interface{}) DB {
	return &db{d.db.Save(value)}
}

func (d *db) Update(attrs ...interface{}) DB {
	return &db{d.db.Update(attrs...)}
}

func (d *db) Where(query interface{}, attrs ...interface{}) DB {
	return &db{d.db.Where(query, attrs...)}
}

// New returns a wrapped gorm DB
func New(gdb *gorm.DB) DB {
	return &db{gdb}
}
