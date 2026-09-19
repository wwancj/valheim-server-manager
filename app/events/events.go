package events

// EventEmitter defines the interface for emitting events to connected clients.
// This replaces the Wails runtime EventsEmit with a WebSocket-based implementation.
type EventEmitter interface {
	Emit(event string, data interface{})
}

// NoopEmitter is a no-op implementation used when no emitter is configured.
type NoopEmitter struct{}

func (n *NoopEmitter) Emit(event string, data interface{}) {}
