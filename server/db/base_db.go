package db

import (
	"gorm.io/gorm"
)

// BaseDB implements CRUDRepo for the user model
type BaseDB[T AllowedModels] struct {
	db   *gorm.DB
	zero T
}

// NewBaseDB returns a new BaseDB instance for the required type,
func NewBaseDB[T AllowedModels](db *gorm.DB) UnsafeCRUDRepo[T] {
	return &BaseDB[T]{db: db}
}

// Add creates a new record of the given object, and returns an occurring error
// the new object is a pointer, so it updates the object's id after creation
func (b *BaseDB[T]) Add(obj *T) error {
	if obj == nil {
		return &ErrNilObject{}
	}

	err := b.db.
		Model(new(T)).
		Create(obj).
		Error

	if err != nil {
		return TryWrapMySqlError(err)
	}

	return nil
}

// AddMany is same as Add but for numerous objects
func (b *BaseDB[T]) AddMany(objs []*T) error {
	if len(objs) == 0 {
		return &ErrNilObject{}
	}

	err := b.db.
		Model(new([]T)).
		Create(&objs).
		Error

	if err != nil {
		return TryWrapMySqlError(err)
	}

	return nil
}

// Exists checks the existence of the given record's id
func (b *BaseDB[T]) Exists(id uint) error {
	if id == 0 { // better to check this before, fetching eh?
		return &ErrRecordNotFound{}
	}

	var obj T
	return b.
		db.
		Select("id").
		Where("id = ?", id).
		First(&obj).
		Error
}

// Get retrieves the object which has the given id
func (b *BaseDB[T]) Get(id uint) (T, error) {
	var obj T

	err := b.db.
		Model(new(T)).
		First(&obj, "id = ?", id).
		Error

	if err != nil {
		return b.zero, err
	}

	return obj, nil
}

// GetByConds is the extended version of Get,
// which uses a given search condition and retrieves every record with the given condition
func (b *BaseDB[T]) GetByConds(conds ...any) ([]T, error) {
	if !checkConds(conds...) {
		return nil, &ErrInvalidWhereConditions{}
	}

	var foundRecords []T

	err := b.db.
		Model(new(T)).
		Find(&foundRecords, conds...).
		Error

	if err != nil || len(foundRecords) == 0 {
		return nil, &ErrRecordNotFound{}
	}

	return foundRecords, nil
}

// GetAll retrieves all the records of the given model
func (b *BaseDB[T]) GetAll() ([]T, error) {
	return b.GetByConds("id != ?", 0)
}

// Count returns the number of records of the given model
func (b *BaseDB[T]) Count() int64 {
	var count int64

	err := b.db.
		Model(new(T)).
		Count(&count).
		Error

	if err != nil {
		return 0
	}

	return count
}

// Update updates the given object/s based on the given condition
// the updated object is a pointer, so it changes the values in it as well,
// and gives it its id(in case searching condition weren't using id)
func (b *BaseDB[T]) Update(obj *T, conds ...any) error {
	if obj == nil {
		return &ErrNilObject{}
	}
	if len(conds) > 1 {
		if !checkConds(conds...) {
			return &ErrInvalidWhereConditions{}
		}
	} else {
		conds = []any{"id = ?", (*obj).GetId()}
	}
	_, err := b.GetByConds(conds...)
	if err != nil {
		return err
	}

	err = b.db.
		Model(new(T)).
		Where(conds[0], conds[1:]...).
		Updates(obj).
		Error
	if err != nil {
		return err
	}

	return nil
}

// Delete deletes the given object/s based on the given object
func (b *BaseDB[T]) Delete(conds ...any) error {
	obj, err := b.GetByConds(conds...)
	if err != nil {
		return err
	}

	err = b.db.
		Model(new(T)).
		Delete(&obj, conds...).
		Error

	if err != nil {
		return err
	}

	return nil
}

// GetDB well, ding...
func (b *BaseDB[T]) GetDB() *gorm.DB {
	return b.db
}

////////////

func checkConds(conds ...any) bool {
	return len(conds) > 1 && checkCondsMeaning(conds...)
}

func checkCondsMeaning(conds ...any) bool {
	ok := false

	switch conds[0].(type) {
	case string:
		ok = true
	default:
		return false
	}

	for _, cond := range conds[1:] {
		switch cond.(type) {
		case bool,
			int, int8, int16, int32, int64,
			uint, uint8, uint16, uint32, uint64,
			float32, float64,
			complex64, complex128,
			string, []uint:
			ok = true
		default:
			return false
		}
	}

	return ok
}
