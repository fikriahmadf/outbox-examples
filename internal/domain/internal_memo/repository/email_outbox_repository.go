package repository

import (
	"context"
	"fmt"

	"github.com/fikriahmadf/outbox-examples/internal/domain/internal_memo/model"
	"github.com/fikriahmadf/outbox-examples/shared/failure"
)

var (
	emailOutboxQueries = struct {
		insertEmailOutbox string
	}{
		insertEmailOutbox: "INSERT INTO \"email_outbox\" %s VALUES %s",
	}
)

func (r *InternalMemoRepositoryPostgres) CreateEmailOutbox(ctx context.Context, emailOutbox *model.EmailOutbox) error {
	insertQuery := fmt.Sprintf(
		emailOutboxQueries.insertEmailOutbox,
		"(id, memo_id, event_type, payload, recipient_email, notification_type, status, retry_count, last_attempt_at, sent_at, error_message, idempotency_key, meta_created_at, meta_updated_at)",
		"($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14)",
	)
	argsList := []any{
		emailOutbox.ID,
		emailOutbox.MemoID,
		emailOutbox.EventType,
		emailOutbox.Payload,
		emailOutbox.RecipientEmail,
		emailOutbox.NotificationType,
		emailOutbox.Status,
		emailOutbox.RetryCount,
		emailOutbox.LastAttemptAt,
		emailOutbox.SentAt,
		emailOutbox.ErrorMessage,
		emailOutbox.IdempotencyKey,
		emailOutbox.MetaCreatedAt,
		emailOutbox.MetaUpdatedAt,
	}
	_, err := r.exec(ctx, insertQuery, argsList)
	if err != nil {
		return failure.AddFuncName(failure.InternalError(err))
	}
	return nil
}

type EmailOutboxRepository interface {
	CreateEmailOutbox(ctx context.Context, emailOutbox *model.EmailOutbox) error
}
