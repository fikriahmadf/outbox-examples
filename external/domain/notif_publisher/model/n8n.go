package model

type SendMemoNotifRequest struct {
	RecipientEmail string `json:"recipientEmail"`
	MemoId         string `json:"memoId"`
	MemoTitle      string `json:"memoTitle"`
	CreatedDate    string `json:"createdDate"`
}

type SendMemoNotifResponse struct {
	Success bool   `json:"success"`
	Message string `json:"message"`
}
