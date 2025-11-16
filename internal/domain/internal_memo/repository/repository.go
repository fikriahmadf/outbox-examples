package repository

import (
	"fmt"

	"github.com/fikriahmadf/outbox-examples/infras"
	"github.com/fikriahmadf/outbox-examples/shared/failure"
	"github.com/jmoiron/sqlx"
	"github.com/rs/zerolog/log"

	"context"
	"database/sql"
)

type InternalMemoRepository interface {
	MemoRepository
}

// InternalMemoRepositoryPostgres is the Postgres-backed implementation of InternalMemoRepository.
type InternalMemoRepositoryPostgres struct {
	DB   *infras.PostgresConn
	dbTx *infras.PostgresTx
}

// ProvideInternalMemoRepositoryPostgres is the provider for this repository.
func ProvideInternalMemoRepositoryPostgres(db *infras.PostgresConn) *InternalMemoRepositoryPostgres {
	s := new(InternalMemoRepositoryPostgres)
	s.DB = db
	return s
}

func (repo *InternalMemoRepositoryPostgres) exec(ctx context.Context, command string, args []any) (sql.Result, error) {
	var (
		stmt *sqlx.Stmt
		err  error
	)
	stmt, err = repo.DB.Write.PreparexContext(ctx, command)
	if err != nil {
		return nil, failure.AddFuncName(failure.InternalError(err))
	}

	defer func() {
		err := stmt.Close()
		if err != nil {
			log.Ctx(ctx).Error().Err(err).Send()
		}
	}()

	result, err := stmt.ExecContext(ctx, args...)
	if err != nil {
		return nil, failure.AddFuncName(failure.InternalError(err))
	}

	return result, nil
}

func (r *InternalMemoRepositoryPostgres) BeginTx(ctx context.Context) (InternalMemoRepository, error) {
	tx, err := r.DB.Write.Beginx(ctx)
	if err != nil {
		log.Error().Err(err).Msg("[InternalMemoRepositoryPostgres][BeginTx] begin transaction failed")
		return nil, failure.InternalError(fmt.Errorf(failure.ErrorInternalSystem))
	}

	return &InternalMemoRepositoryPostgres{
		DB:   r.DB,
		dbTx: tx,
	}, nil
}

func (r *InternalMemoRepositoryPostgres) Rollback(ctx context.Context) error {
	if r.dbTx == nil {
		log.Error().Msg("[InternalMemoRepositoryPostgres][Rollback] not transaction")
		return failure.InternalError(fmt.Errorf(failure.ErrorInternalSystem))
	}
	return r.dbTx.Rollback()
}

func (r *InternalMemoRepositoryPostgres) Commit(ctx context.Context) error {
	if r.dbTx == nil {
		log.Error().Msg("[InternalMemoRepositoryPostgres][Commit] not transaction")
		return failure.InternalError(fmt.Errorf(failure.ErrorInternalSystem))
	}
	return r.dbTx.Commit()
}
