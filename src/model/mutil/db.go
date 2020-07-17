package mutil

import "github.com/jmoiron/sqlx"

func Transaction(dbx *sqlx.DB, fn func(tx *sqlx.Tx) error) error {
	tx, err := dbx.Beginx()
	if err != nil {
		return err
	}

	err = fn(tx)
	if err != nil {
		_ = tx.Rollback()
		return err
	}

	return tx.Commit()
}
