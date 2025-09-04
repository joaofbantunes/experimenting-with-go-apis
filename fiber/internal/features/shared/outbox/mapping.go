package outbox

import (
	"encoding/json"
	"errors"
	"fmt"

	"github.com/joaofbantunes/experimenting-with-go-apis/fiber/internal/features/shared/domain"
)

func MapMsg(o *domain.Order, event any) (*OutboxMessage, error) {
	switch e := event.(type) {
	case *domain.OrderRegistered:
		return mapRegisteredMsg(o, e)
	case *domain.OrderCancelled:
		return mapCancelledMsg(o, e)
	default:
		return nil, errors.New(fmt.Sprintf("unsupported message type: %T", e))
	}
}

// in this case we don't need info from the aggregate, but in other cases we might, to populate the integration event
func mapRegisteredMsg(o *domain.Order, event *domain.OrderRegistered) (*OutboxMessage, error) {
	// serializing the domain event directly, but normally, as the integration event might be different, we would map to a different struct
	payload, err := json.Marshal(event)
	if err != nil {
		return nil, err
	}

	return NewOutboxMessage(
		"OrderRegistered",
		payload,
		event.OccurredAt,
		nil, // TODO: trace context
	), nil
}

// in this case we don't need info from the aggregate, but in other cases we might, to populate the integration event
func mapCancelledMsg(o *domain.Order, event *domain.OrderCancelled) (*OutboxMessage, error) {
	// serializing the domain event directly, but normally, as the integration event might be different, we would map to a different struct
	payload, err := json.Marshal(event)
	if err != nil {
		return nil, err
	}

	return NewOutboxMessage(
		"OrderCancelled",
		payload,
		event.OccurredAt,
		nil, // TODO: trace context
	), nil
}
