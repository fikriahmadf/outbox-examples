package repository

import (
	"context"
	"fmt"

	"github.com/fikriahmadf/outbox-examples/internal/domain/internal_memo/model"
	"github.com/rs/zerolog/log"
)

const MEMO_PREFIX = "MEMO-"

var (
	memoQueries = struct {
		insertMemo string
		countMemo  string
	}{
		insertMemo: "INSERT INTO \"memos\" %s VALUES %s",
		countMemo:  "SELECT COUNT(*) FROM \"memos\"",
	}
)

func (r *InternalMemoRepositoryPostgres) GetCountMemo(ctx context.Context) (int, error) {
	var count int
	err := r.DB.Read.GetContext(ctx, &count, memoQueries.countMemo)
	if err != nil {
		log.Error().Err(err).Msg("[InternalMemoRepositoryPostgres][GetCountMemo] failed to get count memo")
		return 0, err
	}
	return count, nil
}

func (r *InternalMemoRepositoryPostgres) CreateMemo(ctx context.Context, memo *model.Memo) error {

	count, err := r.GetCountMemo(ctx)
	if err != nil {
		log.Error().Err(err).Msg("[InternalMemoRepositoryPostgres][CreateMemo] failed to get count memo")
		return err
	}

	insertQuery := fmt.Sprintf(memoQueries.insertMemo, "(id, memo_number_prefix, memo_number_sequence, department_code, title, purpose, created_at)", "($1, $2, $3, $4, $5, $6, $7)")
	argsList := []any{
		memo.ID,
		MEMO_PREFIX,
		count + 1,
		memo.DepartmentCode,
		memo.Title,
		memo.Purpose,
		memo.CreatedAt,
	}

	_, err = r.exec(ctx, insertQuery, argsList)
	if err != nil {
		log.Error().Err(err).Msg("[InternalMemoRepositoryPostgres][CreateMemo] failed to get count memo")
		return err
	}

	memo.MemoNumberSequence = count + 1
	memo.MemoNumberPrefix = MEMO_PREFIX

	return nil
}

type MemoRepository interface {
	CreateMemo(ctx context.Context, memo *model.Memo) error
	GetCountMemo(ctx context.Context) (int, error)
}
