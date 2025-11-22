package model

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"time"

	"github.com/fikriahmadf/outbox-examples/external/domain/notif_publisher/model"
	"github.com/google/uuid"
)

type MemoEventType int

const (
	MemoEventUnknown MemoEventType = iota
	MemoEventCreated
)

func (e MemoEventType) String() string {
	return []string{"memo.event.unknown", "memo.event.created"}[e]
}

type CreateMemoRequest struct {
	DepartmentCode string `json:"departmentCode" db:"department_code"`
	Title          string `json:"title" db:"title"`
	Purpose        string `json:"purpose" db:"purpose"`
}

func (m *CreateMemoRequest) ToNewModel() Memo {

	id, _ := uuid.NewV7()

	return Memo{
		ID:             id,
		DepartmentCode: m.DepartmentCode,
		Title:          m.Title,
		Purpose:        m.Purpose,
		CreatedAt:      time.Now(),
	}
}

type Memo struct {
	ID                 uuid.UUID `json:"id,omitempty" db:"id"`
	MemoNumberPrefix   string    `json:"memoNumberPrefix" db:"memo_number_prefix"`
	MemoNumberSequence int       `json:"memoNumberSequence" db:"memo_number_sequence"`
	DepartmentCode     string    `json:"departmentCode" db:"department_code"`
	Title              string    `json:"title" db:"title"`
	Purpose            string    `json:"purpose" db:"purpose"`
	CreatedAt          time.Time `json:"createdAt,omitempty" db:"created_at"`
	UpdatedAt          time.Time `json:"updatedAt,omitempty" db:"updated_at"`
}

func (m *Memo) GetIdempotencyKeyForOutbox(eventType MemoEventType, recipientEmail string) string {

	// format: memo.event.created-uuid-Recepient
	return fmt.Sprintf("%s-%s-%s", eventType.String(), m.ID.String(), recipientEmail)
}

func (m *Memo) GetMemoNumber() string {

	return fmt.Sprintf("%s%d", m.MemoNumberPrefix, m.MemoNumberSequence)
}

func (m *Memo) GetPayloadForOutbox(eventType MemoEventType, recipientEmail string) json.RawMessage {

	switch eventType {
	case MemoEventUnknown:
		return json.RawMessage(`{}`)
	case MemoEventCreated:
		payload := model.SendMemoNotifRequest{
			RecipientEmail: recipientEmail,
			MemoId:         m.ID.String(),
			MemoTitle:      m.Title,
			CreatedDate:    m.CreatedAt.Format("2006-01-02 15:04:05"),
			SubjectEmail:   fmt.Sprintf("Internal Memo Notification - New Memo Created - %s - %s", m.Title, m.GetMemoNumber()),
		}

		marsharledPayload, _ := json.Marshal(payload)
		return marsharledPayload
	}

	return json.RawMessage(`{}`)
}

func (m *Memo) ToOutboxModel(eventType MemoEventType, recipientEmail string) EmailOutbox {

	var idempotencyKey string

	switch eventType {
	case MemoEventUnknown:
		idempotencyKey = ""
	case MemoEventCreated:
		idempotencyKey = m.GetIdempotencyKeyForOutbox(eventType, recipientEmail)
	}

	id, _ := uuid.NewV7()

	return EmailOutbox{
		ID:             id,
		MemoID:         m.ID,
		EventType:      eventType.String(),
		Payload:        m.GetPayloadForOutbox(eventType, recipientEmail),
		RecipientEmail: recipientEmail,
		Status:         StatusPending.String(),
		RetryCount:     0,
		IdempotencyKey: sql.NullString{String: idempotencyKey, Valid: idempotencyKey != ""},
		MetaCreatedAt:  time.Now(),
	}
}
