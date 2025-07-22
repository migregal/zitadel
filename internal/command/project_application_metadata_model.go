package command

import (
	"github.com/zitadel/zitadel/internal/eventstore"
	"github.com/zitadel/zitadel/internal/repository/project"
)

type ApplicationMetadataWriteModel struct {
	MetadataWriteModel

	AppID string
}

func NewApplicationMetadataWriteModel(projectID, appID, resourceOwner, key string) *ApplicationMetadataWriteModel {
	return &ApplicationMetadataWriteModel{
		MetadataWriteModel: MetadataWriteModel{
			WriteModel: eventstore.WriteModel{
				AggregateID:   projectID,
				ResourceOwner: resourceOwner,
			},
			Key: key,
		},
		AppID: appID,
	}
}

func (wm *ApplicationMetadataWriteModel) AppendEvents(events ...eventstore.Event) {
	for _, event := range events {
		switch e := event.(type) {
		case *project.ApplicationMetadataSetEvent:
			if e.AppID != wm.AppID {
				continue
			}

			wm.MetadataWriteModel.AppendEvents(&e.SetEvent)
		case *project.ApplicationMetadataRemovedEvent:
			if e.AppID != wm.AppID {
				continue
			}

			wm.MetadataWriteModel.AppendEvents(&e.RemovedEvent)
		case *project.ApplicationMetadataRemovedAllEvent:
			if e.AppID != wm.AppID {
				continue
			}

			wm.MetadataWriteModel.AppendEvents(&e.RemovedAllEvent)
		}
	}
}

func (wm *ApplicationMetadataWriteModel) Query() *eventstore.SearchQueryBuilder {
	return eventstore.NewSearchQueryBuilder(eventstore.ColumnsEvent).
		ResourceOwner(wm.ResourceOwner).
		AddQuery().
		AggregateIDs(wm.MetadataWriteModel.AggregateID).
		AggregateTypes(project.AggregateType).
		EventTypes(
			project.ApplicationMetadataSetType,
			project.ApplicationMetadataRemovedType,
			project.ApplicationMetadataRemovedAllType).
		Builder()
}

type ApplicationMetadataListWriteModel struct {
	MetadataListWriteModel

	AppID string
}

func NewApplicationMetadataListWriteModel(projectID, appID, resourceOwner string) *ApplicationMetadataListWriteModel {
	return &ApplicationMetadataListWriteModel{
		MetadataListWriteModel: MetadataListWriteModel{
			WriteModel: eventstore.WriteModel{
				AggregateID:   projectID,
				ResourceOwner: resourceOwner,
			},
			metadataList: make(map[string][]byte),
		},
		AppID: appID,
	}
}

func (wm *ApplicationMetadataListWriteModel) AppendEvents(events ...eventstore.Event) {
	for _, event := range events {
		switch e := event.(type) {
		case *project.ApplicationMetadataSetEvent:
			if e.AppID != wm.AppID {
				continue
			}

			wm.MetadataListWriteModel.AppendEvents(&e.SetEvent)
		case *project.ApplicationMetadataRemovedEvent:
			if e.AppID != wm.AppID {
				continue
			}

			wm.MetadataListWriteModel.AppendEvents(&e.RemovedEvent)
		case *project.ApplicationMetadataRemovedAllEvent:
			if e.AppID != wm.AppID {
				continue
			}

			wm.MetadataListWriteModel.AppendEvents(&e.RemovedAllEvent)
		}
	}
}

func (wm *ApplicationMetadataListWriteModel) Reduce() error {
	return wm.MetadataListWriteModel.Reduce()
}

func (wm *ApplicationMetadataListWriteModel) Query() *eventstore.SearchQueryBuilder {
	return eventstore.NewSearchQueryBuilder(eventstore.ColumnsEvent).
		ResourceOwner(wm.ResourceOwner).
		AddQuery().
		AggregateIDs(wm.MetadataListWriteModel.AggregateID).
		AggregateTypes(project.AggregateType).
		EventTypes(
			project.ApplicationMetadataSetType,
			project.ApplicationMetadataRemovedType,
			project.ApplicationMetadataRemovedAllType).
		Builder()
}
