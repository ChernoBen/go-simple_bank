package db

import (
	"context"
	"database/sql"
	"fmt"
)

// guarda todas as funções para executar queries e transactions
type Store struct {
	*Queries
	db *sql.DB
}

func NewStore(db *sql.DB) *Store {
	return &Store{
		db:      db,
		Queries: New(db),
	}
}

//executa as funções dentro de uma db transation
//deve retornar um "objeto" ou um erro
func (store *Store) execTx(ctx context.Context, fn func(*Queries) error) error {
	tx, err := store.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	q := New(tx)
	err = fn(q)
	if err != nil {
		// se existe rollback error
		if rbErr := tx.Rollback(); rbErr != nil {
			return fmt.Errorf("tx err:%v, rb err: %v", err, rbErr)
		}
		return err
	}
	return tx.Commit()
}
