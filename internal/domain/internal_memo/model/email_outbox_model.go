package model

import (
	"database/sql"
	"encoding/json"
	"time"

	"github.com/google/uuid"
)

type StatusType string

const (
	StatusPending StatusType = "pending"
	StatusSent    StatusType = "sent"
	StatusFailed  StatusType = "failed"
)

func (d StatusType) String() string {
	return string(d)
}

type EmailOutbox struct {
	ID             uuid.UUID       `json:"id" db:"id"`
	MemoID         uuid.UUID       `json:"memoId" db:"memo_id"`
	EventType      string          `json:"eventType" db:"event_type"`
	Payload        json.RawMessage `json:"payload" db:"payload"`
	RecipientEmail string          `json:"recipientEmail" db:"recipient_email"`
	Status         string          `json:"status" db:"status"`
	RetryCount     int             `json:"retryCount" db:"retry_count"`
	LastAttemptAt  sql.NullTime    `json:"lastAttemptAt" db:"last_attempt_at"`
	SentAt         sql.NullTime    `json:"sentAt" db:"sent_at"`
	ErrorMessage   sql.NullString  `json:"errorMessage" db:"error_message"`
	IdempotencyKey sql.NullString  `json:"idempotencyKey" db:"idempotency_key"`
	MetaCreatedAt  time.Time       `json:"metaCreatedAt" db:"meta_created_at"`
	MetaUpdatedAt  sql.NullTime    `json:"metaUpdatedAt" db:"meta_updated_at"`
}
