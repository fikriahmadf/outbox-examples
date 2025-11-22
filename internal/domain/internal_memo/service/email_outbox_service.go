package service

import (
	"context"
	"database/sql"
	"encoding/json"
	"time"

	"github.com/fikriahmadf/outbox-examples/external/domain/notif_publisher/model"
	model2 "github.com/fikriahmadf/outbox-examples/internal/domain/internal_memo/model"
	"github.com/rs/zerolog/log"
)

func (s *InternalMemoServiceImpl) OutboxProcessor(ctx context.Context) error {
	tx, err := s.InternalMemoRepository.BeginTx(ctx)
	if err != nil {
		log.Warn().Err(err).Msg("[EmailOutboxService][OutboxProcessor] failed to begin transaction")
		return err
	}

	defer func() {
		if err != nil {
			errRollback := tx.Rollback(ctx)
			if errRollback != nil {
				log.Error().Err(errRollback).Msg("[EmailOutboxService][OutboxProcessor][Rollback] rollback transaction failed")
			}
			return
		}
		errCommit := tx.Commit(ctx)
		if errCommit != nil {
			log.Error().Err(errCommit).Msg("[EmailOutboxService][OutboxProcessor][Commit] commit transaction failed")
		}
	}()

	processedEmailOutboxes, err := tx.ResolvePendingEmailOutbox(ctx, 10)
	if err != nil {
		log.Error().Err(err).Msg("[EmailOutboxService][OutboxProcessor] failed to resolve pending email outbox")
		return err
	}

	for _, emailOutbox := range processedEmailOutboxes {

		var memoNotifPayload model.SendMemoNotifRequest
		if err = json.Unmarshal([]byte(emailOutbox.Payload), &memoNotifPayload); err != nil {
			log.Error().Err(err).Msg("[EmailOutboxService][OutboxProcessor] failed to unmarshal memo notif payload")
			return err
		}

		resp, err := s.ExternalNotifPublisherService.SendMemoNotif(ctx, memoNotifPayload)

		if err != nil {
			statusProcessor := model2.StatusPending
			if emailOutbox.RetryCount+1 >= 5 {
				statusProcessor = model2.StatusFailed
			}

			var errMessage string

			if resp.Message != "" {
				errMessage = resp.Message
			} else {
				errMessage = err.Error()
			}

			// update process
			emailOutbox.RetryCount = emailOutbox.RetryCount + 1
			emailOutbox.ErrorMessage = sql.NullString{String: errMessage, Valid: true}
			emailOutbox.Status = statusProcessor.String()
			emailOutbox.LastAttemptAt = sql.NullTime{Time: time.Now(), Valid: true}
			emailOutbox.MetaUpdatedAt = sql.NullTime{Time: time.Now(), Valid: true}
			if err := tx.UpdateErrorProcess(ctx, &emailOutbox); err != nil {
				log.Error().Err(err).Msg("[EmailOutboxService][OutboxProcessor] failed to update outbox process")
				return err
			}
			continue
		}

		emailOutbox.SentAt = sql.NullTime{Time: time.Now(), Valid: true}
		emailOutbox.LastAttemptAt = sql.NullTime{Time: time.Now(), Valid: true}
		emailOutbox.MetaUpdatedAt = sql.NullTime{Time: time.Now(), Valid: true}

		if err := tx.UpdateSentOutboxProcess(ctx, &emailOutbox); err != nil {
			log.Error().Err(err).Msg("[EmailOutboxService][OutboxProcessor] failed to update outbox process")
			return err
		}

	}
	return nil
}

type EmailOutboxService interface {
	OutboxProcessor(ctx context.Context) error
}
