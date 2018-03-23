package db

//go:generate ../../../bin/mockgen.sh storage/discovery/db DB $GOFILE

import (
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
	Find(interface{}, ...interface{}) DB
	First(interface{}, ...interface{}) DB
	GetError() error
	GetErrors() []error
	Preload(string, ...interface{}) DB
	RecordNotFound() bool
	Rollback() DB
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

func (d *db) Preload(column string, conditions ...interface{}) DB {
	return &db{d.db.Preload(column, conditions...)}
}

func (d *db) RecordNotFound() bool {
	return d.db.RecordNotFound()
}

func (d *db) Rollback() DB {
	return &db{d.db.Rollback()}
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
