package model

import "github.com/google/uuid"

type Memo struct {
	ID                 string `json:"id" db:"id"`
	MemoNumberPrefix   string `json:"memoNumberPrefix" db:"memo_number_prefix"`
	MemoNumberSequence int    `json:"memoNumberSequence" db:"memo_number_sequence"`
	DepartmentCode     string `json:"departmentCode" db:"department_code"`
	Title              string `json:"title" db:"title"`
	Purpose            string `json:"purpose" db:"purpose"`
	CreatedAt          string `json:"createdAt" db:"created_at"`
	UpdatedAt          string `json:"updatedAt" db:"updated_at"`
}

func (m *Memo) ToNewModel() Memo {

	id, _ := uuid.NewV7()

	return Memo{
		ID:                 id,
		MemoNumberPrefix:   m.MemoNumberPrefix,
		MemoNumberSequence: m.MemoNumberSequence,
		DepartmentCode:     m.DepartmentCode,
		Title:              m.Title,
		Purpose:            m.Purpose,
		CreatedAt:          m.CreatedAt,
		UpdatedAt:          m.UpdatedAt,
	}
}
