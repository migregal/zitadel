package command

import (
	"context"

	"github.com/zitadel/zitadel/internal/domain"
	"github.com/zitadel/zitadel/internal/eventstore"
	"github.com/zitadel/zitadel/internal/repository/project"
	"github.com/zitadel/zitadel/internal/zerrors"
)

func (c *Commands) SetApplicationMetadata(
	ctx context.Context, projectID, appID, resourceOwner string, metadata *domain.Metadata,
) (_ *domain.Metadata, err error) {
	if projectID == "" || appID == "" {
		return nil, zerrors.ThrowInvalidArgument(nil, "META-983dF", "Errors.IDMissing")
	}

	_, err = c.checkApplicationExists(ctx, projectID, appID, resourceOwner)
	if err != nil {
		return nil, err
	}

	setMetadata := NewApplicationMetadataWriteModel(projectID, appID, resourceOwner, metadata.Key)
	projectAgg := ProjectAggregateFromWriteModel(&setMetadata.WriteModel)
	event, err := c.setApplicationMetadata(ctx, projectAgg, appID, metadata)
	if err != nil {
		return nil, err
	}

	pushedEvents, err := c.eventstore.Push(ctx, event)
	if err != nil {
		return nil, err
	}

	err = AppendAndReduce(setMetadata, pushedEvents...)
	if err != nil {
		return nil, err
	}

	return writeModelToApplicationMetadata(setMetadata), nil
}

func (c *Commands) BulkSetApplicationMetadata(
	ctx context.Context, projectID, appID, resourceOwner string, metadatas ...*domain.Metadata,
) (_ *domain.ObjectDetails, err error) {
	if projectID == "" || appID == "" {
		return nil, zerrors.ThrowInvalidArgument(nil, "META-983dF", "Errors.IDMissing")
	}

	if len(metadatas) == 0 {
		return nil, zerrors.ThrowPreconditionFailed(nil, "META-9mm2d", "Errors.Metadata.NoData")
	}

	_, err = c.checkApplicationExists(ctx, projectID, appID, resourceOwner)
	if err != nil {
		return nil, err
	}

	events := make([]eventstore.Command, len(metadatas))
	setMetadata := NewApplicationMetadataListWriteModel(projectID, appID, resourceOwner)
	projectAgg := ProjectAggregateFromWriteModel(&setMetadata.WriteModel)
	for i, data := range metadatas {
		event, err := c.setApplicationMetadata(ctx, projectAgg, appID, data)
		if err != nil {
			return nil, err
		}
		events[i] = event
	}

	pushedEvents, err := c.eventstore.Push(ctx, events...)
	if err != nil {
		return nil, err
	}

	err = AppendAndReduce(setMetadata, pushedEvents...)
	if err != nil {
		return nil, err
	}
	return writeModelToObjectDetails(&setMetadata.WriteModel), nil
}

func (c *Commands) setApplicationMetadata(
	ctx context.Context, appAgg *eventstore.Aggregate, appID string, metadata *domain.Metadata,
) (command eventstore.Command, err error) {
	if !metadata.IsValid() {
		return nil, zerrors.ThrowInvalidArgument(nil, "META-2ml0f", "Errors.Metadata.Invalid")
	}

	return project.NewApplicationMetadataSetEvent(
		ctx,
		appAgg,
		appID,
		metadata.Key,
		metadata.Value,
	), nil
}

func (c *Commands) RemoveApplicationMetadata(
	ctx context.Context, projectID, appID, resourceOwner, metadataKey string,
) (_ *domain.ObjectDetails, err error) {
	if projectID == "" || appID == "" {
		return nil, zerrors.ThrowInvalidArgument(nil, "META-983dF", "Errors.IDMissing")
	}

	if metadataKey == "" {
		return nil, zerrors.ThrowInvalidArgument(nil, "META-2n0f1", "Errors.Metadata.Invalid")
	}

	_, err = c.checkApplicationExists(ctx, projectID, appID, resourceOwner)
	if err != nil {
		return nil, err
	}

	removeMetadata, err := c.getApplicationMetadataModelByID(
		ctx, projectID, appID, resourceOwner, metadataKey,
	)
	if err != nil {
		return nil, err
	}
	if !removeMetadata.State.Exists() {
		return nil, zerrors.ThrowNotFound(nil, "META-mcnw3", "Errors.Metadata.NotFound")
	}

	projectAgg := ProjectAggregateFromWriteModel(&removeMetadata.WriteModel)
	event, err := c.removeApplicationMetadata(ctx, projectAgg, appID, metadataKey)
	if err != nil {
		return nil, err
	}

	pushedEvents, err := c.eventstore.Push(ctx, event)
	if err != nil {
		return nil, err
	}

	err = AppendAndReduce(removeMetadata, pushedEvents...)
	if err != nil {
		return nil, err
	}

	return writeModelToObjectDetails(&removeMetadata.WriteModel), nil
}

