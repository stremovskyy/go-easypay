package easypay

import "time"

type RecordedExchange struct {
	RequestID string            `json:"request_id"`
	Request   []byte            `json:"request"`
	Response  []byte            `json:"response"`
	Error     error             `json:"error,omitempty"`
	Tags      map[string]string `json:"tags"`
	Timestamp time.Time         `json:"timestamp"`
}
