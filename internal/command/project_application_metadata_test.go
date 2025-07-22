package command

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/zitadel/zitadel/internal/domain"
	"github.com/zitadel/zitadel/internal/eventstore"
	"github.com/zitadel/zitadel/internal/eventstore/v1/models"
	"github.com/zitadel/zitadel/internal/repository/project"
	"github.com/zitadel/zitadel/internal/zerrors"
)

func TestCommandSide_SetApplicationMetadata(t *testing.T) {
	t.Parallel()

	type fields struct {
		eventstore *eventstore.Eventstore
	}
	type (
		args struct {
			ctx           context.Context
			projectID     string
			appID         string
			resourceOwner string
			metadata      *domain.Metadata
		}
	)
	type res struct {
		want *domain.Metadata
		err  func(error) bool
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		res    res
	}{
		{
			name: "app not existing, pre condition error",
			fields: fields{
				eventstore: eventstoreExpect(
					t,
					expectFilter(),
				),
			},
			args: args{
				ctx:           context.Background(),
				projectID:     "prj1",
				appID:         "app1",
				resourceOwner: "org1",
				metadata: &domain.Metadata{
					Key:   "key",
					Value: []byte("value"),
				},
			},
			res: res{
				err: zerrors.IsPreconditionFailed,
			},
		},
		{
			name: "invalid metadata, pre condition error",
			fields: fields{
				eventstore: eventstoreExpect(
					t,
					expectFilter(
						eventFromEventPusher(
							project.NewApplicationAddedEvent(context.Background(),
								&project.NewAggregate("prj1", "org1").Aggregate,
								"app1",
								"key",
							),
						),
					),
				),
			},
			args: args{
				ctx:           context.Background(),
				projectID:     "prj1",
				appID:         "app1",
				resourceOwner: "org1",
				metadata: &domain.Metadata{
					Key: "key",
				},
			},
			res: res{
				err: zerrors.IsErrorInvalidArgument,
			},
		},
		{
			name: "add metadata, ok",
			fields: fields{
				eventstore: eventstoreExpect(
					t,
					expectFilter(
						eventFromEventPusher(
							project.NewApplicationAddedEvent(context.Background(),
								&project.NewAggregate("prj1", "org1").Aggregate,
								"app1",
								"key",
							),
						),
					),
					expectPush(
						project.NewApplicationMetadataSetEvent(context.Background(),
							&project.NewAggregate("prj1", "org1").Aggregate,
							"app1",
							"key",
							[]byte("value"),
						),
					),
				),
			},
			args: args{
				ctx:           context.Background(),
				projectID:     "prj1",
				appID:         "app1",
				resourceOwner: "org1",
				metadata: &domain.Metadata{
					Key:   "key",
					Value: []byte("value"),
				},
			},
			res: res{
				want: &domain.Metadata{
					ObjectRoot: models.ObjectRoot{
						AggregateID:   "prj1", // TODO: fix this behaviour in next commits.
						ResourceOwner: "org1",
					},
					Key:   "key",
					Value: []byte("value"),
					State: domain.MetadataStateActive,
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			r := &Commands{
				eventstore: tt.fields.eventstore,
			}
			got, err := r.SetApplicationMetadata(
				tt.args.ctx, tt.args.projectID, tt.args.appID, tt.args.resourceOwner, tt.args.metadata,
			)
			if tt.res.err == nil {
				assert.NoError(t, err)
			}
			if tt.res.err != nil && !tt.res.err(err) {
				t.Errorf("got wrong err: %v ", err)
			}
			if tt.res.err == nil {
				assert.Equal(t, tt.res.want, got)
			}
		})
	}
}

func TestCommandSide_BulkSetApplicationMetadata(t *testing.T) {
	t.Parallel()

	type fields struct {
		eventstore *eventstore.Eventstore
	}
	type (
		args struct {
			ctx           context.Context
			projectID     string
			appID         string
			resourceOwner string
			metadataList  []*domain.Metadata
		}
	)
	type res struct {
		want *domain.ObjectDetails
		err  func(error) bool
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		res    res
	}{
		{
			name: "empty meta data list, pre condition error",
			fields: fields{
				eventstore: eventstoreExpect(
					t,
				),
			},
			args: args{
				ctx:           context.Background(),
				projectID:     "prj1",
				appID:         "app1",
				resourceOwner: "org1",
			},
			res: res{
				err: zerrors.IsPreconditionFailed,
			},
		},
		{
			name: "app not existing, pre condition error",
			fields: fields{
				eventstore: eventstoreExpect(
					t,
					expectFilter(),
				),
			},
			args: args{
				ctx:           context.Background(),
				projectID:     "prj1",
				appID:         "app1",
				resourceOwner: "org1",
				metadataList: []*domain.Metadata{
					{Key: "key", Value: []byte("value")},
					{Key: "key1", Value: []byte("value1")},
				},
			},
			res: res{
				err: zerrors.IsPreconditionFailed,
			},
		},
		{
			name: "invalid metadata, pre condition error",
			fields: fields{
				eventstore: eventstoreExpect(
					t,
					expectFilter(
						eventFromEventPusher(
							project.NewApplicationAddedEvent(context.Background(),
								&project.NewAggregate("prj1", "org1").Aggregate,
								"app1", "app",
							),
						),
					),
				),
			},
			args: args{
				ctx:           context.Background(),
				projectID:     "prj1",
				appID:         "app1",
				resourceOwner: "org1",
				metadataList: []*domain.Metadata{
					{Key: "key"},
					{Key: "key1"},
				},
			},
			res: res{
				err: zerrors.IsErrorInvalidArgument,
			},
		},
		{
			name: "add metadata, ok",
			fields: fields{
				eventstore: eventstoreExpect(
					t,
					expectFilter(
						eventFromEventPusher(
							project.NewApplicationAddedEvent(context.Background(),
								&project.NewAggregate("prj1", "org1").Aggregate,
								"app1", "app",
							),
						),
					),
					expectPush(
						project.NewApplicationMetadataSetEvent(context.Background(),
							&project.NewAggregate("prj1", "org1").Aggregate,
							"app1",
							"key",
							[]byte("value"),
						),
						project.NewApplicationMetadataSetEvent(context.Background(),
							&project.NewAggregate("prj1", "org1").Aggregate,
							"app1",
							"key1",
							[]byte("value1"),
						),
					),
				),
			},
			args: args{
				ctx:           context.Background(),
				projectID:     "prj1",
				appID:         "app1",
				resourceOwner: "org1",
				metadataList: []*domain.Metadata{
					{Key: "key", Value: []byte("value")},
					{Key: "key1", Value: []byte("value1")},
				},
			},
			res: res{
				want: &domain.ObjectDetails{
					ResourceOwner: "org1",
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			r := &Commands{
				eventstore: tt.fields.eventstore,
			}
			got, err := r.BulkSetApplicationMetadata(
				tt.args.ctx, tt.args.projectID, tt.args.appID, tt.args.resourceOwner, tt.args.metadataList...,
			)
			if tt.res.err == nil {
				assert.NoError(t, err)
			}
			if tt.res.err != nil && !tt.res.err(err) {
				t.Errorf("got wrong err: %v ", err)
			}
			if tt.res.err == nil {
				assertObjectDetails(t, tt.res.want, got)
			}
		})
	}
}

func TestCommandSide_ApplicationRemoveMetadata(t *testing.T) {
	type fields struct {
		eventstore *eventstore.Eventstore
	}
	type (
		args struct {
			ctx           context.Context
			projectID     string
			appID         string
			resourceOwner string
			metadataKey   string
		}
	)
	type res struct {
		want *domain.ObjectDetails
		err  func(error) bool
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		res    res
	}{
		{
			name: "app not existing, pre condition error",
			fields: fields{
				eventstore: eventstoreExpect(
					t,
					expectFilter(),
				),
			},
			args: args{
				ctx:           context.Background(),
				projectID:     "prj1",
				appID:         "app1",
				resourceOwner: "org1",
				metadataKey:   "key",
			},
			res: res{
				err: zerrors.IsPreconditionFailed,
			},
		},
		{
			name: "invalid metadata, pre condition error",
			fields: fields{
				eventstore: eventstoreExpect(
					t,
				),
			},
			args: args{
				ctx:           context.Background(),
				projectID:     "prj1",
				appID:         "app1",
				resourceOwner: "org1",
				metadataKey:   "",
			},
			res: res{
				err: zerrors.IsErrorInvalidArgument,
			},
		},
		{
			name: "empty app id, pre condition error",
			fields: fields{
				eventstore: eventstoreExpect(
					t,
				),
			},
			args: args{
				ctx:           context.Background(),
				projectID:     "prj1",
				appID:         "",
				resourceOwner: "org1",
				metadataKey:   "key",
			},
			res: res{
				err: zerrors.IsErrorInvalidArgument,
			},
		},
		{
			name: "meta data not existing, not found error",
			fields: fields{
				eventstore: eventstoreExpect(
					t,
					expectFilter(
						eventFromEventPusher(
							project.NewApplicationAddedEvent(context.Background(),
								&project.NewAggregate("prj1", "org1").Aggregate,
								"app1", "app",
							),
						),
					),
					expectFilter(),
				),
			},
			args: args{
				ctx:           context.Background(),
				projectID:     "prj1",
				appID:         "app1",
				resourceOwner: "org1",
				metadataKey:   "key",
			},
			res: res{
				err: zerrors.IsNotFound,
			},
		},
		{
			name: "remove metadata, ok",
			fields: fields{
				eventstore: eventstoreExpect(
					t,
					expectFilter(
						eventFromEventPusher(
							project.NewApplicationAddedEvent(context.Background(),
								&project.NewAggregate("prj1", "org1").Aggregate,
								"app1", "app",
							),
						),
					),
					expectFilter(
						eventFromEventPusher(
							project.NewApplicationMetadataSetEvent(context.Background(),
								&project.NewAggregate("prj1", "org1").Aggregate,
								"app1",
								"key",
								[]byte("value"),
							),
						),
					),
					expectPush(
						project.NewApplicationMetadataRemovedEvent(context.Background(),
							&project.NewAggregate("prj1", "org1").Aggregate,
							"app1",
							"key",
						),
					),
				),
			},
			args: args{
				ctx:           context.Background(),
				projectID:     "prj1",
				appID:         "app1",
				resourceOwner: "org1",
				metadataKey:   "key",
			},
			res: res{
				want: &domain.ObjectDetails{
					ResourceOwner: "org1",
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := &Commands{
				eventstore: tt.fields.eventstore,
			}
			got, err := r.RemoveApplicationMetadata(
				tt.args.ctx, tt.args.projectID, tt.args.appID, tt.args.resourceOwner, tt.args.metadataKey,
			)
			if tt.res.err == nil {
				assert.NoError(t, err)
			}
			if tt.res.err != nil && !tt.res.err(err) {
				t.Errorf("got wrong err: %v ", err)
			}
			if tt.res.err == nil {
				assertObjectDetails(t, tt.res.want, got)
			}
		})
	}
}

func TestCommandSide_BulkRemoveApplicationMetadata(t *testing.T) {
	type fields struct {
		eventstore *eventstore.Eventstore
	}
	type (
		args struct {
			ctx           context.Context
			projectID     string
			appID         string
			resourceOwner string
			metadataList  []string
		}
	)
	type res struct {
		want *domain.ObjectDetails
		err  func(error) bool
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		res    res
	}{
		{
			name: "empty meta data list, pre condition error",
			fields: fields{
				eventstore: eventstoreExpect(
					t,
				),
			},
			args: args{
				ctx:           context.Background(),
				projectID:     "prj1",
				appID:         "app1",
				resourceOwner: "org1",
			},
			res: res{
				err: zerrors.IsPreconditionFailed,
			},
		},
		{
			name: "app not existing, pre condition error",
			fields: fields{
				eventstore: eventstoreExpect(
					t,
					expectFilter(),
				),
			},
			args: args{
				ctx:           context.Background(),
				projectID:     "prj1",
				appID:         "app1",
				resourceOwner: "org1",
				metadataList:  []string{"key", "key1"},
			},
			res: res{
				err: zerrors.IsPreconditionFailed,
			},
		},
		{
			name: "remove metadata keys not existing, precondition error",
			fields: fields{
				eventstore: eventstoreExpect(
					t,
					expectFilter(
						eventFromEventPusher(
							project.NewApplicationAddedEvent(context.Background(),
								&project.NewAggregate("prj1", "org1").Aggregate,
								"app1", "app",
							),
						),
					),
					expectFilter(
						eventFromEventPusher(
							project.NewApplicationMetadataSetEvent(context.Background(),
								&project.NewAggregate("prj1", "org1").Aggregate,
								"app1",
								"key",
								[]byte("value"),
							),
						),
					),
				),
			},
			args: args{
				ctx:           context.Background(),
				projectID:     "prj1",
				appID:         "app1",
				resourceOwner: "org1",
				metadataList:  []string{"key", "key1"},
			},
			res: res{
				err: zerrors.IsNotFound,
			},
		},
		{
			name: "invalid metadata, pre condition error",
			fields: fields{
				eventstore: eventstoreExpect(
					t,
					expectFilter(
						eventFromEventPusher(
							project.NewApplicationAddedEvent(context.Background(),
								&project.NewAggregate("prj1", "org1").Aggregate,
								"app1", "app",
							),
						),
					),
					expectFilter(
						eventFromEventPusher(
							project.NewApplicationMetadataSetEvent(context.Background(),
								&project.NewAggregate("prj1", "org1").Aggregate,
								"app1",
								"key",
								[]byte("value"),
							),
						),
						eventFromEventPusher(
							project.NewApplicationMetadataSetEvent(context.Background(),
								&project.NewAggregate("prj1", "org1").Aggregate,
								"app1",
								"key1",
								[]byte("value1"),
							),
						),
					),
				),
			},
			args: args{
				ctx:           context.Background(),
				projectID:     "prj1",
				appID:         "app1",
				resourceOwner: "org1",
				metadataList:  []string{""},
			},
			res: res{
				err: zerrors.IsErrorInvalidArgument,
			},
		},
		{
			name: "remove metadata, ok",
			fields: fields{
				eventstore: eventstoreExpect(
					t,
					expectFilter(
						eventFromEventPusher(
							project.NewApplicationAddedEvent(context.Background(),
								&project.NewAggregate("prj1", "org1").Aggregate,
								"app1", "app",
							),
						),
					),
					expectFilter(
						eventFromEventPusher(
							project.NewApplicationMetadataSetEvent(context.Background(),
								&project.NewAggregate("prj1", "org1").Aggregate,
								"app1",
								"key",
								[]byte("value"),
							),
						),
						eventFromEventPusher(
							project.NewApplicationMetadataSetEvent(context.Background(),
								&project.NewAggregate("prj1", "org1").Aggregate,
								"app1",
								"key1",
								[]byte("value1"),
							),
						),
					),
					expectPush(
						project.NewApplicationMetadataRemovedEvent(context.Background(),
							&project.NewAggregate("prj1", "org1").Aggregate,
							"app1",
							"key",
						),
						project.NewApplicationMetadataRemovedEvent(context.Background(),
							&project.NewAggregate("prj1", "org1").Aggregate,
							"app1",
							"key1",
						),
					),
				),
			},
			args: args{
				ctx:           context.Background(),
				projectID:     "prj1",
				appID:         "app1",
				resourceOwner: "org1",
				metadataList:  []string{"key", "key1"},
			},
			res: res{
				want: &domain.ObjectDetails{
					ResourceOwner: "org1",
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := &Commands{
				eventstore: tt.fields.eventstore,
			}
			got, err := r.BulkRemoveApplicationMetadata(
				tt.args.ctx, tt.args.projectID, tt.args.appID, tt.args.resourceOwner, tt.args.metadataList...,
			)
			if tt.res.err == nil {
				assert.NoError(t, err)
			}
			if tt.res.err != nil && !tt.res.err(err) {
				t.Errorf("got wrong err: %v ", err)
			}
			if tt.res.err == nil {
				assertObjectDetails(t, tt.res.want, got)
			}
		})
	}
}
