package response

type Base struct {
	Data      *interface{} `json:"data,omitempty"`
	Metadata  *interface{} `json:"metadata,omitempty"`
	Error     *string      `json:"error,omitempty"`
	Message   *string      `json:"message,omitempty"`
	RequestID *string      `json:"requestId,omitempty"`
}
