package database

import (
	"context"
	"database/sql"
	"errors"

	"gorm.io/gorm"
)

var ErrBadEntity = errors.New("the provided value does not implement the Entity interface")
var ErrNotFound = errors.New("record not found")

type Repository[T Entity] struct {
	db *gorm.DB
}

func NewRepository[T Entity](db *gorm.DB) (*Repository[T], error) {
	var rawEntity any = new(T)

	entity, ok := rawEntity.(Entity)

	if !ok {
		return nil, ErrBadEntity
	}

	if err := db.AutoMigrate(entity); err != nil {
		return nil, err
	}

	return &Repository[T]{db}, nil
}

func (r *Repository[T]) DB() *gorm.DB {
	return r.db
}

func (r *Repository[T]) Tx(ctx context.Context, txFn func(context.Context) error) error {
	return r.db.Transaction(func(tx *gorm.DB) error {
		return txFn(WithTx(ctx, tx))
	})
}

func (r *Repository[T]) Create(ctx context.Context, entity *T) (*T, error) {
	tx := r.db
	if dbTx, err := FromContext(ctx); err == nil {
		tx = dbTx
	}

	if err := tx.Create(entity).Error; err != nil {
		return nil, err
	}

	return entity, nil
}

func (r *Repository[T]) CreateBatch(ctx context.Context, entities []*T) ([]*T, error) {
	tx := r.db
	if dbTx, err := FromContext(ctx); err == nil {
		tx = dbTx
	}

	if err := tx.Create(entities).Error; err != nil {
		return nil, err
	}

	return entities, nil
}

func (r *Repository[T]) Update(ctx context.Context, entity *T) (*T, error) {
	tx := r.db
	if dbTx, err := FromContext(ctx); err == nil {
		tx = dbTx
	}

	if err := tx.Model(entity).Updates(entity).Error; err != nil {
		return nil, err
	}

	return entity, nil
}

func (r *Repository[T]) FindById(ctx context.Context, id string, opts ...Option) (*T, error) {
	tx := r.db
	if dbTx, err := FromContext(ctx); err == nil {
		tx = dbTx
	}

	for _, opt := range opts {
		tx = opt(tx)
	}

	res := new(T)
	if err := tx.Model(res).Where("id = ?", id).First(res).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrNotFound
		}
		return nil, err
	}

	return res, nil
}

func (r *Repository[T]) FindBy(ctx context.Context, entity *T, opts ...Option) (*T, error) {
	tx := r.db
	if dbTx, err := FromContext(ctx); err == nil {
		tx = dbTx
	}

	for _, opt := range opts {
		tx = opt(tx)
	}

	if err := tx.Model(entity).Where(entity).Take(entity).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrNotFound
		}
		return nil, err
	}

	return entity, nil
}

func (r *Repository[T]) FindAll(ctx context.Context, query *T, opts ...Option) ([]*T, error) {
	tx := r.db
	if dbTx, err := FromContext(ctx); err == nil {
		tx = dbTx
	}

	for _, opt := range opts {
		tx = opt(tx)
	}

	var res []*T
	if err := tx.Where(query).Find(&res).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrNotFound
		}
		return nil, err
	}

	return res, nil
}

func (r *Repository[T]) FindLast(ctx context.Context, query *T, opts ...Option) (*T, error) {
	tx := r.db
	if dbTx, err := FromContext(ctx); err == nil {
		tx = dbTx
	}

	for _, opt := range opts {
		tx = opt(tx)
	}

	if err := tx.Where(query).Last(query).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrNotFound
		}
		return nil, err
	}

	return query, nil
}

func (r *Repository[T]) Query(ctx context.Context, query string, values ...interface{}) (*sql.Rows, error) {
	q := r.db.Raw(query, values)
	rows, err := q.Rows()
	defer rows.Close()

	if err != nil {
		return nil, err
	}

	return rows, nil
}
