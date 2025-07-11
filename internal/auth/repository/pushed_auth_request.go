package repository

import (
	"context"
	"time"

	"github.com/zitadel/oidc/v3/pkg/oidc"
)

type PushedAuthRequestRepository interface {
	CreatePushedAuthRequest(ctx context.Context, requestURI string, request *oidc.AuthRequest, expiresAt time.Time) error

	GetPushedAuthRequest(ctx context.Context, requestURI string) (*oidc.AuthRequest, error)
}
