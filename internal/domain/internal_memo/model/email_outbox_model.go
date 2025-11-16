package model

import (
	"database/sql"
	"encoding/json"
	"time"
)

type EmailOutbox struct {
	ID               string          `json:"id" db:"id"`
	MemoID           string          `json:"memoId" db:"memo_id"`
	EventType        string          `json:"eventType" db:"event_type"`
	Payload          json.RawMessage `json:"payload" db:"payload"`
	RecipientEmail   string          `json:"recipientEmail" db:"recipient_email"`
	NotificationType string          `json:"notificationType" db:"notification_type"`
	Status           string          `json:"status" db:"status"`
	RetryCount       int             `json:"retryCount" db:"retry_count"`
	LastAttemptAt    sql.NullTime    `json:"lastAttemptAt" db:"last_attempt_at"`
	SentAt           sql.NullTime    `json:"sentAt" db:"sent_at"`
	ErrorMessage     sql.NullString  `json:"errorMessage" db:"error_message"`
	IdempotencyKey   sql.NullString  `json:"idempotencyKey" db:"idempotency_key"`
	MetaCreatedAt    time.Time       `json:"metaCreatedAt" db:"meta_created_at"`
	MetaUpdatedAt    sql.NullTime    `json:"metaUpdatedAt" db:"meta_updated_at"`
}
