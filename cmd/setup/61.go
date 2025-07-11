package setup

import (
	"context"
	"database/sql"
	_ "embed"

	"github.com/zitadel/zitadel/internal/eventstore"
)

var (
	//go:embed 61.sql
	createPushAuthRequestCache string
)

type PushedAuthRequestCache struct {
	dbClient *sql.DB
}

func (mig *PushedAuthRequestCache) Execute(ctx context.Context, _ eventstore.Event) error {
	stmt := createPushAuthRequestCache

	_, err := mig.dbClient.ExecContext(ctx, stmt)

	return err
}

func (mig *PushedAuthRequestCache) String() string {
	return "61_pushed_auth_request_cache"
}
