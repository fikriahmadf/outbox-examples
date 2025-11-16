package model

import (
	"time"

	"github.com/google/uuid"
)

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
		UpdatedAt:      time.Now(),
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
