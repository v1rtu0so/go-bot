// pkg/raydium/parsers/event_parser.go

package parsers

import (
	"context"
	"log"

	"corvus_bot/pkg/raydium/listeners"
)

// EventParser handles parsing of raw pool events
type EventParser struct {
	listener *listeners.AMMPoolListener
}

// NewEventParser creates a new parser instance
func NewEventParser(listener *listeners.AMMPoolListener) *EventParser {
	return &EventParser{
		listener: listener,
	}
}

// Start begins processing events from the listener
func (p *EventParser) Start(ctx context.Context) {
	eventChan := p.listener.GetEventChannel()

	go func() {
		for {
			select {
			case <-ctx.Done():
				return
			case event, ok := <-eventChan:
				if !ok {
					return
				}
				p.parseEvent(event)
			}
		}
	}()
}

// parseEvent handles parsing of individual events
func (p *EventParser) parseEvent(event *listeners.RawPoolEvent) {
	// TODO: Implement parsing logic
	log.Printf("Received new pool event - Signature: %s, Slot: %d, BlockTime: %d",
		event.Signature, event.Slot, event.BlockTime)
}
