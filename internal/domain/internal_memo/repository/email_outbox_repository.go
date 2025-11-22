package repository

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/fikriahmadf/outbox-examples/internal/domain/internal_memo/model"
	"github.com/rs/zerolog/log"
)

var (
	emailOutboxQueries = struct {
		insertEmailOutbox         string
		resolvePendingEmailOutbox string
		updateOutboxProcess       string
		UpdateSentOutboxProcess   string
	}{
		insertEmailOutbox:         "INSERT INTO \"email_outbox\" %s VALUES %s",
		resolvePendingEmailOutbox: "SELECT * FROM \"email_outbox\" WHERE status = 'PENDING' ORDER BY meta_created_at FOR UPDATE SKIP LOCKED LIMIT $1",
		updateOutboxProcess:       "UPDATE \"email_outbox\" SET retry_count = $1, error_message = $2, status = $3, last_attempt_at = $4, meta_updated_at =$5 WHERE id = $6",
		UpdateSentOutboxProcess:   "UPDATE \"email_outbox\" SET status = 'SENT', sent_at = $1, last_attempt_at = $2, meta_updated_at =$3, error_message = NULL WHERE id = $4",
	}
)

func (r *InternalMemoRepositoryPostgres) CreateEmailOutbox(ctx context.Context, emailOutbox *model.EmailOutbox) error {
	insertQuery := fmt.Sprintf(
		emailOutboxQueries.insertEmailOutbox,
		"(id, memo_id, event_type, payload, recipient_email, status, retry_count, idempotency_key, meta_created_at)",
		"($1, $2, $3, $4, $5, $6, $7, $8, $9)",
	)
	argsList := []any{
		emailOutbox.ID,
		emailOutbox.MemoID,
		emailOutbox.EventType,
		emailOutbox.Payload,
		emailOutbox.RecipientEmail,
		emailOutbox.Status,
		emailOutbox.RetryCount,
		emailOutbox.IdempotencyKey,
		emailOutbox.MetaCreatedAt,
	}
	_, err := r.exec(ctx, insertQuery, argsList)
	if err != nil {
		log.Error().Err(err).Msg("[InternalMemoRepositoryPostgres][CreateEmailOutbox] failed to create email outbox")
		return err
	}
	return nil
}

func (r *InternalMemoRepositoryPostgres) ResolvePendingEmailOutbox(ctx context.Context, limitProcessor int) ([]model.EmailOutbox, error) {
	var rows *sql.Rows
	var err error

	if r.dbTx != nil {
		rows, err = r.dbTx.QueryContext(ctx, emailOutboxQueries.resolvePendingEmailOutbox, limitProcessor)
		if err != nil {
			log.Error().Err(err).Msg("[InternalMemoRepositoryPostgres][ResolvePendingEmailOutbox][Tx] failed to resolve pending email outbox")
			return nil, err
		}
		defer rows.Close()
	} else {
		rows, err = r.DB.Read.QueryContext(ctx, emailOutboxQueries.resolvePendingEmailOutbox, limitProcessor)
		if err != nil {
			log.Error().Err(err).Msg("[InternalMemoRepositoryPostgres][ResolvePendingEmailOutbox][DB] failed to resolve pending email outbox")
			return nil, err
		}
		defer rows.Close()
	}

	var emailOutboxes []model.EmailOutbox
	for rows.Next() {
		var emailOutbox model.EmailOutbox
		if err := rows.Scan(&emailOutbox.ID, &emailOutbox.MemoID, &emailOutbox.EventType, &emailOutbox.Payload, &emailOutbox.RecipientEmail, &emailOutbox.Status, &emailOutbox.RetryCount, &emailOutbox.IdempotencyKey, &emailOutbox.MetaCreatedAt); err != nil {
			log.Error().Err(err).Msg("[InternalMemoRepositoryPostgres][ResolvePendingEmailOutbox] failed to scan email outbox")
			return nil, err
		}
		emailOutboxes = append(emailOutboxes, emailOutbox)
	}
	if err := rows.Err(); err != nil {
		log.Error().Err(err).Msg("[InternalMemoRepositoryPostgres][ResolvePendingEmailOutbox] failed to iterate email outbox")
		return nil, err
	}
	return emailOutboxes, nil

}

func (r *InternalMemoRepositoryPostgres) UpdateErrorProcess(ctx context.Context, outbox *model.EmailOutbox) error {
	updateQuery := fmt.Sprintf(
		emailOutboxQueries.updateOutboxProcess,
		"($1, $2, $3, $4, $5, $6)",
	)
	argsList := []any{
		outbox.RetryCount,
		outbox.ErrorMessage,
		outbox.Status,
		outbox.LastAttemptAt,
		outbox.MetaUpdatedAt,
		outbox.ID,
	}
	_, err := r.exec(ctx, updateQuery, argsList)
	if err != nil {
		log.Error().Err(err).Msg("[InternalMemoRepositoryPostgres][UpdateErrorProcess] failed to update email outbox")
		return err
	}
	return nil
}

func (r *InternalMemoRepositoryPostgres) UpdateSentOutboxProcess(ctx context.Context, outbox *model.EmailOutbox) error {
	updateQuery := fmt.Sprintf(
		emailOutboxQueries.updateOutboxProcess,
		"($1, $2, $3, $4, $5, $6)",
	)
	argsList := []any{
		outbox.SentAt,
		outbox.LastAttemptAt,
		outbox.MetaUpdatedAt,
		outbox.ID,
	}
	_, err := r.exec(ctx, updateQuery, argsList)
	if err != nil {
		log.Error().Err(err).Msg("[InternalMemoRepositoryPostgres][UpdateSentOutboxProcess] failed to update email outbox")
		return err
	}
	return nil
}

type EmailOutboxRepository interface {
	CreateEmailOutbox(ctx context.Context, emailOutbox *model.EmailOutbox) error
	ResolvePendingEmailOutbox(ctx context.Context, limitProcessor int) ([]model.EmailOutbox, error)
	UpdateErrorProcess(ctx context.Context, outbox *model.EmailOutbox) error
	UpdateSentOutboxProcess(ctx context.Context, outbox *model.EmailOutbox) error
}
