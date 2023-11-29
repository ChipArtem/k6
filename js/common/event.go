package common

import "github.com/ChipArtem/k6/event"

// Events are the event subscriber interfaces for the global event system, and
// the local (per-VU) event system.
type Events struct {
	Global, Local event.Subscriber
}
