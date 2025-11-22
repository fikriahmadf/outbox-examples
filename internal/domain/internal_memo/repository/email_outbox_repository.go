package repository

import (
	"context"
	"fmt"

	"github.com/fikriahmadf/outbox-examples/internal/domain/internal_memo/model"
	"github.com/rs/zerolog/log"
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

type EmailOutboxRepository interface {
	CreateEmailOutbox(ctx context.Context, emailOutbox *model.EmailOutbox) error
}
