package package1

import (
	tsproto "google.golang.org/protobuf/types/known/timestamppb"
	"time"
)

// +Foo=true
// +Bar=123
type AnotherUser struct {
	// +ID=true
	ID          string            `json:"id"`
	DisplayName string            `json:"display_name"`
	Email       string            `json:"email"`
	Time        time.Time         `json:"time"`
	Duration    time.Duration     `json:"duration"`
	Timestamp   tsproto.Timestamp `json:"timestamp"`
}
