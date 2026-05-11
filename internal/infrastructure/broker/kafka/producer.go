package kafka

import (
	"context"
	"errors"
	"fmt"
	"strings"
	"time"

	"api/internal/domain/event"

	kafkago "github.com/segmentio/kafka-go"
)

// Producer publishes domain events to a single configured topic. It is
// safe for concurrent use; the underlying *kafkago.Writer batches and
// retries internally.
type Producer struct {
	writer *kafkago.Writer
	topic  string
}

// NewProducer builds a Writer with sane defaults. Brokers is a comma- or
// space-separated host:port list, matching the PHP DSN format.
func NewProducer(brokers string, topic string) (*Producer, error) {
	if topic == "" {
		return nil, errors.New("kafka: topic is required")
	}
	hosts := splitBrokers(brokers)
	if len(hosts) == 0 {
		return nil, errors.New("kafka: at least one broker is required")
	}

	w := &kafkago.Writer{
		Addr:         kafkago.TCP(hosts...),
		Topic:        topic,
		Balancer:     &kafkago.Hash{},
		RequiredAcks: kafkago.RequireAll,
		Async:        false,
		BatchTimeout: 50 * time.Millisecond,
	}
	return &Producer{writer: w, topic: topic}, nil
}

// Ensure compile-time compliance with domain event bus.
var _ event.Bus = (*Producer)(nil)

// Publish encodes the event and writes one message to the topic.
func (p *Producer) Publish(e event.DomainEvent) error {
	key, body, err := Encode(e)
	if err != nil {
		return fmt.Errorf("encode event: %w", err)
	}
	// 5s upper bound is generous for a single publish; tune via ctx if
	// you propagate cancellation from a request scope.
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	return p.writer.WriteMessages(ctx, kafkago.Message{
		Key:   []byte(key),
		Value: body,
	})
}

// Close flushes the writer.
func (p *Producer) Close() error { return p.writer.Close() }

func splitBrokers(raw string) []string {
	parts := strings.FieldsFunc(raw, func(r rune) bool {
		return r == ',' || r == ' ' || r == ';'
	})
	out := make([]string, 0, len(parts))
	for _, p := range parts {
		if p = strings.TrimSpace(p); p != "" {
			out = append(out, p)
		}
	}
	return out
}
