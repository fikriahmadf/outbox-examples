package repository

import (
	"context"
	"fmt"

	"github.com/fikriahmadf/outbox-examples/internal/domain/internal_memo/model"
	"github.com/fikriahmadf/outbox-examples/shared/failure"
)

var (
	memoQueries = struct {
		insertMemo string
		countMemo  string
	}{
		insertMemo: "INSERT INTO \"memo\" %s VALUES %s",
		countMemo:  "SELECT COUNT(*) FROM \"memo\"",
	}
)

func (r *InternalMemoRepositoryPostgres) GetCountMemo(ctx context.Context) (int, error) {
	var count int
	err := r.DB.Read.GetContext(ctx, &count, memoQueries.countMemo)
	if err != nil {
		return 0, failure.AddFuncName(failure.InternalError(err))
	}
	return count, nil
}

func (r *InternalMemoRepositoryPostgres) CreateMemo(ctx context.Context, memo *model.Memo) error {

	count, err := r.GetCountMemo(ctx)
	if err != nil {
		return failure.AddFuncName(failure.InternalError(err))
	}

	insertQuery := fmt.Sprintf(memoQueries.insertMemo, "(id, memo_number_prefix, memo_number_sequence, department_code, title, purpose, created_at, updated_at)", "($1, $2, $3, $4, $5, $6, $7, $8)")
	argsList := []any{
		memo.ID,
		memo.MemoNumberPrefix,
		count + 1,
		memo.DepartmentCode,
		memo.Title,
		memo.Purpose,
		memo.CreatedAt,
		memo.UpdatedAt,
	}

	_, err = r.exec(ctx, insertQuery, argsList)
	if err != nil {
		return failure.AddFuncName(failure.InternalError(err))
	}
	return nil
}

type MemoRepository interface {
	CreateMemo(ctx context.Context, memo *model.Memo) error
	GetCountMemo(ctx context.Context) (int, error)
}