func (c *Commands) BulkRemoveApplicationMetadata(
	ctx context.Context, projectID, appID, resourceOwner string, metadataKeys ...string,
) (_ *domain.ObjectDetails, err error) {
	if projectID == "" || appID == "" {
		return nil, zerrors.ThrowInvalidArgument(nil, "META-983dF", "Errors.IDMissing")
	}

	if len(metadataKeys) == 0 {
		return nil, zerrors.ThrowPreconditionFailed(nil, "META-9mw2d", "Errors.Metadata.NoData")
	}

	_, err = c.checkApplicationExists(ctx, projectID, appID, resourceOwner)
	if err != nil {
		return nil, err
	}

	events := make([]eventstore.Command, len(metadataKeys))
	removeMetadata, err := c.getApplicationMetadataListModel(ctx, projectID, appID, resourceOwner)
	if err != nil {
		return nil, err
	}

	projectAgg := ProjectAggregateFromWriteModel(&removeMetadata.WriteModel)
	for i, key := range metadataKeys {
		if key == "" {
			return nil, zerrors.ThrowInvalidArgument(nil, "COMMAND-m19ds", "Errors.Metadata.Invalid")
		}
		if _, found := removeMetadata.metadataList[key]; !found {
			return nil, zerrors.ThrowNotFound(nil, "META-2npds", "Errors.Metadata.KeyNotExisting")
		}
		event, err := c.removeApplicationMetadata(ctx, projectAgg, appID, key)
		if err != nil {
			return nil, err
		}
		events[i] = event
	}

	pushedEvents, err := c.eventstore.Push(ctx, events...)
	if err != nil {
		return nil, err
	}

	err = AppendAndReduce(removeMetadata, pushedEvents...)
	if err != nil {
		return nil, err
	}

	return writeModelToObjectDetails(&removeMetadata.WriteModel), nil
}

func (c *Commands) removeApplicationMetadata(
	ctx context.Context, appAgg *eventstore.Aggregate, appID, metadataKey string,
) (command eventstore.Command, err error) {
	return project.NewApplicationMetadataRemovedEvent(
		ctx,
		appAgg,
		appID,
		metadataKey,
	), nil
}

func (c *Commands) getApplicationMetadataModelByID(
	ctx context.Context, projectID, appID, resourceOwner, key string,
) (*ApplicationMetadataWriteModel, error) {
	appMetadataWriteModel := NewApplicationMetadataWriteModel(projectID, appID, resourceOwner, key)

	err := c.eventstore.FilterToQueryReducer(ctx, appMetadataWriteModel)
	if err != nil {
		return nil, err
	}

	return appMetadataWriteModel, nil
}

func (c *Commands) getApplicationMetadataListModel(
	ctx context.Context, projectID, appID, resourceOwner string,
) (*ApplicationMetadataListWriteModel, error) {
	appMetadataWriteModel := NewApplicationMetadataListWriteModel(projectID, appID, resourceOwner)

	err := c.eventstore.FilterToQueryReducer(ctx, appMetadataWriteModel)
	if err != nil {
		return nil, err
	}

	return appMetadataWriteModel, nil
}

func writeModelToApplicationMetadata(wm *ApplicationMetadataWriteModel) *domain.Metadata {
	return &domain.Metadata{
		ObjectRoot: writeModelToObjectRoot(wm.WriteModel),
		Key:        wm.Key,
		Value:      wm.Value,
		State:      wm.State,
	}
}
