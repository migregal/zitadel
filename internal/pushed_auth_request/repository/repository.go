package repository

import (
	"context"

	"github.com/zitadel/oidc/v3/pkg/oidc"
)

type PushedAuthRequestCache interface {
	Health(ctx context.Context) error

	GetPushedAuthRequestByRequestURI(ctx context.Context, requestURI string) (*oidc.AuthRequest, error)
	SavePushedAuthRequest(ctx context.Context, requestURI string, request *oidc.AuthRequest) error
	CachePushedAuthRequest(ctx context.Context, requestURI string, oidcrequest *oidc.AuthRequest)
	DeletePushedAuthRequest(ctx context.Context, requestURI string) error
}
