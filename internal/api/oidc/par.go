package oidc

import (
	"context"
	"time"

	"github.com/zitadel/logging"
	"github.com/zitadel/oidc/v3/pkg/oidc"
	"github.com/zitadel/oidc/v3/pkg/op"

	"github.com/zitadel/zitadel/internal/telemetry/tracing"
)

const PARDefaultLifetime = 1 * time.Minute

type PARConfig struct {
	Lifetime time.Duration
}

// toOPConfig converts DeviceAuthorizationConfig to a [op.DeviceAuthorizationConfig],
// setting sane defaults for empty values.
// Safe to call when c is nil.
func (c *PARConfig) toOPConfig() op.PushedAuthorizationRequestConfig {
	out := op.PushedAuthorizationRequestConfig{
		Lifetime: PARDefaultLifetime,
	}

	if c == nil {
		return out
	}

	if c.Lifetime != 0 {
		out.Lifetime = c.Lifetime
	}

	return out
}

func (o *OPStorage) StorePAR(
	ctx context.Context, requestURI string, request *oidc.AuthRequest, expiresAt time.Time,
) (err error) {
	const logMsg = "store pushed authorization request"
	logger := logging.WithFields("request_uri", requestURI, "expires", expiresAt)

	ctx, span := tracing.NewSpan(ctx)
	defer func() {
		logger.OnError(err).Error(logMsg)
		span.EndWithError(err)
	}()

	err = o.repo.CreatePushedAuthRequest(ctx, requestURI, request, expiresAt)
	if err == nil {
		logger.SetFields("request", request).Debug(logMsg)
	}

	return err
}

func (o *OPStorage) GetPARState(ctx context.Context, requestURI string) (_ *oidc.AuthRequest, err error) {
	const logMsg = "load pushed authorization request"
	logger := logging.WithFields("request_uri", requestURI)

	ctx, span := tracing.NewSpan(ctx)
	defer func() {
		err = oidcError(err)
		logger.OnError(err).Error(logMsg)
		span.EndWithError(err)
	}()

	request, err := o.repo.GetPushedAuthRequest(ctx, requestURI)
	if err != nil {
		return nil, err
	}

	logger.SetFields("request", request).Debug(logMsg)

	return request, nil
}
