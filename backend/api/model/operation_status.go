//go:generate go-enum --marshal
package model

type OperationStatus struct {
	State    OperationStatusState `json:"state"`
	Progress int                  `json:"progress"`
	Error    *HTTPError           `json:"error,omitempty"`
}

// ENUM(scheduled, in_progress, done, error)
type OperationStatusState string
