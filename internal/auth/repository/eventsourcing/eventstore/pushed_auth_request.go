package eventstore

import (
	"context"
	"fmt"
	"time"

	"github.com/zitadel/oidc/v3/pkg/oidc"

	"github.com/zitadel/zitadel/internal/auth/repository"
	cache "github.com/zitadel/zitadel/internal/pushed_auth_request/repository"
	"github.com/zitadel/zitadel/internal/telemetry/tracing"
)

var _ repository.PushedAuthRequestRepository = (*PushedAuthRequestRepo)(nil)

type PushedAuthRequestRepo struct {
	PushedAuthRequests cache.PushedAuthRequestCache
}

func (repo *PushedAuthRequestRepo) Health(ctx context.Context) error {
	return repo.PushedAuthRequests.Health(ctx)
}

func (repo *PushedAuthRequestRepo) CreatePushedAuthRequest(
	ctx context.Context, requestURI string, request *oidc.AuthRequest, expiresAt time.Time,
) (err error) {
	ctx, span := tracing.NewSpan(ctx)
	defer func() { span.EndWithError(err) }()

	return repo.PushedAuthRequests.SavePushedAuthRequest(ctx, requestURI, request, expiresAt)
}

func (repo *PushedAuthRequestRepo) GetPushedAuthRequest(ctx context.Context, requestURI string) (_ *oidc.AuthRequest, err error) {
	ctx, span := tracing.NewSpan(ctx)
	defer func() { span.EndWithError(err) }()

	request, err := repo.PushedAuthRequests.GetPushedAuthRequestByRequestURI(ctx, requestURI)
	if err != nil {
		return nil, fmt.Errorf("get pushed auth request: %w", err)
	}

	err = repo.PushedAuthRequests.DeletePushedAuthRequest(ctx, requestURI)
	if err != nil {
		return nil, fmt.Errorf("delete pushed auth request: %w", err)
	}

	return request, nil
}
