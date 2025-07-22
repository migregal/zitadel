package project

import (
	"context"

	"github.com/zitadel/zitadel/internal/eventstore"
	"github.com/zitadel/zitadel/internal/repository/metadata"
	"github.com/zitadel/zitadel/internal/zerrors"
)

const (
	ApplicationMetadataSetType        = applicationEventTypePrefix + metadata.SetEventType
	ApplicationMetadataRemovedType    = applicationEventTypePrefix + metadata.RemovedEventType
	ApplicationMetadataRemovedAllType = applicationEventTypePrefix + metadata.RemovedAllEventType
)

type ApplicationMetadataSetEvent struct {
	metadata.SetEvent

	AppID string `json:"appId,omitempty"`
}

func (e *ApplicationMetadataSetEvent) Payload() any {
	return e
}

func NewApplicationMetadataSetEvent(ctx context.Context, aggregate *eventstore.Aggregate, appID, key string, value []byte) *ApplicationMetadataSetEvent {
	return &ApplicationMetadataSetEvent{
		SetEvent: *metadata.NewSetEvent(
			eventstore.NewBaseEventForPush(
				ctx,
				aggregate,
				ApplicationMetadataSetType),
			key,
			value),
		AppID: appID,
	}
}

func ApplicationMetadataSetEventMapper(event eventstore.Event) (eventstore.Event, error) {
	me, err := metadata.SetEventMapper(event)
	if err != nil {
		return nil, err
	}

	e := &ApplicationMetadataSetEvent{
		SetEvent: *me,
	}

	err = event.Unmarshal(e)
	if err != nil {
		return nil, zerrors.ThrowInternal(err, "APPLICATION-1m9e3", "unable to unmarshal application")
	}

	return e, nil
}

type ApplicationMetadataRemovedEvent struct {
	metadata.RemovedEvent

	AppID string `json:"appId,omitempty"`
}

func (e *ApplicationMetadataRemovedEvent) Payload() any {
	return e
}

func NewApplicationMetadataRemovedEvent(
	ctx context.Context, aggregate *eventstore.Aggregate, appID, key string,
) *ApplicationMetadataRemovedEvent {
	return &ApplicationMetadataRemovedEvent{
		RemovedEvent: *metadata.NewRemovedEvent(
			eventstore.NewBaseEventForPush(
				ctx,
				aggregate,
				ApplicationMetadataRemovedType),
			key),
		AppID: appID,
	}
}

func ApplicationMetadataRemovedEventMapper(event eventstore.Event) (eventstore.Event, error) {
	me, err := metadata.RemovedEventMapper(event)
	if err != nil {
		return nil, err
	}

	e := &ApplicationMetadataRemovedEvent{
		RemovedEvent: *me,
	}

	err = event.Unmarshal(e)
	if err != nil {
		return nil, zerrors.ThrowInternal(err, "APPLICATION-1m9e3", "unable to unmarshal application")
	}

	return e, nil
}

type ApplicationMetadataRemovedAllEvent struct {
	metadata.RemovedAllEvent

	AppID string `json:"appId,omitempty"`
}

func NewApplicationMetadataRemovedAllEvent(ctx context.Context, aggregate *eventstore.Aggregate) *ApplicationMetadataRemovedAllEvent {
	return &ApplicationMetadataRemovedAllEvent{
		RemovedAllEvent: *metadata.NewRemovedAllEvent(
			eventstore.NewBaseEventForPush(
				ctx,
				aggregate,
				ApplicationMetadataRemovedAllType),
		),
	}
}

func ApplicationMetadataRemovedAllEventMapper(event eventstore.Event) (eventstore.Event, error) {
	me, err := metadata.RemovedAllEventMapper(event)
	if err != nil {
		return nil, err
	}

	e := &ApplicationMetadataRemovedAllEvent{
		RemovedAllEvent: *me,
	}

	err = event.Unmarshal(e)
	if err != nil {
		return nil, zerrors.ThrowInternal(err, "APPLICATION-1m9e3", "unable to unmarshal application")
	}

	return e, nil
}
