package db

import "gorm.io/gorm"

// CreatorRepo is a safe wrapper for records creation for a certain repo
type CreatorRepo[T any] interface {
	// Add creates a new record of the given object, and returns an occurring error
	// the new object is a pointer, so it updates the object's id after creation
	Add(obj *T) error

	// AddMany is same as Add but for numerous objects
	AddMany(objs []*T) error
}

// GetterRepo is a safe wrapper for records retrieval for a certain repo
type GetterRepo[T any] interface {
	// Exists checks the existence of the given record's id
	Exists(id uint) bool

	// Get retrieves the object which has the given id
	Get(id uint) (T, error)

	// GetByConds is the extended version of Get,
	// which uses a given search condition and retrieves every record with the given condition
	GetByConds(conds ...any) ([]T, error)

	// GetAll retrieves all the records of the given model
	GetAll() ([]T, error)

	// Count returns the number of records of the given model
	Count() int64
}

// UpdaterRepo is a safe wrapper for records updating for a certain repo
type UpdaterRepo[T any] interface {
	// Update updates the given object/s based on the given condition
	// the updated object is a pointer, so it changes the values in it as well,
	// and gives it its id(in case searching condition weren't using id)
	Update(obj *T, conds ...any) error
}

// DeleterRepo is a safe wrapper for records deletion for a certain repo
type DeleterRepo[T any] interface {
	// Delete deletes the given object/s based on the given object
	Delete(conds ...any) error
}

// GORMDBGetter is the dark side part of CRUDRepo, which allows direct operations using *gorm.DB instance
type GORMDBGetter interface {
	// GetDB well, ding...
	GetDB() *gorm.DB
}

// CRUDRepo is the whole package of CRUD repos
type CRUDRepo[T any] interface {
	CreatorRepo[T]
	GetterRepo[T]
	UpdaterRepo[T]
	DeleterRepo[T]
}

// UnsafeCRUDRepo is same as CRUDRepo but has the dark method GetDB
type UnsafeCRUDRepo[T any] interface {
	CRUDRepo[T]
	GORMDBGetter
}
